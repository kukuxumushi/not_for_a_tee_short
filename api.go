package main

import (
	"bytes"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

type System struct {
	nonce string
	time  string
	sign  string
}

var seededRand = rand.New(rand.NewSource(time.Now().UnixNano()))

var Id int64 = 0

func MakeRandomString(length int, charset string) string {
	str := make([]byte, length)
	for i := range str {
		str[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(str)
}

func GetNonce() string {
	return MakeRandomString(1, "123456789") + MakeRandomString(31, "0123456789")
}

func GetTime() string {
	return strconv.FormatInt(time.Now().UTC().Unix(), 10)
}

func Sign(param *System) {
	param.sign = fmt.Sprintf("%x", md5.Sum([]byte(fmt.Sprintf("time:%s,nonce:%s,appSecret:%s", param.time, param.nonce, AppSecret))))
}

func GetId() string {
	return strconv.FormatInt(Id, 10)
}

func MakeRequestBody(system System, params, id string) string {
	return fmt.Sprintf(`{
	"system": {
		"ver": "1.0",
    	"sign": "%s",
    	"appId": "<REDACTED>",
    	"time": "%s",
    	"nonce": "%s"
  	},
  	"params": {
		%s
  	},
  	"id": "%s"
}`, system.sign, system.time, system.nonce, params, id)
}

func DoRequest(path string, param string) {
	system := System{GetNonce(), GetTime(), ""}
	Sign(&system)
	id := GetId()
	requestBody := MakeRequestBody(system, param, id)
	if Debug {
		log.Println(requestBody)
	}
	Id++

	if DoNotSendAnything {
		return
	}
	url := "https://example.com/openapi" + path
	resp, err := http.Post(url, "application/json", bytes.NewReader([]byte(requestBody)))
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var format bytes.Buffer
	err = json.Indent(&format, body, "", "\t")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(format.String())
}
