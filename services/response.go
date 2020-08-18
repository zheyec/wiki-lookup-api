package services

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"os"
	"strconv"
)

type Response struct {
	Slots    []Slots       `json:"slots"`
	Messages []ResponseMsg `json:"msg_body"`
	ErrorNo  string        `json:"error_no"`
}

type Slots struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type ResponseMsg struct {
	Data map[string]utf8String `json:"data"`
	Type string                `json:"type"`
}

func NewResponse(code string) *Response {
	return &Response{
		ErrorNo: code,
	}
}

// ToBytes convert to bytes
func (x *Response) ToBytes() []byte {
	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)
	encoder.SetEscapeHTML(false)
	encoder.Encode(x)
	return buffer.Bytes()
}

type utf8String string

// custom marshal method
func (s utf8String) MarshalJSON() ([]byte, error) {
	return []byte(strconv.QuoteToASCII(string(s))), nil
}

// format for read & write
type crawlerResult struct {
	Synp string `json:"synp"`
	Link string `json:"link"`
}

// save synopsis and link as local files
func writeResult(path string, synp string, link string) error {
	data := crawlerResult{synp, link}
	file, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.Write(file)
	if err != nil {
		return err
	}
	return nil
}

// read from json files
func readResult(path string) (*crawlerResult, error) {
	result := crawlerResult{}
	raw, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(raw, &result)
	return &result, err
}
