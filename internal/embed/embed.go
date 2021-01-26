package embed

import (
	"embed"
	"net/http"
)

var (
	//go:embed api.swagger.json
	fs embed.FS

	Handler = http.FileServer(http.FS(fs))
)
