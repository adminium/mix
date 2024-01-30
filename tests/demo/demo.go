package demo

import "context"

type FullAPI interface {
	Ping(ctx context.Context) error
	Login(ctx context.Context, req LoginReq) (resp LoginResp, err error)
	CustomErr(ctx context.Context) error
}

type LoginReq struct {
	Username string
	Password string
}

type LoginResp struct {
	Token string
}
