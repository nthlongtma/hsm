package grpc_server

import (
	context "context"
	"encoding/base64"
	"fmt"
	"hsm/configs"
	hsm_api "hsm/pkg/hsm-api"
	"log"
	"net"

	"github.com/gemalto/pkcs11"
	"google.golang.org/grpc"
)

type (
	Server struct {
		conf *configs.Config
		ctx  *pkcs11.Ctx
		ss   pkcs11.SessionHandle
		UnimplementedHSMServiceServer
	}
)

func NewServer(conf *configs.Config, ctx *pkcs11.Ctx, ss pkcs11.SessionHandle) Server {
	return Server{
		conf: conf,
		ctx:  ctx,
		ss:   ss,
	}
}

func (s Server) Start() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", s.conf.Servers.GRPC.Port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	g := grpc.NewServer()
	RegisterHSMServiceServer(g, s)

	log.Printf("start grpc at port: %s", s.conf.Servers.GRPC.Port)
	if err := g.Serve(lis); err != nil {
		panic(err)
	}
}

func (s Server) Stop() {
	hsm_api.FinishSession(s.ctx, s.ss)
	if s.ctx != nil {
		hsm_api.FinishContext(s.ctx)
	}
}

func (s Server) Encrypt(ctx context.Context, req *EncryptRequest) (*EncryptResponse, error) {
	// decode the plain text
	plainText, err := base64.StdEncoding.DecodeString(req.PlainText)
	if err != nil {
		return nil, fmt.Errorf("failed to decode request plainText: %v", err)
	}
	log.Printf("plain text after decode: %s", string(plainText))

	// encrypt
	cipher, err := s.encrypt(pkcs11.CKO_SECRET_KEY, s.conf.HSM.N2kLabel, plainText)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt: %v", err)
	}

	// encode cipher before response back
	eCipher := base64.StdEncoding.EncodeToString(cipher)

	return &EncryptResponse{
		ErrorCode:    "0000",
		ErrorMessage: "success",
		CipherText:   eCipher,
	}, nil
}

// encrypt with pre-generated iv then append the iv to the output.
func (s Server) encrypt(keyClass uint, keyLabel string, plainText []byte) ([]byte, error) {
	obj, err := hsm_api.FindKeys(s.ctx, s.ss, keyClass, keyLabel)
	if err != nil {
		return nil, fmt.Errorf("failed to find key: %v", err)
	}

	iv := hsm_api.GenIV(s.conf.HSM.IVSize)
	cipher, err := hsm_api.Encrypt(s.ctx, s.ss, obj[0], pkcs11.CKM_AES_CBC_PAD, plainText, iv)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt: %v", err)
	}

	// prepend iv to cipher
	ret := make([]byte, len(iv)+len(cipher))
	copy(ret[:len(iv)], iv)
	copy(ret[len(iv):], cipher)

	return ret, nil
}

func (s Server) Decrypt(ctx context.Context, req *DecryptRequest) (*DecryptResponse, error) {
	// decode cipher text
	cipher, err := base64.StdEncoding.DecodeString(req.CipherText)
	if err != nil {
		return nil, fmt.Errorf("failed to decode request cipherText: %v", err)
	}

	// decrypt
	plainText, err := s.decrypt(pkcs11.CKO_SECRET_KEY, s.conf.HSM.N2kLabel, cipher)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt: %v", err)
	}

	// encode plain text before response
	ePlainText := base64.StdEncoding.EncodeToString(plainText)

	return &DecryptResponse{
		ErrorCode:    "0000",
		ErrorMessage: "success",
		PlainText:    ePlainText,
	}, nil
}

// extract the iv from cipher and decrypt.
func (s Server) decrypt(keyClass uint, keyLabel string, cipher []byte) ([]byte, error) {
	obj, err := hsm_api.FindKeys(s.ctx, s.ss, keyClass, keyLabel)
	if err != nil {
		return nil, err
	}

	// extract iv and cipher
	iv := cipher[:s.conf.HSM.IVSize]
	c := cipher[s.conf.HSM.IVSize:]
	plain, err := hsm_api.Decrypt(s.ctx, s.ss, obj[0], pkcs11.CKM_AES_CBC_PAD, c, iv)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt: %v", err)
	}

	return plain, nil
}
