package arcsv

import (
	"bufio"
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"sync"
)

type Config struct {
	UseBuffer  bool // if need buffer
	BufferSize int  // Bytes for buffer
}

// func Export(filePath string, data [][]string) {
// 	fp, err := os.Create(filePath) // 创建文件句柄
// 	if err != nil {
// 		return
// 	}
// 	defer fp.Close()
// 	w := csv.NewWriter(fp) // 创建一个新的写入文件流
// 	w.WriteAll(data)
// 	w.Flush()
// }

type file struct {
	writeLock *sync.Mutex
	f         *os.File
	bw        *bufio.Writer
	cw        *csv.Writer
	rowCount  int64
	useBuffer bool
}

type File interface {
	addRow(row []string) error

	AddRowString(values ...string) error
	AddRow(values ...interface{}) error
	Export() error
}

func (f *file) addRow(row []string) error {
	return f.cw.Write(row)
}

func (f *file) AddRowString(values ...string) error {
	f.writeLock.Lock()
	defer f.writeLock.Unlock()

	return f.addRow(values)
}

func (f *file) AddRow(values ...interface{}) error {
	f.writeLock.Lock()
	defer f.writeLock.Unlock()

	var row []string
	for _, v := range values {
		row = append(row, fmt.Sprintf("%v", v))
	}

	return f.addRow(row)

	// if err := f.cw.Write(row); err != nil {
	// 	fmt.Println("csv add row error:", err)
	// }

	if f.rowCount += 1; f.rowCount%10000 == 0 {
		f.cw.Flush()
		if err := f.cw.Error(); err != nil {
			return err
		}
	}

	return nil
}

func (f *file) Export() error {
	if f.f == nil {
		return errors.New("file is nil")
	}
	defer func() {
		if err := f.f.Close(); err != nil {
			fmt.Println(err)
		}
		f.f = nil
	}()

	f.cw.Flush()
	if err := f.cw.Error(); err != nil {
		return err
	}
	if f.useBuffer {
		if err := f.bw.Flush(); err != nil {
			return err
		}
	}

	return nil
}

func NewFile(filePath string, config *Config) (File, error) {
	f, err := os.Create(filePath)
	if err != nil {
		return nil, err
	}

	var bw *bufio.Writer
	var cw *csv.Writer

	if config != nil && config.UseBuffer {
		bw = bufio.NewWriter(f)
		cw = csv.NewWriter(bw)
	} else {
		cw = csv.NewWriter(f)
	}

	return &file{
		writeLock: &sync.Mutex{},
		f:         f,
		bw:        bw,
		cw:        cw,
	}, nil
}
