package interfaces

import "net/http"

type Handler interface {
	http.Handler
}
