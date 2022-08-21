package logger

import (
	"context"
	"fmt"
	"github.com/augustazz/camellia/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"strings"
)

var Holder *zap.Logger
var Sugar *zap.SugaredLogger

var defaultLogPath = "./logs"

func SetupLogger(ctx context.Context, appName string, conf config.LogConfig) func() {
	hook := lumberjack.Logger{
		Filename:   buildFileName(conf.Path, appName), // 日志文件路径
		MaxSize:    5 * 1024,                          // 每个日志文件保存的最大尺寸 单位：M
		MaxBackups: 128,                               // 日志文件最多保存多少个备份
		MaxAge:     7,                                 // 文件最多保存多少天
		Compress:   true,                              // 是否压缩
	}

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:  "@timestamp",
		LevelKey: "level",
		NameKey:  "logger",
		//CallerKey:      "class",
		MessageKey:     "message",
		StacktraceKey:  "stack",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,  // 小写编码器
		EncodeTime:     zapcore.ISO8601TimeEncoder,     // ISO8601 UTC 时间格式
		EncodeDuration: zapcore.SecondsDurationEncoder, //
		EncodeCaller:   zapcore.ShortCallerEncoder,
		EncodeName:     zapcore.FullNameEncoder,
	}

	// 设置日志级别
	atomicLevel := zap.NewAtomicLevel()
	if conf.Debug {
		atomicLevel.SetLevel(zap.DebugLevel)
	} else {
		atomicLevel.SetLevel(zap.InfoLevel)
	}

	//writer
	writer := []zapcore.WriteSyncer{zapcore.AddSync(&hook)}
	if conf.Debug {
		writer = append(writer, zapcore.AddSync(os.Stdout))
	}

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		//zapcore.NewConsoleEncoder(encoderConfig),
		zapcore.NewMultiWriteSyncer(writer...),
		atomicLevel, // 日志级别
	)

	// 开启开发模式，堆栈跟踪
	caller := zap.AddCaller()
	// 开启文件及行号
	development := zap.Development()
	// 设置初始化字段
	host, _ := os.Hostname()
	filed := zap.Fields(zap.String("service", appName), zap.String("host", host))

	// 构造日志
	Holder = zap.New(core, caller, development, filed)
	Sugar = Holder.Sugar()

	// 程序推出刷新数据
	return func() {
		e := Holder.Sync()
		if e != nil {
			Sugar.Error("logger sync err", e)
		}
	}
}

func buildFileName(configPath, appName string) string {
	if configPath == "" {
		configPath = defaultLogPath
	}
	if !strings.HasSuffix(configPath, "/") {
		configPath += "/"
	}
	return fmt.Sprintf("%s%s.json", configPath, appName)
}
