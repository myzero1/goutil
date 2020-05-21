package z1log

import (
	"fmt"
	"os"
	"sync"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

/*
https://www.jianshu.com/p/d729c7ec9c85
https://juejin.im/post/5bffa2f15188256693607d7c
https://github.com/uber-go/zap/issues/715
https://my.oschina.net/noevilme/blog/3111521
1、不同级别的日志输出到不同的日志文件中。
2、日志文件按照文件大小或日期进行切割存储，以避免单一日志文件过大。
3、日志使用简单方便，一次定义全局使用。
*/

// Z1logger define
type Z1logger struct {
	zapLogger        *zap.Logger
	zapSugaredLogger *zap.SugaredLogger
	logPath          string //日志文件路径
	maxSize          int    //单个文件大小,MB
	maxBackups       int    //保存的文件个数
	maxAge           int    //保存的天数， 没有的话不删除
	compress         bool   //压缩
	jsonFormat       bool   //是否输出为json格式
	showLine         bool   //显示代码行
	output           string //输出到哪里，console/file/both,default both
	callerSkip       int
	lock             sync.RWMutex
}

func (log *Z1logger) getWriter(filename string) zapcore.WriteSyncer {
	hook := lumberjack.Logger{
		Filename:   filename,       // 日志文件路径
		MaxSize:    log.maxSize,    // megabytes
		MaxBackups: log.maxBackups, // 最多保留300个备份
		Compress:   log.compress,   // 是否压缩 disabled by default
	}

	if log.maxAge > 0 {
		hook.MaxAge = log.maxAge // days
	}

	var syncer zapcore.WriteSyncer

	switch log.output {
	case "console":
		syncer = zapcore.AddSync(os.Stdout)
	case "file":
		syncer = zapcore.AddSync(&hook)
	case "both":
		syncer = zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(&hook))
	}

	return syncer
}

func (log *Z1logger) reNewZapLog() {
	encoderConfig := zapcore.EncoderConfig{
		TimeKey: "time",
		// TimeKey:  "ts",
		LevelKey:  "level",
		NameKey:   "logger",
		CallerKey: "line",
		// CallerKey:     "file",
		MessageKey:    "msg",
		StacktraceKey: "stacktrace",
		LineEnding:    zapcore.DefaultLineEnding,
		EncodeLevel:   zapcore.LowercaseLevelEncoder, // 小写编码器
		EncodeTime:    zapcore.ISO8601TimeEncoder,    // ISO8601 UTC 时间格式
		// EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		// 	enc.AppendString(t.Format("2006-01-02 15:04:05"))
		// },
		EncodeDuration: zapcore.SecondsDurationEncoder,
		// EncodeDuration: func(d time.Duration, enc zapcore.PrimitiveArrayEncoder) {
		// 	enc.AppendInt64(int64(d) / 1000000)
		// },
		EncodeCaller: zapcore.ShortCallerEncoder, // 全路径编码器
		EncodeName:   zapcore.FullNameEncoder,
	}

	var encoder zapcore.Encoder
	if log.jsonFormat {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	} else {
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	}

	// 实现两个判断日志等级的interface
	// 设置日志级别,debug可以打印出info,debug,warn；info级别可以打印warn，info；warn只能打印warn
	// debug->info->warn->error
	debugLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl == zapcore.DebugLevel
	})
	infoLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl == zapcore.InfoLevel
	})
	warnLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl == zapcore.WarnLevel
	})
	errorLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl == zapcore.ErrorLevel
	})

	// 获取 info、error日志文件的io.Writer 抽象 getWriter() 在下方实现
	debugWriter := log.getWriter(fmt.Sprintf("%s/debug/debug.log", log.logPath))
	infoWriter := log.getWriter(fmt.Sprintf("%s/info/info.log", log.logPath))
	warnWriter := log.getWriter(fmt.Sprintf("%s/warn/warn.log", log.logPath))
	errorWriter := log.getWriter(fmt.Sprintf("%s/error/error.log", log.logPath))

	// 最后创建具体的Logger
	core := zapcore.NewTee(
		zapcore.NewCore(encoder, debugWriter, debugLevel),
		zapcore.NewCore(encoder, infoWriter, infoLevel),
		zapcore.NewCore(encoder, warnWriter, warnLevel),
		zapcore.NewCore(encoder, errorWriter, errorLevel),
	)

	log.zapLogger = zap.New(core, zap.AddCaller())
	if log.showLine {
		log.zapLogger = log.zapLogger.WithOptions(zap.AddCaller())
	}

	log.zapLogger = log.zapLogger.WithOptions(zap.AddCallerSkip(log.callerSkip))

	log.zapSugaredLogger = z1logger.zapLogger.Sugar()
}

