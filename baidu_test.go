package baidu

import (
	"log/slog"
	"net/http"

	"github.com/gtkit/llm-baidu/internal/test"
)

func setupBaiduAITestServer() (client *Client, server *test.ServerTest, teardown func()) {
	server = test.NewTestServer()
	ts := server.AITestServer()
	ts.Start()
	teardown = ts.Close
	config := DefaultConfig(ClientId, ClientSecret, true)
	config.BaseURL = ts.URL + "/rpc/2.0/ai_custom/v1/wenxinworkshop"
	config.AuthAPI = ts.URL + "/oauth/2.0/token"
	client, ok := NewClientWithConfig(config).(*Client)
	if !ok {
		slog.Error("NewClientWithConfig is not ok")
	}

	server.RegisterHandler("/oauth/2.0/token", func(w http.ResponseWriter, r *http.Request) {
		data := `{"access_token":"this-is-my-super-token","expires_in": 30}`
		_, _ = w.Write([]byte(data))
	})

	return
}
