package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

const accessToken = "123"

//RawTx struct for handling json post data
type RawTx struct {
	Tx string `json:"tx"`
}

type MyRoundTripper struct {
	r http.RoundTripper
}

func (mrt MyRoundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	r.Header.Add("X-CSRF-Token", accessToken)
	return mrt.r.RoundTrip(r)
}

// GetAddress is for testing
func GetAddress(w http.ResponseWriter, r *http.Request) {
	client := &http.Client{
		Timeout:   time.Second * 10,
		Transport: MyRoundTripper{r: http.DefaultTransport},
	}

	response, err := client.Get("https://taoexplorer.com/ext/getaddress/TmfMKgMG6B8cXr2sewa91VNNrqw6Q9WFMM")

	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(responseData))
	json.NewEncoder(w).Encode(responseData)
}

// DecodeRawTransaction returns a json object of decoded tx
func DecodeRawTransaction(w http.ResponseWriter, r *http.Request) {

	d, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}

	var hash RawTx

	if err := json.Unmarshal(d, &hash); err != nil {
		panic(err)
	}

	var jsonStr = []byte(d)
	url := "http://192.168.0.104:8000/api/decoderawtransaction/"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("X-CSRF-Token", "123")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	json.NewEncoder(w).Encode(string(body))
}

// DecodeRawTransaction returns a json object of decoded tx
func SendRawTransaction(w http.ResponseWriter, r *http.Request) {

	d, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}

	var hash RawTx

	if err := json.Unmarshal(d, &hash); err != nil {
		panic(err)
	}

	var jsonStr = []byte(d)
	url := "http://192.168.0.104:8000/api/sendrawtransaction/"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("X-CSRF-Token", "123")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	json.NewEncoder(w).Encode(string(body))
}

func main() {

	// Init Router
	r := mux.NewRouter()

	// Route Handlers / Endpoints
	r.HandleFunc("/api/getaddress/", GetAddress).Methods("GET")
	r.HandleFunc("/api/sendrawtransaction/", SendRawTransaction).Methods("POST")
	r.HandleFunc("/api/decoderawtransaction/", DecodeRawTransaction).Methods("POST")
	handler := cors.Default().Handler(r)

	//err := http.ListenAndServeTLS(":8001", "./freshmintrecords_com.crt", "./freshmintrecords.key", handler)
	//if err != nil {
	//	log.Fatal("ListenAndServe:", err)
	//}
	log.Fatal(http.ListenAndServe(":8001", handler))
}
