package controller

import (
	"bytes"
	"encoding/json"
	"github.com/gin-contrib/cache/persistence"
	"github.com/gin-gonic/gin"
	"github.com/librarios/go-librarios/app/config"
	"github.com/librarios/go-librarios/app/model"
	"github.com/librarios/go-librarios/app/plugin"
	"github.com/librarios/go-librarios/app/service"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"
)


// RequestLoggerMiddleware returns middleware that logs request body.
func RequestLoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var buf bytes.Buffer
		tee := io.TeeReader(c.Request.Body, &buf)
		body, _ := ioutil.ReadAll(tee)
		c.Request.Body = ioutil.NopCloser(&buf)
		log.Printf("body: %s", string(body))
		log.Printf("header: %s", c.Request.Header)
		c.Next()
	}
}

// TestServer is a gin-gonic server for testing
type TestServer struct {
	r           *gin.Engine
	bookService service.IBookService
}

// Init initializes gin-gonic server
func (t *TestServer) Init() {
	r := gin.New()
	r.Use(RequestLoggerMiddleware())

	bookService := service.NewBookService()

	cacheStore := persistence.NewInMemoryStore(time.Second)
	addEndpoints(r, cacheStore, bookService)

	t.r = r
	t.bookService = bookService
}

// Get sends GET http request to testServer
func (t *TestServer) Get(url string, result interface{}) *httptest.ResponseRecorder {
	return t.request(http.MethodGet, url, result)
}

// Get sends PATCH http request to testServer
func (t *TestServer) Patch(url string, body interface{}, result interface{}) *httptest.ResponseRecorder {
	return t.requestJSON(http.MethodPatch, url, body, result)
}

// Get sends Post http request to testServer
func (t *TestServer) Post(url string, body interface{}, result interface{}) *httptest.ResponseRecorder {
	return t.requestJSON(http.MethodPost, url, body, result)
}

// Get sends Put http request to testServer
func (t *TestServer) Put(url string, body interface{}, result interface{}) *httptest.ResponseRecorder {
	return t.requestJSON(http.MethodPut, url, body, result)
}

// requestJSON sends http request with JSON body and return response.
// JSON type response contents are parsed into result variable.
func (t *TestServer) requestJSON(
	method string,
	url string,
	body interface{},
	result interface{},
) *httptest.ResponseRecorder {
	var reader io.Reader
	switch body.(type) {
	case string:
		reader = strings.NewReader(body.(string))
		break
	default:
		jsonValue, e := json.Marshal(body)
		if e != nil {
			panic(e)
		}
		reader = bytes.NewBuffer(jsonValue)
		break
	}

	req := httptest.NewRequest(method, url, reader)
	req.Header.Set("Content-Type", "application/json")
	rw := httptest.NewRecorder()

	t.r.ServeHTTP(rw, req)

	// parse result JSON data
	if result != nil {
		if err := json.NewDecoder(rw.Body).Decode(result); err != nil {
			panic(err)
		}
	}

	return rw
}

// request sends http request and return response.
// JSON type response contents are parsed into result variable.
func (t *TestServer) request(
	method string,
	url string,
	result interface{},
) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, url, nil)
	rw := httptest.NewRecorder()
	t.r.ServeHTTP(rw, req)

	// parse result JSON data
	if result != nil {
		if err := json.NewDecoder(rw.Body).Decode(result); err != nil {
			panic(err)
		}
	}

	return rw
}

var testServer *TestServer

// TestMain is a test entryPoint
func TestMain(m *testing.M) {
	setup()
	retCode := m.Run()
	teardown()
	os.Exit(retCode)
}

func setup() {
	configFilename := "../../config/librarios-test.yaml"
	c, err := config.LoadConfigFile(configFilename)
	if err != nil {
		log.Panicf("failed to load config file: %s. error: %v", configFilename, err)
	}
	log.Printf("Loaded: %s\n", configFilename)

	// init plugins
	plugin.InitPlugins(c.Plugins)

	// connect to DB
	if err = config.InitDB(c.DB); err != nil {
		log.Panicf("failed to connect DB. err: %v", err)
	}
	if c.DB["autoMigrate"] == true {
		model.AutoMigrate()
	}

	testServer = &TestServer{}
	testServer.Init()
}

func teardown() {
	config.CloseDB()
}
