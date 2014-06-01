package shared

import (
	"github.com/jinzhu/now"
	"strconv"
	"time"
)

func TimeBetweenDates(startTime, endTime string) []time.Time {
	timeStartNeedsInit := true
	timeEndNeedsInit := true

	var timeStart time.Time
	var timeEnd time.Time

	if len(startTime) != 0 {
		parsedInt, err := strconv.ParseInt(startTime, 0, 64)
		if err == nil {
			timeStart = time.Unix(parsedInt, 0)
			timeStartNeedsInit = false
		}
	}

	if len(endTime) != 0 {
		parsedInt, err := strconv.ParseInt(endTime, 0, 64)
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

	return allTimes
}
