// Package bluearchive 插件主体部分
package bluearchive

import (
	"github.com/FloatTech/ZeroBot-Plugin/plugin/bluearchive/utils"
	"github.com/FloatTech/ZeroBot-Plugin/plugin/bluearchive/wiki"

	ctrl "github.com/FloatTech/zbpctrl"
	"github.com/FloatTech/zbputils/control"
	"github.com/FloatTech/zbputils/ctxext"

	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
)

var engine = control.Register("碧蓝档案", &ctrl.Options[*zero.Ctx]{
	DisableOnDefault:  false,
	Brief:             "ba相关信息查询",
	Help:              "bluearchive\n- .ba活动 查询活动信息",
	PrivateDataFolder: "bluearchive",
})

func init() {
	// 完全匹配触发
	// 使用合并消息转发
	engine.OnFullMatch(".ba活动").SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			events, err := wiki.Request()
			if err != nil {
				ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("ERROR: 请求错误"))
				return
			}
			layout := "1月02日 15:04"
			msg := make(message.Message, 0, len(events))
			eventmsg := events.PrintEvent(layout)
			for _, sendmsg := range eventmsg {
				msg = append(msg, ctxext.FakeSenderForwardNode(ctx, utils.Txt2Img(ctx, sendmsg)))
			}
			if id := ctx.Send(msg).ID(); id == 0 {
				ctx.SendChain(message.Text("ERROR: 可能被风控了"))
			}
		})
}
