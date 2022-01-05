package log

import (
	"testing"
	"time"

	"github.com/liangdas/mqant/logv2"
	"github.com/stretchr/testify/assert"
)

func TestBeegoLogger_Log(t *testing.T) {
	//DefaultLogger.Log(logv2.LevelDebug, "test", "1212", "1212121", "121212", "xxxxx")
	//DefaultLogger.Log(logv2.LevelInfo)
	//DefaultLogger.Log(logv2.LevelError)
	//DefaultLogger.Log(logv2.LevelFatal)
	//DefaultLogger.Log(logv2.LevelWarn)
	logger := DefaultLogger
	err := logger.Log(logv2.LevelInfo)
	assert.Error(t, err)

	err = logger.Log(logv2.LevelDebug, "test", "01", "xxx%s%s", "xxxx", "xxxx")
	assert.NoError(t, err)
	time.Sleep(10 * time.Second)
}
