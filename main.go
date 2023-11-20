package main

import (
	"fmt"
	"net/http"
	"path/filepath"
	"sync"
	"text/template"

	"github.com/gin-gonic/gin"
)

// TemplateCache is a simple cache for HTML templates.
type TemplateCache struct {
	templates map[string]*template.Template
	mu        sync.Mutex
}

var cache *TemplateCache = &TemplateCache{
	templates: make(map[string]*template.Template),
}

func main() {
	router := gin.Default()
	// Initialize the template cache
	createTemplateCacheMiddleware(cache)

	router.GET("/", homeHandler)
	router.GET("/about", aboutHandler)
	router.Run("localhost:8080")
}

// TemplateCacheMiddleware is a Gin middleware to handle template caching.
func createTemplateCacheMiddleware(cache *TemplateCache) error {
	layoutFiles, err := filepath.Glob("template/*.html.tmpl")
	if err != nil {
		fmt.Println("Albin error2: ", err)
		return err
	}

	for _, file := range layoutFiles {
		tmplContent, err := template.ParseFiles(file)
		if err != nil {
			fmt.Println("Albin error1: ", err)
			return err
		}

		layoutFiles, err := filepath.Glob("template/*.layout.tmpl")
		if err != nil {
			fmt.Println("Albin error2: ", err)
			return err
		}

		if len(layoutFiles) > 0 {
			tmplContent, err = tmplContent.ParseGlob("template/*.layout.tmpl")
			if err != nil {
				fmt.Println("Albin error3: ", err)
				return err
			}
		}

		_, key := filepath.Split(file)
		fmt.Println("Albin key: ", key)
		cache.set(key, tmplContent)
	}
	return nil
}

func renderTemplate(c *gin.Context, tmpl string, variables gin.H) {
	// Retrieve the template from the context
	tmplContent, ok := cache.get(tmpl)
	if !ok {
		fmt.Println("Albin error4")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	// Execute the template with provided variables
	if err := tmplContent.Execute(c.Writer, variables); err != nil {
		fmt.Println("Albin error5 ", err)
		c.AbortWithStatus(http.StatusInternalServerError)
	}
}

// Template cache methods
func (tc *TemplateCache) get(key string) (*template.Template, bool) {
	tc.mu.Lock()
	defer tc.mu.Unlock()

	tmpl, ok := tc.templates[key]
	return tmpl, ok
}

func (tc *TemplateCache) set(key string, tmpl *template.Template) {
	tc.mu.Lock()
	defer tc.mu.Unlock()

	tc.templates[key] = tmpl
}

func homeHandler(c *gin.Context) {
	renderTemplate(c, "home.html.tmpl", gin.H{
		"title": "Home",
		"name":  "Home ",
	})
}

func aboutHandler(c *gin.Context) {
	renderTemplate(c, "about.html.tmpl", gin.H{
		"title": "About",
		"name":  "About",
	})
}
