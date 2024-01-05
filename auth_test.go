package baidu

import (
	"context"
	"log/slog"
	"testing"
)

// https://cloud.baidu.com/doc/WENXINWORKSHOP/s/Ilkkrb0i5.
func TestCreateAccessToken(t *testing.T) {
	ctx := context.Background()

	client := NewClient(ClientId, ClientSecret, false)

	resp, err := client.CreateAccessToken(ctx)

	if err != nil {
		slog.Error("auth error:" + err.Error())
	}

	slog.Info("", "AccessToken", resp)

}
