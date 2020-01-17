package sample

import(
	_"fmt"
	"net/http"
	"cns"
)

type getProfileObj struct {
    Status bool `json:"status"`
    Msg string `json:"msg"`
}

func GetProfile(req *http.Request,res http.ResponseWriter, cb chan interface{}){ 
	
	responseObject := &getProfileObj{
    	Status:true,
    	Msg:"GetProfile Worked",
    }

	cns.SendMsg(cb, responseObject)
}