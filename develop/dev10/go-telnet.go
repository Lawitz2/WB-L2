package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

/*
Реализовать простейший telnet-клиент.

Примеры вызовов:
go-telnet --timeout=10s host port go-telnet mysite.ru 8080 go-telnet --timeout=3s 1.1.1.1 123


Требования:
Программа должна подключаться к указанному хосту (ip или доменное имя + порт) по протоколу TCP. После подключения STDIN программы должен записываться в сокет, а данные полученные и сокета должны выводиться в STDOUT
Опционально в программу можно передать таймаут на подключение к серверу (через аргумент --timeout, по умолчанию 10s)
При нажатии Ctrl+D программа должна закрывать сокет и завершаться. Если сокет закрывается со стороны сервера, программа должна также завершаться. При подключении к несуществующему сервер, программа должна завершаться через timeout

*/

type client struct {
	destination string
	timeout     time.Duration
	connection  net.Conn
	in          io.Reader
	out         io.Writer
}

func newClient(dest string, timeout time.Duration, in io.Reader, out io.Writer) *client {
	c := &client{
		destination: dest,
		timeout:     timeout,
		in:          in,
		out:         out,
	}
	return c
}

func (cl *client) connect() error {
	var err error
	cl.connection, err = net.DialTimeout("tcp", cl.destination, cl.timeout)
	if err != nil {
		return err
	}
	defer cl.connection.Close()

	fmt.Printf("Connected to %s with timeout %.0f seconds\n", cl.destination, cl.timeout.Seconds())
	fmt.Println("Press Ctrl-D / Ctrl-C to exit")

	errChan := make(chan error, 1)
	receivedChan := make(chan string)
	sendThisOutChan := make(chan string)

	go func() {
		signals := make(chan os.Signal, 1)
		signal.Notify(signals, syscall.SIGTERM, syscall.SIGINT, os.Interrupt)
		<-signals
		errChan <- errors.New("system interrupt")
	}()

	go msgRelay(cl.in, sendThisOutChan, errChan)
	go msgRelay(cl.connection, receivedChan, errChan)

	for {
		select {
		case input := <-receivedChan:
			fmt.Fprintln(cl.out, input)
		case msg := <-sendThisOutChan:
			_, err = cl.connection.Write([]byte(msg))
			if err != nil {
				errChan <- err
			}
		case err = <-errChan:
			return err
		}
	}
}

func msgRelay(in io.Reader, destinationChan chan string, errChan chan error) {
	reader := bufio.NewReader(in)
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			errChan <- err
		}
		destinationChan <- message
	}
}

func main() {
	timeoutStr := flag.String("timeout", "20s", "connection timeout")
	flag.Parse()

	if len(flag.Args()) < 2 {
		fmt.Println("please enter host and port")
		return
	}

	timeout, err := time.ParseDuration(*timeoutStr)
	if err != nil {
		slog.Error(err.Error())
		return
	}

	cl := newClient(net.JoinHostPort(flag.Arg(0), flag.Arg(1)), timeout, os.Stdin, os.Stdout)

	if err = cl.connect(); err != nil {
		fmt.Println("\n\nConnection closed:", err)
	}
}
