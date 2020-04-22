package main

import (
	"bytes"
	"encoding/json"
	"github.com/kataras/golog"
	"io/ioutil"
	"net/http"
	"time"
)

const endpoint = "http://127.0.0.1:6006/v1/models/tfs-model:predict"

type dataResp struct {
	Outputs	[][]float64	`json:"outputs"`
}
type dataPost struct{
	// Model_name	string	`json:"model_name"`
	// Model_version	int	`json:"model_version"`
	// Data	map[string][][]string	`json:"data"`
	Inputs [][]string	`json:"inputs"`
}
type requestBody struct{
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
		arr = append(arr, "<PAD>")
	}

	res := make(map[string]interface{})
	res["inputs"] = [][]string{arr}

	buf, err := json.Marshal(res)
	if err != nil{
		Error.Fatal("转化二进制错误！")
	}
	return buf
}

func requestModel(bodyJson requestBody) map[string]interface{}{
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
