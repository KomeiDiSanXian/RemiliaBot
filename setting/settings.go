package setting

import (
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

type Setting struct {
	vp *viper.Viper
}

var (
	BotSetting *BotSettings
)

// NewSetting 设置conf名为config 路径为 ./configs/config.yaml
func NewSetting() (*Setting, error) {
	vp := viper.New()
	vp.SetConfigName("config")
	vp.AddConfigPath("configs/")
	vp.SetConfigType("yaml")

	if err := vp.ReadInConfig(); err != nil {
		return nil, err
	}

	s := &Setting{vp: vp}
	s.WatchSettingChange()
	return s, nil
}

// WatchSettingChange 监听配置文件更新
func (s *Setting) WatchSettingChange() {
	go func() {
		s.vp.WatchConfig()
		s.vp.OnConfigChange(func(in fsnotify.Event) {
			_ = s.ReloadAllSections()
		})
	}()
}
