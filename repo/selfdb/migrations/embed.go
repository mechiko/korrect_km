package migrations

import "embed"

//go:embed sqlite
var Sqlite embed.FS

//go:embed mssql
var Mssql embed.FS
