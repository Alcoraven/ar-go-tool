package arlog

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
)

const maxFileLines int64 = 100000

type arLogFs struct {
	sync.RWMutex
	logger       *log.Logger
	f            *os.File
	basePath     string
	slice        int
	lineCount    int64
	needsNewline bool
}

func slicePath(base string, slice int) string {
	if slice == 1 {
		return base
	}
	ext := filepath.Ext(base)
	return fmt.Sprintf("%s_slice_%d%s", strings.TrimSuffix(base, ext), slice, ext)
}

// countLines streams the file without imposing a maximum log line length.
func countLines(r io.Reader) (int64, bool, error) {
	buf := make([]byte, 32*1024)
	var lines int64
	var last byte
	var nonempty bool
	for {
		n, err := r.Read(buf)
		if n > 0 {
			lines += int64(bytes.Count(buf[:n], []byte{'\n'}))
			last, nonempty = buf[n-1], true
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			return 0, false, err
		}
	}
	unfinished := nonempty && last != '\n'
	if unfinished {
		lines++
	}
	return lines, unfinished, nil
}

func (f *arLogFs) openFile(filePath string) error {
	f.Lock()
	defer f.Unlock()
	entries, err := os.ReadDir(filepath.Dir(filePath))
	if err != nil {
		return err
	}
	ext := filepath.Ext(filePath)
	prefix := strings.TrimSuffix(filepath.Base(filePath), ext) + "_slice_"
	slice := 1
	for _, entry := range entries {
		name := entry.Name()
		if entry.IsDir() || !strings.HasPrefix(name, prefix) || !strings.HasSuffix(name, ext) {
			continue
		}
		number := strings.TrimSuffix(strings.TrimPrefix(name, prefix), ext)
		n, err := strconv.Atoi(number)
		if err == nil && n > slice && strconv.Itoa(n) == number {
			slice = n
		}
	}
	if err := f.switchFile(slicePath(filePath, slice), false); err != nil {
		return err
	}
	f.basePath, f.slice = filePath, slice
	return nil
}

// switchFile requires the write lock. Failed opens preserve the old state.
func (f *arLogFs) switchFile(path string, exclusive bool) error {
	flags := os.O_APPEND | os.O_CREATE | os.O_RDWR
	if exclusive {
		flags |= os.O_EXCL
	}
	handle, err := os.OpenFile(path, flags, 0644)
	if err != nil {
		return err
	}
	lines, unfinished, err := countLines(handle)
	if err != nil {
		handle.Close()
		return err
	}
	old := f.f
	f.f, f.logger = handle, log.New(handle, "", 0)
	f.lineCount, f.needsNewline = lines, unfinished
	if old != nil {
		old.Close()
	}
	return nil
}

func (f *arLogFs) write(list ...string) {
	f.Lock()
	defer f.Unlock()
	if len(list) == 0 {
		return
	}
	if f.logger == nil {
		fmt.Fprintln(os.Stderr, "write log error: log file is not open")
		return
	}
	if f.lineCount >= maxFileLines {
		if err := f.switchFile(slicePath(f.basePath, f.slice+1), true); err != nil {
			fmt.Fprintln(os.Stderr, "create log slice error:", err)
			return
		}
		f.slice++
	}
	if f.needsNewline {
		if _, err := f.f.WriteString("\n"); err != nil {
			fmt.Fprintln(os.Stderr, "write log error:", err)
			return
		}
		f.needsNewline = false
	}
	// All lines from one log call stay together, even beyond the threshold.
	for _, v := range list {
		if err := f.logger.Output(2, v); err != nil {
			fmt.Fprintln(os.Stderr, "write log error:", err)
			return
		}
		f.lineCount += int64(strings.Count(v, "\n"))
		if !strings.HasSuffix(v, "\n") {
			f.lineCount++
		}
	}
}
