package config

import (
    "fmt"
    "github.com/spf13/viper"
)

type Config struct {
    APIKey      string `mapstructure:"DASHSCOPE_API_KEY"`
    APIEndpoint string `mapstructure:"api_endpoint"`
    DefaultModel string `mapstructure:"default_model"`
}

func Load() (*Config, error) {
    viper.SetConfigName("config")
    viper.SetConfigType("yaml")
    viper.AddConfigPath(".")
    
    // 设置默认值
    viper.SetDefault("api_endpoint", "https://dashscope.aliyuncs.com/compatible-mode/v1/chat/completions")
    viper.SetDefault("default_model", "qwen-vl-max")
    
    // 自动读取环境变量
    viper.AutomaticEnv()
    
    // 读取配置文件（可选）
    if err := viper.ReadInConfig(); err != nil {
        if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
            return nil, fmt.Errorf("读取配置文件失败: %w", err)
        }
    }
    
    var cfg Config
    if err := viper.Unmarshal(&cfg); err != nil {
        return nil, fmt.Errorf("配置解析失败: %w", err)
    }
    
    if cfg.APIKey == "" {
        return nil, fmt.Errorf("缺少必要配置: DASHSCOPE_API_KEY")
    }
    
    return &cfg, nil
}