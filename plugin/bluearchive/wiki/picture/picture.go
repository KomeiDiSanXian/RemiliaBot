// Package picture wiki图片相关
package picture

import (
	"errors"
	"io"
	"os"
	"strings"

	"github.com/FloatTech/ZeroBot-Plugin/plugin/bluearchive/web"
)

// Picture 存储图片信息
type Picture struct {
	Name string
	URL  string
}

// NewPictureByURL 使用url创建图片信息，并基于url命名
func NewPictureByURL(url string) *Picture {
	s := strings.Split(url, "/")
	return &Picture{
		Name: s[len(s)-1],
		URL:  url,
	}
}

// Download 将会下载 Picture到downloadTo 路径
func (p *Picture) Download(downloadTo string) error {
	if p.URL == "" {
		return errors.New("picture url not found")
	}
	link := "http:" + p.URL
	resp, err := web.MakeRequest(link, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	f, err := os.Create(downloadTo + "/" + p.Name)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = io.Copy(f, resp.Body)
	return err
}

// NamedByURL 使用URL的最后一个path 给图片命名
func (p *Picture) NamedByURL() {
	s := strings.Split(p.URL, "/")
	p.Name = s[len(s)-1]
}
