package assets

import (
	"embed"
)

//go:embed repo/*
var Assets embed.FS
