package main

import (
	"github.com/go-rod/rod/lib/input"
	"log"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
)

func main() {
	var err error
	// 启动浏览器
	u := launcher.New().MustLaunch()

	// 创建浏览器实例
	browser := rod.New().ControlURL(u).MustConnect()
	defer browser.MustClose()

	// 创建新页面
	page := browser.MustPage()

	// 打开网页
	page.MustNavigate("https://zh.cppreference.com/w/c/header")

	// 等待页面加载完成
	page.MustWaitLoad()

	// 模拟按下 Ctrl + A 组合键
	err = page.Keyboard.Press(input.ControlLeft)
	if err != nil {
		log.Fatalf("按下 Ctrl 键时出错: %v", err)
	}
	err = page.Keyboard.Press(input.KeyA)
	if err != nil {
		log.Fatalf("按下 A 键时出错: %v", err)
	}
	err = page.Keyboard.Release(input.ControlLeft)
	if err != nil {
		log.Fatalf("释放 Ctrl 键时出错: %v", err)
	}

	// 模拟按下 Ctrl + C 组合键
	err = page.Keyboard.Press(input.ControlLeft)
	if err != nil {
		log.Fatalf("按下 Ctrl 键时出错: %v", err)
	}
	err = page.Keyboard.Press(input.KeyC)
	if err != nil {
		log.Fatalf("按下 C 键时出错: %v", err)
	}
	err = page.Keyboard.Release(input.ControlLeft)
	if err != nil {
		log.Fatalf("释放 Ctrl 键时出错: %v", err)
	}

	// 稍微等待一下，以便观察效果
	time.Sleep(5000 * time.Second)
}
