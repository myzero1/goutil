package z1log

import (
	"fmt"
	"os"
	"strings"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
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
	_, err := os.Stat(log.logPath)
	if err == nil {
		os.MkdirAll(log.logPath, 777)
		if err != nil {
			panic(err)
		}
	}

	fmt.Println(log.logPath)

	// hook := lumberjack.Logger{
	// 	Filename:   filename,       // 日志文件路径
	// 	MaxSize:    log.maxSize,    // megabytes
	// 	MaxBackups: log.maxBackups, // 最多保留300个备份
	// 	Compress:   log.compress,   // 是否压缩 disabled by default
	// }

	// if log.maxAge > 0 {
	// 	hook.MaxAge = log.maxAge // days
	// }

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

	var syncer zapcore.WriteSyncer
	if log.logInConsole {
		syncer = zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(hook))
	} else {
		syncer = zapcore.AddSync(hook)
	}

	return syncer

}

func (log *Z1logger) reNewZapLog() {

	encoderConfig := zapcore.EncoderConfig{
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
	debugWriter := log.getWriter(fmt.Sprintf("%s-debug.log", log.logPath))
	infoWriter := log.getWriter(fmt.Sprintf("%s-info.log", log.logPath))
	warnWriter := log.getWriter(fmt.Sprintf("%s/warn.log", log.logPath))
	errorWriter := log.getWriter(fmt.Sprintf("%s/error.log", log.logPath))

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

func Debug(args ...interface{}) {
	errorLogger := z1logger.zapLogger.Sugar()
	errorLogger.Debug(args...)
}

func Info(args ...interface{}) {
	errorLogger := z1logger.zapLogger.Sugar()
	errorLogger.Info(args...)
}
