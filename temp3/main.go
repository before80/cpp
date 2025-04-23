package main

import (
	"fmt"
	"path/filepath"
	"strings"
	"syscall"
	"time"
	"unsafe"

	"github.com/go-vgo/robotgo"
)

// Windows API 常量定义
const (
	SW_RESTORE = 9
	MAX_PATH   = 260
)

// Windows 结构体定义
type (
	RECT struct {
		Left, Top, Right, Bottom int32
	}
	WINDOWINFO struct {
		cbSize          uint32
		rcWindow        RECT
		rcClient        RECT
		dwStyle         uint32
		dwExStyle       uint32
		dwWindowStatus  uint32
		cxWindowBorders uint32
		cyWindowBorders uint32
		atomWindowType  uint16
		wCreatorVersion uint16
	}
)

// 获取所有可见窗口的标题列表（Windows专用实现）
func getWindowTitles() []string {
	user32 := syscall.NewLazyDLL("user32.dll")
	enumWindows := user32.NewProc("EnumWindows")
	getWindowTextW := user32.NewProc("GetWindowTextW")
	isWindowVisible := user32.NewProc("IsWindowVisible")

	var titles []string

	// 回调函数
	cb := syscall.NewCallback(func(hwnd syscall.Handle, lParam uintptr) uintptr {
		// 检查窗口是否可见
		visible, _, _ := isWindowVisible.Call(uintptr(hwnd))
		if visible == 0 {
			return 1 // 继续枚举
		}

		// 获取窗口标题
		buf := make([]uint16, MAX_PATH)
		length, _, _ := getWindowTextW.Call(
			uintptr(hwnd),
			uintptr(unsafe.Pointer(&buf[0])),
			uintptr(MAX_PATH),
		)

		if length > 0 {
			titles = append(titles, syscall.UTF16ToString(buf[:length]))
		}
		return 1 // 继续枚举
	})

	// 执行枚举
	enumWindows.Call(cb, 0)
	return titles
}

// 改进版窗口定位函数
func findTyporaWindow() (syscall.Handle, error) {
	user32 := syscall.NewLazyDLL("user32.dll")
	enumWindows := user32.NewProc("EnumWindows")

	// 定义目标窗口特征
	targetProcess := "Typora.exe"  // 进程名称
	targetTitlePattern := "Typora" // 标题包含的关键词
	_ = targetProcess
	_ = targetTitlePattern

	var foundWindow syscall.Handle

	// 回调函数
	cb := syscall.NewCallback(func(hwnd syscall.Handle, lParam uintptr) uintptr {
		// 获取进程ID
		var pid uint32
		user32.NewProc("GetWindowThreadProcessId").Call(
			uintptr(hwnd),
			uintptr(unsafe.Pointer(&pid)),
		)
		//fmt.Println("pid=", pid)
		// 根据PID获取进程名称
		name, err := getProcessName(pid)
		_ = err

		fmt.Println("pid=", pid, " name=", name, " title=", getWindowTitle(hwnd))
		if err == nil && name == targetProcess {
			// 验证窗口标题
			title := getWindowTitle(hwnd)
			fmt.Println(title)
			if strings.Contains(title, targetTitlePattern) {
				foundWindow = hwnd
				return 0 // 停止枚举
			}
		}

		return 1 // 继续枚举
	})

	// 执行枚举
	enumWindows.Call(cb, 0)

	if foundWindow == 0 {
		return 0, fmt.Errorf("未找到符合条件的Typora窗口")
	}
	return foundWindow, nil
}

// 定义Windows进程访问权限常量
const (
	PROCESS_QUERY_LIMITED_INFORMATION = 0x1000
)

func getProcessName(pid uint32) (string, error) {
	kernel32 := syscall.NewLazyDLL("kernel32.dll")
	proc := kernel32.NewProc("QueryFullProcessImageNameW")

	// 打开进程（使用修正后的权限）
	hProcess, err := syscall.OpenProcess(
		PROCESS_QUERY_LIMITED_INFORMATION,
		false,
		pid,
	)
	if err != nil {
		return "", fmt.Errorf("打开进程失败(PID:%d): %w", pid, err)
	}
	defer syscall.CloseHandle(hProcess)

	// 获取可执行文件路径
	buf := make([]uint16, MAX_PATH)
	size := uint32(MAX_PATH)
	ret, _, _ := proc.Call(
		uintptr(hProcess),
		0,
		uintptr(unsafe.Pointer(&buf[0])),
		uintptr(unsafe.Pointer(&size)),
	)

	if ret == 0 {
		return "", fmt.Errorf("获取进程路径失败")
	}

	// 提取文件名
	fullPath := syscall.UTF16ToString(buf[:size])
	return filepath.Base(fullPath), nil
}

// 辅助函数：获取窗口标题
func getWindowTitle(hwnd syscall.Handle) string {
	user32 := syscall.NewLazyDLL("user32.dll")
	getWindowTextW := user32.NewProc("GetWindowTextW")

	buf := make([]uint16, MAX_PATH)
	length, _, _ := getWindowTextW.Call(
		uintptr(hwnd),
		uintptr(unsafe.Pointer(&buf[0])),
		uintptr(MAX_PATH),
	)

	if length > 0 {
		return syscall.UTF16ToString(buf[:length])
	}
	return ""
}

func main() {
	// 第一步：查找Typora窗口
	hwnd, err := findTyporaWindow()
	if err != nil {
		fmt.Printf("窗口查找失败: %v\n", err)
		return
	}

	// 第二步：激活窗口
	user32 := syscall.NewLazyDLL("user32.dll")
	user32.NewProc("ShowWindow").Call(
		uintptr(hwnd),
		uintptr(SW_RESTORE),
	)
	user32.NewProc("SetForegroundWindow").Call(uintptr(hwnd))

	// 第三步：获取窗口位置信息
	var rect RECT
	user32.NewProc("GetWindowRect").Call(
		uintptr(hwnd),
		uintptr(unsafe.Pointer(&rect)),
	)

	// 输出窗口信息
	fmt.Printf("窗口位置: (%d, %d)\n", rect.Left, rect.Top)
	fmt.Printf("窗口尺寸: %dx%d\n", rect.Right-rect.Left, rect.Bottom-rect.Top)

	// 第四步：模拟拖动操作（原代码逻辑）
	titleBarX := int(rect.Left) + (int(rect.Right)-int(rect.Left))/2
	titleBarY := int(rect.Top) + 10 // 假设标题栏高度

	robotgo.MoveMouse(titleBarX, titleBarY)
	time.Sleep(500 * time.Millisecond)

	robotgo.Toggle("down")
	time.Sleep(500 * time.Millisecond)

	robotgo.MoveMouseSmooth(200, 200, 0.5, 100)
	time.Sleep(500 * time.Millisecond)

	robotgo.Toggle("up")

	fmt.Println("窗口操作完成")
}
