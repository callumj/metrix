package handlers

import (
	"fmt"
	"github.com/callumj/metrix/shared"
	"io"
	"net/http"
	"time"
)

func IncrementMetricHandler(c http.ResponseWriter, req *http.Request) {
	key := req.FormValue("key")

	if len(key) != 0 {
		recordIncrMetric(key, req.FormValue("subkey"))
	}

	body := "OK"
	c.Header().Add("Content-Type", "text/html")
	c.Header().Add("Content-Length", "2")
	io.WriteString(c, body)
}

func recordIncrMetric(key, subkey string) {
	tPoint := time.Now().UTC().Format("02012006")
	if len(subkey) == 0 {
		subkey = key
		key = tPoint
	} else {
		key = fmt.Sprintf("%v_%v", tPoint, key)
	}

	if len(key) != 0 {
		redis := shared.RedisPool.Get()
		_, err := redis.Do("HINCRBY", key, subkey, "1")
		redis.Close()
		if err != nil {
			shared.HandleError(err)
		}
	}
}
