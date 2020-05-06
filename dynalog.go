package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
)

// LogProperties define a desired logging structure
type LogProperties struct {
	Trace []string
	Debug []string
}

func (lp LogProperties) any() bool {
	return len(lp.Trace) != 0 || len(lp.Debug) != 0
}

func findTemplate(name string) (string, error) {
	path := fmt.Sprintf("templates/%s.log4j.properties", name)
	matches, err := filepath.Glob(path)
	if err != nil {
		return "", err
	}
	if len(matches) == 0 {
		return "", fmt.Errorf("Could not find template `%s`", name)
	}
	return matches[0], nil
}

func buildLog4j(template string, props LogProperties) ([]byte, error) {
	filename, err := findTemplate(template)
	if err != nil {
		fmt.Println("Could not find template: ", err)
		return nil, err
	}
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println("File reading error", err)
		return nil, err
	}

	// If no custom levels are provide, just return the initial file
	if !props.any() {
		return data, nil
	}

	data = append(data, "\n\n## DYNALOG CUSTOM LOGGING LEVELS ENABLED 🚀 \n"...)

	// Append the different logging types
	for _, c := range props.Trace {
		line := fmt.Sprintf("log4j.logger.%s=TRACE\n", c)
		data = append(data, line...)
	}

	// TODO: DRY this up
	for _, c := range props.Debug {
		line := fmt.Sprintf("log4j.logger.%s=DEBUG\n", c)
		data = append(data, line...)
	}

	return data, nil
}

func buildLogProperties(c *gin.Context) LogProperties {
	lp := LogProperties{}
	if len(c.Query("trace")) > 0 {
		lp.Trace = strings.Split(c.Query("trace"), ",")
	}
	if len(c.Query("debug")) > 0 {
		lp.Debug = strings.Split(c.Query("debug"), ",")
	}

	return lp
}

func main() {
	r := gin.Default()

	// ForceSSL in production
	r.Use(ForceSSL(Options{
		SSLProxyHeaders: map[string]string{"X-Forwarded-Proto": "https"},
		IsDevelopment: gin.Mode() != gin.ReleaseMode,
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

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
