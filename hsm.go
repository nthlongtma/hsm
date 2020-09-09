package main

import (
	"bytes"
	"fmt"

	"github.com/gemalto/pkcs11"
)

// Symmetric encryption
func Encryption(ctx *pkcs11.Ctx, ss pkcs11.SessionHandle, key pkcs11.ObjectHandle, mech uint, plainText, iv []byte) ([]byte, error) {
	if err := ctx.EncryptInit(ss, []*pkcs11.Mechanism{pkcs11.NewMechanism(mech, iv)}, key); err != nil {
		return nil, err
	}

	cipher, err := ctx.Encrypt(ss, plainText)
	if err != nil {
		return nil, err
	}

	return cipher, nil
}

// Symmetric encryption
func Decryption(ctx *pkcs11.Ctx, ss pkcs11.SessionHandle, key pkcs11.ObjectHandle, mech uint, cipher, iv []byte) ([]byte, error) {
	if err := ctx.DecryptInit(ss, []*pkcs11.Mechanism{pkcs11.NewMechanism(mech, iv)}, key); err != nil {
		return nil, err
	}

	decrypted, err := ctx.Decrypt(ss, cipher)
	if err != nil {
		return nil, err
	}

	return decrypted, nil
}

func RemoveKey(ctx *pkcs11.Ctx, ss pkcs11.SessionHandle, obj pkcs11.ObjectHandle) error {
	if err := ctx.DestroyObject(ss, obj); err != nil {
		return fmt.Errorf("failed to remove key: %v", err)
	}
	return nil
}

func FindKeys(ctx *pkcs11.Ctx, ss pkcs11.SessionHandle, keyClass uint, tokenLabel string) ([]pkcs11.ObjectHandle, error) {
	searchTemplate := []*pkcs11.Attribute{
		pkcs11.NewAttribute(pkcs11.CKA_LABEL, tokenLabel),
		pkcs11.NewAttribute(pkcs11.CKA_CLASS, keyClass),
	}

	// init
	if err := ctx.FindObjectsInit(ss, searchTemplate); err != nil {
		return nil, fmt.Errorf("failed to init finding key: %v", err)
	}

	// finding
	obj, _, err := ctx.FindObjects(ss, 1)
	if err != nil {
		return nil, fmt.Errorf("failed to find key %v", err)
	}
	if len(obj) == 0 {
		return nil, fmt.Errorf("not found key")
	}

	//final
	if err := ctx.FindObjectsFinal(ss); err != nil {
		return nil, fmt.Errorf("failed to finalize finding key: %v", err)
	}

	return obj, nil
}

// AES key
func CreateKey(ctx *pkcs11.Ctx, ss pkcs11.SessionHandle, tokenLabel string) (pkcs11.ObjectHandle, error) {
	aesKeyTemplate := []*pkcs11.Attribute{
		pkcs11.NewAttribute(pkcs11.CKA_CLASS, pkcs11.CKO_SECRET_KEY), // O
		pkcs11.NewAttribute(pkcs11.CKA_KEY_TYPE, pkcs11.CKK_AES),     // O
		pkcs11.NewAttribute(pkcs11.CKA_LABEL, tokenLabel),
		pkcs11.NewAttribute(pkcs11.CKA_TOKEN, true),
		pkcs11.NewAttribute(pkcs11.CKA_ENCRYPT, true),
		pkcs11.NewAttribute(pkcs11.CKA_DECRYPT, true),
		pkcs11.NewAttribute(pkcs11.CKA_SENSITIVE, true),
		pkcs11.NewAttribute(pkcs11.CKA_EXTRACTABLE, true),
		pkcs11.NewAttribute(pkcs11.CKA_VALUE_LEN, 32),
	}
	return ctx.GenerateKey(ss, []*pkcs11.Mechanism{pkcs11.NewMechanism(pkcs11.CKM_AES_KEY_GEN, nil)}, aesKeyTemplate)
}

func GetContext(modulePath string) (*pkcs11.Ctx, error) {
	ctx := pkcs11.New(modulePath)
	if err := ctx.Initialize(); err != nil {
		return nil, err
	}
	return ctx, nil
}

func FinishContext(ctx *pkcs11.Ctx) {
	ctx.Finalize()
	ctx.Destroy()
}

func GetSession(ctx *pkcs11.Ctx, slotID uint, pin string) (pkcs11.SessionHandle, error) {
	// open session
	ss, err := ctx.OpenSession(slotID, pkcs11.CKF_SERIAL_SESSION|pkcs11.CKF_RW_SESSION)
	if err != nil {
		return 0, fmt.Errorf("failed to open session: %v", err)
	}

	// login
	if err := ctx.Login(ss, pkcs11.CKU_USER, pin); err != nil {
		return 0, fmt.Errorf("failed to login: %v", err)
	}

	return ss, nil
}

func FinishSession(ctx *pkcs11.Ctx, ss pkcs11.SessionHandle) {
	ctx.Logout(ss)
	ctx.CloseSession(ss)
}

// pkcs7Pad right-pads the b slice, so its length becomes the multiply of the blocksize.
func pkcs7Pad(b []byte, blocksize int) ([]byte, error) {
	if blocksize <= 0 {
		panic("invalid block size")
	}
	if b == nil || len(b) == 0 {
		panic("invalid data")
	}
	n := blocksize - (len(b) % blocksize)
	pb := make([]byte, len(b)+n)
	copy(pb, b)
	copy(pb[len(b):], bytes.Repeat([]byte{byte(n)}, n))
	return pb, nil
}

// pkcs7Unpad trims the right most n bytes from the b slice.
func pkcs7Unpad(b []byte, blocksize int) ([]byte, error) {
	if blocksize <= 0 {
		panic("invalid block size")
	}
	if b == nil || len(b) == 0 {
		panic("invalid data")
	}
	if len(b)%blocksize != 0 {
		panic("invalid padding on input")
	}
	c := b[len(b)-1]
	n := int(c)
	if n == 0 || n > len(b) {
		panic("invalid padding on input")
	}
	for i := 0; i < n; i++ {
		if b[len(b)-n+i] != c {
			panic("invalid padding on input")
		}
	}
	return b[:len(b)-n], nil
}
