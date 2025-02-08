package api

import "net/http"

func HealthProbe(w http.ResponseWriter, _ *http.Request) {
	w.Write([]byte("OK"))
}
