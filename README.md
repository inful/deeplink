# deeplink

A small Go library and set of Hugo shortcodes for building deep links
into GitLab, Microsoft Teams, Outlook, and other systems. Designed to
be embedded in [docbuilder](https://github.com/inful/docbuilder) and any
other documentation pipeline that needs to surface "create a thing" or
"open a thing" links.

## What you get

- **A Go library** that produces deep-link URLs as plain strings, with
  typed `Params` structs for each supported operation.
- **A `template.FuncMap`** exposing every builder to Go's `text/template`
  and `html/template` engines (including Hugo's). Use it from docbuilder
  templates, from Hugo's `hugo.toml` `params` evaluation, or anywhere
  else that consumes a `template.FuncMap`.
- **Hugo shortcode HTML templates** embedded in the module via
  `go:embed`, ready to copy into a site's `layouts/shortcodes/`.

## Install

```sh
go get github.com/inful/deeplink
```

## Go API

Every operation takes a typed `Params` struct and returns a URL string
and an error.

```go
import "github.com/inful/deeplink/gitlab"
import "github.com/inful/deeplink/teams"
import "github.com/inful/deeplink/outlook"

// GitLab: pre-filled new work item (issue / incident)
u, err := gitlab.Issue(gitlab.IssueParams{
    BaseURL:     "https://gitlab.com/my-group/my-project",
    Title:       "Login is broken",
    Description: "Steps to reproduce:\n1. ...",
})

// Microsoft Teams: open a new chat
u, err := teams.NewChat(teams.ChatParams{
    TenantID: "f4c2b8e1-7a3d-4e6f-9b12-8c5d0e3a1f47",
    Users:    []string{"joe@contoso.com", "bob@contoso.com"},
    Message:  "Quick question about the rollout",
})

// Outlook: compose a calendar event
u, err := outlook.Event(outlook.EventParams{
    Subject:   "Sprint review",
    StartDate: "2024-01-15T10:00:00.000",
    EndDate:   "2024-01-15T11:00:00.000",
    Location:  "Teams",
})
```

## Use as template functions

Register the entire FuncMap with your template engine in one call:

```go
import "github.com/inful/deeplink"

tpl := template.New("page").Funcs(deeplink.Funcmap())
// ... or pick the system you need:
tpl := template.New("page").Funcs(gitlab.FuncMap())
tpl := template.New("page").Funcs(teams.FuncMap())
tpl := template.New("page").Funcs(outlook.FuncMap())
```

The functions take a `map[string]any` (so they accept the output of
Hugo's `dict` and any plain Go `map`) and return a URL string.

```go-html-template
{{ $u := gitlabIssue (dict
    "BaseURL"     "https://gitlab.com/my-group/my-project"
    "Title"       "Login is broken"
    "Description" "Steps to reproduce:\n1. ..."
) }}
<a href="{{ $u }}">File an issue</a>
```

## Hugo shortcodes

The `shortcodes/` package embeds six shortcode HTML files via
`go:embed`. Copy them into a Hugo site's `layouts/shortcodes/` and
they are ready to use:

| Shortcode | System | Purpose |
|---|---|---|
| `deeplink-gitlab-issue` | GitLab | Pre-filled new work item |
| `deeplink-gitlab-mergerequest` | GitLab | Pre-filled new merge request |
| `deeplink-teams-new-chat` | Teams | Open a new chat compose window |
| `deeplink-teams-channel` | Teams | Open an existing channel |
| `deeplink-outlook-event` | Outlook | Compose a new calendar event |
| `deeplink-outlook-email` | Outlook | `mailto:` compose window |

Example usage in a markdown page:

```markdown
{{< deeplink-gitlab-issue
    baseURL="https://gitlab.com/my-group/my-project"
    title="Login is broken"
    description="Steps to reproduce: 1. Open 2. Click 3. Crash"
    text="Report a bug" >}}
```

To copy the embedded shortcodes into a Hugo site from a Go build tool:

```go
import "github.com/inful/deeplink/shortcodes"

for _, name := range shortcodes.Names() {
    content, _ := shortcodes.Files.ReadFile(name)
    dst := filepath.Join(siteRoot, "layouts/shortcodes", name)
    os.WriteFile(dst, content, 0o644)
}
```

### URL encoding

The Go library and the shortcodes produce **byte-identical URLs**. The
shortcode templates were written to match the Go output exactly:

- Spaces are `%20`, matching the format the vendor documentation uses.
- Namespaced keys keep their brackets — `issue[title]`, not
  `issue%5Btitle%5D`.
- Empty values are omitted entirely.
- Values written literally (the Teams `users` list, `mailto:` recipients)
  keep their commas.

A test in `shortcodes/render_test.go` renders every shortcode with Hugo
and compares the result against the Go function it mirrors, so the two
cannot drift apart silently. It is skipped when Hugo is not installed.

### A note on parallel maintenance

Each shortcode file is a hand-written Hugo template that mirrors the
logic in the corresponding Go function. If you change one, change the
other — the rendering test will tell you if you forget. The doc
comments at the top of each shortcode state which Go function it
mirrors.

### A note on untrusted input

Some parameters are written into the URL without encoding, because the
destination system expects them literally: the Teams `users` list and
the `mailto:` recipient lists. A value containing `&`, `?`, `#` or
whitespace could otherwise add parameters of its own — for example
turning one recipient into a blind copy. The Go library rejects such
values, and so do the shortcodes.

## Architecture

```
deeplink/                  # top-level: combined Funcmap()
├── deeplink.go
├── internal/q/             # shared query/path encoding helpers
├── internal/tmpl/          # shared template-value coercion helpers
├── gitlab/                 # GitLab sub-package (self-contained)
│   ├── issue.go
│   ├── mergerequest.go
│   ├── funcmap.go
│   └── *_test.go
├── teams/                  # Microsoft Teams sub-package
│   ├── chat.go
│   ├── channel.go
│   ├── funcmap.go
│   └── *_test.go
├── outlook/                # Outlook sub-package
│   ├── event.go
│   ├── email.go
│   ├── funcmap.go
│   └── *_test.go
└── shortcodes/             # embedded Hugo shortcode HTML
    ├── embed.go            # go:embed *.html
    ├── deeplink-gitlab-issue.html
    ├── deeplink-gitlab-mergerequest.html
    ├── deeplink-teams-new-chat.html
    ├── deeplink-teams-channel.html
    ├── deeplink-outlook-event.html
    └── deeplink-outlook-email.html
```

Shortcode file names match the shortcode names exactly: Hugo looks up
`{{< deeplink-gitlab-issue >}}` in `deeplink-gitlab-issue.html`, so a
file named anything else would silently never be found.

### Adding a new system

1. Create a new sub-package (e.g. `jira/`) with one Go file per
   operation, each exporting a typed function and (optionally) a
   template-friendly adapter.
2. Add a `func FuncMap() template.FuncMap` to the sub-package that
   returns the namespaced FuncMap for your functions.
3. Add one line to the top-level `Funcmap()` in `deeplink.go` to
   include your FuncMap in the combined result.
4. Add one shortcode HTML file in `shortcodes/`, named after the
   shortcode it defines, that mirrors the Go logic. The `//go:embed`
   directive will pick it up automatically.

That's it — no other plumbing needed. `TestFuncmapExposesExactlyTheExpectedFunctions`
in `deeplink_test.go` will tell you to add your function to the expected
list, and `Funcmap()` panics if two sub-packages claim the same name.

## URL reference

These are the canonical URL formats this library produces. They are
the formats documented by each vendor; refer to the vendor docs for
full parameter lists.

- **GitLab issue (work item)**: `https://<host>/<group>/<project>/-/work_items/new?issue[title]=…&issue[description]=…`
  See: <https://docs.gitlab.com/user/project/issues/create_issues/#using-a-url-with-prefilled-values>
- **GitLab merge request**: `https://<host>/<group>/<project>/-/merge_requests/new?merge_request[title]=…&merge_request[source_branch]=…&merge_request[target_branch]=…`
- **Microsoft Teams new chat**: `https://teams.microsoft.com/l/chat/0/0?tenantId=…&users=user1@x.com,user2@x.com&message=…`
  See: <https://learn.microsoft.com/en-us/microsoftteams/platform/concepts/build-and-test/deep-link-teams>
- **Microsoft Teams channel**: `https://teams.microsoft.com/l/channel/<channelId>/<channelName>?groupId=…&tenantId=…`
- **Outlook calendar event**: `https://outlook.office.com/calendar/deeplink/compose?subject=…&startdt=…&enddt=…&body=…&location=…`
- **Outlook email**: `mailto:user@x.com?cc=…&bcc=…&subject=…&body=…`

## Development

```sh
go test ./...          # run all tests
go test -cover ./...   # with coverage
go vet ./...
```

## License

MIT. See [LICENSE](LICENSE).
