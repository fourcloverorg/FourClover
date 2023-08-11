package logger

import (
	"flag"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	config "fourclover.org/internal/config"
)

// Define your custom logger type.
type logger struct {
	*log.Logger
}

// Optionally make it a interface.
type Logger interface {
	Println(v ...interface{})
}

func init() {
	// Check if the logs directory is defined in the YAML config file, if yes use it, if not use the default logs directory
	var logpath string
	if !config.GetFourCloverSuppressLogs() {
		logpath = config.GetFourCloverLoggerPath()
	}

	if logpath != "" {
		// Get directory path from logpath and create it if it doesn't exist
		dir := filepath.Dir(logpath)
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			os.MkdirAll(dir, 0755)
		}

		flag.Parse()
		var file, err1 = os.Create(logpath)

		if err1 != nil {
			panic(err1)
		}

		// set output of logs to stdout and file
		mw := io.MultiWriter(os.Stdout, file)
		log.SetOutput(mw)
	}
}

// Create a new logger.
func NewLogger() Logger {
	return &logger{log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)}
}

// Set the default logger.
var DefaultLogger = NewLogger()

// // Implement the Logger interface.
// func (l *logger) Println(v ...interface{}) {
// 	l.Output(2, "INFO: FourClover: "+v[0].(string))
// }

// SupressLogs supresses the output of the program
func SupressLogs() {
	log.SetOutput(os.Stdout)
	log.SetOutput(ioutil.Discard)

	// Create a new file descriptor that discards output
	devNull, _ := os.Open(os.DevNull)
	defer devNull.Close()
	// Redirect stdout and stderr to the new file descriptor
	os.Stdout = devNull
	os.Stderr = devNull
}
