package pkg

import (
	"github.com/gin-gonic/gin"
)

func (s *service) routes(r *gin.RouterGroup) {
	r.GET("/", s.home())
	r.GET("/about", s.about())
	r.GET("/login", s.loginGet())
	r.POST("/login", s.loginPost())
	r.GET("/signup", s.signupGet())
	r.POST("/signup", s.signupPost())
}
