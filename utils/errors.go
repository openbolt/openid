package utils

import (
	"log"
	"runtime"
	"strings"
)

func ELog(err error) {
	if err != nil {
		pc, file, line, _ := runtime.Caller(1)
		file = file[strings.LastIndex(file, "/")+1 : len(file)]
		log.Printf("%s:%d %s %s", file, line, runtime.FuncForPC(pc).Name(), err.Error())
	}
}
