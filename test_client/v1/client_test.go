package crypto

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"google.golang.org/grpc"
)

func TestGRPC(t *testing.T) {
	conn, err := grpc.Dial(":9999", grpc.WithInsecure(), grpc.WithBlock(), grpc.WithTimeout(time.Second*3))
	if err != nil {
		t.Error(err)
	}
	defer conn.Close()

	c := NewCryptoClient(conn)

	ctx, _ := context.WithTimeout(context.Background(), time.Second)

	res, err := c.Encrypt(ctx, &EncryptRequest{Type: "N2K_KGK", PlainText: "eyJ1c2VyX3NhZmVfaWQiOiJ0ZXN0X3VzZXJfc2FmZV9pZCIsInR5cGUiOiJpZF9jYXJkIiwiaW1hZ2UiOiJlbmNvZGVkIGJhc2U2NCBpZCBjYXJkIGltYWdlIn0="})
	if err != nil {
		t.Error(err)
	}

	t.Logf("response: %v", res)
}

func TestHTTP(t *testing.T) {
	req := EncryptRequest{Type: "N2K_KGK", PlainText: "eyJ1c2VyX3NhZmVfaWQiOiJ0ZXN0X3VzZXJfc2FmZV9pZCIsInR5cGUiOiJpZF9jYXJkIiwiaW1hZ2UiOiJlbmNvZGVkIGJhc2U2NCBpZCBjYXJkIGltYWdlIn0="}

	b, err := json.Marshal(req)
	if err != nil {
		t.Error(err)
	}

	httpReq, err := http.NewRequest(http.MethodPost, "http://localhost:8888/api/v1/encrypt", bytes.NewBuffer(b))
	if err != nil {
		t.Error(err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	// send request
	httpRes, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		t.Error(err)
	}

	b, err = ioutil.ReadAll(httpRes.Body)
	if err != nil {
		t.Error(err)
	}
	defer httpRes.Body.Close()

	res := EncryptResponse{}
	if err := json.Unmarshal(b, &res); err != nil {
		t.Error(err)
	}

	t.Logf("response: %+v", res)

}
