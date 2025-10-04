package pkg

import (
	"html/template"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

type user struct {
	Username string
	Password []byte
	Fullname string
	Email    string
}

type service struct {
	templateCache map[string]*template.Template
	isProduction  bool
	router        *gin.Engine
	logger        zerolog.Logger
	usersMu       sync.RWMutex
	users         map[string]*user
	basePath      string
	httpServer    *http.Server
	// db         *sql.DB
}
