package goreman

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"sync"
	"time"
	"webcron/app/models"
)

var IsEnd bool
var SuccessTime time.Duration

var wg sync.WaitGroup
var errCh chan procChan
type procChan  struct {
	err error
	name string
	bufErr string
}

var procStatusChan chan procStatus
type procStatus struct {
	name string
	status int8
}

// spawnProc starts the specified proc, and returns any error from running it.
func spawnProc(proc string, errCh chan<- procChan, logger *clogger) {
	procObj := procs[proc]

	procObj.mu.Lock()
	fmt.Fprintf(logger, "Starting %s!!!\n", proc)
	//bufOut := new(bytes.Buffer)
	bufErr := new(bytes.Buffer)

	cs := append(cmdStart, procObj.cmdline)
	cmd := exec.Command(cs[0], cs[1:]...)
	cmd.Stdin = nil
	cmd.Stdout = logger
	cmd.Stderr = bufErr
	cmd.SysProcAttr = procAttrs

	if procObj.setPort {
		cmd.Env = append(os.Environ(), fmt.Sprintf("PORT=%d", procObj.port))
		fmt.Fprintf(logger, "Starting %s on port %d\n", proc, procObj.port)
	}
	if err := cmd.Start(); err != nil {
		select {
		case errCh <- procChan{
			err: err,
			name:proc,
			bufErr:bufErr.String(),
		}:
		default:
		}
		fmt.Fprintf(logger, "Failed to start %s: %s\n", proc, err)
		return
	}

	currentPid := cmd.Process.Pid

	//procObj.Cmd = cmd
	currentNum := len(procObj.CmdList)
	procName := getProcName( proc, currentNum)
	procObj.CmdList[procName] = cmd
	procObj.stoppedBySupervisor = false
	procObj.IsStartSuccess = false
	procObj.mu.Unlock()
	procStatusChan<- procStatus {
		name: proc,
		status:PROC_STATUS_RUNNING,
	}
	now := time.Now()
	err := cmd.Wait()

	elapsed := time.Since(now)
	if elapsed > SuccessTime {
		procObj.IsStartSuccess = true
	}

	//procObj.mu.Lock()

	delete(procObj.CmdList, procName)
	if procObj.IsCmdListEmpty() {
		procStatusChan<- procStatus {
			name: proc,
			status:PROC_STATUS_NOT_START,
		}

		procObj.cond.Broadcast()
	}

	if err != nil && procObj.stoppedBySupervisor == false {
		select {
		case errCh <- procChan{
			err: err,
			name:proc,
			bufErr:fmt.Sprintf("%s, 持续时间：%v", bufErr.String(), elapsed),
		}:
		default:
		}
	} else {
		status := "异常"
		if procObj.IsStartSuccess || procObj.RunStatus == PROC_STATUS_START_STOPPED {
		} else {
			errCh <- procChan{
				err: errors.New(fmt.Sprintf("进程主动退出,持续时间：%v，状态【%s】", elapsed, status)),
				name:proc,
				bufErr:"",
			}
		}

	}
	procObj.waitErr = err
	fmt.Fprintf(logger, "任务【%s】-进程【%d】已停止\n", proc, currentPid)
	if procObj.IsCmdListEmpty() {
		logger.loggerDone <- true
	}

}

// Stop the specified proc, issuing os.Kill if it does not terminate within 10
// seconds. If signal is nil, os.Interrupt is used.
func stopProc(proc string, signal os.Signal) error {
	if signal == nil {
		signal = os.Interrupt
	}
	p, ok := procs[proc]
	if !ok || p == nil {
		return errors.New("unknown proc: " + proc)
	}

	p.mu.Lock()
	defer p.mu.Unlock()

	if p.IsCmdListEmpty() {
		return nil
	}
	p.stoppedBySupervisor = true

	p.RunStatus = PROC_STATUS_START_STOPPED

	err := terminateProc(proc, signal)
	if err != nil {
		return err
	}

	timeout := time.AfterFunc(10*time.Second, func() {
		p.mu.Lock()
		defer p.mu.Unlock()
		if p, ok := procs[proc]; ok && p.CmdList != nil {
			for _,cmd := range p.CmdList {
				err = killProc(cmd.Process)
			}

		}
	})
	p.cond.Wait()
	timeout.Stop()

	if err == nil {
		p.RunStatus = PROC_STATUS_STOPPED
		procStatusChan <- procStatus {
			name: proc,
			status:PROC_STATUS_STOPPED,
		}
	}

	return err
}

