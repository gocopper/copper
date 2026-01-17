package clifecycletest

import (
	"github.com/gocopper/copper/clifecycle"
	"github.com/gocopper/copper/clogger"
)

func New() *clifecycle.Lifecycle {
	return clifecycle.New(clogger.NewNoop())
}
