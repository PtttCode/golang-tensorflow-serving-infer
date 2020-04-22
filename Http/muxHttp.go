package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"time"
)



var (
	Info	*log.Logger
	Warning	*log.Logger
	Error	*log.Logger
)

func timeCal(r *http.Request, start time.Time){
	begin := time.Since(start)
	Info.Println(strings.Join([]string{r.Host, r.URL.Path, " ", begin.String(), " ", r.Method}, ""))
}

func logInit() {
	InfoFile, err := os.OpenFile("log/info.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	WarningFile, err := os.OpenFile("log/warning.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	ErrorFile, err := os.OpenFile("log/error.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalln(err)
	}

	Info = log.New(io.MultiWriter(os.Stdout, InfoFile), "Info: ", log.Ldate|log.Ltime|log.Lshortfile)
	Warning = log.New(io.MultiWriter(os.Stdout, WarningFile), "Warning: ", log.Ldate|log.Ltime|log.Lshortfile)
	Error = log.New(io.MultiWriter(os.Stderr, ErrorFile), "Error: ", log.Ldate|log.Ltime|log.Lshortfile)

}

func inferHandler(w http.ResponseWriter, r *http.Request){
	var bodyJson requestBody
	body, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal([]byte(string(body)), &bodyJson)

	result, _ := json.Marshal(bodyJson)

	w.Write(result)
}

func main() {
	logInit()
	router := mux.NewRouter()
	router.HandleFunc("/infer", inferHandler).Methods("POST")
	//router.Use(app.JwtAuthentication)

	fmt.Println("Start Server!~!~")
	loggedRouter := handlers.LoggingHandler(os.Stdout, router)

	l, err := net.Listen("tcp", "0.0.0.0:9090")
	if err != nil {
		Error.Fatal("Listen: ", err)
	}
	defer l.Close()
	//l = netutil.LimitListener(l, 20)
	http.Serve(l, loggedRouter)
	//err := http.ListenAndServe(":9090", loggedRouter) //设置监听的端口


}
