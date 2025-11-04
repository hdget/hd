package assets

import (
	"embed"
)

//go:embed repo/*
var Manager embed.FS
