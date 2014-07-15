package handlers

import (
	"github.com/callumj/metrix/metric_core"
	"github.com/callumj/metrix/shared"
	"github.com/garyburd/redigo/redis"
	"io"
	"net/http"
)

type VersionGetResponse struct {
	Versions map[string]string `json:"versions"`
}

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

func VersionGetHandler(c http.ResponseWriter, req *http.Request) {
	source := req.FormValue("source")
	redisKey := metric_core.VersionSourcesKey(source)

	versionMaps := make(map[string]string)

	redisConn := shared.RedisPool.Get()
	defer redisConn.Close()

	res, err := redisConn.Do("HGETALL", redisKey)
	if err != nil {
		shared.HandleError(err)
	}

	flatten, err := redis.Strings(res, err)
	if err != nil {
		shared.HandleError(err)
	}

	expectingKey := true
	var lastKey string
	for _, result := range flatten {
		if expectingKey {
			expectingKey = false
			lastKey = result
		} else {
			versionMaps[lastKey] = result
			expectingKey = true
		}
	}

	resp := VersionGetResponse{
		Versions: versionMaps,
	}

	writeJSON(c, resp)
}
