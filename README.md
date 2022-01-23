## 使用
```go
// windows机子

// 编译windows的cmd窗口设置
set GOOS=windows

// 编译linux的cmd窗口设置
set GOOS=linux

// 编译Mac的cmd窗口设置
set GOOS=darwin
```
### 编译
跨平台编译可参考上面的
```go
git clone https://github.com/Rxg1898/rcak.git
cd ./server
go build server.go

cd ./client
go build client.go


// 没有弹窗的exe编译命令
go build -ldflags="-H windowsgui -w -s" client.go
```
### 启动服务端
```go
// 二进制文件 linux启动 默认监听0.0.0.0地址 20221端口
./server  
```
### 启动客户端
```go
// windows
.\client.exe -h                     // 查看帮助
.\client.exe                        // 提前设置好默认值，可以不填参数
.\client.exe -H 127.0.0.1 -p 20221  // 指定IP+端口

// linux
./client                            // 使用默认值
./client -H 127.0.0.1 -p 20221
```
### 简单使用
#### 客户端上线成功
![](https://blog-1301758797.cos.ap-guangzhou.myqcloud.com/%E6%96%87%E6%A1%A3%E5%9B%BE%E7%89%87/%E8%BF%9C%E6%8E%A7/1642949301.png#crop=0&crop=0&crop=1&crop=1&id=SX6gW&originHeight=91&originWidth=473&originalType=binary&ratio=1&rotation=0&showTitle=false&status=done&style=none&title=)
#### 查看有什么命令
![](https://blog-1301758797.cos.ap-guangzhou.myqcloud.com/%E6%96%87%E6%A1%A3%E5%9B%BE%E7%89%87/%E8%BF%9C%E6%8E%A7/1642949481(1).jpg#crop=0&crop=0&crop=1&crop=1&id=ptNr3&originHeight=179&originWidth=534&originalType=binary&ratio=1&rotation=0&showTitle=false&status=done&style=none&title=)
#### 进入客户端
![](https://blog-1301758797.cos.ap-guangzhou.myqcloud.com/%E6%96%87%E6%A1%A3%E5%9B%BE%E7%89%87/%E8%BF%9C%E6%8E%A7/1642949704.jpg#crop=0&crop=0&crop=1&crop=1&id=QpY2h&originHeight=385&originWidth=853&originalType=binary&ratio=1&rotation=0&showTitle=false&status=done&style=none&title=)
#### 挂起客户端，进入另一个客户端
![](https://blog-1301758797.cos.ap-guangzhou.myqcloud.com/%E6%96%87%E6%A1%A3%E5%9B%BE%E7%89%87/%E8%BF%9C%E6%8E%A7/1642949878.jpg#crop=0&crop=0&crop=1&crop=1&id=wHdyC&originHeight=390&originWidth=1007&originalType=binary&ratio=1&rotation=0&showTitle=false&status=done&style=none&title=)
#### 退出客户端，退出程序
因为有重连，所以一推出就新建立连接了
![](https://blog-1301758797.cos.ap-guangzhou.myqcloud.com/%E6%96%87%E6%A1%A3%E5%9B%BE%E7%89%87/%E8%BF%9C%E6%8E%A7/1642950327.jpg#crop=0&crop=0&crop=1&crop=1&id=kIJ2Q&originHeight=327&originWidth=981&originalType=binary&ratio=1&rotation=0&showTitle=false&status=done&style=none&title=)

