// Package service 业务逻辑代码
package service

import (
	zero "github.com/wdvxdr1123/ZeroBot"

	"github.com/FloatTech/ZeroBot-Plugin/plugin/bfhelper/internal/dao"
	"github.com/FloatTech/ZeroBot-Plugin/plugin/bfhelper/pkg/global"
)

// Service 业务
type Service struct {
	ctx *zero.Ctx
	dao *dao.Dao
}

// New 新建业务
func New(ctx *zero.Ctx) *Service {
	svc := Service{ctx: ctx}
	svc.dao = dao.New(global.DB)
	return &svc
}
