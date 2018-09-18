package log

import "go.uber.org/zap"

var rawLogger *zap.Logger
var sugerLogger *zap.SugaredLogger

func Debug(args ...interface{}) {
	sugerLogger.Debug(args...)
}
func Debugf(template string, args ...interface{}) {
	sugerLogger.Debugw(template, args...)
}
func Debugw(msg string, args ...interface{}) {
	sugerLogger.Debugf(msg, args...)
}

func Info(args ...interface{}) {
	sugerLogger.Info(args...)
}
func Infof(template string, args ...interface{}) {
	sugerLogger.Infof(template, args...)
}
func Infow(msg string, args ...interface{}) {
	sugerLogger.Infow(msg, args...)
}

func Warn(args ...interface{}) {
	sugerLogger.Warn(args...)
}
func Warnf(template string, args ...interface{}) {
	sugerLogger.Warnf(template, args...)
}
func Warnw(msg string, args ...interface{}) {
	sugerLogger.Warnw(msg, args...)
}

func Fatal(args ...interface{}) {
	sugerLogger.Fatal(args...)
}
func Fatalf(template string, args ...interface{}) {
	sugerLogger.Fatalf(template, args...)
}
func Fatalw(msg string, args ...interface{}) {
	sugerLogger.Fatalw(msg, args...)
}

func Error(args ...interface{}) {
	sugerLogger.Error(args...)
}
func Errorf(template string, args ...interface{}) {
	sugerLogger.Errorf(template, args...)
}
func Errorw(msg string, args ...interface{}) {
	sugerLogger.Errorw(msg, args...)
}

func Panic(args ...interface{}) {
	sugerLogger.Panic(args...)
}
func Panicf(template string, args ...interface{}) {
	sugerLogger.Panicf(template, args...)
}
func Panicw(msg string, args ...interface{}) {
	sugerLogger.Panicw(msg, args...)
}


func init() {
	var err error
	//rawLogger, err = zap.NewProduction()
	rawLogger, err = zap.NewDevelopment()
	
	if err != nil {
		panic(err)
	}
	
	sugerLogger = rawLogger.Sugar()
	
	Info("[log] zap log initialized")
}