// Package deeplink builds deep links into external systems — GitLab,
// Microsoft Teams, Outlook, and others — and exposes them as a
// template.FuncMap for use with Go's text/template and html/template
// engines (including Hugo's template engine).
//
// Each sub-package (gitlab, teams, outlook, …) is self-contained and
// exports its own FuncMap(). The top-level Funcmap() function returns
// the combined FuncMap of every registered sub-package — convenient
// when you want the entire deep-link surface wired up in one call.
//
// To add a new system, create a sub-package that exposes a FuncMap()
// returning a template.FuncMap, then add one line to deeplink.go to
// register it. Every other consumer of this library picks it up
// automatically.
package deeplink

import (
	"fmt"
	"text/template"

	"github.com/inful/deeplink/gitlab"
	"github.com/inful/deeplink/outlook"
	"github.com/inful/deeplink/teams"
)

// Funcmap combines the FuncMap from every registered sub-package into a
// single template.FuncMap.
//
// It panics if two sub-packages claim the same function name. Each name
// is namespaced by its system (e.g. "gitlabIssue", "teamsNewChat",
// "outlookEvent"), so a collision means two sub-packages were written
// against the same prefix — a bug in this library that would otherwise
// silently shadow one function with another.
func Funcmap() template.FuncMap {
	return merge(
		gitlab.FuncMap(),
		teams.FuncMap(),
		outlook.FuncMap(),
	)
}

// merge combines template FuncMaps into one, panicking if a name is
// declared by more than one source.
func merge(sources ...template.FuncMap) template.FuncMap {
	out := template.FuncMap{}
	for _, src := range sources {
		for name, fn := range src {
			if _, exists := out[name]; exists {
				panic(fmt.Sprintf("deeplink: duplicate template function %q", name))
			}
			out[name] = fn
		}
	}
	return out
}
