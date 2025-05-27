package client

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"qwrunner/config"
	"strings"
	"time"
)

type Client struct {
	cfg        *config.Config
	httpClient *http.Client
}

func New(cfg *config.Config) *Client {
	return &Client{
		cfg: cfg,
		httpClient: &http.Client{
			// 移除全局超时，改为使用context控制
			// Timeout: 30 * time.Second,
		},
	}
}

// 非流式请求
func (c *Client) CreateChatCompletion(ctx context.Context, req *ChatRequest) (*ChatResponse, error) {
	req.Stream = false // 确保非流式
	// 为非流式请求设置30秒超时
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	return c.sendRequest(ctx, req)
}

// 流式请求
func (c *Client) CreateChatCompletionStream(ctx context.Context, req *ChatRequest, callback func(*ChatResponse) error) error {
	req.Stream = true
	if req.StreamOptions == nil {
		req.StreamOptions = &StreamOptions{IncludeUsage: true}
	}

	// 为流式请求设置更长的超时时间（5分钟）
	ctx, cancel := context.WithTimeout(ctx, 5*time.Minute)
	defer cancel()

	jsonData, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("JSON序列化失败: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", c.cfg.APIEndpoint, bytes.NewReader(jsonData))
	if err != nil {
		return fmt.Errorf("创建请求失败: %w", err)
	}

	httpReq.Header.Set("Authorization", "Bearer "+c.cfg.APIKey)
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "text/event-stream")
	httpReq.Header.Set("Cache-Control", "no-cache")

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("请求失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		// 读取错误响应体
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API返回错误状态码: %d, 响应: %s", resp.StatusCode, string(body))
	}

	return c.processStreamResponse(ctx, resp.Body, callback)
}

func (c *Client) sendRequest(ctx context.Context, req *ChatRequest) (*ChatResponse, error) {
	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("JSON序列化失败: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", c.cfg.APIEndpoint, bytes.NewReader(jsonData))
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}

	httpReq.Header.Set("Authorization", "Bearer "+c.cfg.APIKey)
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API返回错误状态码: %d, 响应: %s", resp.StatusCode, string(body))
	}

	var chatResp ChatResponse
	if err := json.NewDecoder(resp.Body).Decode(&chatResp); err != nil {
		return nil, fmt.Errorf("响应解析失败: %w", err)
	}

	return &chatResp, nil
}

func (c *Client) processStreamResponse(ctx context.Context, body io.Reader, callback func(*ChatResponse) error) error {
	scanner := bufio.NewScanner(body)
	
	// 设置更大的缓冲区以处理长响应
	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, 1024*1024)

	for scanner.Scan() {
		// 检查context是否被取消
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		line := strings.TrimSpace(scanner.Text())

		if line == "" || !strings.HasPrefix(line, "data: ") {
			continue
		}

		data := strings.TrimPrefix(line, "data: ")
		if data == "[DONE]" {
			break
		}

		var streamResp ChatResponse
		if err := json.Unmarshal([]byte(data), &streamResp); err != nil {
			// 记录解析错误但继续处理
			fmt.Printf("解析流式数据失败: %v, 数据: %s\n", err, data)
			continue
		}

		if err := callback(&streamResp); err != nil {
			return err
		}
	}

	return scanner.Err()
}
