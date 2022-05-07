package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"runtime"
	"strings"
	"time"
)

// Timer type for debug mode
type Timer struct {
	Begin    time.Time
	Check    time.Duration
	Function string
	File     string
	Line     string
}

// Start() *Timer
// run debug
func (o *Timer) Start() *Timer {
	if debug {
		o.Begin = time.Now()
		o.Function = "-"
		o.Line = "-"
		if pc, file, line, ok := runtime.Caller(1); ok {
			o.Function = runtime.FuncForPC(pc).Name()
			o.File = file
			o.Line = sprintf("%d", line)
		}
	}
	return o
}

// Stop()
// stop debug
func (o *Timer) Stop() {
	o.Check = time.Since(o.Begin)
	if debug {
		toLog(o.Check.String(), o.Line, o.Function)
	}
}

// getCurrentFuncName() string
// PROBABLY NOT USING
func getCurrentFuncName() string {
	pc := make([]uintptr, 10)
	frame, more := runtime.CallersFrames(pc[:runtime.Callers(2, pc)]).Next()
	if !more {
		return ""
	}
	return frame.Function
}

func fmax[T int|int64](i *T, j *T) T {
	if *i > *j {
		return *i
	}
	return *j
}

func fmaxStr(i *string, j *string) string {
	if *i != "" {
		return *i
	}
	return *j
}

// searchFiles(string, string, string) []string
// search foles in folder
func searchFiles(url, folder, ext string) []string {
	mask := regexp.MustCompile(`[\\/]`)
	temp := mask.ReplaceAllString(strings.TrimPrefix(url, folder), "")
	var files []string
	fi, err := ioutil.ReadDir(folder)
	if err != nil {
		return nil
	}
	for i := range fi {
		if strings.Contains(strings.ToLower(fi[i].Name()), strings.ToLower(temp)) && strings.Contains(strings.ToLower(fi[i].Name()), strings.ToLower(ext)) {
			files = append(files, folder+"/"+fi[i].Name())
		}
	}
	return files
}

// delEmpty()
// delete empty folder or trash folder older then 6 hour(too long for working)
func delEmpty() {
	filesInfo, err := ioutil.ReadDir(".")
	if err != nil {
		return
	}
	for i := range filesInfo {
		re := regexp.MustCompile(".{8}_.{4}_.{4}_.{4}_.{12}") //uuid like
		if re.MatchString(filesInfo[i].Name()) && filesInfo[i].IsDir() {
			files, _ := ioutil.ReadDir(filesInfo[i].Name())
			if len(files) == 0 {
				os.Remove(filesInfo[i].Name())
			}
			if time.Since(filesInfo[i].ModTime()).Hours() >= 6 {
				os.RemoveAll(filesInfo[i].Name())
			}
		}
	}
}

// replaceSpecialSymbols(string) string
// without packet 'strings' replace special symbols to '_'
func replaceSpecialSymbols(text string) string {
	n := []rune{}
	for _, v2 := range text {
		val := v2
		for _, v := range pathReplace {
			if v == v2 {
				val = 95 // it's '_'
				break
			}
		}
		n = append(n, val)
	}
	return string(n)
}

// toLog(...any)
// log data
func toLog(val ...any) {
    if debug {
    	text := strings.Trim(fmt.Sprintln(val...), "[]")
    	printf(text)
    	defer fmt.Printf(text)
    }
}

// toLog(any, ...bool)
// log data wrap JsonPls
func (o *JsonPls) toLog(val any, fail ...bool) {
    if debug {
    	symbol := "[+]"
    	for _, v := range fail {
    		if v {
    			symbol = "[-]"
    			break
    		}
    	}
    	text := sprintf("%v[%v] '%v'('%v')\n", symbol, val, o.Artist, o.Song)
    	o.M.RLock()
    	o.Status += text
    	o.M.RUnlock()
    	printf(text)
    	defer fmt.Printf(text)
    }
}

// existFile(string) string
// check exist file
func existFile(url string) string {
	_, err := os.Stat(url)
	if err != nil {
		if os.IsNotExist(err) {
			return ""
		}
	}
	return url
}

//presentation of bool
func sBool(v interface{}) bool {
	switch val := v.(type) {
	case string:
		switch val {
		case "+":
			return true
		case "-":
			return false
		case "🔴":
			return false
		case "🟢":
			return true
		}
	case bool:
		return val
	default:
		return false
	}
	return false
}

//presentation of bool
func ssBool(v interface{}) string {
	switch val := v.(type) {
	case bool:
		switch val {
		case true:
			return "+"
		case false:
			return "-"
		}
	case string:
		return val
	default:
		return "-"
	}
	return "-"
}
