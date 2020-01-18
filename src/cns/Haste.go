package cns

/*
	Import packages
*/

import (
	"context"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	_ "reflect"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
	"math/rand"
	"encoding/json"
	"github.com/gorilla/websocket"
)

var username = "hasteUser"

/*
	Error Pages Static HTML Files
*/

func (v *Http) defaultHeaders(){
	curTime := time.Now()
	v.res.Header().Set("Content-Type", "text/html")
	v.res.Header().Set("Cache-Control", "public,max-age=31536000")
	v.res.Header().Set("Keep-Alive", "timeout=5, max=500")
	v.res.Header().Set("Server", "Go Haste Server")
	v.res.Header().Set("Developed-By", "Pounze It-Solution Pvt Limited")
	v.res.Header().Set("Pragma", "public,max-age=31536000")
	v.res.Header().Set("Expires", curTime.String())
}

// Page Not Found Page 404.html

func (v *Http) pageNotFound() {
	v.defaultHeaders()
	v.res.WriteHeader(http.StatusNotFound)
	pwd, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		return
	}
	_, err = os.Stat(pwd + "/src/error_files/404.html")
	if err == nil {
		var page, err = readErrorFile(pwd + "/src/error_files/404.html")
		if err != nil {
			v.res.Write([]byte("404 Page Not Found"))
		} else {
			v.res.Write(page)
		}
	} else {
		v.res.Write([]byte("404 Page Not Found"))
	}
}

// Access Denies 403 Error File is Served

func (v *Http) accessDenied() {
	v.defaultHeaders()
	v.res.WriteHeader(http.StatusForbidden)
	pwd, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		return
	}

	_, err = os.Stat(pwd + "/src/error_files/403.html")
	if err == nil {
		var page, err = readErrorFile(pwd + "/src/error_files/403.html")
		if err != nil {
			v.res.Write([]byte("403 Access Denied"))
		} else {
			v.res.Write(page)
		}
	} else {
		v.res.Write([]byte("403 Access Denied"))
	}
}

// Internal Server Error 500

func (v *Http) serverError() {
	v.defaultHeaders()
	v.res.WriteHeader(http.StatusInternalServerError)
	pwd, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		return
	}

	_, err = os.Stat(pwd + "/src/error_files/500.html")
	if err == nil {
		var page, err = readErrorFile(pwd + "/src/error_files/500.html")
		if err != nil {
			v.res.Write([]byte("500 Internal Server Error"))
		} else {
			v.res.Write(page)
		}
	} else {
		v.res.Write([]byte("500 Internal Server Error"))
	}
}

// This is where go process is invoked

func init() {
	fmt.Println("Haste framework initiated, GLHF")

	if Config["MAXPROCS"] == "MAX" {
		runtime.GOMAXPROCS(runtime.NumCPU())
	} else {
		var procs, err = strconv.Atoi(Config["MAXPROCS"])
		if err == nil {
			runtime.GOMAXPROCS(1)
		} else {
			runtime.GOMAXPROCS(procs)
		}
	}
}

/*
	Http structure set with http request and response object and current Path
*/

type Http struct {
	res        http.ResponseWriter
	req        *http.Request
	curPath    string
	matchedUrl string
	json       string
}

/*
	chainAndError structure is used to make chaining methods in golang eg: obj.where().middlewares().get()
*/

type ChainAndError struct {
	*Http
	error
}

/*
	globalstruct struct is used to set all routes values
*/

type globalStruct struct {
	method            string
	url               string
	callbackType      string
	callbackFunc      func(req *http.Request, res http.ResponseWriter)
	middlewares       []func(req *http.Request, res http.ResponseWriter, done chan bool)
	globalMiddlewares []func(req *http.Request, res http.ResponseWriter, done chan bool)
}

/*
	Block structure for error handling
*/

type Block struct {
	Try     func()
	Catch   func(Exception)
	Finally func()
}

/*
	Exception interface to get exception
*/

type Exception interface{}

/*
	throw is used to throw exception using panic function and will recieved by recover
*/

func Throw(up Exception) {
	panic(up)
}

/*
	Do function is created to initialize the whole try catch and finally method
*/

