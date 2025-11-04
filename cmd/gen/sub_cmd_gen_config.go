package gen

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/hdget/hd/g"
	"github.com/hdget/hd/pkg/utils"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
)

type configApp struct {
	Name         string
	ExternalPort int
}

type configInput struct {
	Project  string
	Env      string
	RepoHost string
	Apps     []configApp
}

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
    # project name, override HD_NAMESPACE environment variable
    name = "{{.Project}}"
    # running environment
    env = "{{.Env}}"
    # when use --all argument, app start in order, app stop execute reversely
    apps = [{{range $index, $item := .Apps}}{{if $index}},{{end}}"{{ $item.Name }}"{{end}}]

#[[apps]]
#   # app name
#   name = "example-app"
#   # git repo url
#   repo = "https://{{.RepoHost}}/example-app.git"
#	# if external port is set, then app is exposed
#	# external_port = 1000
#	# app plugins 
#	[[apps.plugins]]
#       name = "example-plugin"
#       url = "https://{{.RepoHost}}/{{.Project}}/plugin/example-plugin.git"
{{range $index, $item := .Apps}}
[[apps]]
    name = "{{ .Name }}"
    repo = "https://{{$.RepoHost}}/{{$.Project}}/backend/{{ .Name }}.git"
	{{- if gt .ExternalPort 0}}
	external_port = {{ .ExternalPort }}
	{{- end}}
{{end}}
[[repos]]
    # config repository
    name = "config"
    # repository url
    url = "https://{{.RepoHost}}/{{.Project}}/common/config.git"

[[repos]]
    # proto repository
    name = "proto"
    # repository url
    url = "https://{{.RepoHost}}/{{.Project}}/common/proto.git"

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
	appInput := utils.GetInput("Please input app names", "core,gateway:1000,usercenter")
	repoHost := utils.GetInput("Please input repository host", "github.com")

	appList := make([]configApp, 0)
	apps := strings.Split(appInput, ",")
	for _, appStr := range apps {
		parts := strings.Split(appStr, ":")
		switch len(parts) {
		case 1:
			appList = append(appList, configApp{
				Name: parts[0],
			})
		case 2:
			appList = append(appList, configApp{
				Name:         parts[0],
				ExternalPort: cast.ToInt(parts[1]),
			})
		}
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

	err = tpl.Execute(f, &configInput{
		Project:  project,
		Env:      env,
		RepoHost: repoHost,
		Apps:     appList,
	})
	if err != nil {
		utils.Fatal("error execute template", err)
	}
}
