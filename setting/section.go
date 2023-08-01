// Package setting 配置相关
package setting

// BotSettings Bot配置信息
//
// TODO: web socket settings should surpport multiple accounts
type BotSettings struct {
	Debug          bool     // 开启debug 级别日志
	Warn           bool     // 开启warn 级别日志
	AccessToken    string   // Access token of ws client
	WSClientURL    string   // URL of ws client
	NickNames      []string // 机器人昵称
	CommandPrefix  string   // 命令前缀
	Latency        uint     // 响应延迟 (ms)
	RingSize       uint     // 事件环长度
	MaxProcessTime uint     // 最大处理时间 (min)
	SuperUser      []int64  // 超级管理员
	WebUIURL       string   // webui url
	MarkMessage    bool     // 自动标记消息为已读
}

// 存储设置键值对
var sections = make(map[string]interface{})

// ReadSection 读取设置
func (s *Setting) ReadSection(k string, v interface{}) error {
	if err := s.vp.UnmarshalKey(k, v); err != nil {
		return err
	}
	if _, ok := sections[k]; !ok {
		sections[k] = v
	}
	return nil
}

// ReloadAllSections 重新加载所有设置
func (s *Setting) ReloadAllSections() error {
	for k, v := range sections {
		err := s.ReadSection(k, v)
		if err != nil {
			return err
		}
	}
	return nil
}
