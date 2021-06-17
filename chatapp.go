package chatapp

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
)

func Run(connection string) error {
	l, err := net.Listen("tcp", connection)
	if err != nil {
		fmt.Println("error connecting to chat client", err)
		return err
	}
	r := CreateRoom("chatapp")

	go func() {
		ch := make(chan os.Signal)
		signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
		<-ch
		fmt.Println("closing tcp connection")
		l.Close()
		close(r.quit)
		if r.CLCounts() > 0 {
			<-r.msgch
		}
		os.Exit(0)
	}()

	for {
		con, err := l.Accept()
		if err != nil {
			fmt.Println("error getting connection", err)
			break
		}
		go handleConnection(r, con)

	}
	return err
}

func handleConnection(r *room, c net.Conn) {
	fmt.Println("received request from client", c.RemoteAddr())
	r.Addclient(c)
}
