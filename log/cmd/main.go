package main

import (
	"os"
	"sync"

	"github.com/liangdas/mqant/log"
	"github.com/liangdas/mqant/logv2"
)

func main() {
	defaultLogger := log.DefaultLogger
	// 不传参数， 不输出
	// defaultLogger.Log(logv2.LevelDebug)

	// 传入数据少于等于3位不输出
	// defaultLogger.Log(logv2.LevelDebug, "test", "xxxx", "xxxxx")

	// // 正常输出
	// defaultLogger.Log(logv2.LevelDebug, "test1", "xxxx", "xxxxx%s", "xxxxxxx")

	// defaultLogger.Log(logv2.LevelInfo)

	// // 传入数据少于等于3位不输出
	// defaultLogger.Log(logv2.LevelInfo, "test", "xxxx", "xxxxx")
	// // defaultLogger.Log(logv2.LevelDebug, "test", "xxxx", "xxxxx %s", "xxxxxxx")
	// defaultLogger.Log(logv2.LevelInfo, "test2", "xxxx", "xxxxx %s", "xxxxxxx")
	// // defaultLogger.Log(logv2.LevelInfo, "test", "xxxx", "xxxxx%v", "xxxxxxx")
	// // defaultLogger.Flush()
	// defaultLogger.Log(logv2.LevelWarn)

	// // 传入数据少于等于3位不输出
	// defaultLogger.Log(logv2.LevelWarn, "test", "xxxx", "xxxxx")
	// defaultLogger.Log(logv2.LevelWarn, "test3", "xxxx", "xxxxx")
	// defaultLogger.Log(logv2.LevelWarn, "test3", "xxxx", "xxxxx %s", "xxxxxxx")

	// log.NewLogger(true, log.WithDaily(true))

	wg := &sync.WaitGroup{}
	os.Mkdir("logs", os.ModePerm)
	for i := 0; i < 10000; i++ {
		wg.Add(1)
		go func() {
			defaultLogger.Log(logv2.LevelInfo, "test2", "xxxx", "xxxxx %s", "xxxxxxx")
			wg.Done()
		}()
	}
	wg.Wait()
}
