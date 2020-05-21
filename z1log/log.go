package z1log

import (
	"fmt"
	"os"
	"time"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Z1logger struct {
	zapLogger    *zap.Logger
	logPath      string //日志文件路径
	maxSize      int    //单个文件大小,MB
	maxBackups   int    //保存的文件个数
	maxAge       int    //保存的天数， 没有的话不删除
	compress     bool   //压缩
	jsonFormat   bool   //是否输出为json格式
	showLine     bool   //显示代码行
	logInConsole bool   //是否同时输出到控制台
}

func (log *Z1logger) getWriter(filename string) zapcore.WriteSyncer {
	hook := lumberjack.Logger{
		Filename:   log.logPath,    // 日志文件路径
		MaxSize:    log.maxSize,    // megabytes
		MaxBackups: log.maxBackups, // 最多保留300个备份
		Compress:   log.compress,   // 是否压缩 disabled by default
	}
	if log.maxAge > 0 {
		hook.MaxAge = log.maxAge // days
	}

	var syncer zapcore.WriteSyncer
	if log.logInConsole {
		syncer = zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(&hook))
	} else {
		syncer = zapcore.AddSync(&hook)
	}

	return syncer

}

func (log *Z1logger) reNewZapLog() {

	encoderConfig := zapcore.EncoderConfig{
		// TimeKey:       "time",
		TimeKey:  "ts",
		LevelKey: "level",
		NameKey:  "logger",
		// CallerKey: "line",
		CallerKey:     "file",
		MessageKey:    "msg",
		StacktraceKey: "stacktrace",
		LineEnding:    zapcore.DefaultLineEnding,
		EncodeLevel:   zapcore.LowercaseLevelEncoder, // 小写编码器
		// EncodeTime:     zapcore.ISO8601TimeEncoder,     // ISO8601 UTC 时间格式
		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format("2006-01-02 15:04:05"))
		},
		// EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeDuration: func(d time.Duration, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendInt64(int64(d) / 1000000)
		},
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
	debugWriter := log.getWriter(fmt.Sprintf("%s/debug.log", log.logPath))
	infoWriter := log.getWriter(fmt.Sprintf("%s/inf0.log", log.logPath))
	warnWriter := log.getWriter(fmt.Sprintf("%s/warn.log", log.logPath))
	errorWriter := log.getWriter(fmt.Sprintf("%s/error.log", log.logPath))

	// 最后创建具体的Logger
	core := zapcore.NewTee(
		zapcore.NewCore(encoder, debugWriter, debugLevel),
		zapcore.NewCore(encoder, infoWriter, infoLevel),
		zapcore.NewCore(encoder, warnWriter, warnLevel),
		zapcore.NewCore(encoder, errorWriter, errorLevel),
	)

	log.zapLogger = zap.New(core)
	if log.showLine {
		log.zapLogger = log.zapLogger.WithOptions(zap.AddCaller())
	}
}

// NewZ1logger new
func NewZ1logger() *Z1logger {
	log := Z1logger{
		zapLogger:    nil,
		logPath:      "./logs",
		maxSize:      10,
		maxBackups:   100,
		maxAge:       30,
		compress:     false,
		jsonFormat:   false,
		showLine:     true,
		logInConsole: true,
	}
	return &log

}

// init ===index=== 一旦某一个包被使用，则这个包下边的init函数将会被执行，且只执行一次
var z1logger *Z1logger

func init() {
	z1logger = NewZ1logger()
	z1logger.reNewZapLog()
}

func Info(args ...interface{}) {
	errorLogger := z1logger.zapLogger.Sugar()
	errorLogger.Info(args...)
}
