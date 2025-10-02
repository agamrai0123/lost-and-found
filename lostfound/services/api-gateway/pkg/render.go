package pkg

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func (s *service) renderTemplate(c *gin.Context, tmpl string, data interface{}) {
	var t *template.Template
	var err error

	if !s.isProduction {
		// Rebuild cache in dev mode for hot-reloading
		s.templateCache, err = createTemplateCache()
		if err != nil {
			log.Printf("cannot create template cache: %v", err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}
	}

	t, ok := s.templateCache[tmpl+".page.tmpl"]
	if !ok {
		log.Printf("could not get template from cache: %s", tmpl)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Template not found"})
		return
	}

	// Execute template
	if err := t.ExecuteTemplate(c.Writer, "base", data); err != nil {
		c.String(http.StatusInternalServerError, "Error rendering template")
	}
}

func createTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}
	pages, err := filepath.Glob("./web/templates/*.page.tmpl")
	if err != nil {
		return nil, err
	}

	layouts, err := filepath.Glob("./web/templates/*.layout.tmpl")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)
		ts, err := template.New(name).Funcs(template.FuncMap{}).ParseFiles(page)
		if err != nil {
			return nil, err
		}
		if len(layouts) > 0 {
			ts, err = ts.ParseFiles(layouts...)
			if err != nil {
				return nil, err
			}
		}
		cache[name] = ts
	}
	return cache, nil
}
