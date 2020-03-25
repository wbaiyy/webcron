package goreman

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/mattn/go-colorable"
	"io"
	"os"
	"sync"
	"time"
	"webcron-source/app/libs"
)

type clogger struct {
	idx     int
	proc    string
	writes  chan []byte
	done    chan struct{}
	loggerDone    chan bool
	timeout time.Duration // how long to wait before printing partial lines
	buffers buffers       // partial lines awaiting printing
}
var lineLength int

var Colors = []int{
	32, // green
	36, // cyan
	35, // magenta
	33, // yellow
	34, // blue
	31, // red
}
var mutex = new(sync.Mutex)

var colorOut = colorable.NewColorableStdout()
var outs  map[string]io.Writer

func init()  {
	outs = make(map[string]io.Writer)
}

type buffers [][]byte

func (v *buffers) consume(n int64) {
	for len(*v) > 0 {
		ln0 := int64(len((*v)[0]))
		if ln0 > n {
			(*v)[0] = (*v)[0][n:]
			return
		}
		n -= ln0
		*v = (*v)[1:]
	}
}

func (v *buffers) WriteTo(w io.Writer) (n int64, err error) {
	for _, b := range *v {
		nb, err := w.Write(b)
		n += int64(nb)
		if err != nil {
			v.consume(n)
			return n, err
		}
	}
	v.consume(n)
	return n, nil
}

// write any stored buffers, plus the given line, then empty out
// the buffers.
func (l *clogger) writeBuffers(line []byte) {
	var out io.Writer
	if _, ok := outs[l.proc]; !ok {
		CreateLoggerOutput(l.proc)
	}
	out = outs[l.proc]

	//_, ok := out.(*os.File); ok
	mutex.Lock()
	if  procs[l.proc].OutputFile == "" {
		now := time.Now().Format("15:04:05")
		fmt.Fprintf(out, "\x1b[%dm", Colors[l.idx])
		fmt.Fprintf(out, "%s %*s | ", now, maxProcNameLength, l.proc)
		fmt.Fprintf(out, "\x1b[m")
	}
	l.buffers = append(l.buffers, line)
	l.buffers.WriteTo(out)
	l.buffers = l.buffers[0:0]
	mutex.Unlock()
}


// bundle writes into lines, waiting briefly for completion of lines
func (l *clogger) writeLines() {
	defer func() {
		delete(outs, l.proc)
		if v, ok := outs[l.proc].(*os.File); ok  && procs[l.proc].OutputFile != ""{
			v.Close()
		}
	}()
	var tick <-chan time.Time
	for {
		select {
		case w, ok := <-l.writes:
			if !ok {
				if len(l.buffers) > 0 {
					l.writeBuffers([]byte("\n"))
				}
				return
			}
			buf := bytes.NewBuffer(w)
			for {
				line, err := buf.ReadBytes('\n')
				if len(line) > 0 {
					if line[len(line)-1] == '\n' {
						// any text followed by a newline should flush
						// existing buffers. a bare newline should flush
						// existing buffers, but only if there are any.
						if len(line) != 1 || len(l.buffers) > 0 {
							l.writeBuffers(line)
						}
						tick = nil
					} else {
						l.buffers = append(l.buffers, line)
						tick = time.After(l.timeout)
					}
				}
				if err != nil {
					break
				}
			}
			l.done <- struct{}{}
		case <-tick:
			if len(l.buffers) > 0 {
				l.writeBuffers([]byte("\n"))
			}
			tick = nil
		case  <-l.loggerDone:
			return
		}
	}

}

// write handler of logger.
func (l *clogger) Write(p []byte) (int, error) {
	l.writes <- p
	<-l.done
	return len(p), nil
}

// create logger instance.
func createLogger(proc string, colorIndex int) (l *clogger, err error) {
	defer func() {
		if r := recover(); r != nil {
			err =  r.(error)
		}
	}()

	err = nil
	CreateLoggerOutput(proc)
	mutex.Lock()
	defer mutex.Unlock()
	l = &clogger{
		idx: colorIndex,
		proc: proc,
		writes: make(chan []byte),
		done: make(chan struct{}),
		timeout: 2 * time.Millisecond,
		loggerDone: make(chan bool, 1),
	}
	go l.writeLines()
	return l, err
}

func CreateLoggerOutput(proc string) {
	if _, ok := outs[proc]; !ok{
		out := colorOut
		if procs[proc].OutputFile != "" {
			if !libs.Exist(procs[proc].OutputFile) {
				err := libs.FileCreate(procs[proc].OutputFile)
				if err != nil {
					panic(errors.New(fmt.Sprintf("任务【%s】创建文件失败，原因：%v", proc, err)))
				}
			}
			file, err := os.OpenFile(procs[proc].OutputFile, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
			if  err == nil {
				out = file
			} else {
				panic(err)
			}
		}
		outs[proc] = out
	}
}


