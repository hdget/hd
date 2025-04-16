package sourcecode

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

type routeAnnotation struct {
	Endpoint      string   `json:"endpoint"`      // endpoint
	Methods       []string `json:"methods"`       // http methods
	Origin        string   `json:"origin"`        // 请求来源
	IsRawResponse bool     `json:"isRawResponse"` // 是否返回原始消息
	IsPublic      bool     `json:"isPublic"`      // 是否是公共路由
	Permissions   []string `json:"permissions"`   // 对应的权限列表
}

type Route struct {
	ModuleName    string
	Handler       string
	Endpoint      string
	HttpMethods   []string
	AllowOrigin   string
	IsPublic      int32
	IsRawResponse int32
	Permissions   []string
	Comments      []string
}

const (
	routeAnnotationName = "route"
)

var (
	regModuleName = regexp.MustCompile(`^[vV]([0-9]+)_([a-zA-Z0-9]+)`)
)

// generate 生成相关文件，包括route.json
func (impl *sourceCodeHandlerImpl) generate(scInfo *sourceCodeInfo) error {
	routeItems := make([]*Route, 0)

	for _, h := range scInfo.daprInvocationHandlers {
		for _, m := range scInfo.daprModules {
			if m.pkgRelPath == h.pkgRelPath && m.name == h.moduleName {
				for annKind, annValue := range h.annotations {
					if annKind == routeAnnotationName && strings.TrimSpace(annValue) != "" {
						var ann routeAnnotation
						err := json.Unmarshal([]byte(annValue), &ann)
						if err != nil {
							return err
						}

						// 设置初始值
						routeItem := &Route{
							ModuleName:  h.moduleName,
							Handler:     h.alias,
							Comments:    h.comments,
							Endpoint:    ann.Endpoint,
							HttpMethods: []string{"GET"},
							AllowOrigin: ann.Origin,
						}

						if ann.IsPublic {
							routeItem.IsPublic = 1
						}

						if ann.IsRawResponse {
							routeItem.IsRawResponse = 1
						}

						if len(ann.Methods) > 0 {
							routeItem.HttpMethods = ann.Methods
						}

						routeItems = append(routeItems, routeItem)
					}

				}
			}
		}
	}

	fmt.Println(routeItems)

	return nil
}
