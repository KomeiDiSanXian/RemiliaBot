package service

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"github.com/wdvxdr1123/ZeroBot/message"

	bf1api "github.com/FloatTech/ZeroBot-Plugin/plugin/bfhelper/internal/bf1/api"
	bf1player "github.com/FloatTech/ZeroBot-Plugin/plugin/bfhelper/internal/bf1/player"
	"github.com/FloatTech/ZeroBot-Plugin/plugin/bfhelper/internal/model"
	"github.com/FloatTech/ZeroBot-Plugin/plugin/bfhelper/pkg/global"
	"github.com/FloatTech/ZeroBot-Plugin/plugin/bfhelper/pkg/renderer"
)

// BindAccount 绑定账号
func (s *Service) BindAccount() error {
	id := s.ctx.State["args"].(string)
	// 数据库查询是否绑定
	player, err := s.dao.GetPlayerByQID(s.ctx.Event.UserID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		s.ctx.SendChain(message.At(s.ctx.Event.UserID), message.Text("正在绑定id为 ", id))
		err = s.dao.CreatePlayer(s.ctx.Event.UserID, id)
		if err != nil {
			s.ctx.SendChain(message.At(s.ctx.Event.UserID), message.Text("绑定失败, ERR: 数据库错误"))
			return err
		}
		s.ctx.SendChain(message.At(s.ctx.Event.UserID), message.Text("绑定成功"))
		return nil
	}
	if err != nil {
		return err
	}
	// 绑定的是旧id
	if id == player.DisplayName {
		s.ctx.SendChain(message.At(s.ctx.Event.UserID), message.Text("笨蛋! 你现在绑的就是这个id"))
		return nil
	}
	s.ctx.SendChain(message.At(s.ctx.Event.UserID), message.Text("将原绑定id为 ", player.DisplayName, " 改绑为 ", id))
	err = s.dao.UpdatePlayer(s.ctx.Event.UserID, "", id)
	if err != nil {
		s.ctx.SendChain(message.At(s.ctx.Event.UserID), message.Text("绑定失败, ERR: 数据库错误"))
		return err
	}
	s.ctx.SendChain(message.At(s.ctx.Event.UserID), message.Text("绑定 ", id, " 成功"))
	pid, err := bf1api.GetPersonalID(id)
	if err != nil {
		return err
	}
	err = s.dao.UpdatePlayer(s.ctx.Event.UserID, pid, "")
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) getPlayerID() (string, bool) {
	id := s.ctx.State["regex_matched"].([]string)[1]
	// id 为空就去数据库查
	if id == "" {
		p, err := s.dao.GetPlayerByQID(s.ctx.Event.UserID)
		if err != nil {
			// 查不到或者失败就是无效
			return "", false
		}
		return p.DisplayName, true
	}
	// id 不为空就认为有效
	return id, true
}

func (s *Service) getPlayer(name string) (*model.Player, error) {
	if name == "" {
		return s.dao.GetPlayerByQID(s.ctx.Event.UserID)
	}
	player, err := s.dao.GetPlayerByName(name)
	if err == nil && player.PersonalID != "" {
		return player, nil
	}
	// 有错误或者pid残缺
	pid, err := bf1api.GetPersonalID(name)
	if err != nil {
		return nil, err
	}
	// 残缺
	if player != nil {
		player.PersonalID = pid
		_ = player.Update(global.DB)
		return player, nil
	}

	return &model.Player{DisplayName: name, PersonalID: pid}, nil
}

func (s *Service) sendWeaponInfo(id, class string) error {
	s.ctx.Send("少女折寿中...")
	player, err := s.getPlayer(id)
	if err != nil {
		s.ctx.SendChain(message.At(s.ctx.Event.UserID), message.Text("ERR: 数据库中没有查到该账号! 请检查是否绑定! 如需绑定请使用[.绑定 id], 中括号不需要输入"))
		return err
	}
	log.Println(player == nil)
	weapons, err := bf1player.GetWeapons(player.PersonalID, class)
	if err != nil {
		s.ctx.SendChain(message.At(s.ctx.Event.UserID), message.Text("ERR: 获取武器失败"))
		return err
	}
	txt := "id：" + player.DisplayName + "\n"
	wp := ([]bf1player.Weapons)(*weapons)
	for i := 0; i < 5; i++ {
		txt += fmt.Sprintf("%s\n%s%s\n%s%s\n%s%s\n%s%s\n%s%s\n%s%s\n",
			"---------------",
			"武器名：", wp[i].Name,
			"击杀数：", strconv.FormatFloat(wp[i].Kills, 'f', 0, 64),
			"准度：", wp[i].Accuracy,
			"爆头率：", wp[i].Headshots,
			"KPM：", wp[i].KPM,
			"效率：", wp[i].Efficiency,
		)
	}
	renderer.Txt2Img(s.ctx, txt)
	return nil
}

