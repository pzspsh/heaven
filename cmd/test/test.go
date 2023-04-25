package main

import (
	"heaven/app/APVE/pkg/common/execute"
	"heaven/app/APVE/pkg/controller/runner"
)

func main() {
	runner.Runner()
	execute.Execute()
	//var my store.Mysql = store.Mysql{
	//	User: "root",
	//	Pass: "pzspsh246789",
	//	Ip:   "localhost",
	//	Port: "3306",
	//	DB:   "mysql",
	//}
	//db := my.Mysql()
	//fmt.Println(db.Ping())
}

// 日志文件记录
//import (
//	logger "heaven/logging"
//	"runtime"
//	"strconv"
//	"time"
//)

//func Log(i int) {
//	logger.Debug("Debug>>>>>>>>>>>>>>>>>>>>>>" + strconv.Itoa(i))
//	logger.Info("Info>>>>>>>>>>>>>>>>>>>>>>>>>" + strconv.Itoa(i))
//	logger.Warn("Warn>>>>>>>>>>>>>>>>>>>>>>>>>" + strconv.Itoa(i))
//	logger.Error("Error>>>>>>>>>>>>>>>>>>>>>>>>>" + strconv.Itoa(i))
//	logger.Fatal("Fatal>>>>>>>>>>>>>>>>>>>>>>>>>" + strconv.Itoa(i))
//}
//
//func main() {
//	runtime.GOMAXPROCS(runtime.NumCPU())
//
//	//指定是否控制台打印，默认为true
//	logger.SetConsole(true)
//	//指定日志文件备份方式为文件大小的方式
//	//第一个参数为日志文件存放目录
//	//第二个参数为日志文件命名
//	//第三个参数为备份文件最大数量
//	//第四个参数为备份文件大小
//	//第五个参数为文件大小的单位
//	//logger.SetRollingFile("d:/logtest", "test.log", 10, 5, logger.KB)
//
//	//指定日志文件备份方式为日期的方式
//	//第一个参数为日志文件存放目录
//	//第二个参数为日志文件命名
//	logger.SetRollingDaily("data/log/apve", "apve.log")
//
//	//指定日志级别  ALL，DEBUG，INFO，WARN，ERROR，FATAL，OFF 级别由低到高
//	//一般习惯是测试阶段为debug，生成环境为info以上
//	logger.SetLevel(logger.ERROR)
//
//	for i := 10; i > 0; i-- {
//		go Log(i)
//		time.Sleep(1000 * time.Millisecond)
//	}
//	time.Sleep(15 * time.Second)
//}
