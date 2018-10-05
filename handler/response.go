package handler

import (
	"net/http"
	"net"
	"bufio"
	"errors"
	"encoding/json"
	"bytes"
	"strconv"
	"io"
)

type Map map[string]interface{}
type FillOption map[string]bool

//Response is a wrapper for the http.ResponseWriter
//started set to true if response was written to then don't execute other handler
type Response struct {
	ResponseWriter http.ResponseWriter
	Started bool
	Status  int
	ContentLength int
}

func (r *Response) Reset(rw http.ResponseWriter) {
	r.ResponseWriter = rw
	r.Status = 0
	r.Started = false
}

// Write writes the data to the connection as part of an HTTP reply,
// and sets `started` to true.
// started means the response has sent out.
func (r *Response) Write(p []byte) (int, error) {
	r.Started = true
	return r.ResponseWriter.Write(p)
}

// WriteHeader sends an HTTP response header with status code,
// and sets `started` to true.
func (r *Response) WriteHeader(code int) {
	if r.Status > 0 {
		//prevent multiple response.WriteHeader calls
		return
	}
	r.Status = code
	r.Started = true
	r.ResponseWriter.WriteHeader(code)
}

// Hijack hijacker for http
func (r *Response) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	hj, ok := r.ResponseWriter.(http.Hijacker)
	if !ok {
		return nil, nil, errors.New("webserver doesn't support hijacking")
	}
	return hj.Hijack()
}

// Flush http.Flusher
func (r *Response) Flush() {
	if f, ok := r.ResponseWriter.(http.Flusher); ok {
		f.Flush()
	}
}

// CloseNotify http.CloseNotifier
func (r *Response) CloseNotify() <-chan bool {
	if cn, ok := r.ResponseWriter.(http.CloseNotifier); ok {
		return cn.CloseNotify()
	}
	return nil
}

func (r *Response) Header(key string, value string) {
	r.ResponseWriter.Header().Set(key, value)
}

func (r *Response) Body(content []byte) error {
	//var encoding string
	var buf = &bytes.Buffer{}
	//if output.EnableGzip {
	//	encoding = ParseEncoding(output.Context.Request)
	//}
	
	_, err := buf.Write(content)
	if err != nil {
		return err
	}
	
	r.Header("Content-Length", strconv.Itoa(len(content)))
	r.Header("Server", "MP Server v1.0")
	r.WriteHeader(http.StatusOK)
	r.ContentLength = len(content)
	
	io.Copy(r.ResponseWriter, buf)
	return nil
}

func (r *Response) Error(errCode string, errMsg string) {
	r.ErrorWithCode(500, errCode, errMsg, "")
}

func (r *Response) ErrorWithCode(code int, errCode string, errMsg string, innerErrMsg string) {
	r.Status = code
	r.JSONWithOption(Map{
		"code": code,
		"errCode": errCode,
		"errMsg": errMsg,
		"innerErrMsg": innerErrMsg,
		"data": nil,
	}, true, false)
}


func (r *Response) JSON(data interface{}) {
	r.JSONWithOption(Map{
		"code": http.StatusOK,
		"errCode": "",
		"errMsg": "",
		"innerErrMsg": "",
		"data": data,
	}, false, false)
}

func (r *Response) JSONWithOption(data interface{}, hasIndent bool, coding bool) error {
	r.Header("Content-Type", "application/json; charset=utf-8")
	var content []byte
	var err error
	if hasIndent {
		content, err = json.MarshalIndent(data, "", "  ")
	} else {
		content, err = json.Marshal(data)
	}
	if err != nil {
		http.Error(r.ResponseWriter, err.Error(), http.StatusInternalServerError)
		return err
	}
	if coding {
		content = []byte(stringsToJSON(string(content)))
	}
	return r.Body(content)
}

func stringsToJSON(str string) string {
	var jsons bytes.Buffer
	for _, r := range str {
		rint := int(r)
		if rint < 128 {
			jsons.WriteRune(r)
		} else {
			jsons.WriteString("\\u")
			jsons.WriteString(strconv.FormatInt(int64(rint), 16))
		}
	}
	return jsons.String()
}