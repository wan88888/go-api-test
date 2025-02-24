package client

import (
	"fmt"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/wan/go-api-test/config"
)

// APIClient 封装HTTP客户端
type APIClient struct {
	client  *resty.Client
	config  *config.TestConfig
	Headers map[string]string
}

// NewAPIClient 创建新的API客户端
func NewAPIClient(cfg *config.TestConfig) *APIClient {
	client := resty.New()
	client.SetTimeout(time.Duration(cfg.Timeout) * time.Second)
	client.SetRetryCount(cfg.Retries)
	client.SetBaseURL(cfg.BaseURL)

	return &APIClient{
		client:  client,
		config:  cfg,
		Headers: make(map[string]string),
	}
}

// SetHeader 设置请求头
func (c *APIClient) SetHeader(key, value string) {
	c.Headers[key] = value
}

// Get 发送GET请求
func (c *APIClient) Get(path string) (*resty.Response, error) {
	req := c.client.R()
	for k, v := range c.Headers {
		req.SetHeader(k, v)
	}
	resp, err := req.Get(path)
	if err != nil {
		return nil, fmt.Errorf("GET request failed: %v", err)
	}
	return resp, nil
}

// Post 发送POST请求
func (c *APIClient) Post(path string, body interface{}) (*resty.Response, error) {
	req := c.client.R().SetBody(body)
	for k, v := range c.Headers {
		req.SetHeader(k, v)
	}
	resp, err := req.Post(path)
	if err != nil {
		return nil, fmt.Errorf("POST request failed: %v", err)
	}
	return resp, nil
}