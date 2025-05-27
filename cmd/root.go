package cmd

import (
	"context"
	"fmt"
	"log"
	"os"
	"qwrunner/client"
	"qwrunner/config"
	"strings"

	"github.com/spf13/cobra"
)

var (
	model    string
	imageURL string
	question string
	stream   bool
)

var rootCmd = &cobra.Command{
	Use:   "qwrunner",
	Short: "DashScope多模态AI客户端",
	Long:  "支持图片+文本的多模态AI查询工具，支持qwen-vl-max和qvq-max模型",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.Load()
		if err != nil {
			log.Fatalf("配置加载失败: %v", err)
		}

		// 使用默认模型或命令行指定的模型
		if model == "" {
			model = cfg.DefaultModel
		}

		// 构建请求
		req := &client.ChatRequest{
			Model: model,
			Messages: []client.Message{
				{
					Role: "user",
					Content: []client.Content{
						{Type: "image_url", ImageURL: &client.ImageURL{URL: imageURL}},
						{Type: "text", Text: question},
					},
				},
			},
		}

		c := client.New(cfg)

		if stream {
			// 流式处理
			var fullContent strings.Builder
			err = c.CreateChatCompletionStream(context.Background(), req, func(resp *client.ChatResponse) error {
				if len(resp.Choices) > 0 && resp.Choices[0].Delta != nil {
					content := resp.Choices[0].Delta.Content
					fmt.Print(content)
					fullContent.WriteString(content)
				}
				return nil
			})
			fmt.Println() // 换行
		} else {
			// 非流式处理
			resp, err := c.CreateChatCompletion(context.Background(), req)
			if err != nil {
				log.Fatalf("请求失败: %v", err)
			}

			if len(resp.Choices) > 0 && resp.Choices[0].Message != nil {
				fmt.Println(resp.Choices[0].Message.Content)
			} else {
				fmt.Println("未收到有效响应")
			}
		}

		if err != nil {
			log.Fatalf("请求失败: %v", err)
		}
	},
}

func init() {
	rootCmd.Flags().StringVarP(&model, "model", "m", "", "指定模型名称 (qwen-vl-max, qvq-max等)")
	rootCmd.Flags().StringVarP(&imageURL, "image", "i", "", "图片URL")
	rootCmd.Flags().StringVarP(&question, "question", "q", "", "问题内容")
	rootCmd.Flags().BoolVarP(&stream, "stream", "s", false, "启用流式响应")

	rootCmd.MarkFlagRequired("image")
	rootCmd.MarkFlagRequired("question")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
