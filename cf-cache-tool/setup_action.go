package main

// setupAction 表示首次配置窗口的结束方式。
type setupAction int

const (
	setupOK   setupAction = iota // 点了“开始使用”
	setupQuit                    // 点“退出”或关掉窗口：直接结束软件
)
