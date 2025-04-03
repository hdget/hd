package tools

type sqlboilerTool struct {
	*toolImpl
}

func Sqlboiler() Tool {
	return &sqlboilerTool{
		toolImpl: &toolImpl{
			name: "sqlboiler",
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
		"github.com/volatiletech/sqlboiler/v4@latest",
		"github.com/volatiletech/sqlboiler/v4/drivers/sqlboiler-mysql@latest",
		"go install github.com/volatiletech/sqlboiler/v4/drivers/sqlboiler-sqlite3@latest",
	}
	for _, repo := range repos {
		if err := AllPlatform().GoInstall(repo); err != nil {
			return err
		}
	}
	return nil
}
