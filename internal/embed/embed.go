package embed

import (
	"embed"
	"net/http"
	"net/url"
)

var (
	//go:embed v1.swagger.json
	fs embed.FS

	handler = http.FileServer(http.FS(fs))
)

func HandlerFunc(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/v1/swagger.json" {
		r2 := new(http.Request)
		*r2 = *r
		r2.URL = new(url.URL)
		*r2.URL = *r.URL
		r2.URL.Path = "/v1.swagger.json"
		handler.ServeHTTP(w, r2)
	} else {
		handler.ServeHTTP(w, r)
	}
}
