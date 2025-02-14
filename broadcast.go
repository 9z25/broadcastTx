package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

const accessToken = "123"

const taoNode = "http://192.168.0.100:8000"

//RawTx struct for handling json post data
type RawTx struct {
	Tx string `json:"tx"`
}

// LastTx
type LastTx struct {
	Type      string `json:"type"`
	Addresses string `json:"addresses"`
}

// TaoExplorer
type TaoExplorer struct {
	Address  string   `json:"address"`
	Sent     int      `json:"sent"`
	Received string   `json:"received"`
	Balance  string   `json:"balance"`
	lastTxs  []LastTx `json:"last_txs"`
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

	p := strings.Split(r.URL.Path, "/")

	fmt.Println(p[1])

	response, err := client.Get("https://taoexplorer.com/ext/getaddress/" + p[1])

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

// GetTransaction is for testing
func GetTransaction(w http.ResponseWriter, r *http.Request) {
	client := &http.Client{
		Timeout:   time.Second * 10,
		Transport: MyRoundTripper{r: http.DefaultTransport},
	}

	params := mux.Vars(r)
	txid := params["txid"]
	fmt.Println(txid)

	response, err := client.Get(taoNode + "/api/gettransaction/" + txid)

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

// GetRawTransaction is for testing
func GetRawTransaction(w http.ResponseWriter, r *http.Request) {
	client := &http.Client{
		Timeout:   time.Second * 10,
		Transport: MyRoundTripper{r: http.DefaultTransport},
	}

	params := mux.Vars(r)
	txid := params["txid"]
	fmt.Println(txid)

	response, err := client.Get(taoNode + "/api/getrawtransaction/" + txid)

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
	url := taoNode + "/api/decoderawtransaction/"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("X-CSRF-Token", accessToken)
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
	url := taoNode + "/api/sendrawtransaction/"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("X-CSRF-Token", accessToken)
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

// GetUnspent test
func GetUnspents(w http.ResponseWriter, r *http.Request) {
	var url = "https://taoexplorer.com/ext/getaddress/"
	params := mux.Vars(r)
	addr := params["address"]

	req, _ := http.NewRequest("GET", url+addr, nil)

	client := &http.Client{}

	res, _ := client.Do(req)

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
	}
	json.NewEncoder(w).Encode(body)

}

// GetTxData test
func GetTxData(w http.ResponseWriter, r *http.Request) {
	var url = "https://taoexplorer.com/api/getrawtransaction?txid="
	params := mux.Vars(r)
	txid := params["txid"]

	req, _ := http.NewRequest("GET", url+txid, nil)

	client := &http.Client{}

	res, _ := client.Do(req)

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
	}

	json.NewEncoder(w).Encode(string(body))

}

func main() {

	// Init Router
	r := mux.NewRouter()

	// Route Handlers / Endpoints
	r.HandleFunc("/api/getaddress/{address}", GetAddress).Methods("GET")
	r.HandleFunc("/api/gettransaction/{txid}", GetTransaction).Methods("GET")
	r.HandleFunc("/api/getrawtransaction/{txid}", GetRawTransaction).Methods("GET")
	r.HandleFunc("/api/sendrawtransaction/", SendRawTransaction).Methods("POST")
	r.HandleFunc("/api/decoderawtransaction/", DecodeRawTransaction).Methods("POST")
	r.HandleFunc("/api/getunspents/{address}", GetUnspents).Methods("GET")
	r.HandleFunc("/api/gettxdata/{txid}", GetTxData).Methods("GET")
	handler := cors.Default().Handler(r)

	err := http.ListenAndServeTLS(":8001", "./freshmintrecords_com.crt", "./freshmintrecords.key", handler)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
	//log.Fatal(http.ListenAndServe(":8001", handler))
}
