// Package aipaint 本地部署的ai画图
package aipaint

import (
	ctrl "github.com/FloatTech/zbpctrl"
	"github.com/FloatTech/zbputils/control"
	"github.com/FloatTech/zbputils/ctxext"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"

	"github.com/FloatTech/ZeroBot-Plugin/plugin/aipaint/paintool"
)

var engine = control.Register("ai画图", &ctrl.Options[*zero.Ctx]{
	DisableOnDefault:  false,
	Brief:             "本地部署的ai画图，平均请求时间为10s",
	Help:              "",
	PrivateDataFolder: "aipaint",
}).ApplySingle(ctxext.DefaultSingle)

func init() {
	engine.OnPrefix(".tti").SetBlock(true).Handle(func(ctx *zero.Ctx) {
		prompt := ctx.State["args"].(string)
		ctx.SendChain(message.Text("少女折寿中..."))
		imgs, err := paintool.GetTxt2Img(paintool.NewDefaultTxt2Img(prompt, paintool.DefaultNegtivePrompt))
		if err != nil {
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("绘图失败，请稍后再试"))
			return
		}
		for _, v := range imgs {
			if id := ctx.Send(message.ReplyWithMessage(ctx.Event.MessageID, message.Image("base64://"+*v))); id.ID() == 0 {
				ctx.SendChain(message.At(ctx.Event.UserID), message.Text("ERROR:可能被风控了"))
			}
		}
	})
}
