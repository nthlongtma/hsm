package test_client

import (
	context "context"
	"testing"
	"time"

	grpc "google.golang.org/grpc"
)

func TestEncrypt(t *testing.T) {
	// Set up a connection to the server.
	conn, err := grpc.Dial(":9999", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		t.Errorf("did not connect: %v", err)
	}
	defer conn.Close()
	c := NewHSMServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := c.Encrypt(ctx, &EncryptRequest{Type: "N2K_KGK", PlainText: "eyJ1c2VyX3NhZmVfaWQiOiJ0ZXN0X3VzZXJfc2FmZV9pZCIsInR5cGUiOiJpZF9jYXJkIiwiaW1hZ2UiOiJlbmNvZGVkIGJhc2U2NCBpZCBjYXJkIGltYWdlIn0="})
	if err != nil {
		t.Errorf("could not encrypt: %v", err)
	}
	t.Logf("response: %+v", r)
}
