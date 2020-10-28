package goreman

import (
	"fmt"
	"time"
	"webcron/app/models"
)
var retryChans map[string]chan bool

func InitGoreman(successTime int)  {
	SetSuccessTime(successTime)
	initProcs()
	go func() {
		err := Start()
		if err != nil {
			fmt.Println("start ptask error:", err)
		}
	}()
}

func initProcs() {
	pTasks, _ := models.GetAllGoreman()
	procs = make(map[string]*ProcInfo)
	for _, pTask := range  pTasks {
		if pTask.Status != 1 {
			continue
		}
		procs[pTask.Name] = NewProc(pTask)
	}
}

func SetSuccessTime(successTime int) {
	SuccessTime = time.Duration(successTime) * time.Second
}
