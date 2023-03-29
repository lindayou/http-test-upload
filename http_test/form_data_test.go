package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"testing"
)

//我们首先创建了一个multipart.Writer对象，并向其中添加了一个文件字段和两个其他表单字段。
//然后，我们设置了Content-Type请求头为multipart/form-data格式，并计算了请求体的长度。
//最后，我们创建了一个HTTP请求并发送，同时读取了响应内容并打印。需要注意的是，在使用multipart/form-data格式上传文件时，
//需要调用multipart.Writer.CreateFormFile方法来创建文件字段，然后使用io.Copy方法将文件内容复制到该字段中。

func TestFormData(t *testing.T) {

	// 创建一个新的multipart.Writer  multipart.NewWriter创建了一个写入器
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)

	// 添加文件字段
	file, err := os.Open("example.txt")
	if err != nil {
		fmt.Println("Failed to open file:", err)
		return
	}
	defer file.Close()

	part, err := writer.CreateFormFile("file", "example.txt")
	if err != nil {
		fmt.Println("Failed to create form file:", err)
		return
	}
	_, err = io.Copy(part, file)
	if err != nil {
		fmt.Println("Failed to copy file to form:", err)
		return
	}

	// 添加其他表单字段
	writer.WriteField("username", "john")
	writer.WriteField("password", "password123")

	// 设置Content-Type为multipart/form-data，并计算Content-Length
	contentType := writer.FormDataContentType()
	contentLength := body.Len()

	// 创建HTTP请求并发送
	req, err := http.NewRequest("POST", "http://example.com/upload", body)
	if err != nil {
		fmt.Println("Failed to create request:", err)
		return
	}
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("Content-Length", strconv.Itoa(contentLength))

	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Failed to send request:", err)
		return
	}
	defer resp.Body.Close()

	// 读取响应内容
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Failed to read response:", err)
		return
	}

	fmt.Println(string(respBody))

}
