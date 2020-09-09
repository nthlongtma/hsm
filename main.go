package main

import (
	"fmt"

	"github.com/gemalto/pkcs11"
)

const (
	configPath = "./configs"
	configName = "dev"

	keyTypeSecret = "secret"
	keyTypeRSA    = "rsa-keypair"
)

func main() {
	// read config
	conf := LoadConfig(configPath, configName)
	if conf == nil {
		panic("failed to load config")
	}
	fmt.Printf("config: %+v\n", *conf)

	// init pkcs11
	ctx, err := InitPKCS11(conf.ModulePath)
	if err != nil {
		panic(fmt.Sprintf("failed to init pkcs11 module: %v", err))
	}
	defer FinishPKCS11(ctx)

	// get slot list
	slots, err := ctx.GetSlotList(true)
	if err != nil {
		fmt.Println("failed to get slot list: ", err)
		return
	}
	fmt.Printf("slot list: %+v\n", slots)

	// get info
	info, err := ctx.GetInfo()
	if err != nil {
		fmt.Println("failed to get info: ", err)
	}
	fmt.Printf("info: %+v", info)

	// get session
	ss, err := GetSession(ctx, conf.HSM.SlotID, conf.HSM.Pin)
	if err != nil {
		fmt.Printf("failed to get session: %v\n", err)
		return
	}
	defer FinishSession(ctx, ss)
	fmt.Println("login successfully.")

	// create key
	oh, err := CreateAESKey(ctx, *ss, conf.HSM.TokenLabel)
	if err != nil {
		fmt.Println("failed to create AES key: ", err)
	}

	fmt.Printf("object handler: %+v\n", oh)

}

// func CreateKey(ctx *pkcs11.Ctx, ss pkcs11.SessionHandle, keyType, tokenLabel string) {
// 	switch keyType {
// 	case keyTypeSecret:
// 		CreateAESKey(ctx, ss, tokenLabel)
// 	case keyTypeRSA:
// 		CreateRSAKeyPair(ctx, ss, tokenLabel)
// 	default:
// 		fmt.Println("not support key_type")
// 	}
// }

func CreateAESKey(ctx *pkcs11.Ctx, ss pkcs11.SessionHandle, tokenLabel string) (pkcs11.ObjectHandle, error) {
	// value := make([]byte, 32)
	secretKeyTemplate := []*pkcs11.Attribute{
		pkcs11.NewAttribute(pkcs11.CKA_CLASS, pkcs11.CKO_SECRET_KEY), //
		pkcs11.NewAttribute(pkcs11.CKA_KEY_TYPE, pkcs11.CKK_AES),
		pkcs11.NewAttribute(pkcs11.CKA_LABEL, tokenLabel),
		pkcs11.NewAttribute(pkcs11.CKA_TOKEN, true),
		pkcs11.NewAttribute(pkcs11.CKA_ENCRYPT, true),
		// pkcs11.NewAttribute(pkcs11.CKA_VALUE, value),
		pkcs11.NewAttribute(pkcs11.CKA_VALUE_LEN, 32),
		// pkcs11.NewAttribute(pkcs11.CKA_VERIFY, true),
	}
	return ctx.GenerateKey(ss, []*pkcs11.Mechanism{pkcs11.NewMechanism(pkcs11.CKM_AES_KEY_GEN, nil)}, secretKeyTemplate)
}

func CreateRSAKeyPair(ctx *pkcs11.Ctx, ss pkcs11.SessionHandle, tokenLabel string) {

}

func InitPKCS11(modulePath string) (*pkcs11.Ctx, error) {
	ctx := pkcs11.New(modulePath)
	if err := ctx.Initialize(); err != nil {
		return nil, err
	}
	return ctx, nil
}

func FinishPKCS11(ctx *pkcs11.Ctx) {
	ctx.Finalize()
	ctx.Destroy()
}

func GetSession(ctx *pkcs11.Ctx, slotID uint, pin string) (*pkcs11.SessionHandle, error) {
	// open session
	ss, err := ctx.OpenSession(slotID, pkcs11.CKF_SERIAL_SESSION|pkcs11.CKF_RW_SESSION)
	if err != nil {
		return nil, fmt.Errorf("failed to open session: %v", err)
	}

	// login
	if err := ctx.Login(ss, pkcs11.CKU_USER, pin); err != nil {
		return nil, fmt.Errorf("failed to login: %v", err)
	}

	return &ss, nil
}

func FinishSession(ctx *pkcs11.Ctx, ss *pkcs11.SessionHandle) {
	ctx.Logout(*ss)
	ctx.CloseSession(*ss)
}
