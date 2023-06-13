// Package event 从wiki获取游戏活动
package event

import (
	"fmt"
	"strings"
	"time"

	"github.com/FloatTech/ZeroBot-Plugin/plugin/bluearchive/wiki/picture"
	"github.com/tidwall/gjson"
)

// Event 用于存储活动信息
type Event struct {
	EventName   string
	Description string
	BeginAt     int64
	EndAt       int64
	PictureURL  *picture.Picture
}

// Events 活动信息切片
type Events []*Event

// PrintEvent 输入时间格式，输出字符串切片
//
// 每个字符串都是格式化Event的 开始/结束/剩余时间 等
func (es *Events) PrintEvent(layout string) []string {
	strs := make([]string, 0, len(*es))
	for _, event := range *es {
		h, m, s, isStarted := event.remainingTime()
		if h < 0 {
			continue
		}
		event.fixDescription()
		startTime := time.Unix(event.BeginAt, 0).Format(layout)
		endTime := time.Unix(event.EndAt, 0).Format(layout)
		nonstartfmtstr := "%s\n%s\n开始时间: %s\n结束时间: %s\n距离开始剩余时间: %d 小时 %d 分钟 %d 秒\n"
		startedfmtstr := "%s\n%s\n开始时间: %s\n结束时间: %s\n活动剩余时间: %d 小时 %d 分钟 %d 秒\n"
		if !isStarted {
			strs = append(strs, fmt.Sprintf(nonstartfmtstr, event.EventName, event.Description, startTime, endTime, h, m, s))
		} else {
			strs = append(strs, fmt.Sprintf(startedfmtstr, event.EventName, event.Description, startTime, endTime, h, m, s))
		}

	}
	return strs
}



// fixDescription 删除可能存在的 <br>
func (e *Event) fixDescription() {
	// 删除 <br>
	e.Description = strings.ReplaceAll(e.Description, "<br>", "")
}

// remainingTime 计算活动的剩余时间
//
// 如果活动未开始，输出的剩余时间是距离活动的开始时间
//
// 如果活动进行中，输出的剩余时间是距离活动的结束时间
//
// 剩余时间 比如 3661 会输出 1h1min1s 每个数字单独输出
func (e *Event) remainingTime() (hours, minutes, seconds int64, isStarted bool) {
	now := time.Now().Unix()
	beforeStart := e.BeginAt - now
	// 活动未开始
	if beforeStart > 0 {
		duration := time.Duration(beforeStart) * time.Second
		hours = int64(duration.Hours())
		minutes = int64(duration.Minutes()) % 60
		seconds = beforeStart % 60
		return
	}
	isStarted = true
	remain := e.EndAt - now
	// 活动结束
	if remain < 0 {
		return -1, -1, -1, isStarted
	}
	duration := time.Duration(remain) * time.Second
	hours = int64(duration.Hours())
	minutes = int64(duration.Minutes()) % 60
	seconds = remain % 60
	return
}

// GetEvents 从jsonBytes 中获取活动信息
func (es *Events) GetEvents(jsonBytes []byte) *Events {
	events := gjson.GetBytes(jsonBytes, "data.4.list").Array()
	result := make(Events, 0, len(events))
	for _, value := range events {
		picurl := value.Get("picture").Str
		event := &Event{
			EventName:   value.Get("title").Str,
			Description: value.Get("description").Str,
			BeginAt:     value.Get("begin_at").Int(),
			EndAt:       value.Get("end_at").Int(),
			PictureURL:  picture.NewPictureByURL(picurl),
		}
		result = append(result, event)
	}
	return &result
}