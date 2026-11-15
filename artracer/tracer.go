package tracer

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"math/rand"
	"strconv"
	"sync"
	"time"
)

type sequenceState struct {
	timestamp   int64
	sequence    int
	initialized bool
}

// next is called with lock held. A false result means this timestamp is exhausted.
func (s *sequenceState) next(timestamp int64) (int, bool) {
	if !s.initialized || timestamp > s.timestamp {
		s.timestamp = timestamp
		s.sequence = 0
		s.initialized = true
		return s.sequence, true
	}
	if s.sequence == 999 {
		return 0, false
	}
	s.sequence++
	return s.sequence, true
}

var instanceID uint64
var instanceIDSet bool
var state sequenceState
var lock sync.Mutex

func getTimestampAndSequence() (int64, int) {
	for {
		lock.Lock()
		sequence, ok := state.next(time.Now().UnixMilli())
		timestamp := state.timestamp
		lock.Unlock()
		if ok {
			return timestamp, sequence
		}
		time.Sleep(time.Millisecond)
	}
}

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"

func getRandomStr(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func md5Str(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

// SetInstanceId 设置包内实例 ID，仅首次调用生效，后续调用静默忽略。
// 应在程序启动、开始生成 ID 前调用；调用方需保证 id 在 0~99 范围内，且不同实例不重复。
// 未设置时默认为 0；显式设置为 0 也会消耗唯一的设置机会。
func SetInstanceId(id uint64) {
	lock.Lock()
	defer lock.Unlock()
	if instanceIDSet {
		return
	}
	instanceID = id
	instanceIDSet = true
}

// SetInstanceID 是 SetInstanceId 的兼容入口，两个入口共享一次设置机会。
func SetInstanceID(id uint64) {
	SetInstanceId(id)
}

// GetUniqueTraceId 生成追踪 ID；同一毫秒内序列号耗尽时会等待时间前进。
// 格式为 MD5(随机字符串+实例ID+时间戳+序列号)+实例ID+时间戳+序列号。
// 实例 ID 使用 SetInstanceId 或 SetInstanceID 在启动时首次配置的值。
func GetUniqueTraceId() string {
	lock.Lock()
	id := instanceID
	lock.Unlock()
	timestamp, sequence := getTimestampAndSequence()
	var randStr = getRandomStr(10)
	var instanceStr = fmt.Sprintf("%02d", id)
	var num = fmt.Sprintf("%03d", sequence)
	var nt = strconv.FormatInt(timestamp, 10)
	var prefix = md5Str(randStr + instanceStr + nt + num)
	return prefix + instanceStr + nt + num
}
