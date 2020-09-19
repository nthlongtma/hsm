package client_test

import (
	"context"
	"hsm/test_client/crypto/v1"
	"testing"
	"time"

	"google.golang.org/grpc"
)

func TestGRPC(t *testing.T) {
	conn, err := grpc.Dial(":9999", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		t.Error(err)
	}
	defer conn.Close()

	c := crypto.NewHSMServiceClient(conn)

	ctx, _ := context.WithTimeout(context.Background(), time.Second)

	res, err := c.Encrypt(ctx, &crypto.EncryptRequest{Type: "N2K_KGK", PlainText: "eyJ1c2VyX3NhZmVfaWQiOiJ0ZXN0X3VzZXJfc2FmZV9pZCIsInR5cGUiOiJpZF9jYXJkIiwiaW1hZ2UiOiJlbmNvZGVkIGJhc2U2NCBpZCBjYXJkIGltYWdlIn0="})
	if err != nil {
		t.Error(err)
	}

	t.Logf("response: %v", res)
}
