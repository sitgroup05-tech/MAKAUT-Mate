package teacher

import "net/http"

type Handler interface {
	UploadFile(w http.ResponseWriter, r *http.Request)
}

type handler struct {
}

func NewHanlder() Handler {
	return &handler{}
}

func (h *handler) UploadFile(w http.ResponseWriter, r *http.Request) {

}
