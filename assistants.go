package ragflow

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
)

func (c *Client) CreateAssistant(ctx context.Context, req CreateAssistantRequest) (*Assistant, error) {
	trans := &CreateAssistantRequestTransformed{
		Name: req.Name,
		Prompt: Prompt{
			Prompt: req.Prompt,
			EmptyResponse: req.EmptyResponse,
			Variables: req.Variables,
		},
		LLMModel: TransformLLM{
			ModelName: req.LLMModel,
		},
		DatasetIDs: req.DatasetIDs,
	}
	httpReq, err := c.newRequest(ctx, http.MethodPost, "/api/v1/chats", trans)
	if err != nil {
		return nil, err
	}

	var resp Response[Assistant]
	if err := c.do(httpReq, &resp); err != nil {
		return nil, err
	}

	return &resp.Data, nil
}

func (c *Client) GetAssistant(ctx context.Context, assistantID string) (*Assistant, error) {
	endpoint := fmt.Sprintf("/api/v1/chats/%s", assistantID)
	httpReq, err := c.newRequest(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}

	var resp Response[Assistant]
	if err := c.do(httpReq, &resp); err != nil {
		return nil, err
	}

	return &resp.Data, nil
}

func (c *Client) UpdateAssistant(ctx context.Context, assistantID string, req CreateAssistantRequest) (*Assistant, error) {
	trans := &CreateAssistantRequestTransformed{
		Name: req.Name,
		Prompt: Prompt{
			Prompt: req.Prompt,
			EmptyResponse: req.EmptyResponse,
		},
		LLMModel: TransformLLM{
			ModelName: req.LLMModel,
		},
		DatasetIDs: req.DatasetIDs,
	}
	endpoint := fmt.Sprintf("/api/v1/chats/%s", assistantID)
	httpReq, err := c.newRequest(ctx, http.MethodPut, endpoint, trans)
	if err != nil {
		return nil, err
	}

	var resp Response[Assistant]
	if err := c.do(httpReq, &resp); err != nil {
		return nil, err
	}

	return &resp.Data, nil
}

func (c *Client) DeleteAssistant(ctx context.Context, assistantID string) error {
	endpoint := fmt.Sprintf("/api/v1/chats/%s", assistantID)
	httpReq, err := c.newRequest(ctx, http.MethodDelete, endpoint, nil)
	if err != nil {
		return err
	}

	return c.do(httpReq, nil)
}

type ListAssistantsOptions struct {
	Page     int
	PageSize int
	OrderBy  string
	Desc     bool
	Name     string
	ID       string
}

