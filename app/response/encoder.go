package response

import (
	"bytes"
	"compress/gzip"
	"log"
)

var suppportedEncoding = map[string]func(string) (string, error){
	"gzip": gzipEncoder,
}

func gzipEncoder(body string) (string, error) {
	var buf bytes.Buffer
	zw := gzip.NewWriter(&buf)
	_, err := zw.Write([]byte(body))
	if err != nil {
		return "", err
	}
	if err := zw.Close(); err != nil {
		log.Fatal(err)
	}

	return buf.String(), nil
}
