// Package wiki 从 wiki (https://ba.gamekee.com/v1/wiki/index) 获取信息
package wiki

import (
	"errors"
	"io"

	"github.com/FloatTech/ZeroBot-Plugin/plugin/bluearchive/web"
	"github.com/FloatTech/ZeroBot-Plugin/plugin/bluearchive/wiki/announce"
	"github.com/FloatTech/ZeroBot-Plugin/plugin/bluearchive/wiki/event"

	"github.com/tidwall/gjson"
)

// Data 存储公告和活动
type Data struct {
	Events        *event.Events
	Announcements *announce.Announcements
}

// URL wiki URL
var URL = "https://ba.gamekee.com/v1/wiki/index"

// Headers wiki headers
var Headers = map[string]string{
	"game-alias": "ba",
	"Connection": "close",
}

// NewWikiData 创建空的WikiData，返回其指针
func NewWikiData() *Data {
	return &Data{}
}

// Request 请求获取wiki中的数据
func (w *Data) Request() error {
	resp, err := web.MakeRequest(URL, Headers)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	jsonBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if gjson.GetBytes(jsonBytes, "code").Int() != 0 {
		return errors.New("wiki response code not zero")
	}

	w.Events = w.Events.GetEvents(jsonBytes)
	w.Announcements = w.Announcements.GetAnnouncements(jsonBytes)

	return nil
}
