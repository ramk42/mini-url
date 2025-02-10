// Package api provides the http handlers for the URL shortener service.
package api

import "net/http"

func HealthProbe(w http.ResponseWriter, _ *http.Request) {
	w.Write([]byte("OK"))
}
