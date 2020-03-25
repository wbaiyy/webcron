package main

import (
	"github.com/astaxie/beego"
	"log"
	"os"
	"os/signal"
	"syscall"
	"webcron-source/app/goreman"
	"webcron-source/app/http"
	"webcron-source/app/jobs"
	_ "webcron-source/app/mail"
	"webcron-source/app/models"
)

func main() {
	done := make(chan bool)

	go notifyCh(done)
	defer func() {
		goreman.IsEnd = true
		if err := goreman.RunOpt("stop-all"); err != nil {
			//停止所有进程失败
			log.Fatal(err)
		}

	}()
	models.Init()
	jobs.InitJobs()

	time, error := beego.AppConfig.Int("ptask.start_secs")
	if error != nil {
		time = 60
	}
	goreman.InitGoreman(time)
	http.InitHttp(done)

	<-done
}

func notifyCh(done chan bool) {
	sc := make(chan os.Signal, 10)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGHUP, syscall.SIGTERM, syscall.SIGKILL)

	<-sc
	done <- true
}