package handlers

import (
	"encoding/base64"
	"fmt"
	"github.com/callumj/metrix/shared"
	"github.com/jinzhu/now"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var (
	transGif = "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAQAAAC1HAwCAAAAC0lEQVR42mP4/x8AAwAB/2+Bq7YAAAAASUVORK5CYII="
)

func IncrementMetricHandler(c http.ResponseWriter, req *http.Request) {
	key := req.FormValue("key")

	if (strings.Contains(req.Header.Get("Accept"), "image/") || req.FormValue("image") == "yes") && len(req.FormValue("redirect")) == 0 {
		decoded, err := base64.StdEncoding.DecodeString(transGif)
		if err == nil {
			c.Header().Add("Content-Type", "image/gif")
			c.Header().Add("Content-Length", strconv.Itoa(len(decoded)))
			io.WriteString(c, string(decoded))
		}
	} else {
		if len(req.FormValue("redirect")) != 0 {
			http.Redirect(c, req, req.FormValue("redirect"), 307)
		} else {
			body := "OK"
			c.Header().Add("Content-Type", "text/html")
			c.Header().Add("Content-Length", "2")
			io.WriteString(c, body)
		}
	}

	if len(key) != 0 {
		tPoint := time.Now().UTC()
		subkey := req.FormValue("subkey")
		source := req.FormValue("source")
		recordIncrMetric(key, subkey, source, tPoint)

		headers := req.Header
		sourceIp := headers.Get("X-Real-Ip")
		if len(sourceIp) == 0 {
			sourceIp = headers.Get("X-Forwarded-For")
		}
		if len(sourceIp) == 0 {
			sourceIp = req.RemoteAddr
		}

		lastColon := strings.LastIndex(sourceIp, ":")
		var ipOnly string
		if lastColon != -1 {
			ipOnly = sourceIp[0:lastColon]
		} else {
			ipOnly = sourceIp
		}
		storeIntoMetric(key, subkey, source, tPoint, ipOnly)
	}
}

func recordIncrMetric(key, subkey, source string, tPoint time.Time) {
	if len(key) == 0 {
		return
	}

	start := now.New(tPoint).BeginningOfDay()
	diff := tPoint.Sub(start)

	day := tPoint.Format("02012006")
	if len(source) != 0 {
		day = fmt.Sprintf("%v:%v", source, day)
	}
	totalMinutes := fmt.Sprintf("%v", int(diff.Minutes()))

	var perMinuteKey string
	var kvIncrementKey string
	var kvIncrementSubKey string
	if len(subkey) != 0 {
		perMinuteKey = fmt.Sprintf("%v:%v:%v", day, key, subkey)
		kvIncrementKey = fmt.Sprintf("%v:%v", day, key)
		kvIncrementSubKey = subkey
	} else {
		perMinuteKey = fmt.Sprintf("%v:%v", day, key)
		kvIncrementKey = day
		kvIncrementSubKey = key
	}

	redis := shared.RedisPool.Get()
	defer redis.Close()
	err := redis.Send("HINCRBY", perMinuteKey, totalMinutes, "1")
	if err != nil {
		shared.HandleError(err)
	}

	err = redis.Send("HINCRBY", kvIncrementKey, kvIncrementSubKey, "1")
	if err != nil {
		shared.HandleError(err)
	}

	err = redis.Flush()
	if err != nil {
		shared.HandleError(err)
	}
}

func storeIntoMetric(key, subkey, source string, tPoint time.Time, value string) {
	if len(key) == 0 || len(value) == 0 {
		return
	}
	day := tPoint.Format("02012006")
	if len(source) != 0 {
		day = fmt.Sprintf("%v:%v", source, day)
	}

	var totalKey string
	if len(subkey) != 0 {
		totalKey = fmt.Sprintf("ip:%v:%v:%v", day, key, subkey)
	} else {
		totalKey = fmt.Sprintf("ip:%v:%v", day, key)
	}

	redis := shared.RedisPool.Get()
	defer redis.Close()
	err := redis.Send("SADD", totalKey, value)
	if err != nil {
		shared.HandleError(err)
	}
}
