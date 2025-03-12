package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
)

// 响应结构体
type StreamResponse struct {
	Event          string `json:"event"`
	TaskID         string `json:"task_id"`
	ID             string `json:"id"`
	MessageID      string `json:"message_id"`
	ConversationID string `json:"conversation_id"`
	Mode           string `json:"mode"`
	Answer         string `json:"answer"`
	Metadata       struct {
		Usage struct {
			PromptTokens        int     `json:"prompt_tokens"`
			PromptUnitPrice     string  `json:"prompt_unit_price"`
			PromptPriceUnit     string  `json:"prompt_price_unit"`
			PromptPrice         string  `json:"prompt_price"`
			CompletionTokens    int     `json:"completion_tokens"`
			CompletionUnitPrice string  `json:"completion_unit_price"`
			CompletionPriceUnit string  `json:"completion_price_unit"`
			CompletionPrice     string  `json:"completion_price"`
			TotalTokens         int     `json:"total_tokens"`
			TotalPrice          string  `json:"total_price"`
			Currency            string  `json:"currency"`
			Latency             float64 `json:"latency"`
		} `json:"usage"`
	} `json:"metadata"`
	CreatedAt int64 `json:"created_at"`
}

// 配置结构体
type Config struct {
	ShowThink       bool
	TypewriterDelay time.Duration
	ConversationID  string    `json:"conversation_id"`
	Messages        []Message `json:"messages"`
}

// 消息结构体
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// 历史记录处理函数
func loadHistory() (Config, error) {
	config := Config{
		ShowThink:       true,
		TypewriterDelay: 50,
		Messages:        make([]Message, 0),
	}

	// 读取历史记录文件
	data, err := os.ReadFile("history.json")
	if err != nil {
		if os.IsNotExist(err) {
			return config, nil // 文件不存在时返回默认配置
		}
		return config, err
	}

	// 解析历史记录
	err = json.Unmarshal(data, &config)
	return config, err
}

// 保存历史记录
func saveHistory(config Config) error {
	data, err := json.MarshalIndent(config, "", "    ")
	if err != nil {
		return err
	}
	return os.WriteFile("history.json", data, 0644)
}

func main() {
	// 加载历史记录
	config, err := loadHistory()
	if err != nil {
		fmt.Println("加载历史记录失败:", err)
		return
	}

	// 从环境变量获取配置参数
	apiKey := os.Getenv("DIFY_API_KEY")
	endpoint := os.Getenv("DIFY_ENDPOINT")
	prompt := os.Args[1] // 从命令行参数获取用户输入
	userId := os.Getenv("USER_ID")

	// 检查必要的环境变量
	if apiKey == "" || endpoint == "" || userId == "" {
		fmt.Println("请设置必要的环境变量。参考示例：")
		fmt.Println("export DIFY_API_KEY=your-api-key")
		fmt.Println("export DIFY_ENDPOINT=https://cloud.dify.com/v1/chat-messages")
		fmt.Println("export USER_ID=your-user-id")
		return
	}

	// 检查命令行参数
	if len(os.Args) < 2 {
		fmt.Println("请提供对话内容，例如：")
		fmt.Println("./main \"那你还有什么美白推荐吗\"")
		return
	}

	// 添加新消息到历史记录
	config.Messages = append(config.Messages, Message{
		Role:    "user",
		Content: prompt,
	})

	// 创建HTTP客户端
	client := &http.Client{}

	// 创建请求
	req, err := http.NewRequest("POST", endpoint, strings.NewReader(
		fmt.Sprintf(`{
            "messages": %s,
            "user": "%s",
            "inputs": {},
            "query": "%s",
            "stream": true,
            "conversation_id": "%s"
        }`, marshalMessages(config.Messages), userId, prompt, config.ConversationID))) // 添加conversation_id
	if err != nil {
		fmt.Println("创建请求失败:", err)
		return
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Accept", "text/event-stream")

	// 发送请求
	fmt.Println("正在发送请求...")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("请求失败:", err)
		return
	}
	defer resp.Body.Close()

	// 检查响应状态
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("请求返回非成功状态码: %d\n", resp.StatusCode)
		scanner := bufio.NewScanner(resp.Body)
		for scanner.Scan() {
			fmt.Println(scanner.Text())
		}
		return
	}

	// 处理响应
	fmt.Println("收到响应，状态码:", resp.StatusCode)

	// 设置scanner的最大容量，防止长行被截断
	scanner := bufio.NewScanner(resp.Body)
	buf := make([]byte, 0, 64*1024) // 64KB
	scanner.Buffer(buf, 1024*1024)  // 设置最大扫描长度为1MB

	// 修改响应处理部分，合并两个scanner循环
	fmt.Println("\n回答内容:")
	var currentAnswer strings.Builder

	for scanner.Scan() {
		line := scanner.Text()

		var response StreamResponse
		if err := json.Unmarshal([]byte(line), &response); err != nil {
			fmt.Println("解析JSON失败:", err)
			continue
		}

		// 保存conversation_id
		if config.ConversationID == "" {
			config.ConversationID = response.ConversationID
		}

		// 处理answer
		answer := response.Answer
		if answer != "" {
			// 处理think标签
			if !config.ShowThink && strings.Contains(answer, "<think>") && strings.Contains(answer, "</think>") {
				thinkStart := strings.Index(answer, "<think>")
				thinkEnd := strings.Index(answer, "</think>") + len("</think>")
				if thinkStart >= 0 && thinkEnd > 0 && thinkEnd > thinkStart {
					answer = strings.TrimSpace(answer[thinkEnd:])
				}
			}

			// 打字机效果输出
			for _, char := range answer {
				fmt.Print(string(char))
				time.Sleep(config.TypewriterDelay * time.Millisecond)
			}

			// 累积完整回答
			currentAnswer.WriteString(answer)
		}
	}

	// 保存AI回答到历史记录
	if currentAnswer.Len() > 0 {
		config.Messages = append(config.Messages, Message{
			Role:    "assistant",
			Content: currentAnswer.String(),
		})
	}

	// 保存历史记录到文件
	if err := saveHistory(config); err != nil {
		fmt.Println("保存历史记录失败:", err)
	}
}

/*
示例配置：
- API Key: app-QPjj1x0Gwet7ZtWX0pOVadJv
- Endpoint: https://cloud.dify.com/v1/chat-messages
- User ID: test-user-id
*/

// 添加辅助函数：将消息数组转换为JSON字符串
func marshalMessages(messages []Message) string {
	jsonBytes, err := json.Marshal(messages)
	if err != nil {
		return "[]"
	}
	return string(jsonBytes)
}
