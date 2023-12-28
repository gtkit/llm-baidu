package baidu

import (
	"context"
	"testing"
)

func TestCreateAccessToken(t *testing.T) {
	ctx := context.Background()

	client := NewClient("CqArY3Y4IQMFFkP04U9l4NRI", "BnW6GfYNqMGNmXrvS8sGg2NEu71Oohnr", false)
	resp, err := client.CreateAccessToken(ctx)

	if err != nil {
		println(err.Error())
	}

	println("AccessToken: ", resp.AccessToken)

}
