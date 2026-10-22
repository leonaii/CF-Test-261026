package main

// run.log：把运行过程和完整的账号数据（解密后）落盘，等价于 CLI 的 dump。

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"time"
)

const runLogMaxSize = 5 << 20 // 超过 5MB 轮转成 run.log.old

type runLogger struct {
	mu   sync.Mutex
	f    *os.File
	path string
}

func openRunLogger(path string) (*runLogger, error) {
	if st, err := os.Stat(path); err == nil && st.Size() > runLogMaxSize {
		_ = os.Remove(path + ".old")
		_ = os.Rename(path, path+".old")
	}
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o600)
	if err != nil {
		return nil, err
	}
	return &runLogger{f: f, path: path}, nil
}

// line 写一行带时间戳的日志；nil 安全。
func (l *runLogger) line(format string, args ...any) {
	if l == nil {
		return
	}
	l.write(time.Now().Format("2006-01-02 15:04:05") + "  " + fmt.Sprintf(format, args...) + "\r\n")
}

// block 写一段带分隔线的多行内容（账号明文等）。
func (l *runLogger) block(title, content string) {
	if l == nil {
		return
	}
	var b strings.Builder
	b.WriteString("\r\n===== " + title + " " + time.Now().Format("2006-01-02 15:04:05") + " =====\r\n")
	b.WriteString(content)
	if !strings.HasSuffix(content, "\n") {
		b.WriteString("\r\n")
	}
	b.WriteString("===== end =====\r\n\r\n")
	l.write(b.String())
}

func (l *runLogger) write(s string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.f == nil {
		return
	}
	_, _ = l.f.WriteString(s)
}

func (l *runLogger) close() {
	if l == nil {
		return
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.f != nil {
		_ = l.f.Close()
		l.f = nil
	}
}

// dumpPayload 把 redeem 响应解开后写入日志，其它接口忽略。
func dumpPayload(l *runLogger, action, card, endpoint string, body []byte) {
	if l == nil || endpoint != "api_card_redeem" {
		return
	}
	plain, err := OpenEnvelope(body)
	if err != nil {
		l.line("[DUMP] %s card=%s 解密失败: %v", action, card, err)
		return
	}
	l.block(fmt.Sprintf("%s 账号数据 card=%s", action, card), PrettyJSON(plain))
}
