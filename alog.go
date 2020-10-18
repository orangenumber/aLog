// (c) 2020 Gon Y Yi. <https://gonyyi.com/copyright.txt>

package alog

import (
	"encoding/json"
	"io"
	"strconv"
	"sync"
	"time"
)

const newline = byte('\n')

var unsuppType = []byte("?{unexp}")

const (
	F_TIME uint16 = 1 << iota
	F_MMDD
	F_MICROSEC
	F_PREFIX
	F_UTC
	F_DATE
	F_USE_BUF_1K
	F_USE_BUF_2K
	F_STD = F_MMDD | F_TIME | F_PREFIX
)

// =====================================================================================================================
// A LOGGER
// =====================================================================================================================
type ALogger struct {
	out io.Writer
	// primary buffer
	buf          []byte
	bufUseBuffer bool
	bufSize      int
	// secondary buffer
	buf2    tinyBuffer
	mu      sync.Mutex
	prefix  []byte
	flag    uint16
	jsonEnc *json.Encoder
}

func New(output io.Writer, prefix string, flag uint16) *ALogger {
	if output == nil {
		output = Discard
	}
	l := &ALogger{
		// buf:    make([]byte, 1024),
		out:    output,
		prefix: []byte(prefix),
		flag:   flag,
	}
	if flag&(F_USE_BUF_2K|F_USE_BUF_1K) > 0 {
		if flag&F_USE_BUF_2K > 0 {
			l.bufSize = 2048 * 2
		} else if flag&F_USE_BUF_1K > 0 {
			l.bufSize = 1024 * 2
		}
		l.buf = make([]byte, l.bufSize)
		l.buf = l.buf[:0]
		l.bufUseBuffer = true
	}

	return l
}

// =====================================================================================================================
// A LOGGER / SETTING
// =====================================================================================================================
func (l *ALogger) SetOutput(output io.Writer) {
	l.mu.Lock()
	l.out = output
	l.mu.Unlock()
}
func (l *ALogger) SetPrefix(s string) {
	l.mu.Lock()
	l.prefix = []byte(s)
	l.mu.Unlock()
}
func (l *ALogger) SetFlag(flag uint16) {
	l.mu.Lock()
	l.flag = flag
	l.mu.Unlock()
}

