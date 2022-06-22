package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
)

// Logger variables - 4 types defined
var (
	Trace   *log.Logger // Just about anything
	Info    *log.Logger // Important information
	Warning *log.Logger // Be concerned
	Error   *log.Logger // Critical problem
)

// logInit sets the logger and log type
func logInit() error {
	// setup files with relative location
	base := ".."
	warningsFile := path.Join(base, "log", "warnings.log")
	errorsFile := path.Join(base, "log", "errors.log")

	fmt.Println(warningsFile)
	fmt.Println(errorsFile)

	// Open a file for warnings.
	warnings, err := os.OpenFile(warningsFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("Failed to open warning log file")
		return err
	}
	defer warnings.Close()

	// Open a file for errors.
	errors, err := os.OpenFile(errorsFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("Failed to open errors log file")
		return err
	}
	defer errors.Close()

	// Create a multi writer for errors.
	multi := io.MultiWriter(errors, os.Stderr)

	// Init the log package for each message type.
	initLog(ioutil.Discard, os.Stdout, warnings, multi)

	// Test each log type.
	// Trace.Println("I have something standard to say.")
	// Info.Println("Important Information.")
	// Warning.Println("There is something you need to know about.")
	// Error.Println("Something has failed.")

	return nil
}

// initLog sets the devices for each log type.
func initLog(traceHandle io.Writer, infoHandle io.Writer, warningHandle io.Writer, errorHandle io.Writer) {
	Trace = log.New(traceHandle,
		"TRACE: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Info = log.New(infoHandle,
		"INFO: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Warning = log.New(warningHandle,
		"WARNING: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Error = log.New(errorHandle,
		"ERROR: ",
		log.Ldate|log.Ltime|log.Lshortfile)
}
