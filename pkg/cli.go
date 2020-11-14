package pkg

import (
	"context"
	"github.com/mhchlib/mconfig-api/api/v1/cli"
)

type MConfigCLI struct {
}

func (M *MConfigCLI) PutMconfig(ctx context.Context, request *cli.PutMconfigRequest, response *cli.PutMconfigResponse) error {

	return nil
}

func NewMConfigCLI() *MConfigCLI {
	return &MConfigCLI{}
}
