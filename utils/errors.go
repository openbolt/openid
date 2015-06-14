package utils

import (
	"runtime"
	"strconv"
	"strings"

	log "github.com/Sirupsen/logrus"
)

func init() {
	log.SetLevel(log.DebugLevel)
}

// EDebug prints all debug messages
func EDebug(err error) {
	if err != nil {
		pc, file, line, _ := runtime.Caller(1)
		file = file[strings.LastIndex(file, "/")+1 : len(file)]
		src := file + ":" + strconv.Itoa(line)
		proc := runtime.FuncForPC(pc).Name()
		log.WithFields(log.Fields{"src": src, "proc": proc}).Debug(err.Error())
	}
}

// EInfo prints all info messages
func EInfo(err error) {
	if err != nil {
		pc, file, line, _ := runtime.Caller(1)
		file = file[strings.LastIndex(file, "/")+1 : len(file)]
		src := file + ":" + strconv.Itoa(line)
		proc := runtime.FuncForPC(pc).Name()
		log.WithFields(log.Fields{"src": src, "proc": proc}).Info(err.Error())
	}
}

// ELog prints all error or unclassified messages
func ELog(err error) {
	if err != nil {
		pc, file, line, _ := runtime.Caller(1)
		file = file[strings.LastIndex(file, "/")+1 : len(file)]
		src := file + ":" + strconv.Itoa(line)
		proc := runtime.FuncForPC(pc).Name()
		log.WithFields(log.Fields{"src": src, "proc": proc}).Error(err.Error())
	}
}
