// http请求相关
package paintool

import (
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/tidwall/gjson"
)

// paint struct...
type Paint struct {
	Client   *http.Client
	Response *http.Response
}

// NewPaint creates a paint pointer,which contains client
func NewPaint() *Paint {
	return &Paint{
		Client: http.DefaultClient,
	}
}

// Copy will create a copy of the paint pointer
func (p *Paint) Copy() *Paint {
	pp := *p
	return &pp
}

// Generate a request body
func NewTxt2ImgRequest(reqBody *Txt2ImgReqBody) (*http.Request, error) {
	body, err := anyToJSON(reqBody)
	if err != nil {
		return nil, err
	}
	return http.NewRequest("POST", fmt.Sprintf("%s%s", paintURL, txt2img), body)
}

// response body will be written to Paint.Response
func (p *Paint) SendRequest(req *http.Request) (*Paint, error) {
	resp, err := p.Client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, errors.New("response status code is not 200")
	}
	pp := p.Copy()
	pp.Response = resp
	return pp, nil
}

// Serialize response body
func (p *Paint) ParseRespToGjson(path string) gjson.Result {
	body, _ := io.ReadAll(p.Response.Body)
	defer p.Response.Body.Close()
	return gjson.GetBytes(body, path)
}
