package log

import (
	logs "github.com/liangdas/mqant/log/beego"
)

type BeegoWay int

const (
	Console BeegoWay = iota + 1
	File
	Multifile
	DingtalkWriter
	Conn
	Mail
	Es
	Slack
	jianliao
)

func String(b BeegoWay) string {
	switch b {
	case 1:
		return logs.AdapterConsole
	case 2:
		return logs.AdapterFile
	case 3:
		return logs.AdapterMultiFile
	case 4:
		return logs.AdapterDingtalk
	case 5:
		return logs.AdapterConn
	case 6:
		return logs.AdapterMail
	case 7:
		return logs.AdapterEs
	case 8:
		return logs.AdapterSlack
	case 9:
		return logs.AdapterJianLiao
	}
	return logs.AdapterConsole
}
