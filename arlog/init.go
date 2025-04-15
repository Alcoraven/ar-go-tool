package arlog

import (
	"fmt"
	"github.com/robfig/cron/v3"
	"os"
	"strings"
	"sync"
	"time"
)

// type LogSwitchPlan uint32
//
// const (
// 	LogPlanOneFile LogSwitchPlan = iota
// 	LogPlanDayChange
// )
//
// type FileSwitchConfig struct {
// 	SwitchPlan     uint32
// 	RotateWithSize uint64 // auto rotate log file if FileSize(Bytes) > RotateWithSize(Bytes)
//
// }
//
// func isLogPlan(lp uint32) bool {
// 	if artool.InSlice([]uint32{}, lp) {
// 		return true
// 	}
// 	return false
// }

type manager struct {
	sync.RWMutex
	mode     LogMode
	plan     LogPlan
	savePath string
	depth    int
	timeLoc  *time.Location
	f        *arLogFs

	showPrefix bool
}

func (m *manager) createFile(f *arLogFs, savePath string) {
	if _, err := os.Stat(savePath); err != nil {
		switch {
		case os.IsNotExist(err):
			dir, _ := os.Getwd()
			if e := os.MkdirAll(dir+"/"+savePath, os.ModePerm); e != nil {
				fmt.Println("Save path not exist and created fail", savePath, err.Error())
				panic(e)
			}
		case os.IsPermission(err):
			fmt.Println("Permission denied", savePath, err.Error())
			panic(err)
		default:
			fmt.Println("Unknown error", savePath, err.Error())
			panic(err)
		}
	}

	switch m.plan {
	case LogPlanOnce:
		var filePath = savePath + "out.log"
		f.openFile(filePath)
	case LogPlanDay:
		nf := func() {
			day := time.Now().In(m.timeLoc)
			date := day.Format("20060102")
			var filePath = savePath + date + ".log"
			f.openFile(filePath)
		}
		nf()
		c := cron.New(cron.WithLocation(m.timeLoc))
		c.AddFunc("0 0 * * *", nf)
		go c.Start()
	case LogPlanMonth:
		nf := func() {
			day := time.Now().In(m.timeLoc)
			month := day.Format("200601")
			var filePath = savePath + month + ".log"
			f.openFile(filePath)
		}
		nf()
		c := cron.New(cron.WithLocation(m.timeLoc))
		c.AddFunc("0 0 1 * *", nf)
		go c.Start()
	}
}

func (m *manager) writeLog(level LogLevel, v ...interface{}) {
	var (
		prefix  string
		msg     string
		msgList []string
	)

	if m.showPrefix {
		prefix = getPrefix(m.timeLoc, m.depth, level) + " "
	}
	msg = formatMsg(v...)

	for _, s := range strings.Split(msg, "\n") {
		msgList = append(msgList, fmt.Sprintf("%s%s", prefix, s))
	}

	if m.mode != LogModeProd {
		printStd(level, msgList...)
	}
	m.f.write(msgList...)
}

type PrefixConfig struct {
	ShowPrefix bool
}

type Config struct {
	Mode     LogMode        // log mode
	Plan     LogPlan        // log plan
	Depth    int            // log file depth
	SavePath string         // log save path
	TimeLoc  *time.Location // time location (default for server local location)
	*PrefixConfig
}

func Init(config *Config) (Log, error) {
	var (
		logMode  = LogModeDev
		logPlan  = LogPlanOnce
		savePath = DefaultSavePath
		depth    = 0
		timeLoc  = time.Local

		showPrefix = true
	)

	if config != nil {
		if config.Mode != LogModeDev {
			logMode = config.Mode
		}
		if config.Plan != LogPlanOnce {
			logPlan = config.Plan
		}
		if len(config.SavePath) > 0 {
			savePath = config.SavePath
			if savePath[len(savePath)-1] != '/' {
				savePath += "/"
			}
		}
		if config.Depth > 0 {
			depth = config.Depth
		}
		if config.TimeLoc != nil {
			timeLoc = config.TimeLoc
		}
		if config.PrefixConfig != nil {
			showPrefix = config.ShowPrefix
		}
	}

	var m = &manager{
		mode:       logMode,
		plan:       logPlan,
		depth:      depth,
		timeLoc:    timeLoc,
		f:          new(arLogFs),
		showPrefix: showPrefix,
	}
	m.createFile(m.f, savePath)

	return m, nil
}

type Log interface {
	Debug(...interface{})
	Info(...interface{})
	Warn(...interface{})
	Error(...interface{})
	Fatal(...interface{})

	ChangeMode(mode LogMode)
}

func (m *manager) Debug(v ...interface{}) {
	if m.mode != LogModeDebug {
		return
	}
	m.writeLog(LogLevelDebug, v...)
}

func (m *manager) Info(v ...interface{}) {
	m.writeLog(LogLevelInfo, v...)
}

func (m *manager) Warn(v ...interface{}) {
	m.writeLog(LogLevelWarn, v...)
}

func (m *manager) Error(v ...interface{}) {
	m.writeLog(LogLevelError, v...)
}

func (m *manager) Fatal(v ...interface{}) {
	m.writeLog(LogLevelFatal, v...)
}

func (m *manager) ChangeMode(mode LogMode) {
	m.Lock()
	defer m.Unlock()

	m.mode = mode
}
