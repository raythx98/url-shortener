package reqctx

import (
	"context"

	"github.com/raythx98/gohelpme/tool/reqctx"
)

type IReqCtx interface {
	GetValue(ctx context.Context) *reqctx.Value
}

type ReqCtx struct{}

func New() *ReqCtx {
	return &ReqCtx{}
}

func (r *ReqCtx) GetValue(ctx context.Context) *reqctx.Value {
	return reqctx.GetValue(ctx)
}
