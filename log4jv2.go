package main

import (
	"strings"

	"github.com/gin-gonic/gin"
)

// LogLevel is a single class name/log level pair
type LogLevel struct {
	Name  string
	Level string
}

var logLevels = []string{"trace", "debug", "info", "warn", "error", "fatal"}

func buildLogLevelsFromRequestParams(c *gin.Context) (ll []LogLevel) {
	for _, level := range logLevels {
		if len(c.Query(level)) > 0 {
			for _, name := range strings.Split(c.Query(level), ",") {
				ll = append(ll, LogLevel{name, level})
			}
		}
	}

	return ll
}
