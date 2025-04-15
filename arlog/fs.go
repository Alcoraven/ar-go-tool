package arlog

import (
	"fmt"
	"log"
	"os"
	"sync"
)

type arLogFs struct {
	sync.RWMutex
	logger *log.Logger
	f      *os.File
}

func (f *arLogFs) openFile(filePath string) {
	f.Lock()
	defer f.Unlock()

	handle, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("create log file error:", err.Error())
		return
	}

	if f.f != nil {
		oldFile := f.f
		defer oldFile.Close()
	}

	f.f = handle
	f.logger = log.New(handle, "", 0)
}

func (f *arLogFs) write(list ...string) {
	f.Lock()
	defer f.Unlock()
	for _, v := range list {
		f.logger.Println(v)
	}
}
