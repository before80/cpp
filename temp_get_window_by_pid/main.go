package main

import (
	"fmt"
	"os/exec"
	"strings"
	"syscall"
	"time"
	"unsafe"

	"github.com/go-rod/rod"
	"github.com/go-vgo/robotgo"
	"github.com/shirou/gopsutil/v3/process"
)

func findProcessByName(processName string) ([]int32, error) {
	processes, err := process.Processes()
	if err != nil {
		return nil, err
	}

	var pids []int32
	for _, proc := range processes {
		name, _ := proc.Name()
		if strings.Contains(strings.ToLower(name), strings.ToLower(processName)) {
			pid := proc.Pid
			pids = append(pids, pid)
		}
	}
	return pids, nil
}

func enumWindows(callback func(hwnd uintptr) bool) {
	user32 := syscall.NewLazyDLL("user32.dll")
	enumWindowsProc := user32.NewProc("EnumWindows")

	type enumWindowsCallbackType uintptr

	callbackPtr := syscall.NewCallback(func(hwnd uintptr, lParam enumWindowsCallbackType) uintptr {
		if !callback(hwnd) {
			return 0
		}
		return 1
	})

	enumWindowsProc.Call(callbackPtr, uintptr(unsafe.Pointer(&callbackPtr)))
}

func findWindowByProcessID(pid int32, windowTitlePart string) uintptr {
	var foundHWND uintptr
	enumWindows(func(hwnd uintptr) bool {
		length := getWindowTextLength(hwnd)
		buf := make([]uint16, length+1)
		getWindowText(hwnd, &buf[0], int32(length+1))
		title := syscall.UTF16ToString(buf)

		windowPID := getWindowThreadProcessId(hwnd)
		if strings.Contains(title, windowTitlePart) && int32(windowPID) == pid {
			foundHWND = hwnd
			return false
		}
		return true
	})
	return foundHWND
}

func activateWindow(hwnd uintptr) {
	setForegroundWindow(hwnd)
}

func getWindowTextLength(hwnd uintptr) int32 {
	user32 := syscall.NewLazyDLL("user32.dll")
	getWindowTextLengthProc := user32.NewProc("GetWindowTextLengthW")
	r1, _, _ := getWindowTextLengthProc.Call(hwnd)
	return int32(r1)
}

func getWindowText(hwnd uintptr, buf *uint16, nMaxCount int32) int32 {
	user32 := syscall.NewLazyDLL("user32.dll")
	getWindowTextProc := user32.NewProc("GetWindowTextW")
	r1, _, _ := getWindowTextProc.Call(hwnd, uintptr(unsafe.Pointer(buf)), uintptr(nMaxCount))
	return int32(r1)
}

func getWindowThreadProcessId(hwnd uintptr) uint32 {
	user32 := syscall.NewLazyDLL("user32.dll")
	getWindowThreadProcessIdProc := user32.NewProc("GetWindowThreadProcessId")
	var pid uint32
	getWindowThreadProcessIdProc.Call(hwnd, uintptr(unsafe.Pointer(&pid)))
	return pid
}

func setForegroundWindow(hwnd uintptr) bool {
	user32 := syscall.NewLazyDLL("user32.dll")
	setForegroundWindowProc := user32.NewProc("SetForegroundWindow")
	r1, _, _ := setForegroundWindowProc.Call(hwnd)
	return r1 != 0
}

func main() {
	// 使用 go-rod 打开网页
	url := "https://zh.cppreference.com/w/c/header"
	page := rod.New().MustConnect().MustPage(url).MustWaitLoad()
	_ = page

	// 获取 Chrome 进程 ID
	chromePIDs, err := findProcessByName("chrome.exe")
	if err != nil {
		fmt.Println("无法获取 Chrome 进程 ID:", err)
		return
	}
	if len(chromePIDs) == 0 {
		fmt.Println("未找到 Chrome 进程")
		return
	}

	// 查找 Chrome 窗口 HWND 并激活
	chromeTitle := "cppreference.com" // 根据实际情况调整窗口标题的一部分
	var chromeHWND uintptr
	for _, pid := range chromePIDs {
		chromeHWND = findWindowByProcessID(pid, chromeTitle)
		if chromeHWND != 0 {
			break
		}
	}
	if chromeHWND == 0 {
		fmt.Println("未找到 Chrome 窗口")
		return
	}
	activateWindow(chromeHWND)
	time.Sleep(2 * time.Second) // 等待页面加载完成

	// 使用 robotgo 触发全选并复制
	robotgo.KeyTap("a", "control") // 全选
	robotgo.KeyTap("c", "control") // 复制

	// 获取剪贴板内容
	time.Sleep(2 * time.Second) // 等待复制完成
	text, _ := robotgo.ReadAll()
	if text == "" {
		fmt.Println("无法从剪贴板获取文本")
		return
	}

	// 使用 robotgo 打开 Typora
	filePath := "D:\\Docs\\hugos\\lang\\content\\c\\types\\_index.md"
	cmd := exec.Command("start", "typora.exe", filePath)
	err = cmd.Start()
	if err != nil {
		fmt.Println("无法启动 Typora:", err)
		return
	}
	time.Sleep(5 * time.Second) // 等待 Typora 打开

	// 获取 Typora 进程 ID
	typoraPIDs, err := findProcessByName("typora.exe")
	if err != nil {
		fmt.Println("无法获取 Typora 进程 ID:", err)
		return
	}
	if len(typoraPIDs) == 0 {
		fmt.Println("未找到 Typora 进程")
		return
	}

	// 查找 Typora 窗口 HWND 并激活
	typoraTitle := "_index.md - Typora" // 根据实际情况调整窗口标题的一部分
	var typoraHWND uintptr
	for _, pid := range typoraPIDs {
		typoraHWND = findWindowByProcessID(pid, typoraTitle)
		if typoraHWND != 0 {
			break
		}
	}
	if typoraHWND == 0 {
		fmt.Println("未找到 Typora 窗口")
		return
	}

	activateWindow(typoraHWND)

	// 在 Typora 中全选并粘贴
	robotgo.KeyTap("a", "control") // 全选
	robotgo.KeyTap("v", "control") // 粘贴

	// 保存文件
	robotgo.KeyTap("s", "control")

	// 关闭 Typora
	robotgo.KeyTap("q", "control")

	fmt.Println("操作完成")
}