// start specified proc. if proc is started already, return nil.
func startProc(proc string, wg *sync.WaitGroup, errCh chan<- procChan, isRetry bool) error {
	p, ok := procs[proc]
	if !ok || p == nil {
		return errors.New("unknown proc: " + proc)
	}

	if len(procs[proc].CmdList) >= procs[proc].Num {
		return nil
	}

	//if  p.RunStatus == PROC_STATUS_STOPPED {
	//	return errors.New(fmt.Sprintf("proc[%s]: is stopped status" , proc))
	//}
	//p.mu.Lock()

	if isRetry {
		if  procs[proc].RunStatus == PROC_STATUS_STOPPED {
			//p.mu.Unlock()
			return errors.New(fmt.Sprintf("任务进程【%s】是暂停状态，无需自动启动", proc))
		}

		if procs[proc].FailureTimes < procs[proc].MaxRetryTimes {
			if !procs[proc].IsStartSuccess {
				procs[proc].FailureTimes++
			}
		} else {
			//p.mu.Unlock()
			//todo 邮件通知
			msg := fmt.Sprintf("【任务进程-%s】任务已停止，已经达到最大重启次数:%d", proc, procs[proc].MaxRetryTimes)

			return errors.New(msg)
		}
	}

	procObj := procs[proc]
	logger, error := createLogger(proc, procObj.colorIndex);
	if error != nil {
		//p.mu.Unlock()
		return error
	}

	num := procObj.Num - len(procObj.CmdList)
	for i := num; i > 0 ;i-- {
		if wg != nil {
			wg.Add(1)
		}
		go func() {
			spawnProc(proc, errCh, logger)
			if wg != nil {
				wg.Done()
			}
			//p.mu.Unlock()
		}()
	}

	return nil
}

// restart specified proc.
func restartProc(proc string) error {
	p, ok := procs[proc]
	if !ok || p == nil {
		return errors.New("unknown proc: " + proc)
	}

	stopProc(proc, nil)
	return startProc(proc, &wg, errCh, false)
}

// stopProcs attempts to stop every running process and returns any non-nil
// error, if one exists. stopProcs will wait until all procs have had an
// opportunity to stop.
func stopProcs(sig os.Signal) error {
	var err error
	for proc := range procs {
		stopErr := stopProc(proc, sig)
		if stopErr != nil {
			err = stopErr
		}
	}
	return err
}

// spawn all procs.
func StartProcs(sc <-chan os.Signal, rpcCh <-chan *RpcMessage, exitOnError bool) error {
	errCh = make(chan procChan, 20)
	procStatusChan = make(chan procStatus, 20)

	for name, proc := range procs {
		if proc.RunStatus == PROC_STATUS_RUNNING {
			if error := startProc(name, &wg, errCh,false); error != nil {
				procs[name].RunStatus = PROC_STATUS_NOT_START
				models.SetPtaskRunningStatus(name, procs[name].RunStatus)
				saveExitLog(procChan{
					name: name,
					err: error,
				})
			}
		}
	}
	allProcsDone := make(chan struct{}, 1)
	go func() {
		wg.Wait()
		allProcsDone <- struct{}{}
	}()
	for {
		select {
		case rpcMsg := <-rpcCh:
			//rpc消息通知修改状态
			rpcStatusChange(rpcMsg)
		case procErr := <-errCh:
			//fmt.Println(procErr.name,  procs[procErr.name].MaxRetryTimes, procs[procErr.name].FailureTimes)
			//重启 主动暂停不需要重启
			time.AfterFunc(time.Duration(procs[procErr.name].IntervalTime) * time.Second, func() {
				err := startProc(procErr.name, &wg, errCh, true)
				if err != nil {
					log.Println(err.Error())
				}
			})
			//保存退出日志
			saveExitLog(procErr)


		case <-allProcsDone:
			log.Println("All Procs Done!!")
			//return stopProcs(os.Interrupt)
		case sig := <-sc:
			return stopProcs(sig)
		case procStatus := <- procStatusChan:
			procs[procStatus.name].RunStatus = procStatus.status
			//整个系统退出时，保留运行时状态
			if !IsEnd {
				models.SetPtaskRunningStatus(procStatus.name, procStatus.status)
			}
		}
	}
	return nil
}

func  rpcStatusChange(rpcMsg *RpcMessage) {
	switch rpcMsg.Msg {
	// TODO: add more events here.
	case "stop":
		for _, proc := range rpcMsg.Args {
			if err := stopProc(proc, nil); err != nil {
				rpcMsg.ErrCh <- err
				break
			}
		}
		close(rpcMsg.ErrCh)
	case "stop-all":
		for _, proc := range procs {
			if err := stopProc(proc.name, nil); err != nil {
				rpcMsg.ErrCh <- err
				break
			}
		}
		close(rpcMsg.ErrCh)
	case "start":
		for _, proc := range rpcMsg.Args {
			if err := startProc(proc, nil, errCh, false); err != nil {
				rpcMsg.ErrCh <- err
				break
			}
		}
		close(rpcMsg.ErrCh)
	default:
		panic("unimplemented rpc message type " + rpcMsg.Msg)
	}
}

func saveExitLog(procErr procChan) {
	//写错误日志
	go func() {
		ptask, err := models.GetGoremanByName(procErr.name)
		if err == nil {
			pError := "【BufErr】:" + procErr.bufErr
			if procErr.err != nil{
				pError = pError + ",【Other Error】:" + procErr.err.Error()
			}
			pTaskLog := models.PTaskLog{
				PtaskId : ptask.Id,
				Error : pError,
				CreateTime : time.Now().Unix(),
			}
			_, err := models.PTaskLogAdd(&pTaskLog)
			if err != nil {
				log.Println("insert proc err log error:", err.Error())
			}
		}
	}()
}

func getProcName(name string, num int) string {
	return fmt.Sprintf("%s-%d", name, num)
}
