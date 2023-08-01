// Package bfhelper 战地玩家查询
package bfhelper

import (
	"github.com/FloatTech/ZeroBot-Plugin/plugin/bfhelper/internal/engine"
	"github.com/FloatTech/ZeroBot-Plugin/plugin/bfhelper/internal/handler"
	"github.com/FloatTech/ZeroBot-Plugin/plugin/bfhelper/internal/rule"
)

func init() {
	// QQ绑定ID
	engine.Engine.OnPrefixGroup([]string{".绑定", ".bind"}, rule.Initialized()).SetBlock(true).Handle(handler.BindAccountHandler())
	// bf1个人战绩
	engine.Engine.OnRegex(`\. *1?战绩 *(.*)$`, rule.Initialized()).SetBlock(true).Handle(handler.PlayerStatsHandler())
	// 武器查询，只展示前五个
	engine.Engine.OnRegex(`^\. *1?武器 *(.*)$`, rule.Initialized()).SetBlock(true).Handle(handler.PlayerWeaponHandler())
	// 最近战绩
	engine.Engine.OnRegex(`^\. *1?最近 *(.*)$`, rule.Initialized()).SetBlock(true).Handle(handler.PlayerRecentHandler())
	// 获取所有种类的载具信息
	engine.Engine.OnRegex(`^\. *1?载具 *(.*)$`, rule.Initialized()).SetBlock(true).Handle(handler.PlayerVehicleHandler())
	// 交换查询
	engine.Engine.OnFullMatchGroup([]string{".交换", ".exchange"}, rule.Initialized()).SetBlock(true).Handle(handler.BF1ExchangeHandler())
	// 行动包查询
	engine.Engine.OnFullMatchGroup([]string{".行动", ".行动包", ".pack"}, rule.Initialized()).SetBlock(true).Handle(handler.BF1OpreationPackHandler())
}
