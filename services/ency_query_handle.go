package services

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

var (
	hostport, scheme string
)

// Custom error type 
type encyHandleErr struct {
	Msg string
}

func (err encyHandleErr) Error() string {
	return err.Msg
}

// returnErr - called when an error occurs
func returnErr(e error, w http.ResponseWriter) {
	fmt.Printf("Error: %s\n", e.Error())
	resp := &Response{}
	resp.ErrorNo = e.Error()
	fmt.Printf("[Error]%s \n", resp.ToBytes())
	w.Write(resp.ToBytes())
}

func encyQueryHandle(w http.ResponseWriter, r *http.Request) {
	hostport = r.Host + r.URL.Port()
	scheme = "http://"

	// parse the keyword
	keyword, err := parseKeyword(r)
	if err != nil {
		returnErr(err, w)
		return
	}
	fmt.Printf("[keyword] ==> \"%s\"\n", keyword)

	// crawl and save results
	ec, err := newEncyCrawler(keyword)
	if err != nil {
		returnErr(err, w)
		return
	}

	// response
	resp := &Response{}
	resp.Slots = append(resp.Slots, Slots{"ency_reply", ec.synp})
	resp.Slots = append(resp.Slots, Slots{"ency_cid", ec.cid})
	fmt.Printf("[Response]%s \n", resp.ToBytes())
	w.Write(resp.ToBytes())
}

func parseKeyword(r *http.Request) (string, error) {
	rq := r.URL.RawQuery
	fmt.Println("Raw Query: ", rq)
	if !strings.HasPrefix(rq, "keyword=") {
		return "", encyHandleErr{"Bad URL Format"}
	}
	fields := strings.SplitN(rq, "=", 2)
	keyword, _ := url.QueryUnescape(fields[1])
	if keyword == "" {
		return "", encyHandleErr{"Keyword cannot be empty"}
	}
	return keyword, nil
}
