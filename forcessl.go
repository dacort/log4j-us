package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// Options is a struct for specifying configuration options for the secure.Secure middleware.
type Options struct {
	// When developing, the AllowedHosts, SSL, and STS options can cause some unwanted effects. Usually testing happens on http, not https, and on localhost, not your production domain... so set this to true for dev environment.
	// If you would like your development environment to mimic production with complete Host blocking, SSL redirects, and STS headers, leave this as false. Default if false.
	IsDevelopment bool

	// SSLProxyHeaders is set of header keys with associated values that would indicate a valid https request. Useful when using Nginx: `map[string]string{"X-Forwarded-Proto": "https"}`. Default is blank map.
	SSLProxyHeaders map[string]string
}

// Secure is a middleware that helps setup a few basic security features. A single secure.Options struct can be
// provided to configure which features should be enabled, and the ability to override a few of the default values.
type secure struct {
	// Customize Secure with an Options struct.
	opt Options
}

func (s *secure) process(w http.ResponseWriter, r *http.Request) error {
	// SSL check.
	if s.opt.IsDevelopment == false {
		fmt.Println("IsDevelopment is true, checking for redirect")
		isSSL := false
		if strings.EqualFold(r.URL.Scheme, "https") || r.TLS != nil {
			isSSL = true
		} else {
			for k, v := range s.opt.SSLProxyHeaders {
				if r.Header.Get(k) == v {
					isSSL = true
					break
				}
			}
		}

		if isSSL == false {
			url := r.URL
			url.Scheme = "https"
			url.Host = r.Host

			status := http.StatusMovedPermanently

			http.Redirect(w, r, url.String(), status)
			return fmt.Errorf("Redirecting to HTTPS")
		}
	}

	return nil
}

// ForceSSL redirects the user to https
// Redacted version of https://github.com/gin-gonic/contrib/blob/master/secure/secure.go
func ForceSSL(options Options) gin.HandlerFunc {
	s := secure{options}

	return func(c *gin.Context) {
		err := s.process(c.Writer, c.Request)
		if err != nil {
			if c.Writer.Written() {
				c.AbortWithStatus(c.Writer.Status())
			} else {
				c.AbortWithError(http.StatusInternalServerError, err)
			}
		}
	}

}
