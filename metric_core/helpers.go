package metric_core

import (
	"fmt"
	"github.com/callumj/metrix/shared"
	"github.com/garyburd/redigo/redis"
	"time"
)

func KnownKeysForSource(source string, redisConnPntr *redis.Conn) []string {
	redisKey := KeySourcesKey(source)
	redisConn := *redisConnPntr
	res, err := redisConn.Do("SMEMBERS", redisKey)
	if err != nil {
		shared.HandleError(err)
		return []string{}
	} else {
		conv, err := redis.Strings(res, err)
		if err != nil {
			shared.HandleError(err)
			return []string{}
		} else {
			return conv
		}
	}
}

func FormatDate(date time.Time) string {
	return date.Format("02012006")
}

func BuildKVIncrementKey(date time.Time, source, key string) string {
	day := FormatDate(date)

	return BuildKVIncrementKeyString(day, source, key)
}

func BuildKVIncrementKeyString(day, source, key string) string {
	source = RewriteSource(source)

	if len(source) != 0 {
		day = fmt.Sprintf("%v:%v", source, day)
	}

	return fmt.Sprintf("%s:%s", day, key)
}

func BuildByMinuteKey(date time.Time, source, key, subkey string) string {
	source = RewriteSource(source)
	var day string
	if len(source) == 0 {
		day = FormatDate(date)
	} else {
		day = fmt.Sprintf("%v:%v", source, FormatDate(date))
	}

	if len(subkey) == 0 {
		subkey = "default"
	}

	return fmt.Sprintf("%v:%v:%v", day, key, subkey)
}

func RewriteSource(source string) string {
	if source == "default" {
		source = ""
	}
	return source
}
