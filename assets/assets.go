package assets

import (
	"embed"
)

//go:embed repo/* db/* sql/*
var Manager embed.FS