// formatHeader is modified from builtin logger
func (l *ALogger) formatHeader(buf *[]byte, t time.Time) {
	if l.flag&(F_DATE|F_MMDD|F_TIME|F_MICROSEC) != 0 {
		if l.flag&F_UTC != 0 {
			t = t.UTC()
		}
		if l.flag&(F_DATE|F_MMDD) != 0 {
			year, month, day := t.Date()
			if l.flag&F_DATE != 0 {
				itoa(buf, year, 4)
				*buf = append(*buf, '/')
			}
			if l.flag&(F_DATE|F_MMDD) != 0 {
				itoa(buf, int(month), 2)
				*buf = append(*buf, '/')
				itoa(buf, day, 2)
				*buf = append(*buf, ' ')
			}
		}
		if l.flag&(F_TIME|F_MICROSEC) != 0 {
			hour, min, sec := t.Clock()
			itoa(buf, hour, 2)
			*buf = append(*buf, ':')
			itoa(buf, min, 2)
			*buf = append(*buf, ':')
			itoa(buf, sec, 2)
			if l.flag&F_MICROSEC != 0 {
				*buf = append(*buf, '.')
				itoa(buf, t.Nanosecond()/1e3, 6)
			}
			*buf = append(*buf, ' ')
		}
	}
	if l.flag&F_PREFIX != 0 {
		*buf = append(*buf, l.prefix...)
	}
}
func (l *ALogger) Flush() {
	if l.bufUseBuffer {
		l.mu.Lock()
		defer l.mu.Unlock()
		if len(l.buf) > 0 {
			l.out.Write(l.buf)
			l.buf = l.buf[:0]
		}
	}
}
func (l *ALogger) Printf(format string, a ...interface{}) {
	t := time.Now()

	flagKeyword := false
	var aIdx int = 0
	var aLen = len(a)

	l.mu.Lock()
	defer l.mu.Unlock()

	if !l.bufUseBuffer { // if buffer write is not used, reset buffer each time.
		l.buf = l.buf[:0]
	}

	l.formatHeader(&l.buf, t)

	for _, c := range format {
		if flagKeyword == false {
			if c == '%' {
				flagKeyword = true
			} else {
				l.buf = append(l.buf, byte(c))
			}
		} else {
			// flagKeyword == true
			if c == '%' {
				l.buf = append(l.buf, '%')
				flagKeyword = false
				continue
			}
			if aIdx >= aLen {
				flagKeyword = false
				continue
			}
			switch c {
			case 'd':
				if v, ok := a[aIdx].(int); ok {
					itoa(&l.buf, v, 0)
				} else {
					l.buf = append(l.buf, unsuppType...)
				}
				aIdx++
			case 's':
				if v, ok := a[aIdx].(string); ok {
					l.buf = append(l.buf, []byte(v)...)
				} else {
					l.buf = append(l.buf, unsuppType...)
				}
				aIdx++
			case 'f':
				switch a[aIdx].(type) {
				case float64:
					if v, ok := a[aIdx].(float64); ok {
						ftoa(&l.buf, v, 2)
					} else {
						l.buf = append(l.buf, unsuppType...)
					}
				case float32:
					if v, ok := a[aIdx].(float32); ok {
						ftoa(&l.buf, float64(v), 2)
					} else {
						l.buf = append(l.buf, unsuppType...)
					}
				}
				aIdx++
			case 't':
				if v, ok := a[aIdx].(bool); ok {
					if v {
						l.buf = append(l.buf, []byte("true")...)
					} else {
						l.buf = append(l.buf, []byte("false")...)
					}
				} else {
					l.buf = append(l.buf, unsuppType...)
				}
				aIdx++
			}
			flagKeyword = false
		}
	}
	curBufSize := len(l.buf)
	if curBufSize == 0 || l.buf[curBufSize-1] != '\n' {
		l.buf = append(l.buf, '\n')
	}
	if curBufSize > l.bufSize {
		l.out.Write(l.buf)
		l.buf = l.buf[:0]
	}
}
func (l *ALogger) Print(a ...interface{}) {
	t := time.Now()
	l.mu.Lock()
	defer l.mu.Unlock()
	if !l.bufUseBuffer { // if buffer write is not used, reset buffer each time.
		l.buf = l.buf[:0]
	}

	l.formatHeader(&l.buf, t)
	// -- START

	for _, v := range a {
		switch v.(type) {
		case string:
			l.buf = append(l.buf, []byte(v.(string))...)
		case int:
			l.buf = strconv.AppendInt(l.buf, int64(v.(int)), 10)
		case int8:
			l.buf = strconv.AppendInt(l.buf, int64(v.(int8)), 10)
		case int16:
			l.buf = strconv.AppendInt(l.buf, int64(v.(int16)), 10)
		case int32:
			l.buf = strconv.AppendInt(l.buf, int64(v.(int32)), 10)
		case int64:
			l.buf = strconv.AppendInt(l.buf, v.(int64), 10)
		case bool:
			l.buf = strconv.AppendBool(l.buf, v.(bool))
		case uint:
			l.buf = strconv.AppendUint(l.buf, uint64(v.(uint)), 10)
		case uint8:
			l.buf = strconv.AppendUint(l.buf, uint64(v.(uint8)), 10)
		case uint16:
			l.buf = strconv.AppendUint(l.buf, uint64(v.(uint16)), 10)
		case uint32:
			l.buf = strconv.AppendUint(l.buf, uint64(v.(uint32)), 10)
		case uint64:
			l.buf = strconv.AppendUint(l.buf, v.(uint64), 10)
		case float32:
			l.buf = strconv.AppendFloat(l.buf, float64(v.(float32)), 'f', -1, 32)
		case float64:
			l.buf = strconv.AppendFloat(l.buf, v.(float64), 'f', -1, 64)
		case []byte:
			l.buf = append(l.buf, v.([]byte)...)
		default:
			l.buf = append(l.buf, unsuppType...)
		}
	}

	// -- END
	curBufSize := len(l.buf)
	if curBufSize == 0 || l.buf[curBufSize-1] != '\n' {
		l.buf = append(l.buf, '\n')
	}
	if curBufSize > l.bufSize {
		l.out.Write(l.buf)
		l.buf = l.buf[:0]
	}
}

func (l *ALogger) Printj(addPrefix string, a interface{}) {
	t := time.Now()
	l.mu.Lock()
	defer l.mu.Unlock()

	if !l.bufUseBuffer { // if buffer write is not used, reset buffer each time.
		l.buf = l.buf[:0]
	}

	l.formatHeader(&l.buf, t)
	if addPrefix != "" {
		l.buf = append(l.buf, []byte(addPrefix)...)
	}
	if a == nil {
		l.buf = append(l.buf, []byte("{}")...)
	} else {
		l.buf2 = l.buf2[:0]
		if l.buf2 == nil { // *json.Encode hasn't been initialized until needed.
			l.jsonEnc = json.NewEncoder(&l.buf2)
		}
		if l.jsonEnc.Encode(a) != nil {
			l.buf = append(l.buf, []byte("{}")...)
		} else {
			l.buf = append(l.buf, l.buf2...)
		}
	}

	curBufSize := len(l.buf)
	if curBufSize == 0 || l.buf[curBufSize-1] != '\n' {
		l.buf = append(l.buf, '\n')
	}
	if curBufSize > l.bufSize {
		l.out.Write(l.buf)
		l.buf = l.buf[:0]
	}
}

// just in case when io.Writer has a .Close() method like a file
func (l *ALogger) Close() error {
	l.Flush()
	if c, ok := l.out.(io.Closer); ok && c != nil {
		return c.Close()
	}
	return nil
}
