package pkg

import (
	"html/template"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

type service struct {
	templateCache map[string]*template.Template
	isProduction  bool
	router        *gin.Engine
	logger        zerolog.Logger
	usersMu       sync.RWMutex
	users         map[string]*user
	basePath      string
	// db            *sql.DB
}

func newService() *service {
	if err := readConfiguration(); err != nil {
		panic(err)
	}
	tc, err := createTemplateCache()
	if err != nil {
		panic(err)
	}

	// Initialize Gin router
	gin.SetMode(gin.DebugMode)
	router := gin.Default()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	router.Static("/static", "./web/static")
	router.LoadHTMLGlob("./web/templates/*.tmpl")

	return &service{
		templateCache: tc,
		isProduction:  true,
		router:        router,
		logger:        getLogger(AppConfig.ServiceName),
		users:         make(map[string]*user),
		basePath:      "/api/v1",
	}
}

func StartService() {
	s := newService()
	home := s.router.Group(s.basePath)
	s.routes(home)

	s.logger.Info().Msg("Starting server on :8080")
	if err := s.router.Run(":8080"); err != nil {
		s.logger.Fatal().Err(err).Msg("Server failed to start")
	}
}
