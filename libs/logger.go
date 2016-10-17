package libs

import (
	"os"
	"log"
)

var (
	Log Logging
	log_file *os.File
	err_log_file *os.File
)

type LogLevel int

type Logging struct  {
	Level LogLevel
	NormalLog  *log.Logger
	ErrorLog *log.Logger
}


const (
	VERBOSE LogLevel = iota
	DEBUG
	INFO
	WARNING
	FATAL
)

const (
	SERVER_TAG = " [ SERV ] "
	VERBOSE_TAG =" [ VERB ] "
	DEBUG_TAG = " [ DEBU ] "
	INFO_TAG = " [ INFO ] "
	WARNING_TAG = " [ WARN ] "
	FATAL_TAG = " [ FATA ] "
)

func (logging *Logging) Verbose(v ...interface{})  {
	if(logging.Level <= VERBOSE){
		logging.NormalLog.SetPrefix(VERBOSE_TAG)
		logging.NormalLog.Println(v...)
	}
}

func (logging *Logging) Verbosef(format string, v ...interface{}) {
	if(logging.Level <= VERBOSE){
		logging.NormalLog.SetPrefix(VERBOSE_TAG)
		logging.NormalLog.Printf(format,v...)
	}
}


func (logging *Logging) Debug(v ...interface{})  {
	if(logging.Level <= DEBUG){
		logging.NormalLog.SetPrefix(DEBUG_TAG)
		logging.NormalLog.Println(v...)
	}
}

func (logging *Logging) Debugf(format string, v ...interface{}) {
	if(logging.Level <= DEBUG){
		logging.NormalLog.SetPrefix(DEBUG_TAG)
		logging.NormalLog.Printf(format,v...)
	}
}

func (logging *Logging) Info(v ...interface{})  {
	if(logging.Level <= INFO){
		logging.NormalLog.SetPrefix(INFO_TAG)
		logging.NormalLog.Println(v...)
	}
}

func (logging *Logging) Infof(format string, v ...interface{})  {
	if(logging.Level <= INFO){
		logging.NormalLog.SetPrefix(INFO_TAG)
		logging.NormalLog.Printf(format,v...)
	}
}

func (logging *Logging) Warning(v ...interface{})  {
	if(logging.Level <= WARNING){
		logging.NormalLog.SetPrefix(WARNING_TAG)
		logging.NormalLog.Println(v...)
	}
}

func (logging *Logging) Warningf(format string, v ...interface{})  {
	if(logging.Level <= WARNING){
		logging.NormalLog.SetPrefix(WARNING_TAG)
		logging.NormalLog.Printf(format,v...)
	}
}

func (logging *Logging) Fatal(v ...interface{})  {
	if(logging.Level <= FATAL){
		logging.ErrorLog.SetPrefix(FATAL_TAG)
		logging.ErrorLog.Println(v...)
	}
}

func (logging *Logging) Fatalf(format string, v ...interface{})  {
	if(logging.Level <= FATAL){
		logging.ErrorLog.SetPrefix(FATAL_TAG)
		logging.ErrorLog.Printf(format,v...)
	}
}

func LoadLoggerModule()  {
	if(Config.LogToFile){
		output, err := os.OpenFile(Config.LogFilePath, os.O_WRONLY  | os.O_SYNC | os.O_CREATE | os.O_APPEND, 0666)
		if err != nil{
			log.Fatalln(err)
			os.Exit(1)
		}else{
			log_file = output
		}

		output_err, err := os.OpenFile(Config.ErrLogFilePath, os.O_WRONLY  | os.O_SYNC | os.O_CREATE | os.O_APPEND, 0666)
		if err != nil{
			log.Fatalln(err)
			os.Exit(1)
		}else{
			err_log_file = output_err
		}

	}else{
		log_file = os.Stdout
		err_log_file = os.Stderr
	}
	normalLog := log.New(log_file,"",log.LstdFlags)
	errorLog := log.New(err_log_file,"",log.Lshortfile|log.LstdFlags)

	var level LogLevel
	switch Config.LogLevel {
	case "verbose":
		level = VERBOSE
	case "debug":
		level =  DEBUG
	case "info":
		level =  INFO
	case "warning":
		level =  WARNING
	case "fatal":
		level =  FATAL
	default:
		level = DEBUG
	}

	Log = Logging{level,normalLog,errorLog}

	PrintModuleLoaded("Logger")

}

func ReleaseLoggerModule()  {
	if Config.LogToFile {
		if log_file != nil{
			log_file.Close()
		}

		if err_log_file != nil {
			err_log_file.Close()
		}
	}
	PrintModuleRelease("Logger")
}