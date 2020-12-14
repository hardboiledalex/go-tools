package logging

import (
	"fmt"
	"github.com/gookit/color"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

type LogLevel int

const (
	TRACE LogLevel = iota
	DEBUG
	SUCCESS
	INFO
	WARN
	ERROR
)

const (
	traceStr   = "TRACE"
	debugStr   = "DEBUG"
	successStr = "SUCCESS"
	infoStr    = "INFO"
	warnStr    = "WARN"
	errorStr   = "ERROR"
)

var (
	fileLogger            *log.Logger
	currentLogLevel       LogLevel
	currentLogLevelStr    string // log level string; internal to logging package only
	CurrentLogLevelString string // log level from CLI params
)

func checkFolderExistsAndIsWritable(folderPath string) bool {
	if logsFolderFileInfo, err := os.Stat(folderPath); err == nil {
		if logsFolderFileInfo.IsDir() {
			// check the <folderPath>/t can be written
			if logsFolderFileInfo.Mode().Perm()&(1<<uint(8)|1<<uint(7)|1<<uint(6)) == 0700 {
				var data []byte
				createError := ioutil.WriteFile(folderPath+"/t", data, 0644)
				if createError == nil {
					if deleteError := os.Remove(folderPath + "/t"); deleteError == nil {
						return true
					} else {
						log.Fatal(deleteError)
					}
				} else {
					log.Fatal(createError)
				}
			}
		}
	}
	return false
}

func InitLogging() {
	// If for any reason you need to change the filename format do not change the numbers
	// itself or the TZ name "MST"
	// These value are hardcoded as samples in the time package.
	logFileName := "arcsight-install_" + time.Now().Format("20060102_150405_MST") + ".log"

	var logsFolder string
	// Check for ./arcsight/logs folder existence
	if checkFolderExistsAndIsWritable("arcsight/logs") {
		logsFolder = "arcsight/logs/"
	}
	if logsFolder == "" {
		if checkFolderExistsAndIsWritable("arcsight") {
			if err := os.Mkdir("arcsight/logs", 0755); err == nil {
				// if not present check for ./arcsight folder existence - if exists create subfolder logs
				logsFolder = "arcsight/logs/"
			}
		}
		// if not present do not create anything and store log in current folder
	}

	file, err := os.OpenFile(logsFolder+logFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	fileLogger = log.New(file, "", 0)
	currentLogLevelStr = strings.ToUpper(CurrentLogLevelString)
	switch currentLogLevelStr {
	case traceStr:
		currentLogLevel = TRACE
		break
	case debugStr:
		currentLogLevel = DEBUG
		break
	case warnStr:
		currentLogLevel = WARN
		break
	case errorStr:
		currentLogLevel = ERROR
		break
	case successStr:
		currentLogLevel = SUCCESS
		currentLogLevelStr = "INFO"
	default:
		currentLogLevel = INFO
		currentLogLevelStr = "INFO"
	}
}

func getCaller() string {
	if file, lineno, ok := getCallerFileAndLine(); ok {
		return filepath.Base(file) + ":" + strconv.Itoa(lineno)
	}
	return ""
}

func getCallerFileAndLine() (string, int, bool) {
	programCounters := make([]uintptr, 50)
	framesCount := runtime.Callers(2, programCounters)
	frames := runtime.CallersFrames(programCounters[:framesCount])

	for more := true; more; {
		var frameCandidate runtime.Frame
		frameCandidate, more = frames.Next()
		if !strings.HasSuffix(frameCandidate.File, "utils/logging/log.go") {
			return frameCandidate.File, frameCandidate.Line, true
		}
	}
	return "", 0, false
}

func Log(level LogLevel, msgs ...interface{}) {
	if consoleMessage, logMessage, ok := GetMessages(level, msgs); ok {
		fmt.Println(consoleMessage)
		fileLogger.Println(logMessage)
	}
}

func Logf(level LogLevel, format string, params ...interface{}) {
	Log(level, fmt.Sprintf(format, params...))
}

func LogTrace(msgs ...interface{})                     { Log(TRACE, msgs) }
func LogTracef(format string, params ...interface{})   { Logf(TRACE, format, params...) }
func LogDebug(msgs ...interface{})                     { Log(DEBUG, msgs) }
func LogDebugf(format string, params ...interface{})   { Logf(DEBUG, format, params...) }
func LogSuccess(msgs ...interface{})                   { Log(SUCCESS, msgs) }
func LogSuccessf(format string, params ...interface{}) { Logf(SUCCESS, format, params...) }
func LogInfo(msgs ...interface{})                      { Log(INFO, msgs) }
func LogInfof(format string, params ...interface{})    { Logf(INFO, format, params...) }
func LogWarn(msgs ...interface{})                      { Log(WARN, msgs) }
func LogWarnf(format string, params ...interface{})    { Logf(WARN, format, params...) }
func LogError(msgs ...interface{})                     { Log(ERROR, msgs) }
func LogErrorf(format string, params ...interface{})   { Logf(ERROR, format, params...) }

func LogMessageDirect(consoleMessage string, logMessage string) {
	if consoleMessage != "" {
		fmt.Println(consoleMessage)
	}
	if fileLogger != nil && logMessage != "" {
		fileLogger.Println(logMessage)
	}
}

func LogMessagesDirect(consoleMessages []string, logMessages []string) {
	if len(consoleMessages) == len(logMessages) {
		for index := 0; index < len(consoleMessages); index++ {
			LogMessageDirect(consoleMessages[index], logMessages[index])
		}
	} else {
		for _, message := range consoleMessages {
			LogMessageDirect(message, "")
		}
		if len(logMessages) > 0 {
			for _, message := range logMessages {
				LogMessageDirect("", message)
			}
		} else {
			for _, message := range consoleMessages {
				LogMessageDirect("", message)
			}
		}
	}
}

// returns coloured message for console, plain for logFile
func GetMessages(level LogLevel, msgs ...interface{}) (string, string, bool) {

	if level < currentLogLevel {
		return "", "", false
	}
	caller := getCaller()

	var consoleColor color.Color
	var logLevelStr string

	switch level {
	case TRACE:
		consoleColor = color.LightBlue
		logLevelStr = traceStr
		break
	case DEBUG:
		consoleColor = color.Blue
		logLevelStr = debugStr
		break
	case INFO:
		logLevelStr = infoStr + " "
		break
	case WARN:
		consoleColor = color.Yellow
		logLevelStr = warnStr + " "
		break
	case ERROR:
		consoleColor = color.Red
		logLevelStr = errorStr
		break
	case SUCCESS:
		consoleColor = color.Green
		logLevelStr = infoStr + " "
		break
	default:
		fmt.Println(caller, msgs)
	}

	var consoleMessage strings.Builder
	var fileMessage strings.Builder
	now := time.Now().Format("2006/01/02 15:04:05.000 MST ")
	fileMessage.WriteString(now + "[" + logLevelStr + "] " + caller + " ")
	alignSize := fileMessage.Len()
	for _, msgL1 := range msgs {
		switch msgL1Type := msgL1.(type) {
		default:
			rmString := fmt.Sprintf("%v", msgL1Type)
			consoleMessage.WriteString(rmString)
			writeStringToFile(&fileMessage, rmString, alignSize)
		case []interface{}:
			for _, msgL2 := range msgL1.([]interface{}) {
				switch msgL2Type := msgL2.(type) {
				default:
					rmString := fmt.Sprintf("%v", msgL2Type)
					consoleMessage.WriteString(rmString)
					writeStringToFile(&fileMessage, rmString, alignSize)
				case []interface{}:
					for _, realMessage := range msgL2Type {
						rmString := fmt.Sprintf("%v", realMessage)
						consoleMessage.WriteString(rmString)
						writeStringToFile(&fileMessage, rmString, alignSize)
					}
				}
			}
		}
	}
	return consoleColor.Sprint(consoleMessage.String()), fileMessage.String(), true
}

func writeStringToFile(fileMessage *strings.Builder, message string, alignSize int) {
	messageToWrite := strings.Trim(message, "\n")
	messageToWrite = strings.ReplaceAll(messageToWrite, "\n", "\n"+strings.Repeat(" ", alignSize))
	fileMessage.WriteString(messageToWrite)
}
