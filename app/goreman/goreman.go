package goreman

import (
	"context"
	"flag"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/wbaiyy/goreman"
	"github.com/wbaiyy/webcron-source/app/models"
	"os"
	"regexp"
	"strconv"
)

const version = "0.2.1"

// process informations named with proc.
var procs map[string]*ProcInfo

func usage() {
	fmt.Fprint(os.Stderr, `Tasks:
  goreman check                      # Show entries in Procfile
  goreman help [TASK]                # Show this help
  goreman export [FORMAT] [LOCATION] # Export the apps to another process
                                       (upstart)
  goreman run COMMAND [PROCESS...]   # Run a command
                                       start
                                       stop
                                       stop-all
                                       restart
                                       restart-all
                                       list
                                       status
									   update  #update config file	
  goreman start [PROCESS]            # Start the application
  goreman version                    # Display Goreman version

Options:
`)
	flag.PrintDefaults()
	os.Exit(0)
}
var maxProcNameLength = 0

var re = regexp.MustCompile(`\$([a-zA-Z]+[a-zA-Z0-9_]+)`)


func defaultServer(serverPort uint) string {
	if s, ok := os.LookupEnv("GOREMAN_RPC_SERVER"); ok {
		return s
	}
	return fmt.Sprintf("127.0.0.1:%d", DefaultPort())
}

func defaultAddr() string {
	if s, ok := os.LookupEnv("GOREMAN_RPC_ADDR"); ok {
		return s
	}
	return "0.0.0.0"
}

// default port
func DefaultPort() uint {
	s := os.Getenv("GOREMAN_RPC_PORT")
	if s != "" {
		i, err := strconv.Atoi(s)
		if err == nil {
			return uint(i)
		}
	}
	return 8555
}

func GetProcs() map[string]*ProcInfo {
	return procs
}

func GetProcByName(name string) *ProcInfo {
	_, ok := procs[name]
	if ok {
		return  procs[name]
	}
	return nil
}
func SetProcByName(name string, proc *ProcInfo) {
	procs[name] = proc
}


func NotifyCh() <-chan os.Signal {
	sc := make(chan os.Signal, 10)
	//signal.Notify(sc, os.Interrupt)
	return sc
}

func Start() error {
	ctx, cancel := context.WithCancel(context.Background())
	rpcChan := make(chan *RpcMessage, 10)

	serverPort, err  := beego.AppConfig.Int("ptask.server_port")
	if err != nil {
		return err
	}

	exitOnError, err := beego.AppConfig.Bool("ptask.exit_on_error")
	if err != nil {
		return err
	}

	go StartProcs(NotifyCh(), rpcChan, exitOnError)

	serverErr := StartServer(ctx, rpcChan, uint(serverPort))
	cancel() // If procs have returned/errored, cancel the RPC server.
	return serverErr
}

func RunOpt(cmd string, args ...string) error{
	var optCmd string
	switch cmd {
	case "person-start":
		optCmd = "start"
	default:
		optCmd = cmd
	}

	err := RpcRun(optCmd, args, goreman.DefaultPort())
	if err == nil {
		switch cmd {
		case "person-start":
			ResetRetryTime(args)
			//case "stop":
		}
	}

	return  err
}

func ResetRetryTime(procNames []string) {
	for _, procName := range procNames {
		proc := GetProcByName(procName)
		if proc != nil {
			proc.FailureTimes = 0
			SetProcByName(procName, proc)
		}
	}
}

func PullProcs(pTaslModel *models.PTask) {
	if _, ok :=procs[pTaslModel.Name]; !ok {
		procs[pTaslModel.Name] =  NewProc(pTaslModel)

	}
}

func UpdateProc(pTaslModel *models.PTask) {
	if _, ok :=procs[pTaslModel.Name]; ok {
		procs[pTaslModel.Name].name = pTaslModel.Name
		procs[pTaslModel.Name].cmdline = pTaslModel.Command
		procs[pTaslModel.Name].MaxRetryTimes = pTaslModel.RetryTimes
		procs[pTaslModel.Name].IntervalTime = pTaslModel.IntervalTime
		procs[pTaslModel.Name].OutputFile = pTaslModel.OutputFile
		procs[pTaslModel.Name].NotifyUser = pTaslModel.NotifyUsers
	}
}
