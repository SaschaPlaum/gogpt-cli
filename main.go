package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

type requestBody struct {
	Model       string `json:"model"`
	Prompt      string `json:"prompt"`
	MaxTokens   int    `json:"max_tokens"`
	Temperature int    `json:"temperature"`
}

type Completion struct {
  ID      string   `json:"id"`     
  Object  string   `json:"object"` 
  Created int64    `json:"created"`
  Model   string   `json:"model"`  
  Choices []Choice `json:"choices"`
  Usage   Usage    `json:"usage"`  
}

type Choice struct {
  Text         string      `json:"text"`         
  Index        int64       `json:"index"`        
  Logprobs     interface{} `json:"logprobs"`     
  FinishReason string      `json:"finish_reason"`
}

type Usage struct {
  PromptTokens     int64 `json:"prompt_tokens"`    
  CompletionTokens int64 `json:"completion_tokens"`
  TotalTokens      int64 `json:"total_tokens"`     
}


func main() {
	err := godotenv.Load()
  if err != nil {
    log.Fatal("Error loading .env file")
  }
	
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("...: ")
	text, _ := reader.ReadString('\n')
	answer := callAPI(text)
	fmt.Println(answer)
}



func callAPI(question string) string {
	bearer := "Bearer " + os.Getenv("OPENAI_TOKEN") 
	url := "https://api.openai.com/v1/completions"

	requestBody := requestBody{
		Model:       "text-davinci-003",
		Prompt:      question,
		MaxTokens:  1000,
		Temperature: 0,
	}
	
	reqBody, _ := json.Marshal(requestBody)

	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(reqBody))

	req.Header.Add("Authorization", bearer)
	req.Header.Add("Content-Type", "application/json; charset=UTF-8")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error on response.\n[ERROR] -", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
			fmt.Println("Error while reading the response bytes:", err)
	}

	var completion Completion
	json.Unmarshal(body, &completion)
	return completion.Choices[0].Text
}