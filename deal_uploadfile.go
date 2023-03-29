package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type FileInfo struct {
	Name      string `json:"name"`
	Size      int64  `json:"size"`
	UploadAt  int64  `json:"upload_at"`
	UploadDir string `json:"-"`
}

func main() {
	http.HandleFunc("/upload", uploadHandler)
	http.ListenAndServe(":8080", nil)
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	// 检查请求方法是否为POST
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// 解析multipart/form-data格式请求体
	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "Failed to parse form:", err)
		return
	}

	// 获取上传的文件
	file, header, err := r.FormFile("file")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "Failed to get file:", err)
		return
	}
	defer file.Close()

	// 创建上传目录
	uploadDir := "./uploads/"
	err = os.MkdirAll(uploadDir, os.ModePerm)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Failed to create upload directory:", err)
		return
	}

	// 生成文件名并保存文件
	filename := filepath.Base(header.Filename)
	uploadPath := uploadDir + time.Now().Format("20060102-150405-") + filename
	outFile, err := os.Create(uploadPath)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Failed to create output file:", err)
		return
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, file)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Failed to save file:", err)
		return
	}

	// 构造响应并返回
	fileInfo := FileInfo{
		Name:      filename,
		Size:      header.Size,
		UploadAt:  time.Now().Unix(),
		UploadDir: uploadDir,
	}
	respData, err := json.Marshal(fileInfo)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Failed to marshal response:", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(respData)
}
