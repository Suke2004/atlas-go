// Package migrations provides embedded SQL migration scripts for Goose.
package migrations

import "embed"

// FS embeds all .sql migration files.
//
//go:embed *.sql
var FS embed.FS