func (tcf Block) Do() {
	if tcf.Finally != nil {
		defer tcf.Finally()
	}

	if tcf.Catch != nil {
		defer func() {
			if r := recover(); r != nil {
				tcf.Catch(r)
			}
		}()
		tcf.Try()
	}
}

/*
	global regex structure is used to set regular expressions url
*/

type globalRegex struct {
	regex map[string]string
}

// globalWhereMap is used to collect all regex of the routes

var globalWhereMap = map[string]globalRegex{}

// globalObject is used to collect all global structure

var globalObject = map[string]globalStruct{}

// globalMiscObject is used to collect all global structure

var globalMiscObject = map[string]globalStruct{}

// globalGetObject is used to collect all global structure

var globalGetObject = map[string]globalStruct{}

// globalPostObject is used to collect all global structure

var globalPostObject = map[string]globalStruct{}

// globalPutObject is used to collect all global structure

var globalPutObject = map[string]globalStruct{}

// globalDeleteObject is used to collect all global structure

var globalDeleteObject = map[string]globalStruct{}

// RequestGlobalVariable

var REQUESTMETHOD = ""

type defaultStruct struct {
	callbackFunc func(req *http.Request, res http.ResponseWriter)
}

var defaultObject = map[string]defaultStruct{}

//  block directories list

var dirList []string

// createserver is used to initialize the https server

var mux *http.ServeMux

func CreateHttpsServer(hostPort string, crt string, key string) {
	mux = http.NewServeMux()
	go mux.HandleFunc("/", handHttp)

	fmt.Println("Welcome to Golang Haste Framework ########*****##### Copy-Right 2020, Open Source Project")
	fmt.Println("To get source code visit, github.com/pounze/go_haste or www.pounze.com/pounze/go_haste")
	fmt.Println("Developed by Sudeep Dasgupta(Pounze)")
	fmt.Println("Server started->", hostPort)

	pwd, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		return
	}

	http.ListenAndServeTLS(hostPort, pwd+crt, pwd+key, http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		context := context.WithValue(req.Context(), "username", username)
		mux.ServeHTTP(res, req.WithContext(context))
	}))
}

// createserver is used to initialize the http server

func CreateHttpServer(hostPort string) {
	mux = http.NewServeMux()
	go mux.HandleFunc("/", handHttp)
	fmt.Println("Welcome to Golang Haste Framework ########*****##### Copy-Right 2020, Open Source Project")
	fmt.Println("To get source code visit, github.com/pounze/go_haste or www.pounze.com/pounze/go_haste")
	fmt.Println("Developed by Sudeep Dasgupta(Pounze)")
	fmt.Println("Server started->", hostPort)

	http.ListenAndServe(hostPort, http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		context := context.WithValue(req.Context(), "username", username)
		mux.ServeHTTP(res, req.WithContext(context))
	}))
}

// createserver is used to initialize the http2 server

func CreateHttp2Server(hostPort string, crt string, key string) {
	mux = http.NewServeMux()
	go mux.HandleFunc("/", handHttp)
	fmt.Println("Welcome to Golang Haste Framework ########*****##### Copy-Right 2020, Open Source Project")
	fmt.Println("To get source code visit, github.com/pounze/go_haste or www.pounze.com/pounze/go_haste")
	fmt.Println("Developed by Sudeep Dasgupta(Pounze)")
	fmt.Println("Server started->", hostPort)

	pwd, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		return
	}

	http.ListenAndServeTLS(hostPort, pwd+crt, pwd+key, http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		context := context.WithValue(req.Context(), "username", username)
		mux.ServeHTTP(res, req.WithContext(context))
	}))
}

func (v *Http) DefaultMethod(fn func(req *http.Request, res http.ResponseWriter)) {
	defaultStructObj := defaultStruct{
		callbackFunc: fn,
	}
	defaultObject["default"] = defaultStructObj
}

// handlehttp is used to handle http request

