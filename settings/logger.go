package settings

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/gin-gonic/gin"
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

// 自定义日志输出格式
func LoggerFormateOutput(g *gin.Context) {
	// 请求前
	startTime := time.Now()

	// 复制请求体，以便日志记录后仍可读取
	var requestBody bytes.Buffer
	if g.Request.Body != nil {
		bodyBytes, _ := ioutil.ReadAll(g.Request.Body)
		requestBody.Write(bodyBytes)
		// 重新设置请求体，以供后续使用
		g.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
	}
	fmt.Println("request arguments:", requestBody.String())
	queryParams := g.Request.URL.Query()
	fmt.Println("query arguments:", queryParams)

	// 处理请求
	g.Next()

	// 请求后
	endTime := time.Now()
	latencyTime := endTime.Sub(startTime)
	statusCode := g.Writer.Status()
	clientIP := g.ClientIP()

	// 使用方括号[]格式化日志内容
	zap.L().Info("request details",
		zap.String("method", g.Request.Method),
		zap.String("uri", g.Request.RequestURI),
		zap.Int("status", statusCode),
		zap.String("latency", fmt.Sprintf("[%s]", latencyTime)),
		zap.String("clientIP", fmt.Sprintf("[%s]", clientIP)),
		zap.String("request arguments", fmt.Sprintf("[%s]", requestBody.String())),
		zap.String("queryParams", fmt.Sprintf("[%s]", queryParams)),
		// zap.String("formData", fmt.Sprintf("[%s]", formData)),
	)
}
