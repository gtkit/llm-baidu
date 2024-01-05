package baidu

import (
	"context"
	"net/http"
)

// CreateChatCompletion — API call to Create a completion for the chat message.
func (c *Client) CreateChatCompletion(
	ctx context.Context,
	request ChatCompletionRequest,
	args ...any,
) (response ChatCompletionResponse, err error) {
	if request.Stream {
		err = ErrChatCompletionStreamNotSupported
		return
	}
	model := ""
	if len(args) > 0 {
		m, ok := args[0].(string)
		if !ok {
			err = ErrChatCompletionInvalidModel
			return
		}
		model = m
	}

	var req *http.Request
	if c.config.AutoAuthToken {
		req, err = c.newRequestWithToken(ctx, http.MethodPost, c.fullURL(model), withBody(request))
	} else {
		req, err = c.newRequest(ctx, http.MethodPost, c.fullURL(model), withBody(request))
	}

	if err != nil {
		return
	}

	err = c.sendRequest(req, &response)
	return
}
