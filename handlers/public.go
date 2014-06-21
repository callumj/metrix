package handlers

import (
	"github.com/callumj/metrix/resource_bundle"
	"github.com/gorilla/mux"
	"io"
	"net/http"
	"strconv"
)

func PublicHandler(c http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)

	res, err := resource_bundle.FetchFile(vars["path"])

	if err != nil {
		http.Error(c, err.Error(), http.StatusInternalServerError)
	} else {
		c.Header().Add("Content-Type", res.ContentType)
		c.Header().Add("Content-Length", strconv.Itoa(len(res.Data)))
		io.WriteString(c, string(res.Data))
	}
}
