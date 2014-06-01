package handlers

import (
	"encoding/json"
	"github.com/callumj/metrix/metric_core"
	"github.com/callumj/metrix/shared"
	"github.com/garyburd/redigo/redis"
	"io"
	"net/http"
	"strconv"
)

type SourceListResponse struct {
	Sources []string `json:"sources"`
}

type AvailableKeysResponse struct {
	Keys []string `json:"keys"`
}

type DateKeysResponse struct {
	Count int      `json:"count"`
	Dates []string `json:"dates"`
}

func SourceListHandler(c http.ResponseWriter, req *http.Request) {
	if !verifyAPIKey(c, req) {
		return
	}

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

func AvailableKeysHandler(c http.ResponseWriter, req *http.Request) {
	if !verifyAPIKey(c, req) {
		return
	}

	source := req.FormValue("source")
	if len(source) == 0 {
		http.Error(c, "`source` & `key` param must be provided", http.StatusBadRequest)
		return
	}

	redisConn := shared.RedisPool.Get()
	defer redisConn.Close()

	resp, err := redisConn.Do("SMEMBERS", metric_core.KeySourcesKey(source))
	if err != nil {
		shared.HandleError(err)
	}

	list, err := redis.Strings(resp, err)
	if err != nil {
		shared.HandleError(err)
	}

	json := AvailableKeysResponse{
		Keys: list,
	}
	writeJSON(c, json)
}

func DateKeysHandler(c http.ResponseWriter, req *http.Request) {
	if !verifyAPIKey(c, req) {
		return
	}

	source := req.FormValue("source")
	key := req.FormValue("key")
	if len(source) == 0 || len(key) == 0 {
		http.Error(c, "`source` & `key` param must be provided", http.StatusBadRequest)
		return
	}

	if source == "default" {
		source = ""
	}

	timeStartParam := req.FormValue("start")
	timeEndParam := req.FormValue("end")

	allTimes := shared.TimeBetweenDates(timeStartParam, timeEndParam)

	redisConn := shared.RedisPool.Get()
	defer redisConn.Close()

	redisConn.Send("MULTI")
	for _, date := range allTimes {
		thisKey := metric_core.BuildKVIncrementKey(date, source, key)
		redisConn.Send("EXISTS", thisKey)
	}
	res, err := redisConn.Do("EXEC")

	if err != nil {
		shared.HandleError(err)
	}

	keyStates, err := redis.Values(res, err)
	if err != nil {
		shared.HandleError(err)
	}

	activeKeys := []string{}

	for index, state := range keyStates {
		status := state.(int64)
		if status == 1 {
			matching := allTimes[index]
			activeKeys = append(activeKeys, metric_core.FormatDate(matching))
		}
	}

	resp := DateKeysResponse{
		Count: len(activeKeys),
		Dates: activeKeys,
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
