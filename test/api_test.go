package test

import (
	"testing"
)

func TestJSONPlaceholderAPI(t *testing.T) {
	// 创建测试套件
	suite := NewTestSuite(t, "JSONPlaceholder API Test")

	// 添加GET请求测试用例
	suite.AddTestCase(&TestCase{
		Name:       "Get Post By ID",
		Path:       "/posts/1",
		Method:     "GET",
		ExpectCode: 200,
		ExpectField: map[string]interface{}{
			"id":     float64(1),
			"userId": float64(1),
		},
	})

	// 添加POST请求测试用例
	suite.AddTestCase(&TestCase{
		Name:   "Create New Post",
		Path:   "/posts",
		Method: "POST",
		Body: map[string]interface{}{
			"title":  "Test Post",
			"body":   "This is a test post",
			"userId": 1,
		},
		ExpectCode: 201,
		ExpectField: map[string]interface{}{
			"title": "Test Post",
			"id":    float64(101),
		},
	})

	// 运行测试套件
	suite.Run()

	// 生成测试报告
	suite.GenerateReport()
}