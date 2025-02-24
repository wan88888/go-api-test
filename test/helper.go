package test

import (
	"encoding/json"
	"fmt"
	"html/template"
	"os"
	"testing"
	"time"

	"github.com/fatih/color"
	"github.com/go-resty/resty/v2"
	"github.com/wan/go-api-test/client"
	"github.com/wan/go-api-test/config"
)

// TestCase 定义测试用例结构
type TestCase struct {
	Name        string
	Path        string
	Method      string
	Body        interface{}
	ExpectCode  int
	ExpectField map[string]interface{}
	StartTime   time.Time
	EndTime     time.Time
	Error       error
}

// TestSuite 测试套件结构
type TestSuite struct {
	Name      string
	Client    *client.APIClient
	TestCases []*TestCase
	t         *testing.T
}

// NewTestSuite 创建新的测试套件
func NewTestSuite(t *testing.T, name string) *TestSuite {
	return &TestSuite{
		Name:      name,
		Client:    client.NewAPIClient(config.DefaultConfig()),
		TestCases: make([]*TestCase, 0),
		t:        t,
	}
}

// AddTestCase 添加测试用例到套件
func (s *TestSuite) AddTestCase(tc *TestCase) {
	s.TestCases = append(s.TestCases, tc)
}

// Run 执行测试套件中的所有测试用例
func (s *TestSuite) Run() {
	for _, tc := range s.TestCases {
		s.t.Run(tc.Name, func(t *testing.T) {
			var resp *resty.Response
			tc.StartTime = time.Now()

			var err error
			tc.StartTime = time.Now()

			switch tc.Method {
			case "GET":
				resp, err = s.Client.Get(tc.Path)
			case "POST":
				resp, err = s.Client.Post(tc.Path, tc.Body)
			default:
				err = fmt.Errorf("不支持的HTTP方法: %s", tc.Method)
			}

			tc.EndTime = time.Now()
			tc.Error = err

			if err != nil {
				t.Error(err)
				return
			}

			// 验证响应状态码
			if resp.StatusCode() != tc.ExpectCode {
				err = fmt.Errorf("状态码不匹配，期望 %d，实际 %d", tc.ExpectCode, resp.StatusCode())
				t.Error(err)
				tc.Error = err
				return
			}

			// 验证响应字段
			if tc.ExpectField != nil {
				var respBody map[string]interface{}
				if err := json.Unmarshal(resp.Body(), &respBody); err != nil {
					err = fmt.Errorf("解析响应体失败: %v", err)
					t.Error(err)
					tc.Error = err
					return
				}

				for k, v := range tc.ExpectField {
					if respBody[k] != v {
						err = fmt.Errorf("字段 '%s' 不匹配，期望 %v，实际 %v", k, v, respBody[k])
						t.Error(err)
						tc.Error = err
						return
					}
				}
			}
		})
	}
}

// GenerateReport 生成测试报告
func (s *TestSuite) GenerateReport() {
	totalCases := len(s.TestCases)
	passedCases := 0
	failedCases := 0

	// 控制台输出报告
	fmt.Printf("\n%s\n", color.GreenString("=== Test Report ==="))
	fmt.Printf("Test Suite: %s\n\n", color.CyanString(s.Name))

	// 准备HTML报告数据
	type TestCaseReport struct {
		Name     string
		Duration float64
		Error    error
	}

	type Report struct {
		Name      string
		Total     int
		Passed    int
		Failed    int
		TestCases []TestCaseReport
	}

	report := Report{
		Name:      s.Name,
		Total:     totalCases,
		TestCases: make([]TestCaseReport, 0, totalCases),
	}

	for _, tc := range s.TestCases {
		duration := tc.EndTime.Sub(tc.StartTime)
		if tc.Error == nil {
			passedCases++
			fmt.Printf("%s %s (%.2fs)\n", 
				color.GreenString("✓"),
				tc.Name,
				duration.Seconds())
		} else {
			failedCases++
			fmt.Printf("%s %s (%.2fs)\n%s\n",
				color.RedString("✗"),
				tc.Name,
				duration.Seconds(),
				color.RedString("  Error: %v", tc.Error))
		}

		report.TestCases = append(report.TestCases, TestCaseReport{
			Name:     tc.Name,
			Duration: duration.Seconds(),
			Error:    tc.Error,
		})
	}

	report.Passed = passedCases
	report.Failed = failedCases

	fmt.Printf("\nTotal: %d  Passed: %s  Failed: %s\n",
		totalCases,
		color.GreenString("%d", passedCases),
		color.RedString("%d", failedCases))

	// 生成HTML报告
	tmpl, err := template.ParseFiles("config/template.html")
	if err != nil {
		fmt.Printf("解析HTML模板失败: %v\n", err)
		return
	}

	reportFile, err := os.Create("report/test_report.html")
	if err != nil {
		fmt.Printf("创建报告文件失败: %v\n", err)
		return
	}
	defer reportFile.Close()

	if err := tmpl.Execute(reportFile, report); err != nil {
		fmt.Printf("生成HTML报告失败: %v\n", err)
		return
	}

	fmt.Printf("\nHTML报告已生成: report/test_report.html\n")
}