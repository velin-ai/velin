// Package ui embeds the built frontend so the API binary can serve it.
//
// Run "npm run build" inside ui/api before "go build"; the dist directory is
// produced by Vite and must exist for the embed to succeed.
package ui

import "embed"

//go:embed all:dist
var Dist embed.FS
