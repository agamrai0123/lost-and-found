package pkg

import (
	"fmt"
	"html/template"

	"github.com/gin-gonic/gin"
)

type service struct {
	templateCache map[string]*template.Template
	isProduction  bool
}

func newService() *service {
	tc, err := CreateTemplateCache()
	if err != nil {
		panic(err)
	}
	return &service{
		templateCache: tc,
		isProduction:  true,
	}
}

func StartService() {
	gin.SetMode(gin.DebugMode) // Change to gin.DebugMode for development

	service := newService()
	router := gin.Default()

	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	router.Static("/static", "./web/static")
	router.LoadHTMLGlob("./web/templates/*.tmpl")
	service.routes(router)

	fmt.Println("Starting server on :8080")
	if err := router.Run(":8080"); err != nil {
		fmt.Printf("Server failed to start: %v\n", err)
	}
}
