package main

import (
	"flag"
	"fmt"
	"net"
	"os"

	"github.com/codecrafters-io/http-server-starter-go/app/common"
	"github.com/codecrafters-io/http-server-starter-go/app/request"
	"github.com/codecrafters-io/http-server-starter-go/app/response"
)

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")
	flag.StringVar(&common.FilesFolder, "directory", ".", "Specify the directory to process")
	flag.Parse()
	fmt.Println("fileFolder", common.FilesFolder)
	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}
		go func() {
			req := make([]byte, 1024)
			conn.Read(req)
			reqObj := request.ParseRequest(string(req))
			resp, err := response.ProcessRequest(reqObj)
			if err != nil {
				fmt.Errorf(err.Error())
			}
			conn.Write(resp)
			conn.Close()
		}()
	}

}
