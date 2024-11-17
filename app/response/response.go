package response

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/codecrafters-io/http-server-starter-go/app/common"
	"github.com/codecrafters-io/http-server-starter-go/app/request"
)

var pathMap = map[string](map[string]func(*request.Request) (int, map[string]string, string, error)){
	"GET": {
		"echo/*":  processEcho,
		"files/*": getProcessFiles,
	},
	"POST": {
		"files/*": postProcessFiles,
	},
}

var exactPathMap = map[string](map[string]func(*request.Request) (int, map[string]string, string, error)){
	"GET": {
		"/":           processSlash,
		"/user-agent": processUserAgent,
	},
}

var codes = map[int]string{
	200: "200 OK",
	404: "404 Not Found",
	500: "500 Internal Server Error",
	201: "201 Created",
}

func marshalHeaders(headers map[string]string) string {
	retval := ""
	for key, value := range headers {
		retval = retval + fmt.Sprintf("%s: %s\r\n", key, value)
	}
	return retval
}

func ProcessRequest(req *request.Request) ([]byte, error) {
	var headers map[string]string = nil
	var httpCode int = -1
	var body string = ""
	var err error
	function, ok := exactPathMap[req.Method][req.Path]
	if ok {
		httpCode, headers, body, err = function(req)
	}
	for path, function := range pathMap[req.Method] {
		match, _ := regexp.MatchString(path, req.Path)
		if match {
			httpCode, headers, body, err = function(req)
		}
	}
	if httpCode == -1 {
		httpCode, headers, body, err = processNotFound(req)
	}
	encodings, ok := req.Headers["Accept-Encoding"]
	if ok {
		encodingsList := strings.Split(encodings, ", ")
		for _, encoding := range encodingsList {
			encoder, ok := suppportedEncoding[encoding]
			if ok {
				headers["Content-Encoding"] = encoding
				encodedBody, err := encoder(body)
				if err == nil {
					body = encodedBody
					break
				}
				fmt.Println(fmt.Errorf("error encoding body with %s, with error %s", encoding, err.Error()))
			}
		}
	}
	if body != "" {
		headers["Content-Length"] = fmt.Sprintf("%d", len(body))
	}
	resp := fmt.Sprintf("HTTP/1.1 %s\r\n%s\r\n%s", codes[httpCode], marshalHeaders(headers), body)
	return []byte(resp), err
}

func processSlash(req *request.Request) (int, map[string]string, string, error) {
	httpCode := 200
	headers := map[string]string{}
	body := ""
	return httpCode, headers, body, nil
}

func processNotFound(_ *request.Request) (int, map[string]string, string, error) {
	httpCode := 404
	headers := map[string]string{}
	body := ""
	return httpCode, headers, body, nil

}

func processError(_ *request.Request, err error) (int, map[string]string, string, error) {
	httpCode := 500
	headers := map[string]string{
		"Content-Type": "text/plain",
	}
	body := err.Error()
	return httpCode, headers, body, nil
}

func processEcho(req *request.Request) (int, map[string]string, string, error) {
	slicedPath := strings.Split(req.Path, "/")
	argument := slicedPath[len(slicedPath)-1]
	httpCode := 200
	headers := map[string]string{
		"Content-Type": "text/plain",
	}
	body := argument
	return httpCode, headers, body, nil
}

func processUserAgent(req *request.Request) (int, map[string]string, string, error) {
	httpCode := 200
	headers := map[string]string{
		"Content-Type": "text/plain",
	}
	body := req.Headers["User-Agent"]
	return httpCode, headers, body, nil

}

func getProcessFiles(req *request.Request) (int, map[string]string, string, error) {
	slicedURLPath := strings.Split(req.Path, "/")
	fileName := slicedURLPath[len(slicedURLPath)-1]
	filePath := common.FilesFolder + fileName
	byteData, err := os.ReadFile(filePath)
	if err != nil {
		code, header, body, _ := processNotFound(req)
		return code, header, body, err
	}
	strData := string(byteData)
	httpCode := 200
	headers := map[string]string{
		"Content-Type": "application/octet-stream",
	}
	body := strData
	return httpCode, headers, body, nil
}

func postProcessFiles(req *request.Request) (int, map[string]string, string, error) {
	slicedURLPath := strings.Split(req.Path, "/")
	fileName := slicedURLPath[len(slicedURLPath)-1]
	filePath := common.FilesFolder + fileName
	file, err := os.Create(filePath)
	if err != nil {
		code, header, body, _ := processNotFound(req)
		return code, header, body, err
	}
	_, err = file.Write([]byte(req.Body))
	if err != nil {
		code, header, body, _ := processError(req, err)
		return code, header, body, err
	}
	httpCode := 201
	headers := map[string]string{}
	body := ""
	return httpCode, headers, body, nil

}
