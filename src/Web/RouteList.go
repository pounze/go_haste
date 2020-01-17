package Web

import(
	"cns"
	"fmt"
	"net/http"
	"controllers/sample"
)

func Routes(){
	httpApp := cns.Http{}

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

	httpApp.SetRoutePath("/").Middlewares(func(req *http.Request,res http.ResponseWriter,done chan bool){
		fmt.Println("middleware worked")
		done <- true
	})

	httpApp.SetSocketRoutePath("/realtime").Middlewares(func(req *http.Request,res http.ResponseWriter,done chan bool){
		fmt.Println("middleware worked")
		done <- true
	})

	// httpApp.BlockDirectories([]string{"/UserView/img/","/UserView/js/"})

	// hm := map[string]string{
	//     "$id":"[0-9]{2}",
	//     "$name":"[a-z]+",
	// }

	// httpApp.Put("/$id/$name",func(req *http.Request,res http.ResponseWriter){
	// 	fmt.Println("method invoked")
	// 	fmt.Println(req.URL.Query().Get("$name"))
	// 	fmt.Println(req.URL.RawQuery)
	// 	fmt.Fprintf(res, "Successfully Done")
	// }).Where(hm)


	// httpApp.Get("/",func(req *http.Request,res http.ResponseWriter){
	// 	http.ServeFile(res, req, "src/views/index.html")
	// }).Middlewares(func(req *http.Request,res http.ResponseWriter,done chan bool){
	// 	fmt.Println(req.URL.Query().Get("name"))
	// 	done <- true
	// })

	// httpApp.Post("/",func(req *http.Request,res http.ResponseWriter){
	// 	res.Write([]byte("ok"))
	// })

	// httpApp.GlobalMiddleWares(func(req *http.Request,res http.ResponseWriter,done chan bool){

	// 	fmt.Println("Working Global Middleware 1")

	// 	done <- true
	
	// }).GlobalMiddleWares(func(req *http.Request,res http.ResponseWriter,done chan bool){

	// 	fmt.Println("Working Global Middleware 2")

	// 	done <- true

	// })

	httpApp.Post("/sample", sample.Sample)

	httpApp.Post("/stream",func(req *http.Request,res http.ResponseWriter){
		fmt.Println("method invoked")
	})

	// cns.Block{
	// 	Try:func(){
	// 		fmt.Println("I tried")
	// 		cns.Throw("ohhh")
	// 		fmt.Println("Working")
	// 	},
	// 	Catch:func(e cns.Exception){
	// 		fmt.Println("caught exception",e)
	// 	},
	// 	Finally:func(){
	// 		fmt.Println("Finally")
	// 	},
	// }.Do()
}