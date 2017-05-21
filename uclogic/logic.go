package uclogic

import (
	"fmt"
	"path"
	"path/filepath"
	"strconv"
	"strings"

	"ucloganalyzer/ucfile"
	"ucloganalyzer/ucreg"
	"ucloganalyzer/ucssh"
	"ucloganalyzer/ucutils"
)

const (
	DEFAULT_LOG_DIR = "/var/log/uclog/"
)

func ParseRequest(requestId string) (server string, ip string, timestamp int64) {
	var timestamprand string
	if strings.Contains(requestId, "-") {
		s := strings.Split(requestId, "-")
		if len(s) > 0 {
			server = s[0]
			ip = s[1]
			timestamprand = s[2]
		}
	} else {
		timestamprand = requestId
	}
	if strings.Contains(timestamprand, ".") {
		ts := strings.Split(timestamprand, ".")
		if len(ts) > 0 {
			timestamp, _ = strconv.ParseInt(ts[0], 10, 0)
		}
	}
	return
}

type Invoker struct {
	requestId     string
	serverName    string
	timestamp     int64
	pattern       string
	log           string
	logDir        string
	filterBySever bool
	forceLocal    bool
	remoteBin     string
	remoteBinDir  string
	remoteLog     string
	remoteLogDir  string
	remoteUser    string
	remotePwd     string
	remoteHost    string
	remotePort    int
	output        string
	remoteOutput  string
	parser        *ucreg.Parser
	fileParser    *ucreg.Parser
	requestParser *ucreg.Parser
	keyMap        map[string]string
	offset        int64
	startLine     int64
	endLine       int64
}

func NewInvoker() *Invoker {
	invoker := &Invoker{}
	return invoker
}

func (this *Invoker) FilterBySever(flag bool) {
	this.filterBySever = flag
}

func (this *Invoker) ForceLocal(flag bool) {
	this.forceLocal = flag
}

func (this *Invoker) Log(file string) {
	this.log = file
}

func (this *Invoker) LogDir(dir string) {
	this.logDir = dir
}

func (this *Invoker) RemoteBin(bin string) {
	this.remoteBin = bin
}

func (this *Invoker) RemoteBinDir(dir string) {
	this.remoteBinDir = dir
}

func (this *Invoker) RemoteLog(file string) {
	this.remoteLog = file
}

func (this *Invoker) RemoteLogDir(dir string) {
	this.remoteLogDir = dir
}

func (this *Invoker) RemoteUser(user string) {
	this.remoteUser = user
}

func (this *Invoker) RemotePwd(pwd string) {
	this.remotePwd = pwd
}

func (this *Invoker) RemoteHost(host string) {
	this.remoteHost = host
}

func (this *Invoker) RemotePort(port int) {
	this.remotePort = port
}

func (this *Invoker) RemoteOutput(out string) {
	this.remoteOutput = out
}

func (this *Invoker) Output(out string) {
	this.output = out
}

func (this *Invoker) AddKey(keys ...string) {
	for _, key := range keys {
		if key == "" {
			continue
		}
		if _, ok := this.keyMap[key]; !ok {
			this.keyMap[key] = "1"
		}
	}
}

func (this *Invoker) printFile(lineNum, offset int64, line string) {
	fmt.Println(line)
}

// 1、遍历文件根据lineNum选择输出 2、定位到offset，输出endLine-startLine+1行
func (this *Invoker) filterFile(lineNum, offset int64, line string) {
	if lineNum >= this.startLine && lineNum <= this.endLine {
		fmt.Printf("%d\t%d\t%s\n", lineNum, offset, line)
	}
}

func (this *Invoker) parserFile(lineNum, offset int64, line string) {
	if this.parser != nil {
		result := this.parser.Find(line)
		if len(result) > 0 {
			if len(this.keyMap) > 0 {
				for key, _ := range this.keyMap {
					if ucutils.StrInSlice(key, result) {
						this.offset = offset
						if this.startLine < 0 {
							this.startLine = lineNum
						}
						this.endLine = lineNum
						this.AddKey(result...)
						break
					}
				}
			} else {
				fmt.Printf("%d\t%d\t%s\n", lineNum, offset, line)
			}
		}
	}
}

func (this *Invoker) parserDir(file string) {
	this.HandleLocalFile(file)
}

//日志文件格式：*server-2016-04-20.log、*server.log-20170424、*server.log、*server.log.1和*server.log.1.gz
func (this *Invoker) ValideFile(file string) bool {
	if this.fileParser == nil {
		pattern := fmt.Sprintf("^[a-z]+server[-0123456789]*\\.log(\\.)?[-0123456789]*(\\.%s)?$", ucfile.DEFAULT_COMPRESS_TYPE)
		this.fileParser = ucreg.NewParser(pattern)
	}
	if this.fileParser != nil {
		return this.fileParser.Match(file)
	} else {
		return false
	}
}

