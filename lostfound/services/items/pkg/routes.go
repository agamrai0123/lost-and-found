package pkg

import (
	"github.com/gin-gonic/gin"
)

func (s *ItemService) routes(r *gin.RouterGroup) {
	r.GET("/posts", s.listPosts())
	r.GET("/posts/:id", s.getPost())
	r.POST("/posts", s.createPost())
	// r.PUT("/posts/:id", s.updatePost())
	// r.DELETE("/posts/:id", s.deletePost())
}
