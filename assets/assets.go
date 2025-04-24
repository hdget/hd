package assets

import (
	"embed"
)

//go:embed repo/* db/*
var Assets embed.FS
