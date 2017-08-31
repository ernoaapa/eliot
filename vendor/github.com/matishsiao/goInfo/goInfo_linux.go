package goInfo

import (		
	"strings"	
	"os/exec"
	"os"
	"bytes"
	"runtime"
	"fmt"
	"time"
)

func GetInfo() *GoInfoObject {			
	out := _getInfo()
	for strings.Index(out,"broken pipe") != -1 {
		out = _getInfo()
		time.Sleep(500 * time.Millisecond)
	}
	osStr := strings.Replace(out,"\n","",-1)
	osStr = strings.Replace(osStr,"\r\n","",-1)
	osInfo := strings.Split(osStr," ")
	gio := &GoInfoObject{Kernel:osInfo[0],Core:osInfo[1],Platform:osInfo[2],OS:osInfo[3],GoOS:runtime.GOOS,CPUs:runtime.NumCPU()}
	gio.Hostname,_ = os.Hostname()	
	return gio
}

func _getInfo() string {
	cmd := exec.Command("uname","-srio")
	cmd.Stdin = strings.NewReader("some input")
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {							
		fmt.Println("getInfo:",err)	
	}
	return out.String()
}