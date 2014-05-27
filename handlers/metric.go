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
		subkey := req.FormValue("subkey")
		recordIncrMetric(key, subkey)
		storeIntoMetric(key, subkey, req.RemoteAddr)
	}

	body := "OK"
	c.Header().Add("Content-Type", "text/html")
	c.Header().Add("Content-Length", "2")
	io.WriteString(c, body)
}

func recordIncrMetric(key, subkey string) {
	key, subkey = computeKeys(key, subkey)
	if len(key) != 0 {
		redis := shared.RedisPool.Get()
		_, err := redis.Do("HINCRBY", key, subkey, "1")
		redis.Close()
		if err != nil {
			shared.HandleError(err)
		}
	}
}

func storeIntoMetric(key, subkey, value string) {
	key, subkey = computeKeys(key, subkey)
	if len(key) != 0 {
		join := fmt.Sprintf("%v:%v", key, subkey)
		redis := shared.RedisPool.Get()
		_, err := redis.Do("SADD", join, value)
		redis.Close()
		if err != nil {
			shared.HandleError(err)
		}
	}
}

func computeKeys(key, subkey string) (string, string) {
	if len(key) == 0 {
		return "", ""
	}
	tPoint := time.Now().UTC().Format("02012006")
	if len(subkey) == 0 {
		subkey = key
		key = tPoint
	} else {
		key = fmt.Sprintf("%v_%v", tPoint, key)
	}

	return key, subkey
}
