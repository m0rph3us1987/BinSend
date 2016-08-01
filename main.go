package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"time"
)

type appConfig struct {
	PS4ip       string
	PS4port     uint
	PayloadFile string
	LogPort     uint
	DumpPort    uint
}

var config = appConfig{}

func main() {
	fmt.Printf("                              BinSend by m0rph3us1987\n")
	fmt.Printf("This programm can be used to send a file to a server over a tcp connection.\n\n")
	fmt.Printf("Usage: BinSend [optional filename]\n\n")

	file, err := os.Open("config.json")
	if err != nil {
		fmt.Println("Error: Could not open file config.json")
		return
	}
	defer file.Close()

	jsonConfig, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println("Error: could not read file config.json!")
		return
	}
	file.Close()

	err = json.Unmarshal(jsonConfig, &config)
	if err != nil {
		fmt.Println("Error: Invalid json in config.json!")
		fmt.Println(err.Error())
		return
	}

	if config.PS4ip == "" {
		fmt.Println("No ServerIP speficied. Setting Default to 192.168.1.170")
		config.PS4ip = "192.168.1.170"
	}

	if config.PS4port == 0 {
		fmt.Println("No ServerPort speficied. Setting Default to 9023")
		config.PS4port = 9023
	}

	if config.LogPort == 0 {
		fmt.Println("No Log Port speficied. Setting Default to 9023")
		config.LogPort = 30
	}

	if config.DumpPort == 0 {
		fmt.Println("No Dump Port speficied. Setting Default to 9023")
		config.DumpPort = 31
	}

	if len(os.Args) > 1 {
		if os.Args[1] != "" {
			config.PayloadFile = os.Args[1]
			fmt.Printf("Using file %s\n", config.PayloadFile)
		}
	}

	if config.PayloadFile == "" {
		fmt.Println("No Filename speficied. Setting Default to exploit.bin")
		config.PayloadFile = "exploit.bin"
	}

	go startServer(config.LogPort)
	go startBinServer()

	fmt.Println("Sending", config.PayloadFile, "to", config.PS4ip, "port", config.PS4port)

	file, err = os.Open(config.PayloadFile)
	if err != nil {
		fmt.Printf("Error: Could not open file %s\n", config.PayloadFile)
		return
	}
	defer file.Close()

	fileBuffer, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Printf("Error: could not read file %s\n", config.PayloadFile)
		return
	}
	file.Close()

	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", config.PS4ip, config.PS4port))
	if err != nil {
		fmt.Println("Error: Could not connect to server.")
		return
	}

	conn.Write(fileBuffer)
	conn.Close()
	fmt.Printf("Sent %d bytes\n", len(fileBuffer))

	for 1 == 1 {
		time.Sleep(10 * time.Second)
	}
}

func startServer(Port uint) {
	fmt.Printf("Start logging on port %d\n", Port)
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", Port))
	if err != nil {
		fmt.Println(err.Error())
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println(err.Error())
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	//fmt.Printf("Received connection...\n")
	for 1 == 1 {
		bytes := 0
		buff := make([]byte, 4096)
		bytes, _ = conn.Read(buff)

		if bytes > 0 {
			fmt.Printf("%s\n", buff[:bytes])
		}
		time.Sleep(100 * time.Millisecond)
	}
}

func startBinServer() {
	fmt.Println("Waiting for dump on port 31")
	ln, err := net.Listen("tcp", ":31")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	for {
		conn, err := ln.Accept()

		if err != nil {
			fmt.Println(err.Error())
			return
		}
		go handleBinConnection(conn)
	}
}

func handleBinConnection(c net.Conn) {
	fmt.Println("Writing dump.bin...")

	b := make([]byte, 4096)

	f, _ := os.Create("dump.bin")

	for {
		n, _ := c.Read(b)
		if n == 3 {
			f.Close()
			c.Close()
			fmt.Println("File written...")
			return
		} else {
			f.Write(b[:n])
		}
	}
}
