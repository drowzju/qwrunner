package client

type Content struct {
    Type     string    `json:"type"`
    Text     string    `json:"text,omitempty"`
    ImageURL *ImageURL `json:"image_url,omitempty"`
}

type ImageURL struct {
    URL string `json:"url"`
}

type Message struct {
    Role    string    `json:"role"`
    Content []Content `json:"content"`
}

type ChatRequest struct {
    Model    string    `json:"model"`
    Messages []Message `json:"messages"`
}

type Choice struct {
    Message struct {
        Content string `json:"content"`
    } `json:"message"`
}

type ChatResponse struct {
    Choices []Choice `json:"choices"`
}