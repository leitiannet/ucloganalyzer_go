package main

import (
	"flag"
	"fmt"

	"ucloganalyzer/uclogic"
)

/*
./ucloganalyzer -requestid="msgserver-10.255.0.68-1494550884.897808788.855" -forcelocal=true
./ucloganalyzer -requestid="msgserver-10.255.0.68-1494550884.897808788.855" -pattern="<requestid:([0-9a-z.-]+)>(<gid:[0-9]+>)?"
./ucloganalyzer -requestid="msgserver-10.255.0.68-1494550884.897808788.855" -pattern="<requestid:([0-9a-z.-]+)>(<gid:[0-9]+>)?" -log="msgserver.log" -logdir="." -remotelogdir="/tmp/"
./ucloganalyzer -requestid="msgserver-10.255.0.188-1495233661.344331033.373" -pattern="<requestid:([0-9a-z.-]+)>(<gid:[0-9]+>)?" -log="msgserver.log.1.gz" -forcelocal=true
./ucloganalyzer -uploadfile="upload.sh" -remotehost="10.255.0.68"
*/

const VERSION = "1.0.0"

var requestId *string = flag.String("requestid", "", "filter requestId from file or dir")
var pattern *string = flag.String("pattern", "", "pattern for filter key")
var configFile *string = flag.String("config", "conf/app.conf", "config file")
var log *string = flag.String("log", "msgserver.log", "log file")
var logDir *string = flag.String("logdir", ".", "log dir")
var forceLocal *bool = flag.Bool("forcelocal", false, "force local")
var filterBySever *bool = flag.Bool("filterbysever", true, "filter filename by sever")
var remoteBin *string = flag.String("remotebin", "ucloganalyzer_simple", "remote bin file")
var remoteBinDir *string = flag.String("remotebindir", "/tmp/", "remote bin dir")
var remoteLog *string = flag.String("remotelog", "msgserver.log", "remote log file")
var remoteLogDir *string = flag.String("remotelogdir", "/var/log/uclog", "remote log dir")
var remoteUser *string = flag.String("remoteuser", "yanfa", "remote user")
var remotePwd *string = flag.String("remotepwd", "yanfa", "remote password")
var remoteHost *string = flag.String("remotehost", "", "remote host")
var remotePort *int = flag.Int("remoteport", 22, "remote port")
var remoteOutput *string = flag.String("remoteoutput", "", "remote output")
var output *string = flag.String("output", "", "output")
var uploadFile *string = flag.String("uploadfile", "", "upload file")
var forceUpload *bool = flag.Bool("forceupload", true, "force upload")

func showParams() {
	fmt.Printf(`list params:
	requestId:%s
	pattern:%s
	configFile:%s
	log:%s
	logDir:%s
	forceLocal:%v
	filterBySever:%v
	remoteBin:%s
	remoteBinDir:%s
	remoteLog:%s
	remoteLogDir:%s
	remoteUser:%s
	remotePwd:%s
	remoteHost:%s
	remotePort:%d
	remoteOutput:%s
	output:%s
	uploadFile:%s
	forceUpload:%v
`, *requestId, *pattern, *configFile, *log, *logDir, *forceLocal, *filterBySever,
		*remoteBin, *remoteBinDir, *remoteLog, *remoteLogDir, *remoteUser, *remotePwd, *remoteHost, *remotePort,
		*remoteOutput, *output, *uploadFile, *forceUpload)
}

func main() {
	flag.Parse()
	showParams()

	var err error
	invoker := uclogic.NewInvoker()
	invoker.Log(*log)
	invoker.LogDir(*logDir)
	invoker.ForceLocal(*forceLocal)
	invoker.FilterBySever(*filterBySever)
	invoker.RemoteLog(*remoteLog)
	invoker.RemoteLogDir(*remoteLogDir)
	invoker.RemoteBin(*remoteBin)
	invoker.RemoteBinDir(*remoteBinDir)
	invoker.RemoteUser(*remoteUser)
	invoker.RemotePwd(*remotePwd)
	invoker.RemoteHost(*remoteHost)
	invoker.RemotePort(*remotePort)
	invoker.RemoteOutput(*remoteOutput)
	invoker.Output(*output)
	if *requestId != "" || *pattern != "" {
		err = invoker.FilterLog(*requestId, *pattern)
	} else if *uploadFile != "" && *remoteHost != "" {
		err = invoker.Upload(*uploadFile, *forceUpload)
	}
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("invoker execute finish")
	}
}
