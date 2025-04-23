package pg

import (
	"encoding/json"
	"fmt"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
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

type Hn2Header map[string][]Header
type Hn2Url map[string]string

func GetAllHeader(page *rod.Page) (hn2Url Hn2Url, err error) {
	page.MustNavigate("https://zh.cppreference.com/w/c/header")
	page.MustWaitLoad()
	var result *proto.RuntimeRemoteObject
	result, err = page.Eval(`() => {
	const pattern = /<([^<>./]+)\.h>/;
    const table = document.querySelector('table.t-dsc-begin');
    const linkData = [];
    // 用于记录已经存在的 header
    const existingHeaders = {};
    if (table) {
        // 获取表格的所有行
        const rows = table.querySelectorAll('tr.t-dsc');
        rows.forEach(row => {
            // 获取第一列的所有链接
            const firstColumn = row.querySelector('td:first-child');
            const links = firstColumn.querySelectorAll('a');
            links.forEach(link => {
                const fullHeader = link.textContent.trim();
                const headerMatch  = fullHeader.match(pattern);
                const url = link.href;
                const header = headerMatch ? headerMatch[1] : "";
                if (!existingHeaders[header]) {
                    linkData.push({
                        header: header,
                        fullHeader: fullHeader,
                        url: url
                    });
                    // 标记该 header 已经存在
                    existingHeaders[header] = true;
                }                
            });
        });
        // 打印结果
        console.log(linkData);
    }
    return linkData
	}`)

	if err != nil {
		return nil, err
	}
	// 将结果序列化为 JSON 字节
	jsonBytes, err := json.Marshal(result.Value)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal result: %v", err)
	}

	// 定义一个结构体来解析 JSON 数据
	var data []struct {
		Header     string `json:"header"`
		FullHeader string `json:"fullHeader"`
		Url        string `json:"url"`
	}

	// 将 JSON 数据反序列化到结构体中
	err = json.Unmarshal(jsonBytes, &data)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal result: %v", err)
	}
	hn2Url = make(Hn2Url)
	for _, st := range data {
		hn2Url[st.Header] = st.Url
	}

	return hn2Url, nil
}

func GetHeaderIdentifier(page *rod.Page, url string) {
	
}
