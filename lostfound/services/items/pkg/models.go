package pkg

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

type PostStatus string

const (
	StatusLost    PostStatus = "lost"
	StatusFound   PostStatus = "found"
	StatusClaimed PostStatus = "claimed"
)

type ContactInfo struct {
	Name  string `json:"name,omitempty"`  // display name for contact person
	Phone string `json:"phone,omitempty"` // E.164 preferred (e.g. +919876543210)
	Email string `json:"email,omitempty"` // RFC-5322 / standard email string
}

type Post struct {
	ID         string      `json:"id"`                    // UUID string (unique post id)
	Title      string      `json:"title"`                 // short title/summary of the post
	ImageURL   string      `json:"image_url,omitempty"`   // URL to an image (optional)
	AuthorName string      `json:"author_name,omitempty"` // public display name of the author (optional)
	CreatedBy  string      `json:"created_by"`            // creator user ID (UUID)
	CreatedAt  time.Time   `json:"created_at"`            // creation timestamp (RFC3339)
	UpdatedAt  *time.Time  `json:"updated_at,omitempty"`  // last update timestamp (nullable)
	Status     PostStatus  `json:"status"`                // one of: "lost", "found", "claimed"
	Location   string      `json:"location,omitempty"`    // human-readable location text
	Details    string      `json:"details,omitempty"`     // full description / notes
	Contact    ContactInfo `json:"contact,omitempty"`     // contact information object
}

type ItemService struct {
	storePath  string
	imageDir   string
	router     *gin.Engine
	mu         sync.RWMutex
	posts      map[string]*Post
	logger     zerolog.Logger
	basePath   string
	httpServer *http.Server
}
type createPostRequest struct {
	Title      string      `json:"title" binding:"required"`
	ImageURL   string      `json:"image_url"`
	AuthorName string      `json:"author_name"`
	CreatedBy  string      `json:"created_by" binding:"required"`
	Status     PostStatus  `json:"status" binding:"required"`
	Location   string      `json:"location"`
	Details    string      `json:"details"`
	Contact    ContactInfo `json:"contact"`
}

type updatePostRequest struct {
	Title      *string      `json:"title"`
	ImageURL   *string      `json:"image_url"`
	AuthorName *string      `json:"author_name"`
	Status     *PostStatus  `json:"status"`
	Location   *string      `json:"location"`
	Details    *string      `json:"details"`
	Contact    *ContactInfo `json:"contact"`
}
