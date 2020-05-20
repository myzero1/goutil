package z1log

import (
	"io"
	"strings"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var errorLogger *zap.SugaredLogger

/*
https://www.jianshu.com/p/d729c7ec9c85
https://juejin.im/post/5bffa2f15188256693607d7c
https://github.com/uber-go/zap/issues/715
1、不同级别的日志输出到不同的日志文件中。
2、日志文件按照文件大小或日期进行切割存储，以避免单一日志文件过大。
3、日志使用简单方便，一次定义全局使用。
*/
func init() {
	// 设置一些基本日志格式 具体含义还比较好理解，直接看zap源码也不难懂
	encoder := zapcore.NewConsoleEncoder(zapcore.EncoderConfig{
		MessageKey:  "msg",
		LevelKey:    "level",
		EncodeLevel: zapcore.CapitalLevelEncoder,
		TimeKey:     "ts",
		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format("2006-01-02 15:04:05"))
		},
		CallerKey:    "file",
		EncodeCaller: zapcore.ShortCallerEncoder,
		EncodeDuration: func(d time.Duration, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendInt64(int64(d) / 1000000)
		},
	})

	// 实现两个判断日志等级的interface
	infoLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.InfoLevel
	})

	errorLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.ErrorLevel
	})

	// 获取 info、error日志文件的io.Writer 抽象 getWriter() 在下方实现
	infoWriter := getWriter("./logs/demo_info.log")
	errorWriter := getWriter("./logs/demo_error.log")

	// 最后创建具体的Logger
	core := zapcore.NewTee(
		zapcore.NewCore(encoder, zapcore.AddSync(infoWriter), infoLevel),
		zapcore.NewCore(encoder, zapcore.AddSync(errorWriter), errorLevel),
	)

	// log := zap.New(core, zap.AddCaller()) // 需要传入 zap.AddCaller() 才会显示打日志点的文件名和行数, 有点小坑
	log := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1)) // 需要传入 zap.AddCaller() 才会显示打日志点的文件名和行数, 有点小坑
	errorLogger = log.Sugar()
}

func getWriter(filename string) io.Writer {
	// 生成rotatelogs的Logger 实际生成的文件名 demo.log.YYmmddHH
	// demo.log是指向最新日志的链接
	// 保存7天内的日志，每1小时(整点)分割一次日志
	hook, err := rotatelogs.New(
		strings.Replace(filename, ".log", "", -1) + "-%Y%m%d%H.log", // 没有使用go风格反人类的format格式
		//rotatelogs.WithLinkName(filename),
		//rotatelogs.WithMaxAge(time.Hour*24*7),
		//rotatelogs.WithRotationTime(time.Hour),
	)

	if err != nil {
		panic(err)
	}
	return hook
}
func Debug(args ...interface{}) {
	errorLogger.Debug(args...)
}

func Debugf(template string, args ...interface{}) {
	errorLogger.Debugf(template, args...)
}

func Info(args ...interface{}) {
	errorLogger.Info(args...)
}

func Infof(template string, args ...interface{}) {
	errorLogger.Infof(template, args...)
}

func Warn(args ...interface{}) {
	errorLogger.Warn(args...)
}

func Warnf(template string, args ...interface{}) {
	errorLogger.Warnf(template, args...)
}

func Error(args ...interface{}) {
	errorLogger.Error(args...)
}

func Errorf(template string, args ...interface{}) {
	errorLogger.Errorf(template, args...)
}

func DPanic(args ...interface{}) {
	errorLogger.DPanic(args...)
}

func DPanicf(template string, args ...interface{}) {
	errorLogger.DPanicf(template, args...)
}

func Panic(args ...interface{}) {
	errorLogger.Panic(args...)
}

func Panicf(template string, args ...interface{}) {
	errorLogger.Panicf(template, args...)
}

func Fatal(args ...interface{}) {
	errorLogger.Fatal(args...)
}

func Fatalf(template string, args ...interface{}) {
	errorLogger.Fatalf(template, args...)
}
