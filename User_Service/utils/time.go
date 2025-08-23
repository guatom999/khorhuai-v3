package utils

import "time"

func GetLocalTime() time.Time {

	now := time.Now()

	location, _ := time.LoadLocation("Asia/Bangkok")

	return now.In(location)

}
