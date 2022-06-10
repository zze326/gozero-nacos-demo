package logic

import (
	"context"
	"demo/microservice/user/rpc/user"

	"demo/microservice/demo/api/internal/svc"
	"demo/microservice/demo/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserLogic {
	return &GetUserLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUserLogic) GetUser(req *types.GetUserReq) (resp *types.GetUserReply, err error) {
	u, err := l.svcCtx.UserRpc.GetUser(l.ctx, &user.GetUserReq{
		Id: uint32(req.ID),
	})

	if err != nil {
		return nil, err
	}

	return &types.GetUserReply{
		ID:   int(u.Id),
		Name: u.Name,
	}, nil
}
