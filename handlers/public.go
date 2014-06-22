package handlers

import (
	"github.com/callumj/metrix/resource_bundle"
	"github.com/gorilla/mux"
	"io"
	"net/http"
	"strconv"
	"strings"
)

func PublicDevHandler(c http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)

	processFile(vars["path"], c, req)
}

func PublicProdHandler(c http.ResponseWriter, req *http.Request) {
	path := strings.TrimPrefix(req.RequestURI, "/")
	if len(path) == 0 {
		path = "index.html"
	}
	processFile(path, c, req)
}

func processFile(key string, c http.ResponseWriter, req *http.Request) {
	res, err := resource_bundle.FetchFile(key)

	if err == resource_bundle.ErrNotExist {
		http.Error(c, err.Error(), http.StatusNotFound)
	} else if err == resource_bundle.ErrLoading || err == resource_bundle.ErrReading {
		http.Error(c, err.Error(), http.StatusInternalServerError)
	} else {
		c.Header().Add("Content-Type", res.ContentType)
		c.Header().Add("Content-Length", strconv.Itoa(len(res.Data)))
		c.Header().Add("eTag", res.Hash)
		io.WriteString(c, string(res.Data))
	}
}