func handHttp(res http.ResponseWriter, req *http.Request) {

	defer req.Body.Close()

	var mutex = &sync.Mutex{}

	if defaultObject != nil {
		defaultObject["default"].callbackFunc(req, res)
	}

	// if url match with global

	mutex.Lock()

	httpObj := Http{
		req: req,
		res: res,
	}

	if req.Method == "GET" {

		globalObject = globalGetObject

	} else if req.Method == "POST" {

		globalObject = globalPostObject

	} else if req.Method == "PUT" {

		globalObject = globalPutObject

	} else {

		globalObject = globalDeleteObject

	}

	var routesLen = len(globalObject)

	if routesLen == 0 {
		fmt.Println("Please create a route")
		return
	}

	var dirFailed = 0

	var dirListLen = len(dirList)

	for _, val := range dirList {

		val = strings.Replace(val, "/", "\\/", -1)

		match, _ := regexp.MatchString("(?i)"+val+"[a-z0-9A-Z\\.]*", req.URL.Path)

		if match {

			httpObj.accessDenied()
			break
			return

		} else {

			dirFailed += 1

		}
	}

	if dirFailed == dirListLen {

		var ext = filepath.Ext(req.URL.Path)

		if ext == "" {

			var notMatched = 0

			var urlJoin string

			urlJoin = ""

			for key, val := range globalObject {

				val.url = strings.Replace(val.url, "/", "\\/", -1)

				match, _ := regexp.MatchString("^"+val.url+"$", req.URL.Path)

				if match && req.Method == val.method {

					httpObj.matchedUrl = key

					key = strings.Replace(key, "[\\", "", 1)

					urlArr := strings.Split(req.URL.Path, "/")

					urlLen := len(urlArr)

					valArr := strings.Split(key, "/")

					valLen := len(valArr)

					if urlLen == valLen {

						for key, val := range urlArr {

							if val != "" {

								urlJoin += "&" + valArr[key] + "=" + val

							}
						}

						if req.URL.RawQuery == "" {

							req.URL.RawQuery = urlJoin

						} else {

							req.URL.RawQuery = req.URL.RawQuery + urlJoin

						}
					}

					var parseChan = make(chan bool)

					if httpObj.invokeGlobalMiddlewares(parseChan) {
						
						go httpObj.parseHttpRequest(parseChan)

						if <-parseChan == true {

							if httpObj.invokeMiddlewares(parseChan) {

								httpObj.invokeMethod()

								defer close(parseChan)

								break

							} else {
								break
							}

						}

					} else {
						break
					}

				} else {

					notMatched += 1

				}
			}

			if notMatched == routesLen {

				httpObj.pageNotFound()

			}
		} else {

			httpObj.serveStaticFiles(ext)

		}
	}

	mutex.Unlock()
}

func (v *Http) parseHttpRequest(parseChan chan bool) {

	if v.req.Method == "POST" || v.req.Method == "PUT" || v.req.Method == "DELETE" {

		if v.req.Header.Get("Content-Type") == "" || v.req.Header.Get("Content-Type") == "multipart/form-data" {

			parseChan <- true

		} else {

			if v.req.Header.Get("Content-Type") == "application/json" {

				parseChan <- true

			} else {

				v.req.ParseForm()
				parseChan <- true

			}

		}

	} else {

		parseChan <- true

	}
}

func (v *Http) SaveFile(w http.ResponseWriter, file multipart.File, path string, fileChan chan bool) {

	data, err := ioutil.ReadAll(file)

	if err != nil {

		fmt.Fprintf(w, "%v", err)

		fileChan <- false

		return
	}

	err = ioutil.WriteFile(path, data, 0666)

	if err != nil {

		fmt.Fprintf(w, "%v", err)

		fileChan <- false

		return
	}

	fileChan <- true
}

