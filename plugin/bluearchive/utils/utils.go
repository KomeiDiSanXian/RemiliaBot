// Package utils 工具函数
package utils

import (
	"github.com/FloatTech/zbputils/img/text"
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
