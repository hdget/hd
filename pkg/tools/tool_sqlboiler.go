package tools

import "github.com/hdget/hd/g"

type sqlboilerTool struct {
	*toolImpl
}

func Sqlboiler() Tool {
	return &sqlboilerTool{
		toolImpl: &toolImpl{
			&g.ToolConfig{
				Name:    "sqlboiler",
				Version: "4.19.5",
			},
		},
	}
}

func (t *sqlboilerTool) IsInstalled() bool {
	return t.success("sqlboiler --version")
}

func (t *sqlboilerTool) LinuxInstall() error {
	return t.install()
}

func (t *sqlboilerTool) WindowsInstall() error {
	return t.install()
}

func (t *sqlboilerTool) install() error {
	repos := []string{
		"github.com/aarondl/sqlboiler/v4@latest",
		"github.com/aarondl/sqlboiler/v4/drivers/sqlboiler-mysql@latest",
		"github.com/aarondl/sqlboiler/v4/drivers/sqlboiler-psql@latest",
		"github.com/aarondl/sqlboiler/v4/drivers/sqlboiler-sqlite3@latest",
	}
	for _, repo := range repos {
		if err := AllPlatform().GoInstall(repo); err != nil {
			return err
		}
	}
	return nil
}
