package main

import(
	"cns"
	"Web"
	"net/http"
	"fmt"
)

func main(){
	
	httpApp := cns.Http{}

	Web.Routes()

	go func(){
		defer cns.CreateHttpServer(":8000")
	}()

	defer cns.CreateHttpStreaming(":8100", "./ssl/https-server.crt","./ssl/https-server.key", "/streaming")

	defer httpApp.DefaultMethod(func(req *http.Request,res http.ResponseWriter){
		res.Header().Set("Name", "Sudeep Dasgupta")
		res.Header().Set("Content-Type", "application/json")
		res.Header().Set("Access-Control-Allow-Origin", "*")
		res.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")
		fmt.Println("Default header executed")
	})
}