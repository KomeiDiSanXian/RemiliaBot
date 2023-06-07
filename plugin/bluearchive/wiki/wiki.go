// Package wiki 从 wiki (https://ba.gamekee.com/v1/wiki/index) 获取信息
package wiki

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/FloatTech/ZeroBot-Plugin/plugin/bluearchive/web"

	"github.com/tidwall/gjson"
)

// Event 用于存储活动信息
type Event struct {
	EventName   string
	Description string
	BeginAt     int64
	EndAt       int64
	PictureURL  string
}

// Events 活动信息切片
type Events []*Event

// URL wiki URL
var URL = "https://ba.gamekee.com/v1/wiki/index"

// Headers wiki headers
var Headers = map[string]string{
	"game-alias": "ba",
	"Connection": "close",
}

// PrintEvent 输入时间格式，输出字符串切片
//
// 每个字符串都是格式化Event的 开始/结束/剩余时间 等
func (es *Events) PrintEvent(layout string) []string {
	strs := make([]string, 0, len(*es))
	for _, event := range *es {
		h, m, s, isStarted := event.RemainingTime()
		if h < 0 {
			continue
		}
		event.FixDescription()
		startTime := time.Unix(event.BeginAt, 0).Format(layout)
		endTime := time.Unix(event.EndAt, 0).Format(layout)
		nonstartfmtstr := "%s\n%s\n开始时间: %s\n结束时间: %s\n距离开始剩余时间: %d 小时 %d 分钟 %d 秒"
		startedfmtstr := "%s\n%s\n开始时间: %s\n结束时间: %s\n活动剩余时间: %d 小时 %d 分钟 %d 秒"
		if !isStarted {
			strs = append(strs, fmt.Sprintf(nonstartfmtstr, event.EventName, event.Description, startTime, endTime, h, m, s))
		} else {
			strs = append(strs, fmt.Sprintf(startedfmtstr, event.EventName, event.Description, startTime, endTime, h, m, s))
		}
	}
	return strs
}

// FixDescription 删除可能存在的 <br>
func (e *Event) FixDescription() {
	// 删除 <br>
	e.Description = strings.ReplaceAll(e.Description, "<br>", "")
}

// RemainingTime 计算活动的剩余时间
//
// 如果活动未开始，输出的剩余时间是距离活动的开始时间
//
// 如果活动进行中，输出的剩余时间是距离活动的结束时间
//
// 剩余时间 比如 3661 会输出 1h1min1s 每个数字单独输出
func (e *Event) RemainingTime() (hours, minutes, seconds int64, isStarted bool) {
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

// DownloadPicture 将会下载 Event.PictureURL到downloadTo 路径
//
// 文件名为最后一个 path
func (e *Event) DownloadPicture(downloadTo string) error {
	if e.PictureURL == "" {
		return errors.New("picture not found")
	}
	link := "http:" + e.PictureURL
	resp, err := web.MakeRequest(link, Headers)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	s := strings.Split(link, "/")
	name := s[len(s)-1]
	f, err := os.Create(downloadTo + "/" + name)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = io.Copy(f, resp.Body)
	return err
}

// Request 请求获取wiki中的活动数据，返回Event切片
func Request() (Events, error) {
	resp, err := web.MakeRequest(URL, Headers)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	jsonBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if gjson.GetBytes(jsonBytes, "code").Int() != 0 {
		return nil, errors.New("wiki response code not zero")
	}

	events := gjson.GetBytes(jsonBytes, "data.4.list")
	result := make(Events, 0, len(events.Array()))
	for _, value := range events.Array() {
		wiki := &Event{
			EventName:   value.Get("title").Str,
			Description: value.Get("description").Str,
			BeginAt:     value.Get("begin_at").Int(),
			EndAt:       value.Get("end_at").Int(),
			PictureURL:  value.Get("picture").Str,
		}
		result = append(result, wiki)
	}

	return result, nil
}
