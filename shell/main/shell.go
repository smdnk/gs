package main

import (
	"bytes"
	"fmt"
	"golang.org/x/crypto/ssh"
	"io"
	"log"
	"time"
)

type SSHTerm struct {
	In       io.Writer
	Out      *bytes.Buffer
	Session  *ssh.Session
	SendLock chan string
	User     User
}

type User struct {
	name string
	pwd  string
	addr string
	port string
}

func InitSSHTerm(name string, pwd string, addr string, port string) (sshTerm SSHTerm) {
	config := &ssh.ClientConfig{
		User:            name,
		Auth:            []ssh.AuthMethod{ssh.Password(pwd)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	client, err := ssh.Dial("tcp", addr+":"+port, config)
	if err != nil {
		fmt.Printf("dial ssh error :%v\n", err)

	}
	session, err := client.NewSession()
	sshTerm.Session = session

	modes := ssh.TerminalModes{
		ssh.ECHO:          0,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}
	if err = session.RequestPty("xterm", 80, 40, modes); err != nil {
		fmt.Printf("get pty error:%v\n", err)
	}
	stdinBuf, err := session.StdinPipe()
	if err != nil {
		log.Printf("get stdin pipe error%v\n", err)
	}
	sshTerm.In = stdinBuf
	sshTerm.Out = bytes.NewBuffer(make([]byte, 0))
	session.Stdout = sshTerm.Out

	err = session.Shell()
	if err != nil {
		fmt.Printf("shell session error%v", err)
	}
	fmt.Println(sshTerm.Out.String())

	sshTerm.SendLock = make(chan string, 1)

	sshTerm.User = User{
		name: name,
		port: port,
		addr: addr,
		pwd:  pwd,
	}
	// 初始化完成后 把链接成功的缓存信息打印出来
	for {
		buf := make([]byte, 8192)
		n, _ := sshTerm.Out.Read(buf)
		if n > 0 {
			fmt.Printf(string(buf[0:n]))
			break
		}
		time.Sleep(time.Millisecond * 200)
	}
	return sshTerm
}

func SendCmd(sshTerm SSHTerm, cmd string) {
	sshTerm.In.Write([]byte(fmt.Sprintf("%v\n", cmd)))
	fmt.Printf("%v\n", cmd)
	go func() {
		terminator := "$"
		if sshTerm.User.name == "root" {
			terminator = "#"
		}
		for {
			buf := make([]byte, 8192)
			n, err := sshTerm.Out.Read(buf)
			if err != nil && err != io.EOF {
				fmt.Printf("read out buffer err:%v", err)
				sshTerm.SendLock <- "send ok"
				break
			}
			if n > 0 {
				split := bytes.Split(buf, []byte(terminator))
				if len(split) < 2 {
					fmt.Printf(string(buf))
					break
				}
				fmt.Printf(string(split[len(split)-2]))
				sshTerm.SendLock <- "send ok"
				break
			}
			time.Sleep(time.Millisecond * 200)
		}
		return // 方法结束后会销毁携程，但是有问题。还是规范点，结束了就return掉
	}()
	<-sshTerm.SendLock

}

func main() {
	sshTerm := InitSSHTerm("smdnk", "199802273612Lw@", "192.168.0.7", "22")

	SendCmd(sshTerm, "cd /home")

	SendCmd(sshTerm, "cd /")

	SendCmd(sshTerm, "ls -l")

	SendCmd(sshTerm, "docker ps -a")

}
