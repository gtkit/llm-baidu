package baidu

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"testing"

	"github.com/gtkit/llm-baidu/internal/test/checks"
)

func TestCreateChatCompletionRealServer(t *testing.T) {
	client := NewClient(ClientId, ClientSecret, true)
	resp, err := client.CreateChatCompletion(context.Background(), ChatCompletionRequest{
		Messages: []ChatCompletionMessage{
			{
				Role:    ChatMessageRoleUser,
				Content: "Hello!",
			},
		},
		Stream: false,
	}, "eb-instant")
	checks.NoError(t, err, "CreateCompletionStream returned error")

	// println(resp.ErrorMsg)
	// println(resp.Result)

	t.Log(resp.ErrorMsg)
	t.Log(resp.Result)
}

// TestCreateChatCompletionStreamOnRealServer.
func TestCreateChatCompletionStreamOnRealServer(t *testing.T) {
	client := NewClient(ClientId, ClientSecret, true)
	// fmt.Printf("---client---- %+v\n", client)
	slog.Info("new client info")

	prompt := []ChatCompletionMessage{
		{
			Role:    ChatMessageRoleUser,
			Content: "Go 语言发展前景!",
		},
	}

	stream, err := client.CreateChatCompletionStream(context.Background(), ChatCompletionRequest{
		Messages:    prompt,
		Temperature: 1,
		Stream:      true,
	})
	checks.NoError(t, err, "CreateCompletionStream returned error")
	defer stream.Close()

	for {
		response, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			// fmt.Printf("\nStream finished: %d %s\n", response.ErrorCode, response.ErrorMsg)
			return
		}

		if err != nil {
			t.Logf("\nStream error: %v\n", err)
			return
		}

		t.Logf("error: %s\n", response.ErrorMsg)
		t.Logf("resp: %s\n", response.Result)
	}
}

func TestCreateChatCompletionStream(t *testing.T) {
	client, server, teardown := setupBaiduAITestServer()
	defer teardown()
	server.RegisterHandler("/rpc/2.0/ai_custom/v1/wenxinworkshop/chat/completions", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")

		// Send test responses
		dataBytes := []byte{}
		dataBytes = append(dataBytes, []byte("event: message\n")...)

		data := `{"id":"1","object":"chat.completion","created":1598069254,"result":"response1"}`
		dataBytes = append(dataBytes, []byte("data: "+data+"\n\n")...)

		dataBytes = append(dataBytes, []byte("event: message\n")...)

		data = `{"id":"2","object":"chat.completion","created":1598069255,"model":"gpt-3.5-turbo","result":"response2"}`
		dataBytes = append(dataBytes, []byte("data: "+data+"\n\n")...)

		dataBytes = append(dataBytes, []byte("event: done\n")...)
		dataBytes = append(dataBytes, []byte("data: [DONE]\n\n")...)

		_, err := w.Write(dataBytes)
		checks.NoError(t, err, "Write error")
	})

	stream, err := client.CreateChatCompletionStream(context.Background(), ChatCompletionRequest{
		Messages: []ChatCompletionMessage{
			{
				Role:    ChatMessageRoleUser,
				Content: "Hello!",
			},
		},
		Stream: true,
	})
	checks.NoError(t, err, "CreateCompletionStream returned error")
	defer stream.Close()

	expectedResponses := []ChatCompletionResponse{
		{
			ID:      "1",
			Object:  "chat.completion",
			Created: 1598069254,
			Result:  "response1",
		},
		{
			ID:      "2",
			Object:  "chat.completion",
			Created: 1598069255,
			Result:  "response2",
		},
	}

	for ix, expectedResponse := range expectedResponses {
		b, _ := json.Marshal(expectedResponse)
		t.Logf("%d: %s", ix, string(b))

		receivedResponse, streamErr := stream.Recv()
		checks.NoError(t, streamErr, "stream.Recv() failed")
		if !compareChatResponses(expectedResponse, receivedResponse) {
			t.Errorf("Stream response %v is %v, expected %v", ix, receivedResponse, expectedResponse)
		}
	}

	_, streamErr := stream.Recv()
	if !errors.Is(streamErr, io.EOF) {
		t.Errorf("stream.Recv() did not return EOF in the end: %v", streamErr)
	}

	_, streamErr = stream.Recv()

	checks.ErrorIs(t, streamErr, io.EOF, "stream.Recv() did not return EOF when the stream is finished")
	if !errors.Is(streamErr, io.EOF) {
		t.Errorf("stream.Recv() did not return EOF when the stream is finished: %v", streamErr)
	}
}

func TestCreateChatCompletionStreamError(t *testing.T) {
	client, server, teardown := setupBaiduAITestServer()
	defer teardown()
	server.RegisterHandler("/v1/chat/completions", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")

		// Send test responses
		dataBytes := []byte{}
		dataStr := []string{
			`{"error_code": 1, "error_msg": "Unknown error"}`,
		}
		for _, str := range dataStr {
			dataBytes = append(dataBytes, []byte(str+"\n")...)
		}

		_, err := w.Write(dataBytes)
		checks.NoError(t, err, "Write error")
	})

	stream, err := client.CreateChatCompletionStream(context.Background(), ChatCompletionRequest{
		Messages: []ChatCompletionMessage{
			{
				Role:    ChatMessageRoleUser,
				Content: "Hello!",
			},
		},
		Stream: true,
	})
	checks.NoError(t, err, "CreateCompletionStream returned error")
	defer stream.Close()

	_, streamErr := stream.Recv()
	checks.HasError(t, streamErr, "stream.Recv() did not return error")

	var apiErr *APIError
	if !errors.As(streamErr, &apiErr) {
		t.Errorf("stream.Recv() did not return APIError")
	}
	t.Logf("%+v\n", apiErr)
}

// Helper funcs.
func compareChatResponses(r1, r2 ChatCompletionResponse) bool {
	if r1.ID != r2.ID || r1.Object != r2.Object || r1.Created != r2.Created || r1.Result != r2.Result {
		return false
	}

	return true
}
