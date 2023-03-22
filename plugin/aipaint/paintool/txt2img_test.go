package paintool

import (
	"log"
	"testing"
)

func TestGetImg(t *testing.T) {
	p := "masterpiece, best quality, (masterpiece:1.21), (best quality:1.331), (ultra-detailed:1.20), (illustration:1.21), (disheveled hair:1.21), (frills:1.21), (1girl:1.1), (solo:1.1), dynamic angle, big top sleeves, floating, beautiful detailed sky, beautiful detailed water, beautiful detailed eyes, overexposure, (fist:1.1), expressionless, side blunt bangs, hair between eyes, ribbons, bowties, buttons, bare shoulders, (small breast:1.331), detailed wet clothes, blank stare, pleated skirt, flowers"
	np := "lowres, bad anatomy, bad hands, text, error, missing fingers, extra digit, fewer digits, cropped, worst quality, low quality, normal quality, jpeg artifacts, signature, watermark, username, blurry, lowres, text, cropped, worst quality, low quality, normal quality, jpeg artifacts, signature, watermark, username, blurry, nsfw, lowres, bad anatomy, bad hands, text, error, missing fingers, extra digit, fewer digits, cropped, worst quality, low quality, normal quality, jpeg artifacts, signature, watermark, username, blurry, missing arms, long neck, humpbacked"
	s, err := GetTxt2Img(NewDefaultTxt2Img(p, np))
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(*(s[0]))
}
