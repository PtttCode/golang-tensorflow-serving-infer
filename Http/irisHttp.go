package main

import (
	"github.com/kataras/golog"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/hero"
	"github.com/kataras/iris/v12/middleware/logger"
	"github.com/kataras/iris/v12/middleware/recover"
	"math/rand"
	"net"
	"sort"
)

type sortBody struct {
	Array []int
}

func requestDeal(bodyJson sortBody) map[string][]int{
	var xrange []int
	if bodyJson.Array != nil{
		xrange = bodyJson.Array
	}else{
		for i := 0;i < 5000; i++{
			xrange = append(xrange, i)
		}
		rand.Shuffle(len(xrange), func(i, j int) {
			xrange[i], xrange[j] = xrange[j], xrange[i]
		})
	}
	sort.Ints(xrange)

	return map[string][]int{
		"data": xrange,
	}
}

func main() {
	golog.SetTimeFormat("")
	golog.SetLevel("debug")
	app := iris.New()
	app.Logger().SetLevel("debug")
	app.Use(recover.New())
	app.Use(logger.New())
	//goroutines := New(40)

	hero.Register(func (ctx iris.Context)(jsonData requestBody){
		ctx.ReadJSON(&jsonData)
		return
	})
	inferHandler := hero.Handler(requestModel)
	app.Post("/infer", inferHandler)
	//app.Post("/infer",  func(ctx iris.Context){
	//	goroutines.Add(1)
	//	go func() {
	//		var jsonData requestBody
	//		golog.Debug(runtime.NumGoroutine())
	//		ctx.ReadJSON(&jsonData)
	//		result := requestModel(jsonData)
	//		ctx.JSON(result)
	//		goroutines.Done()
	//	}()
	//	goroutines.Wait()
	//})

	hero.Register(func (ctx iris.Context)(array sortBody){
		ctx.ReadJSON(&array)
		return
	})
	sortHandler := hero.Handler(requestDeal)
	app.Post("/sort", sortHandler)

	app.Get("/")

	l, err := net.Listen("tcp", "0.0.0.0:9090")
	if err != nil{
		golog.Error("Listener Error", err)
	}
	defer l.Close()
	//l = netutil.LimitListener(l, 1)

	app.Run(iris.Listener(l), iris.WithoutServerError(iris.ErrServerClosed))

}