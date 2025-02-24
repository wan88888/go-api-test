package config

// TestConfig 存储测试相关的配置信息
type TestConfig struct {
	BaseURL string
	Timeout int
	Retries int
}

// DefaultConfig 返回默认的测试配置
func DefaultConfig() *TestConfig {
	return &TestConfig{
		BaseURL: "https://jsonplaceholder.typicode.com",
		Timeout: 10,
		Retries: 3,
	}
}