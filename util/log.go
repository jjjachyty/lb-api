package util

import (
	"os"

	log "github.com/sirupsen/logrus"
)

var Glog *log.Logger

func init() {
	Glog = log.New()
	// Glog.Out = os.Stdout

	file, err := os.OpenFile("lb.log", os.O_CREATE|os.O_WRONLY, 0666)
	if err == nil {
		Glog.Out = os.Stdout
	} else {
		log.Info("Failed to log to file, using default stderr")
	}

	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.JSONFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(file)

	// Only log the warning severity or above.
	Glog.SetLevel(log.DebugLevel)
	log.Debug("日志启动成功")

}