// GetPlayerRecent 获取玩家最近游玩
func (s *Service) GetPlayerRecent() error {
	id, isVaild := s.getPlayerID()
	if !isVaild {
		s.ctx.SendChain(message.At(s.ctx.Event.UserID), message.Text("ERR: 数据库中没有查到该账号! 请检查是否绑定! 如需绑定请使用[.绑定 id], 中括号不需要输入"))
		return nil
	}
	s.ctx.Send("少女折寿中...")
	recent, err := bf1player.GetBF1Recent(id)
	if err != nil {
		s.ctx.SendChain(message.At(s.ctx.Event.UserID), message.Text("ERR: 获取最近战绩失败"))
		return err
	}
	// 发送最近战绩
	// TODO: 修改为卡片发送
	msg := "id：" + id + "\n"
	for i := range *recent {
		msg += "服务器：" + (*recent)[i].Server[:24] + "\n"
		msg += "地图：" + (*recent)[i].Map + "   (" + (*recent)[i].Mode + ")\n"
		msg += "kd：" + strconv.FormatFloat((*recent)[i].Kd, 'f', -1, 64) + "\n"
		msg += "kpm：" + strconv.FormatFloat((*recent)[i].Kpm, 'f', -1, 64) + "\n"
		msg += "游玩时长：" + strconv.FormatFloat(float64((*recent)[i].Time/60), 'f', -1, 64) + "分钟"
		msg += "\n---------------\n"
	}
	renderer.Txt2Img(s.ctx, msg)
	return nil
}

// GetPlayerStats 获取玩家战绩
func (s *Service) GetPlayerStats() error {
	id, isVaild := s.getPlayerID()
	if !isVaild {
		s.ctx.SendChain(message.At(s.ctx.Event.UserID), message.Text("ERR: 数据库中没有查到该账号! 请检查是否绑定! 如需绑定请使用[.绑定 id], 中括号不需要输入"))
		return nil
	}
	s.ctx.Send("少女折寿中...")
	stat, err := bf1player.GetStats(id)
	if err != nil {
		s.ctx.SendChain(message.At(s.ctx.Event.UserID), message.Text("ERR: 获取玩家战绩失败, 请自行检查id是否正确"))
		return err
	}
	if stat.Rank == "" {
		s.ctx.SendChain(message.At(s.ctx.Event.UserID), message.Text("获取到的部分数据为空，请检查id是否有效"))
		return errors.Errorf("%s stat.Rank is blank", id)
	}
	// 发送战绩
	// TODO: 修改为卡片发送, 部分数据不准确，等待更改
	txt := "id：" + id + "\n" +
		"等级：" + stat.Rank + "\n" +
		"技巧值：" + stat.Skill + "\n" +
		"游玩时长：" + stat.TimePlayed + "\n" +
		"总kd：" + stat.TotalKD + "(" + stat.Kills + "/" + stat.Deaths + ")" + "\n" +
		"总kpm：" + stat.KPM + "\n" +
		"准度：" + stat.Accuracy + "\n" +
		"爆头数：" + stat.Headshots + "\n" +
		"胜率：" + stat.WinPercent + "(" + stat.Wins + "/" + stat.Losses + ")" + "\n" +
		"场均击杀：" + stat.KillsPerGame + "\n" +
		"步战kd：" + stat.InfantryKD + "\n" +
		"步战击杀：" + stat.InfantryKills + "\n" +
		"步战kpm：" + stat.InfantryKPM + "\n" +
		"载具击杀：" + stat.VehicleKills + "\n" +
		"载具kpm：" + stat.VehicleKPM + "\n" +
		"近战击杀：" + stat.DogtagsTaken + "\n" +
		"最高连杀：" + stat.HighestKillStreak + "\n" +
		"最远爆头：" + stat.LongestHeadshot + "\n" +
		"MVP数：" + stat.MVP + "\n" +
		"作为神医拉起了 " + stat.Revives + " 人" + "\n" +
		"开棺材车创死了 " + stat.CarriersKills + " 人"
	renderer.Txt2Img(s.ctx, txt)
	return nil
}

