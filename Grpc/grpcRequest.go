package main

import (
	"golang.org/x/net/context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	pb "github.com/tensorflow/serving/tensorflow_serving/apis"
	tf_core_framework "github.com/tensorflow/tensorflow/tensorflow/go/core/framework"
	"google.golang.org/grpc"
	"net"
	"net/http"
	"runtime"
)

type requestBody struct{
	Question string
	Length int
	Pid string
}

const (
	modelName = "tfs-model"
	REQUESTURL = "127.0.0.1:6007"
)

type responseBody struct {
	Data	[]map[string]interface{}
	Code 	int
}

var rpcClients = make(map[string]pb.PredictionServiceClient)

func HasItemInClients(dic map[string]pb.PredictionServiceClient, item string) bool{
	for k := range dic{
		if k == item{
			return true
		}
	}
	return false
}

func getModelInput(sentence string, totalLength int)	[]byte{
	arr := []string{"bos"}
	for _, i := range sentence{
		arr = append(arr, string(i))
	}
	arr = append(arr, "eos")
	length := len(arr)

	for i:=0;i<totalLength-length;i++{
		arr = append(arr, "<PAD>")
	}


	buf, err := json.Marshal(arr)
	if err != nil{
		fmt.Println("转化二进制错误！")
	}
	return buf
}

func findMaxNum(floatArray []float32) int{
	var max float32 = 0
	var pos int
	for idx, i := range floatArray{
		if i > max{
			max = i
			pos = idx
		}
	}
	return pos
}

func grpcRequestModel(w http.ResponseWriter, r *http.Request){
	var bodyJson requestBody
	var res responseBody
	var client pb.PredictionServiceClient
	fmt.Println(runtime.NumGoroutine())
	body, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(body, &bodyJson)

	x := getModelInput(bodyJson.Question, bodyJson.Length)

	request := &pb.PredictRequest{
		ModelSpec: &pb.ModelSpec{
			Name:	modelName,
			// SignatureName: "add_fn",
		},
		Inputs: map[string]*tf_core_framework.TensorProto{
			"inputs": {
				Dtype:    tf_core_framework.DataType_DT_STRING,
				TensorShape:	&tf_core_framework.TensorShapeProto{
					Dim:	[]*tf_core_framework.TensorShapeProto_Dim{
						{
							Size: int64(len(bodyJson.Question)),
						},
						{
							Size: int64(bodyJson.Length),
						},
					},
				},
				StringVal:	[][]byte{x},
			},
		},
	}

	mark := HasItemInClients(rpcClients, modelName)
	if mark{
		client = rpcClients[modelName]
	}else{
		conn, err := grpc.Dial(REQUESTURL, grpc.WithInsecure())
		if err != nil {
			log.Fatalln(err)
		}
		//defer conn.Close()
		client = pb.NewPredictionServiceClient(conn)
		rpcClients[modelName] = client
	}

	resp, err := client.Predict(context.Background(), request)

	if err != nil {
		res.Code = 1
		log.Fatalln(err)
	}

	for _, v := range resp.Outputs{
		idx := findMaxNum(v.FloatVal)
		i2byte := map[string]interface{}{
			"probility": v.FloatVal[idx],
			"intent": idx}
		res.Data = append(res.Data, i2byte)
	}

	data, _ := json.Marshal(res)

	w.Write(data)

}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/infer", grpcRequestModel).Methods("POST")

	l, err := net.Listen("tcp", "0.0.0.0:9090")
	if err != nil{
		fmt.Println(err)
	}
	//l = netutil.LimitListener(l, 1)
	//loggingRouter := handlers.LoggingHandler(os.Stdout, router)
	fmt.Println("ServerStart!!!!!!!!!!!!!!!")
	http.Serve(l, router)
}
