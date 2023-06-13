// Package utils 工具函数
package utils

import (
	"encoding/json"
	"os"

	"github.com/FloatTech/zbputils/img/text"
	"github.com/tidwall/gjson"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
	"github.com/wdvxdr1123/ZeroBot/utils/helper"
)

// Txt2Img 文字转图片
func Txt2Img(ctx *zero.Ctx, txt string) message.MessageSegment {
	data, err := text.RenderToBase64(txt, text.FontFile, 400, 20)
	if err != nil {
		ctx.Send(message.ReplyWithMessage(ctx.Event.MessageID, message.Text("将文字转换成图片时发生错误")))
	}
	return message.Image("base64://" + helper.BytesToString(data))
}

// L11nJSON 本地化json
func L11nJSON(l11nJSONPath, toBeTranslatedJSONPath string) ([]byte, error) {
	// 读取本地化json
	l11nBytes, err := ReadJSONFile(l11nJSONPath)
	if err != nil {
		return nil, err
	}
	// 读取待翻译json
	jsonBytes, err := ReadJSONFile(toBeTranslatedJSONPath)
	if err != nil {
		return nil, err
	}

	// 解析待翻译json
	var transJSON []map[string]interface{}
	err = json.Unmarshal(jsonBytes, &transJSON)
	if err != nil {
		return nil, err
	}
	// 解析 l11n json
	l11n := gjson.ParseBytes(l11nBytes)

	// 创建字段翻译映射
	translationMap := make(map[string]map[string]string)

	// 遍历l11n json并创建映射
	l11n.ForEach(func(field, translations gjson.Result) bool {
		translation := make(map[string]string)
		translations.ForEach(func(key, value gjson.Result) bool {
			translation[key.String()] = value.String()
			return true
		})
		translationMap[field.String()] = translation
		return true
	})

	// 遍历待翻译json 进行翻译
	for _, trans := range transJSON {
		// 翻译每个字段
		for field, translation := range translationMap {
			value, ok := trans[field].(string)
			if !ok {
				// 如果字段不是字符串类型，跳过翻译
				continue
			}

			translatedValue, found := translation[value]
			if found {
				trans[field] = translatedValue
			} else {
				// 在找不到翻译时，使用原始值作为默认值
				trans[field] = value
			}
		}
	}
	return json.Marshal(transJSON)
}

// ReadJSONFile 从data 中读取json
func ReadJSONFile(path string) ([]byte, error) {
	return os.ReadFile(path)
}