func (v *Http) serveStaticFiles(ext string) {

	v.res.Header().Set("Content-Type", MimeList[ext])

	curTime := time.Now()

	v.res.Header().Set("Keep-Alive", "timeout=5, max=500")
	v.res.Header().Set("Server", "Go Haste Server")
	v.res.Header().Set("Developed-By", "Pounze It-Solution Pvt Limited")
	v.res.Header().Set("Expires", curTime.String())

	pwd, err := os.Getwd()

	if err != nil {

		fmt.Println(err)

		return
	}
	_, err = os.Stat(pwd + v.req.URL.Path)

	if err == nil {

		var page, err = readErrorFile(pwd + v.req.URL.Path)

		if err != nil {

			v.pageNotFound()

		} else {

			lastModifiedDate := v.req.Header.Get("if-modified-since")

			file, err := os.Stat(pwd + v.req.URL.Path)

			if err != nil {

				v.pageNotFound()

			} else {

				modifiedtime := file.ModTime()

				if lastModifiedDate == "" {

					v.res.Header().Set("Last-Modified", modifiedtime.Format("2006-01-02 15:04:05"))

					v.res.WriteHeader(http.StatusOK)

				} else {

					time, err := time.Parse("2006-01-02 15:04:05", lastModifiedDate)

					if err != nil {

						v.pageNotFound()

					} else {

						v.res.Header().Set("Last-Modified", modifiedtime.Format("2006-01-02 15:04:05"))

						if time.Unix() < modifiedtime.Unix() {

							v.res.WriteHeader(http.StatusOK)

						} else {

							v.res.WriteHeader(http.StatusNotModified)

						}
					}
				}
			}

			v.res.Write(page)
		}
	} else {

		v.pageNotFound()

	}
}

func (v *Http) invokeGlobalMiddlewares(parseChan chan bool) bool {

	var middlewareCount int
	middlewareCount = 0

	if globalMiscObject["GlobalMiddleWares"].globalMiddlewares != nil {
		for index, _ := range globalMiscObject["GlobalMiddleWares"].globalMiddlewares {

			go globalMiscObject["GlobalMiddleWares"].globalMiddlewares[index](v.req, v.res, parseChan)

			if <-parseChan == false {
				middlewareCount += 1
			}
		}

		if middlewareCount > 0 {
			return false
		}
	}

	return true
}

func (v *Http) invokeMiddlewares(parseChan chan bool) bool {

	var middlewareCount int
	middlewareCount = 0

	if globalObject[v.matchedUrl].middlewares != nil {
		for index, _ := range globalObject[v.matchedUrl].middlewares {
			
			go globalObject[v.matchedUrl].middlewares[index](v.req, v.res, parseChan)

			if <-parseChan == false {
				middlewareCount += 1
			}
		}

		if middlewareCount > 0 {
			return false
		}
	}

	return true
}

func (v *Http) invokeMethod() {

	globalObject[v.matchedUrl].callbackFunc(v.req, v.res)

}

func (v *Http) BlockDirectories(dir []string) {

	dirList = dir

}

func readErrorFile(path string) ([]byte, error) {

	b, err := ioutil.ReadFile(path) // just pass the file name

	if err != nil {

		var emptyByte []byte
		return emptyByte, err

	}

	return b, nil
}

// getFunc is used to handle get request with func parameters

func (v *Http) Get(url string, fn func(req *http.Request, res http.ResponseWriter)) ChainAndError {

	globalStructObj := globalStruct{
		method:       "GET",
		url:          url,
		callbackType: "func",
		callbackFunc: fn,
	}

	v.curPath = url
	globalGetObject[url] = globalStructObj
	REQUESTMETHOD = "GET"
	return ChainAndError{v, nil}
}

// PostFunc is used to handle post request with func parameters

func (v *Http) Post(url string, fn func(req *http.Request, res http.ResponseWriter)) ChainAndError {
	globalStructObj := globalStruct{
		method:       "POST",
		url:          url,
		callbackType: "func",
		callbackFunc: fn,
	}
	v.curPath = url
	globalPostObject[url] = globalStructObj
	REQUESTMETHOD = "POST"
	return ChainAndError{v, nil}
}

// PutFunc is used to handle post request with func parameters

func (v *Http) Put(url string, fn func(req *http.Request, res http.ResponseWriter)) ChainAndError {
	globalStructObj := globalStruct{
		method:       "PUT",
		url:          url,
		callbackType: "func",
		callbackFunc: fn,
	}
	v.curPath = url
	globalPutObject[url] = globalStructObj
	REQUESTMETHOD = "PUT"
	return ChainAndError{v, nil}
}

// DeleteFunc is used to handle post request with func parameters