//requestid格式：*server-10.255.0.68-1494032304.356359546.431和1494032304.356359546.431
func (this *Invoker) ValideRequest(requestId string) bool {
	if this.requestParser == nil {
		pattern := "^([a-z]+server-[0-9]+\\.[0-9]+\\.[0-9]+\\.[0-9]+-)?[0-9]+\\.[0-9]+\\.[0-9]+$"
		this.requestParser = ucreg.NewParser(pattern)
	}
	if this.requestParser != nil {
		return this.requestParser.Match(requestId)
	} else {
		return false
	}
}

func (this *Invoker) FilterLog(requestId, pattern string) error {
	if !this.ValideRequest(requestId) {
		return fmt.Errorf("%s not support\n", requestId)
	}

	var err error
	this.requestId = requestId
	this.pattern = pattern
	this.keyMap = make(map[string]string, 0)
	this.startLine = -1
	this.endLine = -1
	this.offset = 0
	this.AddKey(requestId)
	if pattern != "" {
		this.parser = ucreg.NewParser(pattern)
	} else {
		this.parser = ucreg.NewParser(requestId)
	}

	localIps := ucutils.GetLocalIp()
	serverName, servrIp, timestamp := ParseRequest(this.requestId)
	fmt.Println(serverName, servrIp, timestamp, ucutils.FormatTimestamp(timestamp), localIps)
	this.timestamp = timestamp
	this.serverName = serverName
	if this.remoteHost == "" {
		this.remoteHost = servrIp
	}

	if this.forceLocal || this.remoteHost == "" || ucutils.StrInSlice(this.remoteHost, localIps) {
		err = this.HandleLocal()
	} else {
		err = this.HandleRemote()
	}
	return err
}

func (this *Invoker) HandleLocal() error {
	var err error
	var logDir string = this.logDir
	if logDir == "" {
		logDir = DEFAULT_LOG_DIR
	}
	if this.log == "" {
		if ucfile.IsDir(logDir) {
			err = this.HandleLocalDir(logDir)
		}
	} else {
		logFile := filepath.Join(logDir, this.log)
		if ucfile.IsFile(logFile) {
			err = this.HandleLocalFile(logFile)
		}
	}
	return err
}

func (this *Invoker) HandleLocalFile(file string) error {
	if !this.ValideFile(path.Base(file)) {
		return fmt.Errorf("%s not support\n", file)
	}
	var err error
	err = ucfile.ScanFile(file, this.parserFile)
	if err != nil {
		return err
	}
	if this.requestId != "" {
		fmt.Println(this.keyMap, this.offset, this.startLine, this.endLine)
		err = ucfile.ScanFile(file, this.filterFile)
	}
	return err
}

func (this *Invoker) HandleLocalDir(dir string) error {
	var prefix string
	if this.filterBySever {
		prefix = this.serverName
	}
	return ucfile.ScanDir(dir, prefix, this.parserDir)
}

func (this *Invoker) HandleRemote() error {
	var err error
	var result string
	sshClient := ucssh.NewSSHClient()
	succ := sshClient.Dial(this.remoteHost, this.remotePort, this.remoteUser, this.remotePwd)
	if !succ {
		return fmt.Errorf("dial %s:%d fail", this.remoteHost, this.remotePort)
	}
	defer sshClient.Close()

	excuteCmd := fmt.Sprintf("cd %s;chmod a+x %s;./%s -requestid=\"%s\" -log=\"%s\" -logdir=\"%s\" -forcelocal=true",
		this.remoteBinDir, this.remoteBin, this.remoteBin, this.requestId, this.remoteLog, this.remoteLogDir)

	outputName := fmt.Sprintf("%s.log", this.requestId)
	remoteOutput := filepath.Join(this.remoteBinDir, outputName)
	if this.remoteOutput == "file" {
		excuteCmd = fmt.Sprintf("%s > %s", excuteCmd, remoteOutput)
	}
	fmt.Println(excuteCmd)
	result, err = sshClient.Run(excuteCmd)
	if err != nil {
		return err
	}
	fmt.Println(result)

	if this.remoteOutput == "file" {
		err = sshClient.Download(remoteOutput, ".")
		if err != nil {
			return err
		}
		fmt.Println("==========================================list result==========================================")
		err = ucfile.ScanFile(outputName, this.printFile)
	}
	return err
}

// 上传大文件时比较慢，推荐使用rz命令
func (this *Invoker) Upload(uploadFile string, force bool) error {
	sshClient := ucssh.NewSSHClient()
	succ := sshClient.Dial(this.remoteHost, this.remotePort, this.remoteUser, this.remotePwd)
	if !succ {
		return fmt.Errorf("dial %s:%d fail", this.remoteHost, this.remotePort)
	}
	defer sshClient.Close()

	if force && len(uploadFile) > 0 {
		remoteFile := filepath.Join(this.remoteBinDir, path.Base(uploadFile))
		// 无法使用[ ] && rm -rf?
		removeCmd := fmt.Sprintf("filename=%s;if [ -e \"$filename\" ]; then echo remove $filename;rm -rf $filename;fi", remoteFile)
		fmt.Println(removeCmd)

		result, err := sshClient.Run(removeCmd)
		if err != nil {
			return err
		}
		fmt.Println(result)

		err = sshClient.Upload(uploadFile, this.remoteBinDir)
		if err != nil {
			return err
		}
	}
	return nil
}
