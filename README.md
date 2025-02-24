# Go API 测试框架

这是一个基于 Go 语言的 API 测试框架，提供了简单易用的接口测试功能，支持 HTTP 请求测试、响应验证和测试报告生成。

## 目录结构

```
├── client/         # HTTP 客户端实现
├── config/         # 配置相关
├── report/         # 测试报告输出
├── test/           # 测试用例和测试助手
└── go.mod          # Go 模块依赖文件
```

## 核心功能

- HTTP 请求测试支持（GET、POST 等方法）
- 灵活的测试用例配置
- 自动化测试执行
- 美观的 HTML 测试报告
- 控制台彩色输出
- 请求重试机制

## 快速开始

### 安装

```bash
go get github.com/wan/go-api-test
```

### 使用示例

```go
package test

import "testing"

func TestAPI(t *testing.T) {
    // 创建测试套件
    suite := NewTestSuite(t, "API Test Suite")

    // 添加 GET 请求测试
    suite.AddTestCase(&TestCase{
        Name:       "Get User",
        Path:       "/users/1",
        Method:     "GET",
        ExpectCode: 200,
        ExpectField: map[string]interface{}{
            "id": float64(1),
        },
    })

    // 添加 POST 请求测试
    suite.AddTestCase(&TestCase{
        Name:   "Create User",
        Path:   "/users",
        Method: "POST",
        Body: map[string]interface{}{
            "name":  "Test User",
            "email": "test@example.com",
        },
        ExpectCode: 201,
    })

    // 运行测试
    suite.Run()

    // 生成报告
    suite.GenerateReport()
}
```

## 配置说明

在 `config/config.go` 中可以配置测试框架的基本参数：

```go
type TestConfig struct {
    BaseURL string  // API 基础 URL
    Timeout int     // 请求超时时间（秒）
    Retries int     // 失败重试次数
}
```

默认配置：
- BaseURL: "https://jsonplaceholder.typicode.com"
- Timeout: 10 秒
- Retries: 3 次

## 测试用例结构

```go
type TestCase struct {
    Name        string                 // 测试用例名称
    Path        string                 // 请求路径
    Method      string                 // HTTP 方法
    Body        interface{}            // 请求体
    ExpectCode  int                    // 期望的状态码
    ExpectField map[string]interface{} // 期望的响应字段
}
```

## 测试报告

框架会生成两种格式的测试报告：

1. 控制台输出：彩色显示测试结果，包含成功/失败状态、执行时间和错误信息
2. HTML 报告：位于 `report/test_report.html`，包含详细的测试统计和结果展示

## 依赖项

- github.com/cucumber/godog: BDD 测试支持
- github.com/go-resty/resty/v2: HTTP 客户端
- github.com/fatih/color: 控制台彩色输出
- github.com/stretchr/testify: 测试断言

## 注意事项

1. 确保测试环境能够访问目标 API
2. 配置适当的超时时间和重试次数
3. 测试用例的 ExpectField 支持嵌套的 JSON 结构
4. HTML 报告需要 template.html 模板文件

## 许可证

MIT License