func (v *Http) Delete(url string, fn func(req *http.Request, res http.ResponseWriter)) ChainAndError {
	globalStructObj := globalStruct{
		method:       "DELETE",
		url:          url,
		callbackType: "func",
		callbackFunc: fn,
	}
	v.curPath = url
	globalDeleteObject[url] = globalStructObj
	REQUESTMETHOD = "DELETE"
	return ChainAndError{v, nil}
}

// middlewares method is used to handle middlewares

func (v *Http) Middlewares(middlewares func(req *http.Request, res http.ResponseWriter, done chan bool)) ChainAndError {

	if REQUESTMETHOD == "GET" {

		var tmp = globalGetObject[v.curPath]
		tmp.middlewares = append(tmp.middlewares, middlewares)
		globalGetObject[v.curPath] = tmp

	} else if REQUESTMETHOD == "POST" {

		var tmp = globalPostObject[v.curPath]
		tmp.middlewares = append(tmp.middlewares, middlewares)
		globalPostObject[v.curPath] = tmp

	} else if REQUESTMETHOD == "PUT" {

		var tmp = globalPutObject[v.curPath]
		tmp.middlewares = append(tmp.middlewares, middlewares)
		globalPutObject[v.curPath] = tmp

	} else {

		var tmp = globalDeleteObject[v.curPath]
		tmp.middlewares = append(tmp.middlewares, middlewares)
		globalDeleteObject[v.curPath] = tmp

	}

	return ChainAndError{v, nil}
}

// middlewares method is used to handle middlewares

func (v *Http) GlobalMiddleWares(middlewares func(req *http.Request, res http.ResponseWriter, done chan bool)) ChainAndError {

	var tmp = globalMiscObject["GlobalMiddleWares"]

	tmp.globalMiddlewares = append(tmp.globalMiddlewares, middlewares)

	globalMiscObject["GlobalMiddleWares"] = tmp

	return ChainAndError{v, nil}
}

// where method is used to parse url matching with the regular expression

func (v *Http) Where(regex map[string]string) ChainAndError {

	if REQUESTMETHOD == "GET" {
		var tmp = globalGetObject[v.curPath]

		tmp.url = globalGetObject[v.curPath].url

		for key, value := range regex {
			tmp.url = strings.Replace(tmp.url, key, value, -1)
		}

		globalGetObject[v.curPath] = tmp

	} else if REQUESTMETHOD == "POST" {

		var tmp = globalPostObject[v.curPath]

		tmp.url = globalPostObject[v.curPath].url

		for key, value := range regex {
			tmp.url = strings.Replace(tmp.url, key, value, -1)
		}

		globalPostObject[v.curPath] = tmp
	} else if REQUESTMETHOD == "PUT" {

		var tmp = globalPutObject[v.curPath]

		tmp.url = globalPutObject[v.curPath].url

		for key, value := range regex {
			tmp.url = strings.Replace(tmp.url, key, value, -1)
		}

		globalPutObject[v.curPath] = tmp
	} else {
		var tmp = globalDeleteObject[v.curPath]

		tmp.url = globalDeleteObject[v.curPath].url

		for key, value := range regex {
			tmp.url = strings.Replace(tmp.url, key, value, -1)
		}

		globalDeleteObject[v.curPath] = tmp
	}

	return ChainAndError{v, nil}
}

/*
	Projection Api For Golang Haste Framework.
	To check for node implementation go to the node projection api
*/

type Projection struct {
	Method func(input map[string]interface{}, req *http.Request, res http.ResponseWriter, cb chan interface{})
	Timeout int32
}

// creating hashmap for projection schema

var projectHM = map[string]Projection{}

// setting schema

func (v *Http) CreateSchema(projectMap map[string]Projection){
	projectHM = projectMap
}

// default response structure

type DefaultResponse struct {
    Status bool `json:"status"`
    Msg string `json:"msg"`
}

// method to check if channel is closed before sending message

func IsClosed(ch chan interface{}) bool {
	select {
	case <-ch:
		return true
	default:
	}

	return false
}

// send message method is used in projection api to send message over the channels

func SendMsg(ch chan interface{}, class interface{}){
	if(!IsClosed(ch)){
		ch <- class
	}
}

