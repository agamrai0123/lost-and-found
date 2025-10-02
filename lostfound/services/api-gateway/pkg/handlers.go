package pkg

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *service) home() gin.HandlerFunc {
	return func(c *gin.Context) {
		data := map[string]any{
			"Title":  "Home - LostAndFound",
			"Active": "home",
		}
		s.renderTemplate(c, "home", data)
	}
}

func (s *service) about() gin.HandlerFunc {
	return func(c *gin.Context) {
		data := map[string]any{
			"Title":  "About - LostAndFound",
			"Active": "about",
		}
		s.renderTemplate(c, "about", data)
	}
}

func (s *service) login() gin.HandlerFunc {
	return func(c *gin.Context) {
		data := map[string]any{
			"Title":  "Login - LostAndFound",
			"Active": "login",
		}
		s.renderTemplate(c, "login", data)
	}
}

func (s *service) logout(c *gin.Context) {
	c.Redirect(http.StatusSeeOther, "/")
}
