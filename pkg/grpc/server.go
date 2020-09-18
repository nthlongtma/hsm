package grpc

import (
	"hsm/configs"

	"github.com/gemalto/pkcs11"
)

type (
	Server struct {
		conf *configs.Config
		ctx  *pkcs11.Ctx
		ss   pkcs11.SessionHandle
	}
)
