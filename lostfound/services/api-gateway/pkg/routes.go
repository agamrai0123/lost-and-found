package pkg

import "github.com/gin-gonic/gin"

func (s *service) routes(r *gin.Engine) {
	api := r.Group("/api")
	v1 := api.Group("/v1")
	v1.GET("/", s.home())
	v1.GET("/about", s.about())
	v1.Any("/login", s.login())
	v1.POST("/logout", s.logout)
}
