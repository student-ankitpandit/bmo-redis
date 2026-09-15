package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
)

func main() {
	fmt.Println("Listening on port :6379")
	
	l, err := net.Listen("tcp", ":6379")

	if err != nil {
		fmt.Println(err)
		return 
	}

	conn, err := l.Accept()

	if err != nil {
		fmt.Println(err)
		return
	}

	defer conn.Close()

	for {
		buf := make([]byte, 1024)

		_, err := conn.Read(buf)

		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println("error reading from client", err.Error())
			return
		}

		conn.Write([]byte("+OK\r\n"))

		input := "$5\r\nankit\r\n"

		r := bufio.NewReader(strings.NewReader(input))

		b, _ := r.ReadByte() //consuming '$' imp here

		if b != '$' {
			fmt.Println("invalid type, expecting bulk string only")
			os.Exit(1)
		}

		size, _ := r.ReadByte() //"5"

		strSize, _ := strconv.ParseInt(string(size), 10, 64) //5 

		r.ReadByte() // /r
		r.ReadByte() // /n

		name := make([]byte, strSize)

		r.Read(name)

		fmt.Println(string(name))
	}
}
