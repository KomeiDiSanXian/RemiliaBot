// Package kfccrazythursday 疯狂星期四
package kfccrazythursday

import (
	"github.com/FloatTech/floatbox/web"
	ctrl "github.com/FloatTech/zbpctrl"
	"github.com/FloatTech/zbputils/control"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
	"github.com/wdvxdr1123/ZeroBot/utils/helper"
)

const (
	crazyURL = "http://api.jixs.cc/api/wenan-fkxqs/index.php"
)

func init() {
	engine := control.Register("疯狂星期四", &ctrl.Options[*zero.Ctx]{
		DisableOnDefault: false,
		Brief:            "疯狂星期四文案",
		Help:             "疯狂星期四\n",
	})
	engine.OnFullMatch("疯狂星期四").SetBlock(true).Handle(func(ctx *zero.Ctx) {
		data, err := web.GetData(crazyURL)
		if err != nil {
			ctx.SendChain(message.Text("ERROR: ", err))
			return
		}
		ctx.SendChain(message.Text(helper.BytesToString(data)))
	})
}
