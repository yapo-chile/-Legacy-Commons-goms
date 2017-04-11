package sources

import (
	"github.com/Yapo/goutils"
	"github.com/Yapo/logger"
	"net/http"
)

//  MyGOMSHandler is an example handler
func MyGOMSHandler(w http.ResponseWriter, r *http.Request) {
	logger.Info("Request: [%s] %s from: %s", r.Method, r.URL, r.RemoteAddr)
	var response goutils.Response
	defer goutils.WriteJSONResponse(w, &response)
	defer goutils.CreateJSON(&response)

	response.Body = struct {
		Response string
	}{Response: "HOLA MUNDO!"}
	response.Code = http.StatusOK
}
