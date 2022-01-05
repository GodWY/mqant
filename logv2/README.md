# 日志接口

> 目前日志库的一个困扰 是框架层自己使用了一个日志库，无法和外部通用，导致业务层无法控制

新建日志接口，使用传入日志组件接口的方式调用日志库。 业务可以随便选型，只要实现了日志组件的接口就可以。框架内部完全按照业务侧进行rpc接口的输出。

# Logger

## Usage

### Structured logging

```go
logger := log.NewStdLogger(os.Stdout)
// fields & valuer
logger = log.With(logger,
    "service.name", "hellworld",
    "service.version", "v1.0.0",
    "ts", log.DefaultTimestamp,
    "caller", log.DefaultCaller,
)
logger.Log(log.LevelInfo, "key", "value")

// helper
helper := log.NewHelper(logger)
helper.Log(log.LevelInfo, "key", "value")
helper.Info("info message")
helper.Infof("info %s", "message")
helper.Infow("key", "value")

// filter
log := log.NewHelper(log.NewFilter(logger,
	log.FilterLevel(log.LevelInfo),
	log.FilterKey("foo"),
	log.FilterValue("bar"),
	log.FilterFunc(customFilter),
))
log.Debug("debug log")
log.Info("info log")
log.Warn("warn log")
log.Error("warn log")
```
