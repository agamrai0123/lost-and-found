package pkg

import (
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func (s *service) home() gin.HandlerFunc {
	return func(c *gin.Context) {
		data := map[string]any{
			"Title":    "Home - LostAndFound",
			"Active":   "home",
			"BasePath": s.basePath,
		}
		s.renderTemplate(c, "home", data)
	}
}

func (s *service) about() gin.HandlerFunc {
	return func(c *gin.Context) {
		data := map[string]any{
			"Title":    "About - LostAndFound",
			"Active":   "about",
			"BasePath": s.basePath,
		}
		s.renderTemplate(c, "about", data)
	}
}

// GET: show login page
func (s *service) loginGet() gin.HandlerFunc {
	return func(c *gin.Context) {
		data := map[string]any{
			"Title":    "Login - LostAndFound",
			"Active":   "login",
			"Flash":    c.Query("msg"),
			"BasePath": s.basePath,
		}
		s.renderTemplate(c, "login", data)
	}
}

// POST: process login
func (s *service) loginPost() gin.HandlerFunc {
	return func(c *gin.Context) {
		username := c.PostForm("username")
		password := c.PostForm("password")

		if username == "" || password == "" {
			c.Redirect(http.StatusSeeOther, s.basePath+"/login?msg="+url.QueryEscape("missing fields"))
			return
		}
		s.usersMu.RLock()
		u, ok := s.users[username]
		s.usersMu.RUnlock()
		if !ok {
			c.Redirect(http.StatusSeeOther, s.basePath+"/login?msg="+url.QueryEscape("invalid credentials"))
			return
		}

		if err := bcrypt.CompareHashAndPassword(u.Password, []byte(password)); err != nil {
			c.Redirect(http.StatusSeeOther, s.basePath+"/login?msg="+url.QueryEscape("invalid credentials"))
			return
		}

		// create cookie (1 hour)
		c.SetCookie("session_user", username, 3600, "/", "", false, true) // set Secure: true in prod
		c.Redirect(http.StatusSeeOther, s.basePath+"/")
	}
}

// GET: show signup page
func (s *service) signupGet() gin.HandlerFunc {
	return func(c *gin.Context) {
		data := map[string]any{
			"Title":    "Sign up - LostAndFound",
			"Active":   "signup",
			"Flash":    c.Query("msg"),
			"BasePath": s.basePath,
		}
		s.renderTemplate(c, "signup", data)
	}
}

// POST: process signup form
func (s *service) signupPost() gin.HandlerFunc {
	return func(c *gin.Context) {
		fullname := c.PostForm("fullname")
		email := c.PostForm("email")
		username := c.PostForm("username")
		password := c.PostForm("password")

		// basic validation
		if username == "" || password == "" || email == "" || fullname == "" {
			c.Redirect(http.StatusSeeOther, s.basePath+"/signup?msg="+url.QueryEscape("missing fields"))
			return
		}
		if len(password) < 6 {
			c.Redirect(http.StatusSeeOther, s.basePath+"/signup?msg="+url.QueryEscape("password too short"))
			return
		}

		// check duplicate username
		s.usersMu.Lock()
		if _, exists := s.users[username]; exists {
			s.usersMu.Unlock()
			c.Redirect(http.StatusSeeOther, s.basePath+"/signup?msg="+url.QueryEscape("username already exists"))
			return
		}

		// create password hash
		hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			s.usersMu.Unlock()
			s.logger.Error().Err(err).Msg("failed to hash password")
			c.String(http.StatusInternalServerError, "internal error")
			return
		}

		// save to in-memory store
		s.users[username] = &user{
			Username: username,
			Password: hash,
			Fullname: fullname,
			Email:    email,
		}
		s.usersMu.Unlock()

		// set session cookie (demo)
		c.SetCookie("session_user", username, 3600, "/", "", false, true)
		c.Redirect(http.StatusSeeOther, s.basePath+"/")
	}
}
