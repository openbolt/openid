package utils

import (
	"net/http"
	"runtime"
	"strconv"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/context"
)

const REQUEST_UUID = "request-uuid"

func init() {
	log.SetLevel(log.DebugLevel)
}

// EDebug prints all debug messages
func EDebug(err error, r *http.Request) {
	if err != nil {
		pc, file, line, _ := runtime.Caller(1)
		file = file[strings.LastIndex(file, "/")+1 : len(file)]
		src := file + ":" + strconv.Itoa(line)
		proc := runtime.FuncForPC(pc).Name()
		log.WithFields(log.Fields{"src": src, "proc": proc, "request": context.Get(r, REQUEST_UUID)}).Debug(err.Error())
	}
}

// EInfo prints all info messages
func EInfo(err error, r *http.Request) {
	if err != nil {
		pc, file, line, _ := runtime.Caller(1)
		file = file[strings.LastIndex(file, "/")+1 : len(file)]
		src := file + ":" + strconv.Itoa(line)
		proc := runtime.FuncForPC(pc).Name()
		log.WithFields(log.Fields{"src": src, "proc": proc, "request": context.Get(r, REQUEST_UUID)}).Info(err.Error())
	}
}

// ELog prints all error or unclassified messages
func ELog(err error, r *http.Request) {
	if err != nil {
		pc, file, line, _ := runtime.Caller(1)
		file = file[strings.LastIndex(file, "/")+1 : len(file)]
		src := file + ":" + strconv.Itoa(line)
		proc := runtime.FuncForPC(pc).Name()
		log.WithFields(log.Fields{"src": src, "proc": proc, "request": context.Get(r, REQUEST_UUID)}).Error(err.Error())
	}
}
