package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

const (
	ORGNAME     = "org1"
	CHANNELID   = "mychannel"
	USERNAME    = "Admin"
	CHAINCODEID = "zcc"
	PUT_FCN     = "put"
	GET_FCN     = "get"
	CONFIG_PATH = "/opt/gopath/src/github.com/learnergo/fabric-performance-test/fabric-client-go/fixtures/config.yaml"
	TARGETPEER  = "peer0.org1.example.com"
)

var client = Client{}

func GetHandler(w http.ResponseWriter, r *http.Request) {

	args := [][]byte{[]byte("a")}

	_, err := client.Query(InterRequest{
		CHAINCODEID,
		GET_FCN,
		args,
		nil,
		[]string{TARGETPEER},
		nil,
	})
	if err != nil {
		fmt.Println("Failed to query")
		http.Error(w, err.Error(), 500)
		return

	}
}

func PutHandler(w http.ResponseWriter, r *http.Request) {

	kv := []byte(strconv.FormatInt(time.Now().UnixNano(), 10) + strconv.Itoa(rand.Int()))

	args := [][]byte{kv, kv}
	_, err := client.Invoke(InterRequest{
		CHAINCODEID,
		PUT_FCN,
		args,
		nil,
		[]string{TARGETPEER},
		nil,
	})
	if err != nil {
		fmt.Println("Failed to invoke")
		http.Error(w, err.Error(), 500)
		return

	}
}

func HelloWorldHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("hello world!!")
}

func main() {
	var err error
	client, err = NewClient(CreateChannelClient{
		CONFIG_PATH,
		CHANNELID,
		USERNAME,
		ORGNAME,
	})
	if err != nil {
		fmt.Println("Failed to create new SDK: %s", err)
		return
	}

	http.HandleFunc("/gettest", GetHandler)
	http.HandleFunc("/puttest", PutHandler)
	http.HandleFunc("/helloworld", HelloWorldHandler)
	err = http.ListenAndServe("127.0.0.1:8080", nil)
	if err != nil {
		fmt.Println("ListenAndServe:", err)
	}

}
