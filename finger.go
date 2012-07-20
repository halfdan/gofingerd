package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"errors"
	"net"
	"os"
	"strings"
	"flag"
)

var port *int = flag.Int("port", 79, "Port to listen on")

func main() {
	flag.Parse()
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		panic(err)
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	user, _, _ := reader.ReadLine()

	info, err := getUserInfo(string(user))
	if err != nil {
		conn.Write([]byte(err.Error()))
	} else {
		conn.Write(info)
	}
}

func getUserInfo(user string) (info []byte, err error) {
	home, err := getHomeDir(user)
	if err != nil {
		return nil, err
	}
	plan := fmt.Sprintf("%s/.plan", home)
	data, err := ioutil.ReadFile(plan)
	if err != nil {
		return data, errors.New("User doesn't have a .plan file!\n")
	}
	return data, nil
}

func getHomeDir(user string) (string, error) {
	passwd, err := os.Open("/etc/passwd")
	if err != nil {
		return "", err
	}
	defer passwd.Close()
	reader := bufio.NewReader(passwd)
	var line []byte
	for err != io.EOF {
		line, _, err = reader.ReadLine()
		data := strings.FieldsFunc(string(line), func(r rune) bool {
			return r == ':'
		})
		if len(data) == 7 {
			if data[0] == user {
				return data[5], nil
			}
		}
	}
	return "", errors.New("User does not exist!\n")
}
