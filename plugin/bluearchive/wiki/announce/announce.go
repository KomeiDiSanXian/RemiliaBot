// Package announce 从wiki获取游戏公告
package announce

import (
	"fmt"
	"strings"

	"github.com/tidwall/gjson"

	"github.com/FloatTech/ZeroBot-Plugin/plugin/bluearchive/wiki/picture"
)

// Announcement 公告
type Announcement struct {
	Title       string
	Summary     string
	PictureURLs []*picture.Picture
}

// Announcements 公告切片
type Announcements []*Announcement

// PrintAnnouncements 返回每条公告的格式化信息
func (as *Announcements) PrintAnnouncements() []string {
	strs := make([]string, 0, len(*as))
	for _, a := range *as {
		str := fmt.Sprintf("%s\n\n%s...", a.Title, a.Summary)
		strs = append(strs, str)
	}
	return strs
}

// GetAnnouncements 从jsonBytes 中获取公告信息
func (as *Announcements) GetAnnouncements(jsonBytes []byte) *Announcements {
	var announcements []gjson.Result
	gjson.ParseBytes(jsonBytes).Get("data").ForEach(func(_, value gjson.Result) bool {
		if value.Get("module.name").Str == "公告" {
			announcements = value.Get("list").Array()
			return false
		}
		return true
	})
	result := make(Announcements, 0, len(announcements))
	for _, value := range announcements {
		announcement := &Announcement{
			Title:   value.Get("title").Str,
			Summary: value.Get("summary").Str,
		}
		picurls := value.Get("thumb").Str
		if picurls != "" {
			urls := strings.Split(picurls, ",")
			pics := make([]*picture.Picture, 0, len(urls))
			for _, url := range urls {
				pics = append(pics, picture.NewPictureByURL(url))
			}
			announcement.PictureURLs = pics
		}
		result = append(result, announcement)
	}
	return &result
}
