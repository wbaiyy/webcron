package goreman

import (
	"webcron/app/models"
	"os/exec"
	"sync"
)

const (
	PROC_STATUS_NOT_START = 1
	PROC_STATUS_RUNNING = 2
	PROC_STATUS_STOPPED = 3
)

var colorIndex int

// -- process information structure.
type ProcInfo struct {
	name       string
	cmdline    string
	Cmd        *exec.Cmd
	CmdList      map[string]*exec.Cmd
	port       uint
	setPort    bool
	colorIndex int

	// True if we called stopProc to kill the process, in which case an
	// *os.ExitError is not the fault of the subprocess
	stoppedBySupervisor bool

	mu      sync.Mutex
	cond    *sync.Cond
	waitErr error

	//failure times
	FailureTimes int
	//运行状态
	RunStatus int8   //1未开启，2运行中,3暂停中

	MaxRetryTimes int
	IntervalTime  int64
	OutputFile string
	NotifyUser string
	IsStartSuccess bool
	Num int  //启动数量
}


func NewProc(pTaslModel *models.PTask) *ProcInfo {
	this := &ProcInfo {
		name : pTaslModel.Name,
		cmdline : pTaslModel.Command,
		colorIndex : getColorIndex(),
		RunStatus : pTaslModel.RunStatus,
		MaxRetryTimes : pTaslModel.RetryTimes,
		IntervalTime : pTaslModel.IntervalTime,
		OutputFile : pTaslModel.OutputFile,
		NotifyUser : pTaslModel.NotifyUsers,
		Num: pTaslModel.Num,
		CmdList: make(map[string] *exec.Cmd),
	}
	this.cond = sync.NewCond(&this.mu)

	return this
}

func getColorIndex() int{
	colorIndex++
	if colorIndex >= len(Colors) {
		colorIndex = 0
	}
	return colorIndex
}


