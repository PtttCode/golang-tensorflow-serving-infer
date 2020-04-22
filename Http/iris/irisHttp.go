package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/kataras/golog"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/hero"
	"github.com/kataras/iris/v12/middleware/logger"
	"github.com/kataras/iris/v12/middleware/recover"
	"io/ioutil"
	"net"
	"net/http"
	"time"
)

const (
	endpoint = "http://127.0.0.1:6006/v1/models/tfs-model:predict"
	padding = "<PAD>"
)

type dataResp struct {
	Outputs	[][]float64	`json:"outputs"`
}
type dataPost struct{
	// Model_name	string	`json:"model_name"`
	// Model_version	int	`json:"model_version"`
	// Data	map[string][][]string	`json:"data"`
	Inputs [][]string	`json:"inputs"`
}
type RequestBody struct{
	Question string
	Length int
	Pid string
}

func maxNum(array []float64) int{
	maxDig := 0.0
	var pos int
	for i := 0; i< len(array); i++{
		if array[i] > maxDig {
			maxDig = array[i]
			pos = i
		}
	}
	// fmt.Println(pos)
	return pos
}

func buildModelInput(sentence string, totalLength int)	[]byte{
	arr := []string{"bos"}
	for _, i := range sentence{
		arr = append(arr, string(i))
	}
	arr = append(arr, "eos")
	length := len(arr)

	for i:=0;i<totalLength-length;i++{
		arr = append(arr, padding)
	}

	res := make(map[string]interface{})
	res["inputs"] = [][]string{arr}

	buf, err := json.Marshal(res)
	if err != nil{
		fmt.Println("转化二进制错误！")
	}
	return buf
}

func RequestModel(bodyJson RequestBody) map[string]interface{}{
	//defer timeCal(r, time.Now())
	//dataByte := `{"inputs": [["bos", "这", "东", "西", "太", "垃", "圾", "了", "eos", "<PAD>", "<PAD>", "<PAD>", "<PAD>", "<PAD>", "<PAD>", "<PAD>", "<PAD>", "<PAD>", "<PAD>", "<PAD>", "<PAD>", "<PAD>", "<PAD>", "<PAD>", "<PAD>", "<PAD>", "<PAD>", "<PAD>", "<PAD>", "<PAD>", "<PAD>", "<PAD>", "<PAD>", "<PAD>", "<PAD>", "<PAD>", "<PAD>", "<PAD>", "<PAD>", "<PAD>", "<PAD>", "<PAD>", "<PAD>", "<PAD>", "<PAD>", "<PAD>", "<PAD>", "<PAD>", "<PAD>", "<PAD>", "<PAD>", "<PAD>", "<PAD>", "<PAD>", "<PAD>", "<PAD>", "<PAD>", "<PAD>", "<PAD>", "<PAD>", "<PAD>", "<PAD>", "<PAD>", "<PAD>", "<PAD>", "<PAD>", "<PAD>", "<PAD>", "<PAD>", "<PAD>", "<PAD>", "<PAD>", "<PAD>", "<PAD>", "<PAD>", "<PAD>", "<PAD>", "<PAD>", "<PAD>", "<PAD>", "<PAD>", "<PAD>", "<PAD>", "<PAD>", "<PAD>", "<PAD>", "<PAD>", "<PAD>", "<PAD>", "<PAD>", "<PAD>", "<PAD>", "<PAD>", "<PAD>", "<PAD>", "<PAD>", "<PAD>", "<PAD>", "<PAD>", "<PAD>"]]}`

	dataByte := buildModelInput(bodyJson.Question, bodyJson.Length)
	var dp dataPost
	var dr dataResp

	json.Unmarshal(dataByte, &dp)
	dataJson, _ := json.Marshal(dp)

	s := time.Now()
	resp, err := http.Post(endpoint, "application/json", bytes.NewBuffer(dataJson))
	golog.Info(time.Now().Sub(s).String())
	if err != nil {
		golog.Error("出错", err)
	}
	res, _ := ioutil.ReadAll(resp.Body)

	json.Unmarshal(res, &dr)
	predArray := dr.Outputs

	pos := maxNum(predArray[0])

	result := map[string]interface{}{
		"probility": predArray[0][pos],
		"intent": pos}


	return result
}

func main() {
	golog.SetTimeFormat("")
	golog.SetLevel("debug")
	app := iris.New()
	app.Logger().SetLevel("debug")
	app.Use(recover.New())
	app.Use(logger.New())
	//goroutines := New(40)

	hero.Register(func (ctx iris.Context)(jsonData RequestBody){
		ctx.ReadJSON(&jsonData)
		return
	})
	inferHandler := hero.Handler(RequestModel)
	app.Post("/infer", inferHandler)

	l, err := net.Listen("tcp", "0.0.0.0:9090")
	if err != nil{
		golog.Error("Listener Error", err)
	}
	defer l.Close()
	//l = netutil.LimitListener(l, 1)

	app.Run(iris.Listener(l), iris.WithoutServerError(iris.ErrServerClosed))

}