# go_haste
GoLang Web Framework 

This is Go Haste Framework 1.0.0, developed in 10th Feb 2018 by Pounze It Solutions Pvt Limited.
This framework is open source project.

This is basic tutorials of goLang haste framework.

For any queries visit www.pounze.com/dev/help
For complete documentation visit: www.pounze.com/projects/goHaste

For consulting query contact:

sudeepdasgupta25@gmail.com | sudeep.dasgupta@pounze.com

Features:

1) Middlewares
2) Go default routes are modified
3) 401 Authentication support
4) URL Binding 
5) Handles all types of request and return object
6) Handles errors like 500 errors, 404 not found etc.

Upcoming Features:

1) Images crop and scaling lib.
2) ORM for SQL
3) ODM for Mongodb
4) Session handling for files and redis support
5) Parallelism support for multiple CPU cores.
6) Inbuilt Messaging queue with streaming architechture for real time streams.
7) JSON parser with complete object, no need for static scheme structure.
8) Dynamic request routes to controllers.
9) Monitoring tool with mail support

###########################################################################################################################

Basic Routes  

Server.go

package main

import(
	"cns"
	"Web"
)

func main(){
	
	Web.Routes()

	defer cns.CreateHttpServer(":8000")
}

/cns/Haste.go

It is the main library file of the whole framework.

/Web/RouteList.go

Its contains all routes

package Web

import(
	"cns"
	"fmt"
	"net/http"
)

func Routes(){
	httpApp := cns.Http{}

	//httpApp.BlockDirectories([]string{"/UserView/img/","/UserView/js/"})

	hm := map[string]string{
	    "$id":"[0-9]{2}",
	    "$name":"[a-z]+",
	}

	httpApp.Post("/$id/$name[\\/]*",func(req *http.Request,res http.ResponseWriter){
		fmt.Println("method invoked")
		fmt.Println(req.URL.Query().Get("$name"))
		fmt.Fprintf(res, "Successfully Done")
	}).Middlewares(func(req *http.Request,res http.ResponseWriter,done chan bool){

		err := req.ParseMultipartForm(200000)
        if err != nil {
            fmt.Println("Unable to parse form data")
            return
        }
        _,handle,_ := req.FormFile("name")
        fmt.Println(handle)

		done <- true
	}).Middlewares(func(req *http.Request,res http.ResponseWriter,done chan bool){
		fmt.Println("working both")
		done <- true
	}).Where(hm)


	httpApp.Get("/",func(req *http.Request,res http.ResponseWriter){
		//cns.Push(res,"/UserView/js/test.js")
		stat,_,result := cns.Authorization(req,res,"Enter username and password to Authenticate")
		if stat{
			fmt.Println(result)
			http.ServeFile(res, req, "src/views/index.html")
		}else{
			fmt.Fprintf(res, "Unauthorized")
		}
	})

	cns.Block{
		Try:func(){
			fmt.Println("I tried")
			cns.Throw("ohhh")
			fmt.Println("Working")
		},
		Catch:func(e cns.Exception){
			fmt.Println("caught exception",e)
		},
		Finally:func(){
			fmt.Println("Finally")
		},
	}.Do()
}


In the above example it inclues all examples for 401, middlewares, url binding, routes, request handling like : json, multipart etc.
