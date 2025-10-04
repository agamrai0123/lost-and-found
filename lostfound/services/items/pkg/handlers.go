package pkg

import (
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (s *ItemService) listPosts() gin.HandlerFunc {
	return func(c *gin.Context) {
		statusQ := c.Query("status")
		locationQ := c.Query("location")

		s.mu.RLock()
		list := make([]*Post, 0, len(s.posts))
		for _, p := range s.posts {
			// filter
			if statusQ != "" && string(p.Status) != statusQ {
				continue
			}
			if locationQ != "" && !containsIgnoreCase(p.Location, locationQ) {
				continue
			}
			list = append(list, p)
		}
		s.mu.RUnlock()
		c.JSON(http.StatusOK, list)
	}
}

// // POST /api/posts
const maxUploadSize = 10 << 20 // 10 MiB

func (s *ItemService) createPost() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Protect from huge bodies
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxUploadSize)

		// Parse form (multipart)
		if err := c.Request.ParseMultipartForm(maxUploadSize); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "failed to parse multipart form: " + err.Error()})
			return
		}

		// Read required fields from the form
		req := createPostRequest{
			Title:      c.PostForm("title"),
			ImageURL:   "", // will fill later
			AuthorName: c.PostForm("author_name"),
			CreatedBy:  c.PostForm("created_by"),
			Location:   c.PostForm("location"),
			Details:    c.PostForm("details"),
			Contact: ContactInfo{
				Name:  c.PostForm("contact_name"),
				Phone: c.PostForm("contact_phone"),
				Email: c.PostForm("contact_email"),
			},
		}
		// status is required
		statusStr := c.PostForm("status")
		req.Status = PostStatus(statusStr)

		// validate required fields (created_by, title, status)
		if req.CreatedBy == "" || req.Title == "" || req.Status == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "missing required fields: created_by, title, status"})
			return
		}
		if !validStatus(req.Status) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid status (must be lost|found|claimed)"})
			return
		}

		// Handle optional image upload (form key "image")
		var savedURL string
		fileHeader, err := c.FormFile("image")
		if err == nil && fileHeader != nil {
			// Open uploaded file to inspect mime
			f, err := fileHeader.Open()
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "cannot open uploaded file"})
				return
			}
			defer f.Close()

			// Read first 512 bytes for content sniffing
			buf := make([]byte, 512)
			n, _ := f.Read(buf)
			mimeType := http.DetectContentType(buf[:n])

			// Validate MIME
			switch mimeType {
			case "image/jpeg", "image/png", "image/gif", "image/webp":
				// allowed
			default:
				c.JSON(http.StatusBadRequest, gin.H{"error": "unsupported image type: " + mimeType})
				return
			}

			// reset reader to start so SaveUploadedFile can read whole file (open again)
			// Gin's SaveUploadedFile expects the file header, so we use SaveUploadedFile directly.
			// Generate unique filename preserving extension if available.
			ext := filepath.Ext(fileHeader.Filename)
			if ext == "" {
				// fallback from MIME
				switch mimeType {
				case "image/png":
					ext = ".png"
				case "image/gif":
					ext = ".gif"
				case "image/webp":
					ext = ".webp"
				default:
					ext = ".jpg"
				}
			}
			filename := uuid.New().String() + ext
			uploadDir := s.imageDir // see service struct change below; default "uploads"
			if err := os.MkdirAll(uploadDir, 0o755); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create upload dir"})
				return
			}
			dst := filepath.Join(uploadDir, filename)

			// Use Gin helper to save the uploaded file
			if err := c.SaveUploadedFile(fileHeader, dst); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save uploaded file: " + err.Error()})
				return
			}

			// savedURL should be the publicly-accessible URL you serve the file from
			// if you use s.router.Static("/uploads", uploadDir) then:
			savedURL = "/uploads/" + filename
		}

		now := time.Now().UTC()
		id := uuid.New().String()
		post := &Post{
			ID:         id,
			Title:      req.Title,
			ImageURL:   savedURL, // path/URL to file
			AuthorName: req.AuthorName,
			CreatedBy:  req.CreatedBy,
			CreatedAt:  now,
			UpdatedAt:  nil,
			Status:     req.Status,
			Location:   req.Location,
			Details:    req.Details,
			Contact:    req.Contact,
		}

		s.mu.Lock()
		s.posts[id] = post
		s.mu.Unlock()

		if err := s.saveToFile(); err != nil {
			s.logger.Printf("warning: saveToFile failed: %v", err)
		}

		c.JSON(http.StatusCreated, post)
	}
}

// // GET /api/posts/:id
func (s *ItemService) getPost() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		s.mu.RLock()
		post, ok := s.posts[id]
		s.mu.RUnlock()
		if !ok {
			c.JSON(http.StatusNotFound, gin.H{"error": "post not found"})
			return
		}
		c.JSON(http.StatusOK, post)
	}
}

// // PUT /api/posts/:id
// func (s *ItemService) handleUpdatePost() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		id := c.Param("id")
// 		var req updatePostRequest
// 		if err := c.ShouldBindJSON(&req); err != nil {
// 			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 			return
// 		}

// 		s.mu.Lock()
// 		post, ok := s.posts[id]
// 		if !ok {
// 			s.mu.Unlock()
// 			c.JSON(http.StatusNotFound, gin.H{"error": "post not found"})
// 			return
// 		}

// 		// apply updates if set
// 		updated := false
// 		if req.Title != nil {
// 			post.Title = *req.Title
// 			updated = true
// 		}
// 		if req.ImageURL != nil {
// 			post.ImageURL = *req.ImageURL
// 			updated = true
// 		}
// 		if req.AuthorName != nil {
// 			post.AuthorName = *req.AuthorName
// 			updated = true
// 		}
// 		if req.Status != nil {
// 			if !validStatus(*req.Status) {
// 				s.mu.Unlock()
// 				c.JSON(http.StatusBadRequest, gin.H{"error": "invalid status (must be lost|found|claimed)"})
// 				return
// 			}
// 			post.Status = *req.Status
// 			updated = true
// 		}
// 		if req.Location != nil {
// 			post.Location = *req.Location
// 			updated = true
// 		}
// 		if req.Details != nil {
// 			post.Details = *req.Details
// 			updated = true
// 		}
// 		if req.Contact != nil {
// 			post.Contact = *req.Contact
// 			updated = true
// 		}
// 		if updated {
// 			now := time.Now().UTC()
// 			post.UpdatedAt = &now
// 		}
// 		s.mu.Unlock()

// 		if updated {
// 			if err := s.saveToFile(); err != nil {
// 				s.logger.Printf("warning: saveToFile failed: %v", err)
// 			}
// 		}

// 		c.JSON(http.StatusOK, post)
// 	}
// }

// // DELETE /api/posts/:id
// func (s *ItemService) handleDeletePost() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		id := c.Param("id")
// 		s.mu.Lock()
// 		_, ok := s.posts[id]
// 		if !ok {
// 			s.mu.Unlock()
// 			c.JSON(http.StatusNotFound, gin.H{"error": "post not found"})
// 			return
// 		}
// 		delete(s.posts, id)
// 		s.mu.Unlock()

// 		if err := s.saveToFile(); err != nil {
// 			s.logger.Printf("warning: saveToFile failed: %v", err)
// 		}
// 		c.Status(http.StatusNoContent)
// 	}
// }
