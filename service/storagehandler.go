package service

import (
	"io"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/yangb8/webservice/common/api"
	"github.com/yangb8/webservice/common/config"
)

// StorageHandler ...
type StorageHandler struct {
	api.Handler
	storageServ Storage
}

// NewtSorageHandler ...
func NewStorageHandler(c *config.Config) http.Handler {
	res := &StorageHandler{
		storageServ: GetStorage(c),
	}
	res.Handler.Router = mux.NewRouter()

	res.Methods("PUT").Path("/objects/{key:.+}").
		Name("put_object").
		HandlerFunc(res.Upload)
	res.Methods("GET").Path("/objects/{key:.+}").
		Name("get_object").
		HandlerFunc(res.Download)
	res.Methods("DELETE").Path("/objects/{key:.+}").
		Name("delete_object").
		HandlerFunc(res.Delete)

	return res
}

func (h *StorageHandler) Upload(w http.ResponseWriter, r *http.Request) {
	var b struct {
		Key string `mux:"key" preprocess:"urldecode"`
	}
	if err := api.BindSkipBody(r, &b); err != nil {
		panic(err)
	}

	defer r.Body.Close()
	h.storageServ.Upload(b.Key, r.Body)
}

func (h *StorageHandler) Download(w http.ResponseWriter, r *http.Request) {
	var b struct {
		Key string `mux:"key" preprocess:"urldecode"`
	}
	if err := api.Bind(r, &b); err != nil {
		panic(err)
	}

	res, err := h.storageServ.Download(b.Key)
	if err != nil {
		panic(err)
	}
	defer res.Close()
	io.Copy(w, res)
}

func (h *StorageHandler) Delete(w http.ResponseWriter, r *http.Request) {
	var b struct {
		Key string `mux:"key" preprocess:"urldecode"`
	}
	if err := api.Bind(r, &b); err != nil {
		panic(err)
	}
	h.storageServ.Delete(b.Key)
}
