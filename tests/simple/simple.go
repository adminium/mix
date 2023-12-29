package simple

import "context"

type SimpleAPI interface {
	Login(ctx context.Context, req LoginReq) (resp LoginResp, err error)
}

type LoginReq struct {
	Username string
	Password string
}

type LoginResp struct {
	Token string
}
