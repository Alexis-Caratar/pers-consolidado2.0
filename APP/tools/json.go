package tools

import (
	"encoding/json"
	"net/http"
	"time"
)

//JSONEncode struct a json([]byte)
func JSONEncode(data interface{}) (js []byte, err error) {
	js, err = json.Marshal(data)
	return js, err
}

//Responder Reponse
func Responder(w http.ResponseWriter, data interface{}, json bool) (err error) {
	if json {
		js, err := JSONEncode(data)
		if err == nil {
			w.Header().Set("Expires: Mon", "5 Jan 1993 05:00:00 GMT")
			w.Header().Set("Last-Modified", time.Now().Format(time.RFC1123)+" GMT")
			w.Header().Set("Cache-Control", "no-cache, must-revalidate")
			w.Header().Set("Pragma", "no-cache")
			w.Header().Set("Content-Type", "application/json")
			w.Write(js)
		}
	} else {
		w.Write([]byte(data.(string)))
	}
	return err
}
