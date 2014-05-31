package handlers

import (
	"encoding/json"
	"github.com/callumj/metrix/metric_core"
	"github.com/callumj/metrix/shared"
	"github.com/garyburd/redigo/redis"
	"github.com/jinzhu/now"
	"io"
	"net/http"
	"strconv"
	"time"
)

type SourceListResponse struct {
	Sources []string `json:"sources"`
}

type DateKeysResponse struct {
	Count int `json:"count"`
}

func SourceListHandler(c http.ResponseWriter, req *http.Request) {

	redisConn := shared.RedisPool.Get()
	defer redisConn.Close()

	var sourceList []string
	sources, err := redisConn.Do("SMEMBERS", metric_core.SourcesKey)
	if err != nil {
		shared.HandleError(err)
	} else {
		sourceList, err = redis.Strings(sources, err)
		if err != nil {
			shared.HandleError(err)
		}
	}

	resp := SourceListResponse{
		Sources: sourceList,
	}

	writeJSON(c, resp)
}

func DateKeysHandler(c http.ResponseWriter, req *http.Request) {
	source := req.FormValue("source")
	if len(source) == 0 {
		http.Error(c, "`source` param must be provided", http.StatusBadRequest)
		return
	}

	if source == "default" {
		source = ""
	}

	timeStartNeedsInit := true
	timeEndNeedsInit := true

	var timeStart time.Time
	var timeEnd time.Time

	timeStartParam := req.FormValue("start")
	timeEndParam := req.FormValue("end")

	if len(timeStartParam) != 0 {
		parsedInt, err := strconv.ParseInt(timeStartParam, 0, 64)
		if err == nil {
			timeStart = time.Unix(parsedInt, 0)
			timeStartNeedsInit = false
		}
	}

	if len(timeEndParam) != 0 {
		parsedInt, err := strconv.ParseInt(timeEndParam, 0, 64)
		if err == nil {
			timeEnd = time.Unix(parsedInt, 0)
			timeEndNeedsInit = false
		}
	}

	if timeStartNeedsInit {
		timeStart = now.New(time.Now().UTC()).BeginningOfDay()
	}

	if timeEndNeedsInit {
		timeEnd = timeStart.Add((24 * time.Hour) * -7)
	}

	allTimes := []time.Time{timeEnd}
	lastTime := timeEnd.Add(24 * time.Hour)

	for lastTime.Unix() <= timeStart.Unix() {
		allTimes = append(allTimes, lastTime)
		lastTime = lastTime.Add(24 * time.Hour)
	}

	resp := DateKeysResponse{
		Count: len(allTimes),
	}
	writeJSON(c, resp)
}

func writeJSON(c http.ResponseWriter, resp interface{}) {
	json, err := json.Marshal(resp)
	if err != nil {
		shared.HandleError(err)
	}

	jsonString := string(json)

	if len(json) != 0 {
		c.Header().Add("Content-Type", "application/json")
		c.Header().Add("Content-Length", strconv.Itoa(len(jsonString)))
		io.WriteString(c, jsonString)
	}
}
