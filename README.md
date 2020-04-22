# golang-tensorflow-serving-infer
Communicate with Tensorflow-serving via Grpc and Http.Programmed by Goalng
## Grpc
### requirements
```
go get github.com/gorilla/mux
go get golang.org/x/net
# if you fail to get "net" package, you can run this command:
# git clone https://github.com/golang/net.git $GOPATH/src/golang.org/x/net
go get google.golang.org/grpc
```
You can get Tensorflow and Tensorflow-serving-api packages with their guide :
```
# For quickstart, you can copy ./tensorflow to $GOPATH/src/github.com/
# Necessary requirements
https://github.com/tensorflow/tensorflow/tree/master/tensorflow/go
https://github.com/tensorflow/serving
```
## Http
### requirements
```
go get github.com/gorilla/handlers  # v1.4.2
go get github.com/gorilla/mux	# v1.7.4
go get github.com/kataras/golog	# v0.0.10
```

## Request
### Body
```
{
    "question": "小红今天怎么样啊",	# 问句
    "length": 30	# 模型的maxlen
}
```
### Constant Variables
```
modelName = "tfs-model"	# Tensorflow-serving --model_name
REQUESTURL = "127.0.0.1:6007"	# Tensorflow-serving ip:{port or rest_api_port}
padding = "<PAD>"	# Model padding strings
```
