package sample

import(
	"fmt"
	"net/http"
	"cns"
)

type getProfileObj struct {
    Status bool `json:"status"`
    Msg string `json:"msg"`
}

func GetProfile(input map[string]interface{}, req *http.Request,res http.ResponseWriter, cb chan interface{}){ 

	fmt.Println(input)
	
	responseObject := &getProfileObj{
    	Status:true,
    	Msg:"GetProfile Worked",
    }

	cns.SendMsg(cb, responseObject)
}