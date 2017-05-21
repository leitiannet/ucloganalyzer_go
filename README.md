**目的：**

1. 根据requestid从本地日志文件/日志目录中提取相关请求的日志

2. 自动登录到远程机器上根据requestid提取相关请求的日志  

3. 上传文件到远程机器/从远程主机下载执行结果   

4. 跨平台，支持在Window和Linux下查看服务器日志 

5. 通过正则表达式指定匹配规则 

6. 支持普通文件和.gz压缩文件  

**说明：**
1. requestid格式：	msgserver-10.255.0.68-1494032304.356359546.431和1494032304.356359546.431
2. 日志文件格式：	msgserver-2016-04-20.log、msgserver.log-20170424、msgserver.log、msgserver.log.1和msgserver.log.1.gz
3. 依赖以下包  
    https://github.com/golang/crypto.git  
    https://github.com/pkg/sftp  

**使用：**  
./ucloganalyzer -h
./ucloganalyzer -requestid="msgserver-10.255.0.68-1494550884.897808788.855" -forcelocal=true  
./ucloganalyzer -requestid="msgserver-10.255.0.68-1494550884.897808788.855" -pattern="<requestid:([0-9a-z.-]+)>(<gid:[0-9]+>)?" -log="msgserver.log" -logdir="." -remotelogdir="/tmp/"  
./ucloganalyzer -uploadfile="upload.sh" -remotehost="10.255.0.68"  
