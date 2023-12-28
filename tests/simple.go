package tests

import "context"

type SimpleAPI interface {
	Login(ctx context.Context, req LoginReq) (resp LoginResp, err error)
}

type LoginReq struct {
}

type LoginResp struct {
}
