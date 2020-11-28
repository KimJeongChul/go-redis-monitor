package logger

import (
	"fmt"
	"log"
)

func LogE(packageName string, funcName string, logMsg ...interface{}) {
	logStr := "[E] ( " + packageName + " ) < " + funcName + " > "
	for _, item := range logMsg {
		logStr += fmt.Sprint(item)
	}
	log.Println(logStr)
}
func LogI(packageName string, funcName string, logMsg ...interface{}) {
	logStr := "[I] ( " + packageName + " ) < " + funcName + " > "
	for _, item := range logMsg {
		logStr += fmt.Sprint(item)
	}
	log.Println(logStr)
}
func LogD(packageName string, funcName string, logMsg ...interface{}) {
	//todo :: Debug Log가 켜져 있는 경우만 출력하도록
	logStr := "[D] ( " + packageName + " ) < " + funcName + " > "
	for _, item := range logMsg {
		logStr += fmt.Sprint(item)
	}
	log.Println(logStr)
}