// project method for http/1.1

func (v *Http) ProjectionMethod(req *http.Request,res http.ResponseWriter){

	// creating hashmap for response

	var responseHashMap = make(map[string]interface{})

	// creating another hashmap to parse request json 

	requestHM := make(map[string]interface{})

	// decoding json request

	err := json.NewDecoder(req.Body).Decode(&requestHM)

	if err != nil {
        responseObject := &DefaultResponse{
	    	Status:false,
	    	Msg:"Oops, something went wrong",
	    }

	    jsonData, _ := json.Marshal(responseObject)
		res.Write([]byte(jsonData))
        return
    }

    // iterating over all the request schema in the request payload

	for key, _ := range requestHM{

		if _, ok := projectHM[key]; ok {

			// creating channels for callback

		    callbackChan := make(chan interface{})

		    // setting timeout for context

			timeoutTime := rand.Int31n(projectHM[key].Timeout)

			ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeoutTime) * time.Millisecond)

			defer cancel()

			// invoking method with go routines

			md, ok := requestHM[key].(map[string]interface{})

			if !ok{

				responseObject := &DefaultResponse{
			    	Status:false,
			    	Msg:"Method "+key+" does not contain valid json",
			    }
				responseHashMap[key] = responseObject

				continue
			}

	        go projectHM[key].Method(md, req, res, callbackChan)

	        // checking if context is done

	        select{
			  case <-ctx.Done():
			    responseObject := &DefaultResponse{
			    	Status:false,
			    	Msg:"Method "+key+" timed out",
			    }
				responseHashMap[key] = responseObject
			  case callback := <-callbackChan:
			    responseHashMap[key] = callback
		  	}

		  	// closing the channel

		  	close(callbackChan)

		}else{
			 responseObject := &DefaultResponse{
		    	Status:false,
		    	Msg:"Method "+key+" does not exists",
		    }
			responseHashMap[key] = responseObject
		}
    }

    // sending the http response

    jsonData, _ := json.Marshal(responseHashMap)

    res.Write([]byte(jsonData))
}

// set route path to set projection url

func (v *Http) SetRoutePath(path string) ChainAndError{
	v.Post(path, v.ProjectionMethod)
	return ChainAndError{v, nil}
}

// creating upgrader object for the gorilla websocket

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// method to handle socket projection request

func (v *Http) SocketProjection(res http.ResponseWriter, req *http.Request){

	// upgrade method to upgrade http to websocket

	con, err := upgrader.Upgrade(res, req, nil)
	
	if err != nil {
		return
	}

	// closing the socket connection

	defer con.Close()
	
	for {


		// reading messages from sockets

		mt, message, err := con.ReadMessage()
		
		if err != nil {
			break
		}

		// creating hashmap  for parsing the request

		requestHM := make(map[string]interface{})

		// creating object from json string

		jsonError := json.Unmarshal(message, &requestHM)

		if jsonError != nil {
	        responseObject := &DefaultResponse{
		    	Status:false,
		    	Msg:"Oops, something went wrong",
		    }

		    jsonData, _ := json.Marshal(responseObject)

	        writeError := con.WriteMessage(mt, []byte(jsonData))

	        if writeError != nil {
				break
			}
	    }

	    // iterating over the json object

	    for key, _ := range requestHM{

	    	var responseHashMap = make(map[string]interface{})
	    	
			if _, ok := projectHM[key]; ok {

				// creating channel to get the response

			    callbackChan := make(chan interface{})

			    // setting context timeout

				timeoutTime := rand.Int31n(projectHM[key].Timeout)

				ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeoutTime) * time.Millisecond)

				defer cancel()

				// invoking method for matched method from request payload

				md, ok := requestHM[key].(map[string]interface{})

				if !ok{
					
					responseObject := &DefaultResponse{
				    	Status:false,
				    	Msg:"Method "+key+" does not contain valid json",
				    }
					responseHashMap[key] = responseObject

					continue
				}

		        go projectHM[key].Method(md, req, res, callbackChan)

		        select{
				  case <-ctx.Done():
				    responseObject := &DefaultResponse{
				    	Status:false,
				    	Msg:"Method "+key+" timed out",
				    }
					responseHashMap[key] = responseObject
				  case callback := <-callbackChan:
				    responseHashMap[key] = callback
			  	}

			  	// closing the channel

			  	close(callbackChan)

			}else{
				 responseObject := &DefaultResponse{
			    	Status:false,
			    	Msg:"Method "+key+" does not exists",
			    }
				responseHashMap[key] = responseObject
			}

			// creating json string and sending as response

			jsonData, _ := json.Marshal(responseHashMap)

    		writeError := con.WriteMessage(mt, []byte(jsonData))

	        if writeError != nil {
				break
			}
	    }
		
	}
}

