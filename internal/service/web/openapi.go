package web

import (
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/thirdfort/thirdfort-go-code-review/internal"
	"github.com/wI2L/fizz/openapi"
)

func (s *WebService) setOpenAPI() {
	infos := &openapi.Info{
		Title:       "Consumer API",
		Description: `Consumer API web service.`,
		Version:     internal.Version,
	}

	var auth []*openapi.SecurityRequirement
	var authval openapi.SecurityRequirement = map[string][]string{"authentication": {"bearerAuth", "deviceFingerprint"}}
	auth = append(auth, &authval)

	s.Fizz.Generator().SetSecurityRequirement(auth)
	s.Fizz.Generator().SetSecuritySchemes(map[string]*openapi.SecuritySchemeOrRef{
		"bearerAuth": {
			SecurityScheme: &openapi.SecurityScheme{
				Type:         "http",
				Scheme:       "bearer",
				BearerFormat: "JWT",
			},
		},
		"deviceFingerprint": {
			SecurityScheme: &openapi.SecurityScheme{
				In:   "header",
				Name: "device-fingerprint",
			},
		},
	})

	s.Fizz.GET("/openapi.yaml", nil, s.Fizz.OpenAPI(infos, "yaml"))
	s.Fizz.GET("/openapi.json", nil, s.Fizz.OpenAPI(infos, "json"))
	s.Fizz.GET("/swagger/*any", nil,
		ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.URL("/openapi.yaml")))

}
