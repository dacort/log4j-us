package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

// LogProperties define a desired logging structure
type LogProperties struct {
	Trace []string
	Debug []string
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

func main() {
	r := gin.Default()
	r.GET("/templates/:name", func(c *gin.Context) {
		n := c.Param("name")
		lp := LogProperties{
			Trace: strings.Split(c.Query("trace"), ","),
		}
		logFile, err := buildLog4j(n, lp)
		if err != nil {
			c.String(http.StatusBadRequest, fmt.Sprintf("bad request: %s", err))
		}
		c.String(http.StatusOK, string(logFile))
	})
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
