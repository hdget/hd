package assets

import (
	"embed"
)

//go:embed repo/* db/*
var Manager embed.FS
