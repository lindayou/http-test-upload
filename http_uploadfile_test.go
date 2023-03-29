package main

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"testing"
)

//打开要上传的文件。
//创建一个multipart.Writer实例，它用于将表单数据编码为multipart格式。
//将要上传的文件添加到multipart.Writer实例中。
//关闭multipart.Writer实例，这样它会在内存中生成multipart格式的数据。
//创建一个HTTP POST请求，将multipart格式的数据作为请求体发送。
//发送HTTP请求，并输出响应状态码。
func TestUpLoad(t *testing.T) {
	file, err := os.Open("./test") // 替换为你要上传的文件路径
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file", file.Name())
	if err != nil {
		fmt.Println(err)
		return
	}

	_, err = io.Copy(part, file)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = writer.Close()
	if err != nil {
		fmt.Println(err)
		return
	}

	req, err := http.NewRequest("POST", "http://127.0.0.1:8080/upload", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	fmt.Println(resp.Status)
}
