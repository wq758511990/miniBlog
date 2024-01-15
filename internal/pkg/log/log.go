package log

import (
	"context"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"myMiniBlog/internal/pkg/known"
	"sync"
	"time"
)

type zapLogger struct {
	z *zap.Logger
}

var (
	mu  sync.Mutex
	std = NewLogger(NewOptions())
)

type Logger interface {
	Debugw(msg string, keysAndValues ...interface{})
	Infow(msg string, keysAndValues ...interface{})
	Warnw(msg string, keysAndValues ...interface{})
	Errorw(msg string, keysAndValues ...interface{})
	Panicw(msg string, keysAndValues ...interface{})
	Fatalw(msg string, keysAndValues ...interface{})
	Sync()
}

func Init(opts *Options) {
	mu.Lock()
	defer mu.Unlock()
	std = NewLogger(opts)
}

func NewLogger(opts *Options) *zapLogger {
	if opts == nil {
		opts = NewOptions()
	}
	var zapLevel zapcore.Level
	if err := zapLevel.UnmarshalText([]byte(opts.Level)); err != nil {
		zapLevel = zapcore.InfoLevel
	}
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.MessageKey = "message"
	encoderConfig.TimeKey = "timestamp"
	encoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format(time.DateTime + ".000"))
	}
	encoderConfig.EncodeDuration = func(duration time.Duration, encoder zapcore.PrimitiveArrayEncoder) {
		encoder.AppendFloat64(float64(duration) / float64(time.Millisecond))
	}
	cfg := &zap.Config{
		// 是否在日志中显示调用日志所在的文件和行号，例如：`"caller":"miniblog/miniblog.go:75"`
		DisableCaller: opts.DisableCaller,
		// 是否禁止在 panic 及以上级别打印堆栈信息
		DisableStacktrace: opts.DisableStacktrace,
		// 指定日志级别
		Level: zap.NewAtomicLevelAt(zapLevel),
		// 指定日志显示格式，可选值：console, json
		Encoding:      opts.Format,
		EncoderConfig: encoderConfig,
		// 指定日志输出位置
		OutputPaths: opts.OutputPaths,
		// 设置 zap 内部错误输出位置
		ErrorOutputPaths: []string{"stderr"},
	}
	z, err := cfg.Build(zap.AddStacktrace(zapcore.PanicLevel), zap.AddCallerSkip(1))
	if err != nil {
		panic(err)
	}
	logger := &zapLogger{z: z}
	zap.RedirectStdLog(z)

	return logger
}
func Sync() { std.Sync() }

func (l *zapLogger) Sync() {
	_ = l.z.Sync()
}

func C(ctx context.Context) *zapLogger {
	return std.C(ctx)
}
func (l *zapLogger) C(ctx context.Context) *zapLogger {
	lc := l.clone()
	if requestID := ctx.Value(known.XRequestIDKey); requestID != nil {
		lc.z = lc.z.With(zap.Any(known.XRequestIDKey, requestID))
	}
	if userId := ctx.Value(known.XUsernameKey); userId != nil {
		lc.z = lc.z.With(zap.Any(known.XUsernameKey, userId))
	}
	return lc
}

func (l *zapLogger) clone() *zapLogger {
	lc := *l
	return &lc
}

// Debugw 这个Debugw是包级别的，引用了全局的std实例
func Debugw(msg string, keysAndValues ...interface{}) {
	std.z.Sugar().Debugw(msg, keysAndValues...)
}
func Infow(msg string, keysAndValues ...interface{}) {
	std.z.Sugar().Infow(msg, keysAndValues...)
}
func Warnw(msg string, keysAndValues ...interface{}) {
	std.z.Sugar().Warnw(msg, keysAndValues...)
}
func Errorw(msg string, keysAndValues ...interface{}) {
	std.z.Sugar().Errorw(msg, keysAndValues...)
}
func Panicw(msg string, keysAndValues ...interface{}) {
	std.z.Sugar().Panicw(msg, keysAndValues...)
}
func Fatalw(msg string, keysAndValues ...interface{}) {
	std.z.Sugar().Fatalw(msg, keysAndValues...)
}

// Debugw 这个是局部变量级别的
func (l *zapLogger) Debugw(msg string, keysAndValues ...interface{}) {
	l.z.Sugar().Debugw(msg, keysAndValues...)
}
func (l *zapLogger) Infow(msg string, keysAndValues ...interface{}) {
	l.z.Sugar().Infow(msg, keysAndValues...)
}
func (l *zapLogger) Warnw(msg string, keysAndValues ...interface{}) {
	l.z.Sugar().Warnw(msg, keysAndValues...)
}
func (l *zapLogger) Errorw(msg string, keysAndValues ...interface{}) {
	l.z.Sugar().Errorw(msg, keysAndValues...)
}
func (l *zapLogger) Panicw(msg string, keysAndValues ...interface{}) {
	l.z.Sugar().Panicw(msg, keysAndValues...)
}
func (l *zapLogger) Fatalw(msg string, keysAndValues ...interface{}) {
	l.z.Sugar().Fatalw(msg, keysAndValues...)
}
