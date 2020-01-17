package sample

import(
	_"fmt"
	"net/http"
	"time"
	"cns"
)

type getLoginObj struct {
    Status bool `json:"status"`
    Msg string `json:"msg"`
}

func Login(req *http.Request, res http.ResponseWriter, cb chan interface{}){ 
	
	time.Sleep(1 * time.Second)

	responseObject := &getLoginObj{
    	Status:true,
    	Msg:"Login Worked",
    }

	cns.SendMsg(cb, responseObject)

}