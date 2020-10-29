package main

import (
	"bytes"
	"fmt"
	"net/http"
	"strings"
	"text/template"

	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
)

func main() {
	funcMap := template.FuncMap{"ToUpper": strings.ToUpper}
	r := gin.Default()

	// ForceSSL in production
	r.Use(ForceSSL(Options{
		SSLProxyHeaders: map[string]string{"X-Forwarded-Proto": "https"},
		IsDevelopment:   gin.Mode() != gin.ReleaseMode,
	}))

	// Basic site placeholder
	r.Use(static.Serve("/", static.LocalFile("./public", true)))

	// Dynamic log builder
	r.GET("/templates/:name", func(c *gin.Context) {
		n := c.Param("name")
		lp := buildLogProperties(c)
		logFile, err := buildLog4j(n, lp)
		if err != nil {
			c.String(http.StatusBadRequest, fmt.Sprintf("bad request: %s", err))
		}
		c.String(http.StatusOK, string(logFile))
	})

	// Now supporting log4jv2
	r.GET("/v2/templates/:name", func(c *gin.Context) {
		n := c.Param("name")
		ll := buildLogLevelsFromRequestParams(c)

		templateName := fmt.Sprintf("%s.log4j.xml", n)
		t, err := template.New(templateName).Funcs(funcMap).ParseFiles("./templates/" + templateName)
		if err != nil {
			c.String(http.StatusBadRequest, fmt.Sprintf("template not found"))
			return
		}

		var result bytes.Buffer
		t.Execute(&result, ll)
		c.String(http.StatusOK, result.String())

	})

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
