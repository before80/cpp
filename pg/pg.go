package pg

import (
	"cppreference/contants"
	"cppreference/js"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Itu struct {
	Id  string
	Typ string
	Url string
}

type Header struct {
	Ids    []string
	I2ItuS map[string][]Itu
}

type HInfo struct {
	Header     string `json:"header"`
	FullHeader string `json:"fullHeader"`
	Url        string `json:"url"`
}

type Hn2Header map[string][]Header
type Hn2HInfo map[string]HInfo

func GetAllHeader(page *rod.Page) (hn2HInfo Hn2HInfo, err error) {
	page.MustNavigate("https://zh.cppreference.com/w/c/header")
	page.MustWaitLoad()
	var result *proto.RuntimeRemoteObject

	result, err = page.Eval(js.InHeaderPageGetAllHeaderInfoJs)

	if err != nil {
		return nil, err
	}
	// 将结果序列化为 JSON 字节
	jsonBytes, err := json.Marshal(result.Value)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal result: %v", err)
	}

	// 定义一个结构体切片来解析 JSON 数据
	var data []HInfo

	// 将 JSON 数据反序列化到结构体中
	err = json.Unmarshal(jsonBytes, &data)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal result: %v", err)
	}
	hn2HInfo = make(Hn2HInfo)
	for _, st := range data {
		hn2HInfo[st.Header] = st
	}

	return hn2HInfo, nil
}

func GetHeaderIdentifier(page *rod.Page, url string) {

}

var InitSpeacialHeaderMdContent = `
+++
title = "%s"
date = %s
weight = %d
type = "docs"
description = "%s"
isCJKLanguage = true
draft = false

+++
`

func InitMdFile(index int, hInfo HInfo, page *rod.Page) (err error) {

	page.MustNavigate(hInfo.Url)
	page.MustWaitLoad()

	// 获取h1标签的内容
	h1Content := page.MustElement("#firstHeading").MustText()
	h1Content = strings.TrimSpace(strings.Replace(h1Content, "标准库标头", "", -1))

	err = os.MkdirAll(filepath.Join(contants.OutputFolderName, contants.CStdFolderName), 0777)
	if err != nil {
		return fmt.Errorf("无法创建%s目录：%v\n", filepath.Join(contants.OutputFolderName, contants.CStdFolderName), err)
	}

	newMdFp := filepath.Join(contants.OutputFolderName, contants.CStdFolderName, hInfo.Header)
	var newMd *os.File
	_, err1 := os.Stat(newMdFp)
	// 当文件不存在的情况下，新建文件并初始化该文件
	if err1 != nil && errors.Is(err1, fs.ErrNotExist) {
		//fmt.Println("err=", err1)
		newMd, err = os.OpenFile(newMdFp, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil {
			return fmt.Errorf("创建文件 %s 时出错: %w", newMdFp, err)
		}
		defer newMd.Close()
		_, err = newMd.WriteString(fmt.Sprintf(`
+++
title = "%s"
date = %s
weight = %d
type = "docs"
description = "%s"
isCJKLanguage = true
draft = false

+++
`, h1Content, time.Now().Format(time.RFC3339), index*10, ""))
		if err != nil {
			return fmt.Errorf("写入文件时出错: %w", err)
		}
		// 初始化文件内容

	}
	//
	//_, err := os.Stat(newMdFp)
	//if err == nil {
	//	// 文件存在
	//	fmt.Println("文件已存在:", filePath)
	//}
	//
	//errors.Is(err，fs.ErrNotExist)

	//if !os.IsNotExist(err) {
	//	// 出现了其他错误
	//	return err
	//}
	//
	//// 文件不存在，尝试创建
	//file, err := os.Create(filePath)
	//if err != nil {
	//	return err
	//}
	//defer file.Close()
	return

}
