package models

import (
	"crypto/md5"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

// UnixToTime 时间戳转换成日期
func UnixToTime(timestamp int) string {
	t := time.Unix(int64(timestamp), 0)
	return t.Format("2006-01-02 15:04:05")
}

// DateToUnix 日期转换成时间戳 2020-05-02 15:04:05
func DateToUnix(str string) int64 {
	template := "2006-01-02 15:04:05"
	t, err := time.ParseInLocation(template, str, time.Local)
	if err != nil {
		return 0
	}
	return t.Unix()
}

// GetUnix 获取时间戳
func GetUnix() int64 {
	return time.Now().Unix()
}

// GetDate 获取当前的日期
func GetDate() string {
	template := "2006-01-02 15:04:05"
	return time.Now().Format(template)
}

// GetDay 获取年月日
func GetDay() string {
	template := "20060102"
	return time.Now().Format(template)
}

func Md5(str string) string {
	h := md5.New()
	io.WriteString(h, str)
	return fmt.Sprintf("%x", h.Sum(nil))
}

// Int 表示把string转换成int
func Int(str string) (int, error) {
	n, err := strconv.Atoi(str)
	return n, err
}

// String 表示把int转换成string
func String(n int) string {
	str := strconv.Itoa(n)
	return str
}

// UploadImg
// @Description 上传图片
// @Author xYuan 2024-04-17 15:43:02
// @Param c
// @Param picName
// @Return string
// @Return error
func UploadImg(c *gin.Context, picName string) (string, error) {
	// 1、获取上传的图片
	file, err := c.FormFile(picName)
	if file == nil {
		return "", nil
	}
	if err != nil {
		return "", fmt.Errorf("failed to get uploaded file: %w", err)
	}
	// 2、获取后缀名 判断是否正确 .jpg .png .gif .jpeg
	extName := filepath.Ext(file.Filename)
	allowExtMap := map[string]bool{
		".jpg":  true,
		".png":  true,
		".gif":  true,
		".jpeg": true,
	}
	if _, ok := allowExtMap[extName]; !ok {
		return "", errors.New("文件后缀名不合法")
	}

	// 3、创建图片保存目录 static/upload/images/20200417/
	today := GetDay()
	dirPath := fmt.Sprintf("static/upload/images/%s/", today)
	if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
		return "", fmt.Errorf("failed to create directory: %w", err)
	}
	// 4、生成文件名称和保存的目录使用图片GetUnix()+后缀名
	timestamp := time.Now().Unix()
	filename := fmt.Sprintf("%d%s", timestamp, extName)
	finalImgPath := filepath.Join(dirPath, filename)
	// 5、执行上传
	if err := c.SaveUploadedFile(file, finalImgPath); err != nil {
		return "", fmt.Errorf("failed to save uploaded file: %w", err)
	}
	urlPath := fmt.Sprintf("%s", finalImgPath)
	fmt.Println("urlPath:", urlPath)
	return urlPath, nil
}
