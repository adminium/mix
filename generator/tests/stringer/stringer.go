package tests_stringer

import (
	"context"
	"github.com/shopspring/decimal"
	"time"
)

type StringerAPI interface {
	Test(ctx context.Context, in *Stringer) (out Stringer, err error)
}

type Stringer struct {
	Time     time.Time
	Duration time.Duration
	Decimal  decimal.Decimal
}
