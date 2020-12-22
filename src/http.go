package main

import (
	"net/http"
	"strings"
	"time"
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"encoding/json"
	"fmt"
)


type AdvHandler struct {
	h	func(http.ResponseWriter, *http.Request) int
}


func (handler AdvHandler)ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var body []byte
	var bodyStr string
	var err error
	now := time.Now()
	retval := 0

	defer writeAccessLog(r, now, &retval, &bodyStr)
	if strings.EqualFold(r.Method, "POST") {
		b := bytes.NewBuffer(make([]byte, 0))
		reader := io.TeeReader(r.Body, b)
		body, err = ioutil.ReadAll(reader)
		bodyStr = string(body)
		if nil != err {
			log.Printf("read data from body failed, err: %s", err.Error())
			retval = ErrorReadBodyFailed
			writeResponseMsg(w, ErrorReadBodyFailed, nil)
			return
		}
		r.Body = ioutil.NopCloser(b)
	}

	if !cfg.NoAuth {
		if !APIAuth(r, bodyStr) {
			retval = ErrorAuthentication
			writeResponseMsg(w, retval, nil)
			return
		}
	}
	retval = handler.h(w, r)
}

func writeResponseMsg(w http.ResponseWriter, code int, data interface{}) int {
	var msg string
	var apiResp APIResponse

	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	msg, ok := StatusText[code]
	if !ok {
		msg = "Unknown error"
	}

	apiResp.Code = code
	apiResp.Msg = msg
	apiResp.Module = "wd-test"
	if data != nil {
		apiResp.Data = data
	}
	if code == StatusSuccess {
		apiResp.Success = true
	}else {
		apiResp.Success = false
	}
	retStr, err := json.Marshal(apiResp)
	if err != nil {
		fmt.Fprintf(w, "{\"code\":300,\"module\":\"wd-test\",\"message\":\"system error\"}")
		return ErrorMarshalJSON
	}
	w.Write(retStr)
	return code
}

func writeAccessLog(r *http.Request, requestTime time.Time, retval *int, rawbody *string) {
	

	bt := requestTime.UnixNano()
	et := time.Now().UnixNano()
	remote := strings.SplitN(r.RemoteAddr, ":", 2)
	useragent := ""
	
	if _, ok := r.Header["User-Agent"]; ok {
		useragent = r.Header["User-Agent"][0]
	}
	referer := ""
	if _, ok := r.Header["Referer"]; ok {
		referer = r.Header["Referer"][0]
	}

	SourceCode := r.FormValue("sourceCode")
	bid := r.FormValue("bid")
	identity := r.FormValue("identity")
	qtime := r.FormValue("qtime")
	sign := r.FormValue("sign")

	log.Printf("[ACCESS] Remote: %s, RequestTime:%s, Method:%s, URL:%s, Proto:%s, Retval:%d, Time:%d, UserAgent:%s, Referer:%s, SourceCode:%s, Bid:%s, Identity:%s, Qtime:%s, Sign:%s", remote[0], requestTime.Format("2006-01-02 15:04:05"), r.Method, r.URL, r.Proto, *retval, int((et-bt)/(1000*1000)), useragent, referer, SourceCode, bid, identity, qtime, sign)

	

}