// txt to img
package paintool

import "github.com/tidwall/gjson"

// txt to image request body structure
type Txt2ImgReqBody struct {
	EnableHires                       bool     `json:"enable_hr"`
	DenoisingStrength                 int      `json:"denoising_strength"`
	HrScale                           int      `json:"hr_scale"`
	HrUpscaler                        string   `json:"hr_upscaler"`
	HrSecondPassSteps                 int      `json:"hr_second_pass_steps"`
	Prompt                            string   `json:"prompt"`
	Seed                              int      `json:"seed"`
	SamplerName                       string   `json:"sampler_name"`
	Steps                             int      `json:"steps"`
	CfgScale                          int      `json:"cfg_scale"`
	Width                             int      `json:"width"`
	Height                            int      `json:"height"`
	RestoreFaces                      bool     `json:"restore_faces"`
	Tiling                            bool     `json:"tiling"`
	NegativePrompt                    string   `json:"negative_prompt"`
	Eta                               int      `json:"eta"`
	OverrideSettings                  settings `json:"override_settings"`
	OverrideSettingsRestoreAfterwards bool     `json:"override_settings_restore_afterwards"`
	SamplerIndex                      string   `json:"sampler_index"`
}

type settings struct {
	FilterNSFW bool `json:"filter_nsfw"`
}

// New a Txt2ImgReqBody,which sampler is DPM++ 2M Karras,steps 20, cfg 7, 512*512px.
// nsfw will be filtered out
func NewDefaultTxt2Img(prompt, negprompt string) *Txt2ImgReqBody {
	return &Txt2ImgReqBody{
		EnableHires:                       false,
		Prompt:                            prompt,
		NegativePrompt:                    negprompt,
		Seed:                              -1,
		SamplerName:                       samplerDPMPP2MK,
		Steps:                             20,
		CfgScale:                          7,
		Width:                             512,
		Height:                            512,
		RestoreFaces:                      false,
		Tiling:                            false,
		Eta:                               0,
		OverrideSettings:                  settings{FilterNSFW: true},
		OverrideSettingsRestoreAfterwards: false,
		SamplerIndex:                      "Euler",
	}
}

// Send a request body then get an array of base64 string.
func GetTxt2Img(body *Txt2ImgReqBody) (b64Img []*string, err error) {
	paint := NewPaint()
	reqBody, err := NewTxt2ImgRequest(body)
	if err != nil {
		return nil, err
	}
	paint, err = paint.SendRequest(reqBody)
	if err != nil {
		return nil, err
	}
	paint.ParseRespToGjson("images").ForEach(func(_, value gjson.Result) bool {
		b64Img = append(b64Img, &value.Str)
		return true
	})
	return
}
