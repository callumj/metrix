package handlers

import (
	"github.com/callumj/metrix/metric_core"
	"github.com/callumj/metrix/shared"
	"io"
	"net/http"
)

func VersionSetHandler(c http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		http.Error(c, "Must be POST", http.StatusBadRequest)
		return
	}

	source := req.FormValue("source")
	key := req.FormValue("key")
	if len(key) == 0 {
		http.Error(c, "`key` must be present", http.StatusBadRequest)
		return
	}
	version := req.FormValue("version")
	if len(version) == 0 {
		http.Error(c, "`version` must be present", http.StatusBadRequest)
		return
	}

	redisConn := shared.RedisPool.Get()
	defer redisConn.Close()
	err := redisConn.Send("HSET", metric_core.VersionSourcesKey(source), key, version)
	if err != nil {
		shared.HandleError(err)
	}

	body := "OK"
	c.Header().Add("Content-Type", "text/html")
	c.Header().Add("Content-Length", "2")
	io.WriteString(c, body)
}