// NewZ1logger new
func NewZ1logger() *Z1logger {
	log := Z1logger{
		zapLogger:        nil,
		zapSugaredLogger: nil,
		logPath:          "./logs",
		// maxSize:      100,
		// maxBackups:   100,
		// maxAge:       30,
		compress:   false,
		jsonFormat: false,
		showLine:   true,
		output:     "both",
		callerSkip: 1,
	}
	return &log

}

// init ===index=== 一旦某一个包被使用，则这个包下边的init函数将会被执行，且只执行一次
var z1logger *Z1logger

func init() {
	z1logger = NewZ1logger()
	z1logger.reNewZapLog()
}

// Debug define
func Debug(args ...interface{}) {
	z1logger.zapSugaredLogger.Debug(args...)
}

// Debugf defin
func Debugf(template string, args ...interface{}) {
	z1logger.zapSugaredLogger.Debugf(template, args...)
}

// Info defin
func Info(args ...interface{}) {
	z1logger.zapSugaredLogger.Info(args...)
}

// Infof defin
func Infof(template string, args ...interface{}) {
	z1logger.zapSugaredLogger.Infof(template, args...)
}

// Warn defin
func Warn(args ...interface{}) {
	z1logger.zapSugaredLogger.Warn(args...)
}

// Warnf defin
func Warnf(template string, args ...interface{}) {
	z1logger.zapSugaredLogger.Warnf(template, args...)
}

// Error defin
func Error(args ...interface{}) {
	z1logger.zapSugaredLogger.Error(args...)
}

// Errorf defin
func Errorf(template string, args ...interface{}) {
	z1logger.zapSugaredLogger.Errorf(template, args...)
}

// DPanic defin
func DPanic(args ...interface{}) {
	z1logger.zapSugaredLogger.DPanic(args...)
}

// DPanicf defin
func DPanicf(template string, args ...interface{}) {
	z1logger.zapSugaredLogger.DPanicf(template, args...)
}

// Panic defin
func Panic(args ...interface{}) {
	z1logger.zapSugaredLogger.Panic(args...)
}

// Panicf defin
func Panicf(template string, args ...interface{}) {
	z1logger.zapSugaredLogger.Panicf(template, args...)
}

// Fatal defin
func Fatal(args ...interface{}) {
	z1logger.zapSugaredLogger.Fatal(args...)
}

// Fatalf defin
func Fatalf(template string, args ...interface{}) {
	z1logger.zapSugaredLogger.Fatalf(template, args...)
}

// SetLogPath default ./logs
func SetLogPath(logPath string) {
	z1logger.lock.Lock()
	defer z1logger.lock.Unlock()
	z1logger.logPath = logPath
	z1logger.reNewZapLog()
}

// SetMaxSize defin
func SetMaxSize(maxSize int) {
	z1logger.lock.Lock()
	defer z1logger.lock.Unlock()
	z1logger.maxSize = maxSize
	z1logger.reNewZapLog()
}

// SetMaxBackups defin
func SetMaxBackups(maxBackups int) {
	z1logger.lock.Lock()
	defer z1logger.lock.Unlock()
	z1logger.maxBackups = maxBackups
	z1logger.reNewZapLog()
}

// SetMaxAge defin
func SetMaxAge(maxAge int) {
	z1logger.lock.Lock()
	defer z1logger.lock.Unlock()
	z1logger.maxAge = maxAge
	z1logger.reNewZapLog()
}

// SetCompress default false
func SetCompress(compress bool) {
	z1logger.lock.Lock()
	defer z1logger.lock.Unlock()
	z1logger.compress = compress
	z1logger.reNewZapLog()
}

// SetJsonFormat default false
func SetJsonFormat(jsonFormat bool) {
	z1logger.lock.Lock()
	defer z1logger.lock.Unlock()
	z1logger.jsonFormat = jsonFormat
	z1logger.reNewZapLog()
}

// SetShowLine default true
func SetShowLine(showLine bool) {
	z1logger.lock.Lock()
	defer z1logger.lock.Unlock()
	z1logger.showLine = showLine
	z1logger.reNewZapLog()
}

// SetOutput console/file/both,default both
func SetOutput(output string) {
	z1logger.lock.Lock()
	defer z1logger.lock.Unlock()
	z1logger.output = output
	z1logger.reNewZapLog()
}

// SetCallerSkip default 1
func SetCallerSkip(skip int) {
	z1logger.lock.Lock()
	defer z1logger.lock.Unlock()
	z1logger.callerSkip = skip
	z1logger.reNewZapLog()
}
