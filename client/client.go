package client

import (
    "bytes"
    "context"
    "encoding/json"
    "fmt"
    "net/http"
    "time"
    "qwrunner/config"
)

type Client struct {
    cfg        *config.Config
    httpClient *http.Client
}

func New(cfg *config.Config) *Client {
    return &Client{
        cfg: cfg,
        httpClient: &http.Client{
            Timeout: 30 * time.Second,
        },
    }
}

func (c *Client) CreateChatCompletion(ctx context.Context, req *ChatRequest) (*ChatResponse, error) {
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
        return nil, fmt.Errorf("API返回错误状态码: %d", resp.StatusCode)
    }
    
    var chatResp ChatResponse
    if err := json.NewDecoder(resp.Body).Decode(&chatResp); err != nil {
        return nil, fmt.Errorf("响应解析失败: %w", err)
    }
    
    return &chatResp, nil
}