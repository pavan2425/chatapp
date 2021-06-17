package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	name := fmt.Sprintf("anonymus %d", rand.Intn(400))
	fmt.Println("starting hrdrachat client")
	fmt.Println("whats your name?")
	fmt.Scanln(&name)

	fmt.Printf("hellos %s, conecting to hydra chat system   ... \n", name)
	conn, err := net.Dial("tcp", "127.0.0.1:2310")
	if err != nil {
		log.Fatal("error connecting to net")
	}
	fmt.Println("connected to  chat app system")
	name += ":"
	defer conn.Close()
	go func() {
		scanner := bufio.NewScanner(conn)
		for scanner.Scan() {
			fmt.Println(scanner.Text())
		}
	}()

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() && err == nil {
		msg := scanner.Text()
		_, err = fmt.Fprintf(conn, name+msg+"\n")

	}
}
