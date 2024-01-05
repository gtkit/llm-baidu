package baidu

import (
	"context"
	"testing"
)

func TestChatCompletion(t *testing.T) {
	ctx := context.Background()
	client := NewClient(ClientId, ClientSecret, true)
	prompt := []ChatCompletionMessage{
		{
			Role:    ChatMessageRoleUser,
			Content: "Go 语言发展前景!",
		},
	}
	resp, err := client.CreateChatCompletion(ctx, ChatCompletionRequest{
		Messages:    prompt,
		Temperature: 0.7,
		Stream:      false,
		UserId:      "",
	})

	if err != nil {
		t.Log(err.Error())
	}

	t.Log(resp.ErrorMsg)
	if resp.ErrorCode != 0 {
		t.Log(resp.ErrorMsg)
	}

	t.Log(resp.Result)
}
