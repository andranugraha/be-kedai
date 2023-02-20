package test

import (
	"encoding/json"
	"kedai/backend/be-kedai/internal/server"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/gin-gonic/gin"
)

func ServeReq(opts *server.RouterConfig, req *http.Request) (*gin.Engine, *httptest.ResponseRecorder) {
	router := server.NewRouter(opts)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	return router, rec
}

func MakeRequestBody(dto interface{}) *strings.Reader {
	payload, _ := json.Marshal(dto)
	return strings.NewReader(string(payload))
}