// creating socket path for websocket streaming

func (v *Http) SetSocketRoutePath(path string) ChainAndError{
	
	go func(){
		for true {
		
			if mux != nil{
				break
			}

	        time.Sleep(1 * time.Second)
	    }

		go mux.HandleFunc(path, v.SocketProjection)
	}()

	return ChainAndError{v, nil}
}

// declaring variable for http2 multiplexing

var muxStream *http.ServeMux

// creating http2 streaming, it has only support with TLS 

func CreateHttpStreaming(hostPort string, crt string, key string, path string){
	muxStream = http.NewServeMux()
	go muxStream.HandleFunc(path, handleHttpStreaming)
	fmt.Println("Welcome to Golang Haste Framework ########*****##### Copy-Right 2020, Open Source Project")
	fmt.Println("To get source code visit, github.com/pounze/go_haste or www.pounze.com/pounze/go_haste")
	fmt.Println("Developed by Sudeep Dasgupta(Pounze)")
	fmt.Println("Server started->", hostPort)

	pwd, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		return
	}

	http.ListenAndServeTLS(hostPort, pwd+crt, pwd+key, http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		context := context.WithValue(req.Context(), "username", username)
		muxStream.ServeHTTP(res, req.WithContext(context))
	}))
}

// handling the http2 streaming

func handleHttpStreaming(res http.ResponseWriter, req *http.Request) {

	// creating hashmap for request parsing

	requestHM := make(map[string]interface{})

	// decoding json string to structure

	err := json.NewDecoder(req.Body).Decode(&requestHM)

	if err != nil {
		fmt.Println(err)
        responseObject := &DefaultResponse{
	    	Status:false,
	    	Msg:"Oops, something went wrong",
	    }

	    jsonData, _ := json.Marshal(responseObject)
		res.Write([]byte(jsonData))
        return
    }

    // iterating over the request

    for key, _ := range requestHM{

    	// creating response hashmap

    	var responseHashMap = make(map[string]interface{})

		if _, ok := projectHM[key]; ok {

			// creating channel for getting response from the invoked methods

		    callbackChan := make(chan interface{})

		    // setting context timeout 

			timeoutTime := rand.Int31n(projectHM[key].Timeout)

			ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeoutTime) * time.Millisecond)

			defer cancel()

			// invoking methods with goroutines

			md, ok := requestHM[key].(map[string]interface{})

			if !ok{
				
				responseObject := &DefaultResponse{
			    	Status:false,
			    	Msg:"Method "+key+" does not contain valid json",
			    }
				responseHashMap[key] = responseObject

				continue
			}

	        go projectHM[key].Method(md, req, res, callbackChan)

	        select{
			  case <-ctx.Done():
			    responseObject := &DefaultResponse{
			    	Status:false,
			    	Msg:"Method "+key+" timed out",
			    }
				responseHashMap[key] = responseObject
			  case callback := <-callbackChan:
			    responseHashMap[key] = callback
		  	}

		  	// closing the channel

		  	close(callbackChan)

		}else{
			 responseObject := &DefaultResponse{
		    	Status:false,
		    	Msg:"Method "+key+" does not exists",
		    }
			responseHashMap[key] = responseObject
		}

		// creating json string for streaming response with a splitter \r\n\r\n

		jsonData, _ := json.Marshal(responseHashMap)
		res.Write([]byte(string(jsonData)+"\r\n\r\n"))
		res.(http.Flusher).Flush()
    } 
}