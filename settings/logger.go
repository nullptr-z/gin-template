package settings

import (
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func InitLogger() error {
	// 构建一个生产环境下推荐使用的配置
	config := zap.NewProductionConfig()

	// 修改日志时间格式为ISO8601
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	// 设置日志级别，例如Debug级别
	config.Level = zap.NewAtomicLevelAt(zap.DebugLevel)

	// 创建Logger实例
	logger, err := config.Build()
	if err != nil {
		return err
	}
	// 创建SugaredLogger
	// sugar := logger.Sugar()

	defer logger.Sync() // 退出时缓冲区的日志都刷到磁盘里

	// 替换全局的 Logger; 其他地方直接 zap.L就能访问这个 logger 了
	zap.ReplaceGlobals(logger)
	fmt.Println("Logger initialized.......")
	return nil
}

// func main() {
// 	writesyncer := &zapcore.WriteSyncer{}

// 	level := new(zapcore.Level)
// 	if err := level.UnmarshalText([]byte(viper.GetString("log.level"))); err != nil {
// 		return
// 	}
// 	core := zapcore.NewCore(nil, writesyncer, level)

// 	log := zap.New(core, zap.AddCaller())
// 	// 替换全局的 Logger
// 	zap.ReplaceGlobals(log)

// 	// level: debug
// 	// filename: app_log.log
// 	// max_size: 200
// 	// max_age: 30
// 	// max_backups: 7

// }
