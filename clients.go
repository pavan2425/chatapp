package chatapp

import (
	"bufio"
	"fmt"
	"io"
)

type client struct {
	*bufio.Reader
	*bufio.Writer
	wc chan string
}

func StartClient(cn io.ReadWriteCloser, msgch chan<- string, quit chan struct{}) (chan<- string, chan struct{}) {
	c := new(client)

	c.Reader = bufio.NewReader(cn)
	c.Writer = bufio.NewWriter(cn)
	c.wc = make(chan string)
	done := make(chan struct{})

	go func() {
		scanner := bufio.NewScanner(c.Reader)
		for scanner.Scan() {
			fmt.Println("banthu")
			fmt.Println(scanner.Text())
			msgch <- scanner.Text()
		}
		done <- struct{}{}
	}()
	c.writeMonitor()

	go func() {
		select {
		case <-quit:
			cn.Close()
		case <-done:
		}
	}()
	fmt.Println("2")
	return c.wc, done
}

func (c *client) writeMonitor() {
	go func() {
		for s := range c.wc {
			fmt.Println("3")
			c.WriteString(s + "\n")
			c.Flush()
		}
	}()
}
