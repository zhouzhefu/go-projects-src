package logging

import (
	log "github.com/cihub/seelog"
	// customLogger "github.com/zhouzhefu/util/logging"
	"testing"
)


func Test_logWithCustomLogger(t *testing.T) {
	var1 := "Variable_1"
	Logger.Info("Log INFO some variable:%v", var1)
	Logger.Warn("Log WARN some message")
	Logger.Critical("Log Critical will send email")
}

/*
* Without custom conf, the Default logger is created 
* by "<seelog />"
*/
func Test_sayHello(t *testing.T) {
	defer log.Flush()
	log.Info("Say Hello to seelog!")
}