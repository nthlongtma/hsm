package main

import (
	"testing"
)

const (
	testModulePath = "./module/libsofthsm2.so"
	testSlotID     = 866444829
	testPin        = "12345678"
	testKeyType    = "secret"
	testKeyLabel   = "master-key"
)

func TestInitPKCS11(t *testing.T) {
	ctx, err := InitPKCS11(testModulePath)
	if err != nil {
		t.Error(err)
	}
	FinishPKCS11(ctx)
}

func TestGetSession(t *testing.T) {
	ctx, err := InitPKCS11(testModulePath)
	if err != nil {
		t.Error(err)
	}
	defer FinishPKCS11(ctx)
	ss, err := GetSession(ctx, testSlotID, testPin)
	if err != nil {
		t.Error(err)
	}
	defer FinishSession(ctx, ss)
}
