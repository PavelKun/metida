package logger

import (
	"github.com/Dsmit05/metida/internal/config"
	"go.uber.org/zap"
	//"go.uber.org/zap/zapcore"
)

var L *zap.SugaredLogger
var Prod *zap.Logger

//https://stackoverflow.com/questions/57745017/how-to-initialize-a-zap-logger-once-and-reuse-it-in-other-go-files
// Todo: нормально сконфигурить логгер
func InitLogger(flagCmd config.CommandLineI) error {
	if flagCmd.IfDebagOn() {

	}

	newLogger, err := zap.NewProduction()
	if err != nil {
		return err
	}
	L = newLogger.Sugar()

	L.Info(" logger is initialized")
	Prod, err = zap.NewProduction()
	if err != nil {
		return err
	}

	return nil
}

//func Info(message string, text string) {
//	zapLog.Info(message, zap.String("text", text))
//}
//
//func Debug(message string, text string) {
//	zapLog.Debug(message, zap.String("text", text))
//}
//
//func Error(message string, err error) {
//	zapLog.Error(message, zap.Error(err))
//}
//
//func Fatal(message string,  err error) {
//	zapLog.Fatal(message, zap.Error(err))
//}
