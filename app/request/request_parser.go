package request

import (
	"strings"
)

type Request struct {
	Method      string
	Path        string
	HTTPVersion string
	Headers     map[string]string
	Body        string
}

func ParseRequest(request_str string) *Request {
	reqObj := Request{}
	reqObj.Headers = map[string]string{}
	data := strings.Split(request_str, "\r\n")
	baseData := strings.Split(data[0], " ")
	reqObj.Method, reqObj.Path, reqObj.HTTPVersion = baseData[0], baseData[1], baseData[2]
	var index int
	var header string
	for index, header = range data[1 : len(data)-1] {
		if strings.TrimSpace(header) == "" {
			break
		}
		headerSplit := strings.Split(header, ":")
		reqObj.Headers[headerSplit[0]] = strings.TrimSpace(headerSplit[1])
	}
	reqObj.Body = ""
	if index+2 < len(data) {
		for _, char := range data[index+2] {
			if char == 0 {
				break
			}
			reqObj.Body = reqObj.Body + string(char)
		}
	}

	return &reqObj
}
