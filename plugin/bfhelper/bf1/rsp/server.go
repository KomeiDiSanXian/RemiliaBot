// Package bf1rsp 战地服务器操作
package bf1rsp

import (
	"errors"
	"fmt"

	bf1api "github.com/FloatTech/ZeroBot-Plugin/plugin/bfhelper/bf1/api"
	"github.com/tidwall/gjson"
)

type server struct {
	SID  string
	GID  string
	PGID string
}

type m struct {
	MapName  string
	ModeName string
}

type maps []m

// NewServer 服务器
func NewServer(sid, gid, pgid string) *server {
	return &server{
		SID:  sid,
		GID:  gid,
		PGID: pgid,
	}
}

// Kick player, reason needs BIG5, return reason and err
func (s *server) Kick(pid, reason string) (string, error) {
	reason = fmt.Sprintf("%s%s", "Remi:", reason)
	if len(reason) > 32 {
		return "", errors.New("理由过长")
	}
	post := NewPostKick(pid, s.GID, reason)
	data, err := bf1api.ReturnJSON(bf1api.NativeAPI, "POST", post)
	if err != nil {
		return "", err
	}
	return gjson.Get(data, "result.reason").Str, err
}

// Ban player, check returned id
func (s *server) Ban(pid string) error {
	post := NewPostBan(pid, s.SID)
	data, err := bf1api.ReturnJSON(bf1api.NativeAPI, "POST", post)
	if err != nil {
		return err
	}
	if gjson.Get(data, "id").Str == "" {
		return errors.New("服务器未发出正确的响应，请稍后再试")
	}
	return nil
}

// Unban player
func (s *server) Unban(pid string) error {
	post := NewPostRemoveBan(pid, s.SID)
	data, err := bf1api.ReturnJSON(bf1api.NativeAPI, "POST", post)
	if err != nil {
		return err
	}
	if gjson.Get(data, "id").Str == "" {
		return errors.New("服务器未发出正确的响应，请稍后再试")
	}
	return nil
}

// ChangeMap will change the map for players
func (s *server) ChangeMap(index int) error {
	post := NewPostChangeMap(s.PGID, index)
	data, err := bf1api.ReturnJSON(bf1api.NativeAPI, "POST", post)
	if err != nil {
		return err
	}
	if gjson.Get(data, "id").Str == "" {
		return errors.New("服务器未发出正确的响应，请稍后再试")
	}
	return nil
}

// GetMaps returns maps
func (s *server) GetMaps() (*maps, error) {
	post := NewPostGetServerInfo(s.GID)
	data, err := bf1api.ReturnJSON(bf1api.NativeAPI, "POST", post)
	if err != nil {
		return nil, err
	}
	if gjson.Get(data, "result").String() == "" {
		return nil, errors.New("服务器gameid可能无效，请更新服务器信息")
	}
	result := gjson.Get(data, "result.rotation").Array()
	if result == nil {
		return nil, errors.New("获取到的地图池为空")
	}
	mp := make(maps, len(result))
	for i, v := range result {
		mp[i] = m{MapName: v.Get("mapPrettyName").Str, ModeName: v.Get("modePrettyName").Str}
	}
	return &mp, nil
}

// GetAdminspid returns pids of admins
func (s *server) GetAdminspid() ([]string, error) {
	post := NewPostRSPInfo(s.SID)
	data, err := bf1api.ReturnJSON(bf1api.NativeAPI, "POST", post)
	if err != nil {
		return nil, err
	}
	result := gjson.Get(data, "result.adminList.#.personaId").Array()
	result = append(result, gjson.Get(data, "result.owner.personaId"))
	strs := make([]string, len(result))
	for i, v := range result {
		strs[i] = v.Str
	}
	return strs, bf1api.Exception(gjson.Get(data, "error.code").Int())
}

// input keywords for map id
/* not compiled
func (s *server) GetMapidByKeywords(keyword string) (int, error) {
	switch keyword{
		case
	}
}
*/
