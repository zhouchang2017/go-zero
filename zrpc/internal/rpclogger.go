package internal

import (
	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/grpclog"
)

func init() {
	grpclog.SetLoggerV2(logx.GlobalLogger())
}
