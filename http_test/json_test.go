package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestJson(t *testing.T) {
	// 要访问的URL
	url := "https://example.com/api"

	// 要伪造的JSON数据
	data := map[string]string{
		"name":    "Alice",
		"age":     "25",
		"country": "USA",
		"fake":    "data",
	}

	// 将JSON数据编码为字节数组
	jsonData, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	//Json.NewEncoder.Encode
	//buf := new(bytes.Buffer)
	//err := json.NewEncoder(buf).Encode(data)
	//req, err := http.NewRequest("POST", url, buf)
	// 创建HTTP请求

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		panic(err)
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3")

	// 发送HTTP请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// 读取响应内容
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	// 输出响应内容
	fmt.Println(string(body))
}
