//go:build !js
// +build !js

package reader

import (
	"bufio"
	"io"
	"os"
	"strings"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

func GetReader(filename string) (io.Reader, error) {
	csvFile, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	return csvFile, nil
}

func GetGBKReader(filename string) (io.Reader, error) {
	csvFile, err := GetReader(filename)
	if err != nil {
		return nil, err
	}

	gbkReader := transform.NewReader(csvFile, simplifiedchinese.GBK.NewDecoder())
	return gbkReader, nil
}

// GetWechatBillReader 处理微信支付账单文件，将制表符替换为逗号
func GetWechatBillReader(filename string) (io.Reader, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	// 创建一个管道
	pr, pw := io.Pipe()

	// 在后台处理文件
	go func() {
		defer file.Close()
		defer pw.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()
			// 替换制表符为逗号，保留双引号以正确处理带逗号的字段
			line = strings.ReplaceAll(line, "\t", ",")
			// 写入处理后的行
			pw.Write([]byte(line + "\n"))
		}
	}()

	return pr, nil
}
