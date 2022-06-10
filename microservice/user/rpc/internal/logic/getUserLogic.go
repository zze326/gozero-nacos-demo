package logic

import (
	"context"
	"fmt"

	"demo/microservice/user/rpc/internal/svc"
	"demo/microservice/user/rpc/types/app"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserLogic {
	return &GetUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetUserLogic) GetUser(in *app.GetUserReq) (*app.GetUserReply, error) {
	name := fmt.Sprintf("zze%d", in.Id)
	return &app.GetUserReply{
		Id:   in.Id,
		Name: name,
	}, nil
}
