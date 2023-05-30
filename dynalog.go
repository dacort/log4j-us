package main

import (
	"bytes"
	"embed"
	"fmt"
	"io/fs"
	"net/http"
	"strings"
	"text/template"

	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
)

//go:embed templates/*
var resources embed.FS

//go:embed public/*
var staticFS embed.FS

var funcMap = template.FuncMap{"ToUpper": strings.ToUpper}

func renderTemplate(c *gin.Context, templateName string) {
	ll := buildLogLevelsFromRequestParams(c)

	templatePath := fmt.Sprintf("templates/%s", templateName)
	t, err := template.New(templateName).Funcs(funcMap).ParseFS(resources, templatePath)
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("template not found"))
		return
	}

	var result bytes.Buffer
	t.Execute(&result, ll)
	c.String(http.StatusOK, result.String())
}

func main() {
	
	r := gin.Default()

	// ForceSSL in production
	r.Use(ForceSSL(Options{
		SSLProxyHeaders: map[string]string{"X-Forwarded-Proto": "https"},
		IsDevelopment:   gin.Mode() != gin.ReleaseMode,
	}))

	// Basic site placeholder
	r.Use(static.Serve("/", EmbedFolder(staticFS, "public", true)))

	// Dynamic log builder
	r.GET("/templates/:name", func(c *gin.Context) {
		n := c.Param("name")
		templateName := fmt.Sprintf("%s.log4j.properties", n)
		renderTemplate(c, templateName)
	})

	// Now supporting log4jv2
	r.GET("/v2/templates/:name", func(c *gin.Context) {
		n := c.Param("name")
		templateName := fmt.Sprintf("%s.log4j.xml", n)
		renderTemplate(c, templateName)
	})

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

type embedFileSystem struct {
	http.FileSystem
	indexes bool
}

func (e embedFileSystem) Exists(prefix string, path string) bool {
	f, err := e.Open(path)
	if err != nil {
		return false
	}

	// check if indexing is allowed
	s, _ := f.Stat()
	if s.IsDir() && !e.indexes {
		return false
	}

	return true
}

func EmbedFolder(fsEmbed embed.FS, targetPath string, index bool) static.ServeFileSystem {
	subFS, err := fs.Sub(fsEmbed, targetPath)
	if err != nil {
		panic(err)
	}
	return embedFileSystem{
		FileSystem: http.FS(subFS),
		indexes:    index,
	}
}
