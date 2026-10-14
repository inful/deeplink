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
	"maps"
	"text/template"

	"github.com/inful/deeplink/gitlab"
	"github.com/inful/deeplink/outlook"
	"github.com/inful/deeplink/teams"
)

// Funcmap combines the FuncMap from every registered sub-package into a
// single template.FuncMap. Sub-package keys never collide (each is
// namespaced by its system, e.g. "gitlabIssue", "teamsNewChat",
// "outlookEvent"); if they ever do, the last writer wins and the
// top-level TestNoFuncNameCollisions test fails.
func Funcmap() template.FuncMap {
	out := template.FuncMap{}
	for _, src := range []template.FuncMap{
		gitlab.FuncMap(),
		teams.FuncMap(),
		outlook.FuncMap(),
	} {
		maps.Copy(out, src)
	}
	return out
}