func (c *Client) ListAssistants(ctx context.Context, opts *ListAssistantsOptions) (*Response[[]Assistant], error) {
	params := make(map[string]string)

	if opts != nil {
		if opts.Page > 0 {
			params["page"] = strconv.Itoa(opts.Page)
		}
		if opts.PageSize > 0 {
			params["page_size"] = strconv.Itoa(opts.PageSize)
		}
		if opts.OrderBy != "" {
			params["orderby"] = opts.OrderBy
		}
		if opts.Desc {
			params["desc"] = "true"
		}
		if opts.Name != "" {
			params["name"] = opts.Name
		}
		if opts.ID != "" {
			params["id"] = opts.ID
		}
	}

	url := c.buildURL("/api/v1/chats", params)
	httpReq, err := c.newRequest(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	var resp Response[[]Assistant]
	if err := c.do(httpReq, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

func (c *Client) CreateSession(ctx context.Context, assistantID string, req CreateSessionRequest) (*Session, error) {
	endpoint := fmt.Sprintf("/api/v1/chats/%s/sessions", assistantID)
	httpReq, err := c.newRequest(ctx, http.MethodPost, endpoint, req)
	if err != nil {
		return nil, err
	}

	var resp Response[Session]
	if err := c.do(httpReq, &resp); err != nil {
		return nil, err
	}

	return &resp.Data, nil
}

func (c *Client) GetSession(ctx context.Context, assistantID, sessionID string) (*Session, error) {
	endpoint := fmt.Sprintf("/api/v1/chats/%s/sessions/%s", assistantID, sessionID)
	httpReq, err := c.newRequest(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}

	var resp Response[Session]
	if err := c.do(httpReq, &resp); err != nil {
		return nil, err
	}

	return &resp.Data, nil
}

func (c *Client) UpdateSession(ctx context.Context, assistantID, sessionID string, req UpdateSessionRequest) (*Session, error) {
	endpoint := fmt.Sprintf("/api/v1/chats/%s/sessions/%s", assistantID, sessionID)
	httpReq, err := c.newRequest(ctx, http.MethodPut, endpoint, req)
	if err != nil {
		return nil, err
	}

	var resp Response[Session]
	if err := c.do(httpReq, &resp); err != nil {
		return nil, err
	}

	return &resp.Data, nil
}

func (c *Client) DeleteSession(ctx context.Context, assistantID, sessionID string) error {
	endpoint := fmt.Sprintf("/api/v1/chats/%s/sessions/%s", assistantID, sessionID)
	httpReq, err := c.newRequest(ctx, http.MethodDelete, endpoint, nil)
	if err != nil {
		return err
	}

	return c.do(httpReq, nil)
}

type ListSessionsOptions struct {
	Page     int
	PageSize int
	OrderBy  string
	Desc     bool
	Name     string
	ID       string
}

func (c *Client) ListSessions(ctx context.Context, assistantID string, opts *ListSessionsOptions) (*ListResponse[Session], error) {
	params := make(map[string]string)

	if opts != nil {
		if opts.Page > 0 {
			params["page"] = strconv.Itoa(opts.Page)
		}
		if opts.PageSize > 0 {
			params["page_size"] = strconv.Itoa(opts.PageSize)
		}
		if opts.OrderBy != "" {
			params["orderby"] = opts.OrderBy
		}
		if opts.Desc {
			params["desc"] = "true"
		}
		if opts.Name != "" {
			params["name"] = opts.Name
		}
		if opts.ID != "" {
			params["id"] = opts.ID
		}
	}

	endpoint := fmt.Sprintf("/api/v1/chats/%s/sessions", assistantID)
	url := c.buildURL(endpoint, params)
	httpReq, err := c.newRequest(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	var resp ListResponse[Session]
	if err := c.do(httpReq, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// ChatWithSession converses with a chat assistant using a session
func (c *Client) ChatWithSession(ctx context.Context, chatID string, req ChatWithSessionRequest) (string, error) {
	req.Stream = true
	stream, errChan := c.ChatWithSessionStream(ctx, chatID, req)

	var full string

	for {
		select {
		case ev, ok := <-stream:
			if !ok || ev.Done {
				return full, nil
			}
			full += ev.Delta

		case err := <-errChan:
			if err != nil {
				return "", err
			}
		}
	}
}

func (c *Client) ChatWithSessionStream(ctx context.Context, chatID string, req ChatWithSessionRequest) (<-chan ChatStreamEvent, <-chan error) {

	respChan := make(chan ChatStreamEvent)
	errChan := make(chan error, 1)

	go func() {
		defer close(respChan)
		defer close(errChan)

		endpoint := fmt.Sprintf("/api/v1/chats/%s/completions", chatID)
		req.Stream = true

		httpReq, err := c.newRequest(ctx, http.MethodPost, endpoint, req)
		if err != nil {
			errChan <- err
			return
		}

		resp, err := c.HTTPClient.Do(httpReq)
		if err != nil {
			errChan <- fmt.Errorf("error making request: %w", err)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			bodyBytes, _ := io.ReadAll(resp.Body)
			errChan <- c.handleErrorResponse(resp.StatusCode, bodyBytes)
			return
		}

		scanner := bufio.NewScanner(resp.Body)

		// ⭐ 非常关键：扩大 buffer，避免长 token 被截断
		scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)
		var fullAnswer string
		for scanner.Scan() {
			select {
			case <-ctx.Done():
				errChan <- ctx.Err()
				return
			default:
			}

			line := strings.TrimSpace(scanner.Text())
			if line == "" {
				continue
			}

			// 只处理 SSE data 行
			if !strings.HasPrefix(line, "data:") {
				continue
			}

			payload := strings.TrimSpace(strings.TrimPrefix(line, "data:"))
			var envelope ChatStreamEnvelope

			if err := json.Unmarshal([]byte(payload), &envelope); err != nil {
				continue
			}

			// data: true —— 结束信号
			if string(envelope.Data) == "true" {
				respChan <- ChatStreamEvent{Done: true}
				return
			}
			// 1️⃣ 先判断是否是 error 包
			var errResp struct {
				Code    int    `json:"code"`
				Message string `json:"message"`
			}
			if err := json.Unmarshal([]byte(payload), &errResp); err == nil && errResp.Code != 0 {
				errChan <- fmt.Errorf("API error %d: %s", errResp.Code, errResp.Message)
				return
			}

			// 2️⃣ 正常流式响应
			var streamResp ChatStreamData
			if err := json.Unmarshal([]byte(envelope.Data), &streamResp); err != nil {
				// 流式里偶尔有非结构 JSON，直接跳过
				continue
			}

			if streamResp.Answer != "" {
				delta := streamResp.Answer[len(fullAnswer):]
				fullAnswer = streamResp.Answer

				respChan <- ChatStreamEvent{
					Delta:     delta,
					Full:      fullAnswer,
					ID:        streamResp.ID,
					SessionID: streamResp.SessionID,
					Reference: streamResp.Reference,
				}
			}
		}

		if err := scanner.Err(); err != nil {
			errChan <- fmt.Errorf("error reading stream: %w", err)
		}
	}()

	return respChan, errChan
}