// GetPlayerWeapon 获取玩家武器
func (s *Service) GetPlayerWeapon() error {
	str := strings.Split(s.ctx.State["regex_matched"].([]string)[1], " ")
	var id string
	// 相当于只输入 .武器
	if str[0] == "" {
		return s.sendWeaponInfo(id, bf1player.ALL)
	}
	if len(str) > 1 {
		id = str[1]
	}
	switch str[0] {
	// 除default 相当于输入 .武器 class id
	case "半自动", "semi":
		return s.sendWeaponInfo(id, bf1player.Semi)
	case "冲锋枪", "冲锋":
		return s.sendWeaponInfo(id, bf1player.SMG)
	case "轻机枪", "机枪":
		return s.sendWeaponInfo(id, bf1player.LMG)
	case "步枪", "狙击枪", "狙击":
		return s.sendWeaponInfo(id, bf1player.Bolt)
	case "霰弹枪", "散弹枪", "霰弹", "散弹":
		return s.sendWeaponInfo(id, bf1player.Shotgun)
	case "配枪", "手枪", "副手":
		return s.sendWeaponInfo(id, bf1player.Sidearm)
	case "近战", "刀":
		return s.sendWeaponInfo(id, bf1player.Melee)
	case "手榴弹", "手雷", "雷":
		return s.sendWeaponInfo(id, bf1player.Grenade)
	case "驾驶员", "坦克兵", "载具":
		return s.sendWeaponInfo(id, bf1player.Dirver)
	case "配备", "装备":
		return s.sendWeaponInfo(id, bf1player.Gadget)
	case "精英", "精英兵":
		return s.sendWeaponInfo(id, bf1player.Elite)
	default:
		// 相当于 .武器 id
		if regexp.MustCompile(`\w+`).MatchString(str[0]) {
			id = str[0]
			return s.sendWeaponInfo(id, bf1player.ALL)
		}
		s.ctx.SendChain(message.At(s.ctx.Event.UserID), message.Text("ERR: 获取玩家武器失败, 不能识别的输入格式"))
		return nil
	}
}

// GetPlayerVehicle 获取玩家载具信息
func (s *Service) GetPlayerVehicle() error {
	id := s.ctx.State["regex_matched"].([]string)[1]
	s.ctx.Send("少女折寿中...")
	player, err := s.getPlayer(id)
	if err != nil {
		s.ctx.SendChain(message.At(s.ctx.Event.UserID), message.Text("ERR: 数据库中没有查到该账号! 请检查是否绑定! 如需绑定请使用[.绑定 id], 中括号不需要输入"))
		return err
	}
	car, err := bf1player.GetVehicles(player.PersonalID)
	if err != nil {
		s.ctx.SendChain(message.At(s.ctx.Event.UserID), message.Text("ERR: 获取玩家载具失败, 请自行检查id是否正确"))
		return err
	}
	msg := "id：" + player.DisplayName + "\n"
	for i := range *car {
		msg += "------------\n"
		msg += (*car)[i].Name + "\n"
		msg += fmt.Sprintf("%s%6.0f\t", "击杀数：", (*car)[i].Kills)
		msg += "kpm：" + (*car)[i].KPM + "\n"
		msg += fmt.Sprintf("%s%6.0f\t", "击毁数：", (*car)[i].Destroyed)
		msg += "游玩时间：" + (*car)[i].Time + " 小时\n"
	}
	renderer.Txt2Img(s.ctx, msg)
	return nil
}

// GetBF1Exchange 获取BF1本期交换信息
func (s *Service) GetBF1Exchange() error {
	exchange, err := bf1api.GetExchange()
	if err != nil {
		s.ctx.SendChain(message.Reply(s.ctx.Event.MessageID), message.Text("ERR: 获取交换失败"))
		return err
	}
	var msg string
	for i, v := range exchange {
		msg += i + ": \n"
		for _, skin := range v {
			msg += "\t" + skin + "\n"
		}
	}
	renderer.Txt2Img(s.ctx, msg)
	return nil
}

// GetBF1OpreationPack 获取本期行动包信息
func (s *Service) GetBF1OpreationPack() error {
	pack, err := bf1api.GetCampaignPacks()
	if err != nil {
		s.ctx.SendChain(message.Reply(s.ctx.Event.MessageID), message.Text("ERR: 获取行动包失败"))
		return err
	}
	var msg string
	msg += "行动名：" + pack.Name + "\n"
	msg += "剩余时间：" + fmt.Sprintf("%.2f", float64(pack.RemainTime)/60/24) + " 天\n"
	msg += "箱子重置时间：" + fmt.Sprintf("%.2f", float64(pack.ResetTime)/60) + " 小时\n"
	msg += "行动地图：" + pack.Op1Name + " 与 " + pack.Op2Name + "\n"
	msg += "行动简介：" + pack.Desc
	renderer.Txt2Img(s.ctx, msg)
	return nil
}
