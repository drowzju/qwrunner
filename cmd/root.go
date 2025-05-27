package cmd

import (
    "context"
    "fmt"
    "log"
    "os"
    "github.com/spf13/cobra"
    "qwrunner/client"
    "qwrunner/config"
)

var (
    model    string
    imageURL string
    question string
)

var rootCmd = &cobra.Command{
    Use:   "qwrunner",
    Short: "DashScope多模态AI客户端",
    Long:  "支持图片+文本的多模态AI查询工具",
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
                    Role: "system",
                    Content: []client.Content{
                        {Type: "text", Text: "You are a helpful assistant."},
                    },
                },
                {
                    Role: "user",
                    Content: []client.Content{
                        {Type: "image_url", ImageURL: &client.ImageURL{URL: imageURL}},
                        {Type: "text", Text: question},
                    },
                },
            },
        }
        
        // 发送请求
        c := client.New(cfg)
        resp, err := c.CreateChatCompletion(context.Background(), req)
        if err != nil {
            log.Fatalf("请求失败: %v", err)
        }
        
        // 输出结果
        if len(resp.Choices) > 0 {
            fmt.Println(resp.Choices[0].Message.Content)
        } else {
            fmt.Println("未收到有效响应")
        }
    },
}

func init() {
    rootCmd.Flags().StringVarP(&model, "model", "m", "", "指定模型名称")
    rootCmd.Flags().StringVarP(&imageURL, "image", "i", "", "图片URL")
    rootCmd.Flags().StringVarP(&question, "question", "q", "", "问题内容")
    
    rootCmd.MarkFlagRequired("image")
    rootCmd.MarkFlagRequired("question")
}

func Execute() {
    if err := rootCmd.Execute(); err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
}