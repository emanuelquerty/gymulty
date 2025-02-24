package gymulty

import "embed"

//go:embed postgres/migrations/*.sql
var EmbedMigrations embed.FS
