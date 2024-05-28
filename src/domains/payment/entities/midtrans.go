package entities

import (
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/coreapi"
	_ "github.com/midtrans/midtrans-go/iris"
	_ "github.com/midtrans/midtrans-go/snap"
	"main.go/configs/app"
	"strings"
)

var (
	Midtrans *midtran
)

type midtran struct {
	CoreApi *coreapi.Client
	Env     midtrans.EnvironmentType
}

func NewMidtrans() *midtran {
	if Midtrans == nil {
		// Midtrans Production or Development
		var env midtrans.EnvironmentType
		switch strings.ToLower(app.GetConfig().Midtrans.Environment) {
		case "sandbox":
			env = midtrans.Sandbox
		case "production":
			env = midtrans.Production
		}
		// client
		client := &coreapi.Client{}
		client.New(app.GetConfig().Midtrans.ServerKey, env)
		Midtrans = &midtran{
			CoreApi: client,
			Env:     env,
		}
	}
	return Midtrans
}
