package gen

import (
	"fmt"
	"github.com/hdget/hd/g"
	"github.com/hdget/hd/pkg/utils"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"text/template"
)

var (
	subCmdGenConfig = &cobra.Command{
		Use:   "config",
		Short: "new project config",
		Run: func(cmd *cobra.Command, args []string) {
			genConfig()
		},
	}
)

const (
	configTemplate = `[project]
    # project name, will be used as HD_NAMESPACE environment variable
    name = "{{.Project.Name}}"
    # running environment
    env = "{{.Project.Env}}"
    # gateway app listen port
    gateway_port = {{.Project.GatewayPort}}
    # app lists with order when using '--all' flag
    apps = [{{range $index, $app := .Project.Apps}}{{if $index}},{{end}}"{{ $app }}"{{end}}]

# repos
#[[repos]]
#    # ususally it is the same as app name
#    name = "example_repo"
#    # git repo url
#    url = "https://github.com/repo/example"

# 3rd party tools
#[[tools]]
#    # tool name
#    name = "example_tool"
#    # tool version
#    version = "1.0"
#    # tool windows binary release url
#    url_win_release = ""
#    # tool linux binary release url
#    url_linux_release = ""`
)

func genConfig() {
	baseDir, _ := os.Getwd()
	possibleProject := filepath.Base(baseDir)

	project := utils.GetInput("Please input project name", possibleProject)
	env := utils.GetInput("Please input running environment", "test")

	exampleConfig := &g.RootConfig{
		Project: g.ProjectConfig{
			Name:        project,
			Env:         env,
			GatewayPort: g.DefaultGatewayPort,
			Apps:        []string{},
		},
	}

	tpl, err := template.New("toml").Parse(configTemplate)
	if err != nil {
		utils.Fatal("error parse template", err)
	}

	if utils.ExistsFile(g.ConfigFile) {
		fmt.Printf("%s exists, automatically saved as %s.bak\n", g.ConfigFile, g.ConfigFile)
		_ = os.Rename(g.ConfigFile, g.ConfigFile+".bak")
	}

	f, err := os.Create(g.ConfigFile)
	if err != nil {
		utils.Fatal("error create config file", err)
	}
	defer func() {
		_ = f.Close()
	}()

	err = tpl.Execute(f, exampleConfig)
	if err != nil {
		utils.Fatal("error execute template", err)
	}
}
