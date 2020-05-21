package z1log

import (
	"os"

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

func (log *Z1logger) reNewZapLog() {
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

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "line",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,  // 小写编码器
		EncodeTime:     zapcore.ISO8601TimeEncoder,     // ISO8601 UTC 时间格式
		EncodeDuration: zapcore.SecondsDurationEncoder, //
		EncodeCaller:   zapcore.ShortCallerEncoder,     // 全路径编码器
		EncodeName:     zapcore.FullNameEncoder,
	}

	var encoder zapcore.Encoder
	if log.jsonFormat {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	} else {
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	}

	core := zapcore.NewCore(
		encoder,
		syncer,
		zap.InfoLevel,
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
