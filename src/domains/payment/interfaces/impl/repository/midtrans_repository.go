package repository

import (
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/coreapi"
	"main.go/src/configs/app"
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
		default:
			env = midtrans.Sandbox // default to sandbox if not specified
		}
		// client
		client := &coreapi.Client{
			//ServerKey:  app.GetConfig().Midtrans.ServerKey,
			//ClientKey:  app.GetConfig().Midtrans.ClientKey,
			//Env:        env,
			//HttpClient: midtrans.GetHttpClient(env),
			//Options:    &midtrans.ConfigOptions{},
		}
		client.New(app.GetConfig().Midtrans.ServerKey, env)

		Midtrans = &midtran{
			CoreApi: client,
			Env:     env,
		}
	}
	return Midtrans
}

func (s *midtran) CreateChargeRequest(req *coreapi.ChargeReq) (*coreapi.ChargeResponse, *midtrans.Error) {
	return s.CoreApi.ChargeTransaction(req)
}
