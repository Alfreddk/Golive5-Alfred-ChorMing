package logger

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"log"
	"os"
	"path"

	"github.com/joho/godotenv"
)

// Loggers that log errors.
var (
	Trace *log.Logger // Just about anything
	Info  *log.Logger // Authentication log
)

// Logs and Log Hashes files.
var (
	logDir           string
	InfoLogFile      string
	InfoLogHashFile  string
	TraceLogFile     string
	TraceLogHashFile string
)

// init() initialized all customised type loggers to allow logging of logs.
// It also performs hash verification of the content of the log files.
func init() {

	envFile := path.Join("config", ".env")

	err := godotenv.Load(envFile)
	if err != nil {
		log.Fatalln("Error loading .env file: ", err)
	}

	// Retrieving env variables LOG_DIR, INFO_LOG, INFO_LOG_HASH, TRACE_LOG and TRACE_LOG_HASH.
	logDir = os.Getenv("LOG_DIR")
	InfoLogFile = path.Join(logDir, os.Getenv("INFO_LOG"))
	InfoLogHashFile = path.Join(logDir, os.Getenv("INFO_LOG_HASH"))
	TraceLogFile = path.Join(logDir, os.Getenv("TRACE_LOG"))
	TraceLogHashFile = path.Join(logDir, os.Getenv("TRACE_LOG_HASH"))

	VerifyLogHash(TraceLogFile, TraceLogHashFile)
	traceLog, err := os.OpenFile(TraceLogFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0700)
	if err != nil {
		log.Fatalln("Failed to open trace log file: ", err)
	}
	Trace = log.New(io.MultiWriter(traceLog, os.Stdout), "TRACE: ", log.Ldate|log.Ltime|log.Lshortfile)

	VerifyLogHash(InfoLogFile, InfoLogHashFile)
	infoLog, err := os.OpenFile(InfoLogFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0700)
	if err != nil {
		log.Fatalln("Failed to open info log file: ", err)
	}
	Info = log.New(io.MultiWriter(infoLog, os.Stdout), "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
}

// LogHashing hashes the content of the log file and stores the hash values onto a .txt file.
func LogHashing(logFileName string, hashFileName string) {
	f, err := os.Open(logFileName)
	defer f.Close()
	if err != nil {
		if os.IsNotExist(err) {
			Trace.Fatalln("Log file not found!", err)
		}
		Trace.Fatalln("Failed to log file: ", err)
	}
	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		Trace.Println("Error occured while copying from log file: ", err)
	}

	hash := hex.EncodeToString(h.Sum(nil))

	err = os.WriteFile(hashFileName, []byte(hash), 0700)
	if err != nil {
		Trace.Println("Error occured while writing to log file hash file: ", err)
	}
}

// VerifyLogHash performs hash verification of the content of the log files.
func VerifyLogHash(logFileName string, hashFileName string) {
	f, err := os.Open(logFileName)
	defer f.Close()
	if err != nil {
		if os.IsNotExist(err) {
			log.Fatalln("Log file not found!", err)
		}
		log.Fatalln("Failed to open log file: ", err)
	}

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		log.Println("Error occured while copying from log file: ", err)
	}

	hash := hex.EncodeToString(h.Sum(nil))

	checkSum, err := os.ReadFile(hashFileName)
	if err != nil {
		if os.IsNotExist(err) {
			log.Println("Log file hash file not found!", err)
		}
		log.Fatalln("Error occured while reading records from Log file hash file: ", err)
	}

	if string(checkSum) != hash {
		log.Printf("%s has been tampered with!", logFileName)
	}
}
