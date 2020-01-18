package sample

import(
	"fmt"
	"net/http"
	"time"
	"cns"
)

type getLoginObj struct {
    Status bool `json:"status"`
    Msg string `json:"msg"`
}

func Login(input map[string]interface{}, req *http.Request, res http.ResponseWriter, cb chan interface{}){ 

	fmt.Println(input)
	fmt.Println(input["username"])
	
	time.Sleep(1 * time.Second)

	responseObject := &getLoginObj{
    	Status:true,
    	Msg:"Login Worked",
    }

	cns.SendMsg(cb, responseObject)

}