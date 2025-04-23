package main

import (
	"cppreference/bs"
	"cppreference/cfg"
	"cppreference/js"
	"cppreference/lg"
	"cppreference/myf"
	"cppreference/pg"
	"cppreference/wind"
	_ "embed"
	"fmt"
	"github.com/go-rod/rod"
	"github.com/go-vgo/robotgo"
	"github.com/tailscale/win"
	"strconv"
	"strings"
	"time"
)

func main() {
	var err error
	_ = err
	var browser *rod.Browser
	var page *rod.Page
	var chromeHwnd, typoraHwnd win.HWND
	var HadReplace bool

	mdFilepath := cfg.Default.UniqueMdFilepath

	// 获取文件名
	spSlice := strings.Split(mdFilepath, "\\")
	mdFilename := spSlice[len(spSlice)-1]

	// 打开浏览器
	browser, err = bs.GetBrowser(strconv.Itoa(0))
	defer browser.MustClose()
	// 创建新页面
	page = browser.MustPage()
	chromeHwnd = robotgo.GetHWND()
	//wind.SetChromeWindowsName(chromeHwnd, "MyWin")

	res, _ := pg.GetAllHeader(page)
	fmt.Println(res)
	return

	urls := []string{
		"https://zh.cppreference.com/w/c/string/byte/strcpy",
		"https://zh.cppreference.com/w/c/types/ptrdiff_t",
	}
	for _, url := range urls {
		err = myf.TruncFileContent(mdFilepath)
		if err != nil {
			lg.ErrorToFileAndStdOutWithSleepSecond(fmt.Sprintf("url=%s，清空%q文件内容出现错误：%v\n", url, mdFilepath, err), 3)
			return
		}

		// 打开 markdown文件
		err = wind.OpenTypora(mdFilepath)
		if err != nil {
			lg.ErrorToFileAndStdOutWithSleepSecond(fmt.Sprintf("url=%s，打开窗口名为%q时出现错误：%v\n", url, mdFilename, err), 3)
			return
		}
		// 保证能打开 typora
		time.Sleep(3 * time.Second)
		typoraWindowName := fmt.Sprintf("%s - Typora", mdFilename)
		typoraHwnd, err = wind.FindWindowHwndByWindowTitle(typoraWindowName)
		if err != nil {
			lg.ErrorToFileAndStdOutWithSleepSecond(fmt.Sprintf("url=%s，找不到%q窗口：%v\n", url, typoraWindowName, err), 3)
			return
		}

		// 打开网页
		page.MustNavigate(url)
		// 等待页面加载完成
		page.MustWaitLoad()

		// 获取h1中的内容
		h1 := page.MustElement("#firstHeading").MustText()
		// 判断是否存在多个标识符
		ids := strings.Split(h1, ", ")
		_ = ids

		_, err = page.Eval(fmt.Sprintf(`() => { %s }`, js.WaitRunJs))
		if err != nil {
			lg.ErrorToFileAndStdOutWithSleepSecond(fmt.Sprintf("url=%s，执行js错误：%v\n", url, err), 3)
			return
		}

		wind.SelectAllAndCtrlC(chromeHwnd)
		wind.SelectAllAndDelete(typoraHwnd)
		wind.CtrlV(typoraHwnd)

		wind.CtrlS(typoraHwnd)
		robotgo.CloseWindow()
		time.Sleep(1 * time.Second)
		HadReplace, err = myf.ReplaceMarkdownFileContent(mdFilepath)
		if err != nil {
			lg.ErrorToFileAndStdOutWithSleepSecond(fmt.Sprintf("替换出现错误%v\n", err), 3)
			return
		}
		lg.InfoToFileAndStdOut(fmt.Sprintf("url=%s，HadReplace=%v\n", url, HadReplace))

		//将md的内容放到ids对应的三级菜单下

	}

	//wind.OpenDevToolToConsole(chromeHwnd)
	//time.Sleep(2 * time.Second)
	//_ = robotgo.WriteAll(js.WaitRunJs)

	//if err = robotgo.WriteAll(js.WaitRunJs); err != nil {
	//	fmt.Println("✘ 复制到剪贴板失败：", err)
	//}
	//
	//robotgo.SetForeg(chromeHwnd)
	//robotgo.TypeStr("allow pasting", 200)
	////time.Sleep(time.Second)
	//_ = robotgo.KeyTap("enter")
	//time.Sleep(time.Second)
	//wind.CtrlV(chromeHwnd)
	////wind.RunJsCodeInConsole(js.WaitRunJs)
	//time.Sleep(3 * time.Second)

	//curHwnd, err := wind.FindWindowHwndByWindowTitle("MyWin")
	//if err != nil {
	//	fmt.Println(err)
	//} else {
	//	_ = curHwnd
	//	fmt.Println("找到 MyWin 窗口")
	//	//wind.OpenDevToolToConsole(curHwnd)
	//	//fmt.Println("关闭DevTool")
	//}

	//n, _ := wind.GetWindowText(typoraHwnd)
	//fmt.Println("n=", n)
	//time.Sleep(1 * time.Second)
	//robotgo.CloseWindow()
	//time.Sleep(time.Second)

	time.Sleep(2000 * time.Second)
}
