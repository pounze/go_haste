# go_haste
GoLang Web Framework 

This is Go Haste Framework 1.0.0, developed in 10th Feb 2018 by Pounze It Solutions Pvt Limited.
This framework is open source project.

This is basic tutorials of goLang haste framework.

For any queries visit www.pounze.com/dev/help
For complete documentation visit: www.pounze.com/projects/goHaste

For consulting query contact:

sudeepdasgupta25@gmail.com | sudeep@pounze.com

Features:

1) Middlewares
2) Go default routes are modified
3) 401 Authentication support
4) URL Binding 
5) Handles all types of request and return object
6) Handles errors like 500 errors, 404 not found etc.
7) Projection Api inspired from Google GRPC. It reduces http calls to single payload and send messages one after another when they are completed.

###########################################################################################################################

For Projection api Websocket Gorilla is needed. It is open source and very powerfull websocket library in golang.
To download write the command in the root directory
go get github.com/gorilla/websocket

To set configuration for database and everything
/src/cns/Config.go

Basic Routes  

Server.go
    
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
        
        // we can also create https server 
        
        defer cns.CreateHttpsServer(":8300", "./ssl/https-server.crt","./ssl/https-server.key", "/")
        
        // for http/2 server we write the following code
        
        defer cns.CreateHttp2Server(":8300", "./ssl/https-server.crt","./ssl/https-server.key", "/")

        defer httpApp.DefaultMethod(func(req *http.Request,res http.ResponseWriter){
            res.Header().Set("Name", "Sudeep Dasgupta")
            res.Header().Set("Content-Type", "application/json")
            res.Header().Set("Access-Control-Allow-Origin", "*")
            res.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")
            fmt.Println("Default header executed")
        })
    }

/cns/Haste.go

It is the main library file of the whole framework.

/Web/RouteList.go

Its contains all routes

package Web

Web/RouteList.go

    import(
        "cns"
        "fmt"
        "net/http"
    )

    func Routes(){
        httpApp := cns.Http{}
        
        // for projection api create schema 
        
        httpApp.CreateSchema(map[string]cns.Projection{
            "Login": cns.Projection{
            sample.Login,
            10,
            },
            "GetProfile": cns.Projection{
            sample.GetProfile,
            10,
            },
        });
        
        // set route path for projection schema add middleware method to it
        
        httpApp.SetRoutePath("/").Middlewares(func(req *http.Request,res http.ResponseWriter,done chan bool){
            fmt.Println("middleware worked")
            done <- true
        })
        
        // projection api for socket support 
        
        httpApp.SetSocketRoutePath("/realtime").Middlewares(func(req *http.Request,res http.ResponseWriter,done chan bool){
            fmt.Println("middleware worked")
            done <- true
        })
        
        // block directories to get access
        
        httpApp.BlockDirectories([]string{"/UserView/img/","/UserView/js/"})

        hm := map[string]string{
            "$id":"[0-9]{2}",
            "$name":"[a-z]+",
        }
        
        // to save file in case of multipart form data
        
        httpApp.SaveFile(w http.ResponseWriter, file multipart.File, path string, fileChan chan bool)

        // url matching using regular expression and calling multiple middleware with chaining
        
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

        // calling Get request method same POST, PUT and DELETE available
        
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
        
        // TO call global middleware
        
        httpApp.GlobalMiddleWares(func(req *http.Request,res http.ResponseWriter,done chan bool){
            fmt.Println("working both")
            done <- true
        })
    
        // creating try catch block to handle exception
        
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
