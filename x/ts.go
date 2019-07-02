package x

import (
	"time"
	"wcore_old/x"
)

// GetNowTimeByMilli 获取当前时间的毫秒时间戳
func GetNowTimeByMilli() int64 {
	return time.Now().UnixNano() / 1e6
}

// GetWeekDay 获取本周凌晨时间戳
func GetWeekDay() int64 {
	now := time.Now()
	offset := int(time.Monday - now.Weekday())
	if offset > 0 {
		offset = -6
	}
	weekStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local).AddDate(0, 0, offset)
	return weekStart.UnixNano() / 1e6
}

// GetByZeroMorningTs 获取指定 时间戳当天的零点时间戳
func GetByZeroMorningTs(ts int64) int64 {
	return ts - (ts+8*3600*1000)%x.OneDay
}

// GetCurrDay 获取当天的凌晨时间戳
func GetCurrDay() int64 {
	t := time.Now()
	now := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	return now.UnixNano() / 1e6
}
