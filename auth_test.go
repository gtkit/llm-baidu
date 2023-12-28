package baidu

import (
	"context"
	"testing"
)

func TestCreateAccessToken(t *testing.T) {
	ctx := context.Background()

	client := NewClient("", "", false)
	resp, err := client.CreateAccessToken(ctx)

	if err != nil {
		println(err.Error())
	}

	println("AccessToken: ", resp.AccessToken)

}
