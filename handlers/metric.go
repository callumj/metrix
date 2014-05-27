package handlers

import (
	"fmt"
	"github.com/callumj/metrix/shared"
	"github.com/jinzhu/now"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

func IncrementMetricHandler(c http.ResponseWriter, req *http.Request) {
	key := req.FormValue("key")

	if len(key) != 0 {
		tPoint := time.Now().UTC()
		subkey := req.FormValue("subkey")
		source := req.FormValue("source")
		recordIncrMetric(key, subkey, source, tPoint)

		headers := req.Header
		log.Printf("%v", headers)
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

	body := "OK"
	c.Header().Add("Content-Type", "text/html")
	c.Header().Add("Content-Length", "2")
	io.WriteString(c, body)
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
	var daySubKey string
	if len(subkey) != 0 {
		perMinuteKey = fmt.Sprintf("%v:%v:%v", day, key, subkey)
		daySubKey = fmt.Sprintf("%v:%v", key, subkey)
	} else {
		perMinuteKey = fmt.Sprintf("%v:%v", day, key)
		daySubKey = key
	}

	redis := shared.RedisPool.Get()
	defer redis.Close()
	err := redis.Send("HINCRBY", perMinuteKey, totalMinutes, "1")
	if err != nil {
		shared.HandleError(err)
	}

	err = redis.Send("HINCRBY", day, daySubKey, "1")
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
