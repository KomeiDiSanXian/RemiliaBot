package paintool

import (
	"fmt"
	"io"
	"net/http"

	"github.com/tidwall/gjson"
)

type Paint struct {
	Client   *http.Client
	Response *http.Response
}

func NewPaint() *Paint {
	return &Paint{
		Client: http.DefaultClient,
	}
}

func (p *Paint) Copy() *Paint {
	pp := *p
	return &pp
}

func NewTxt2ImgRequest(reqBody *Txt2ImgReqBody) (*http.Request, error) {
	body, err := anyToJSON(reqBody)
	if err != nil {
		return nil, err
	}
	return http.NewRequest("POST", fmt.Sprintf("%s%s", paintURL, txt2img), body)
}

func (p *Paint) SendRequest(req *http.Request) (*Paint, error) {
	resp, err := p.Client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("response status code is not 200: %d", resp.StatusCode)
	}
	pp := p.Copy()
	pp.Response = resp
	return pp, nil
}

func (p *Paint) ParseRespToGjson(path string) gjson.Result {
	body, _ := io.ReadAll(p.Response.Body)
	return gjson.GetBytes(body, path)
}
