// Package dailynews 今日早报
package dailynews

import (
	"github.com/FloatTech/floatbox/binary"
	"github.com/FloatTech/floatbox/web"
	ctrl "github.com/FloatTech/zbpctrl"
	"github.com/FloatTech/zbputils/control"
	"github.com/tidwall/gjson"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
)

const api = "http://dwz.2xb.cn/zaob"

func init() {
	engine := control.Register("早报", &ctrl.Options[*zero.Ctx]{
		DisableOnDefault:  false,
		Brief:             "今日早报",
		Help:              "- 今日早报",
		PrivateDataFolder: "dailynews",
	})

	engine.OnFullMatch(`今日早报`).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			data, err := web.GetData(api)
			if err != nil {
				ctx.SendChain(message.Text("ERROR: ", err))
				return
			}
			picURL := gjson.Get(binary.BytesToString(data), "imageUrl").String()
			ctx.SendChain(message.Image(picURL))
		})
}
