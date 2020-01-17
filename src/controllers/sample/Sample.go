package sample

import(
	"fmt"
	"net/http"
)

func Sample(req *http.Request, res http.ResponseWriter){ 
	fmt.Println("Working");
}