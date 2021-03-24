package rpc

import (
	"github.com/mhchlib/mconfig-api/api/v1/server"
	"google.golang.org/grpc"
)

// InitRpc ...
func InitRpc(s *grpc.Server) {
	server.RegisterMConfigServer(s, NewMConfigServer())
}
