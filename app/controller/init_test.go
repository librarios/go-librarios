package controller

import (
	"bytes"
	"encoding/json"
	"github.com/gin-contrib/cache/persistence"
	"github.com/gin-gonic/gin"
	"github.com/librarios/go-librarios/app/config"
	"github.com/librarios/go-librarios/app/plugin"
	"github.com/librarios/go-librarios/app/service"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"
)

type TestServer struct {
	r           *gin.Engine
	bookService service.IBookService
}

func (t *TestServer) Init() {
	t.r = gin.New()
	t.bookService = service.NewBookService()
	cacheStore := persistence.NewInMemoryStore(time.Second)
	addEndpoints(t.r, cacheStore, t.bookService)
}

func (t *TestServer) Get(url string, result interface{}) *httptest.ResponseRecorder {
	return t.request(http.MethodGet, url, result)
}

func (t *TestServer) Patch(url string, body interface{}, result interface{}) *httptest.ResponseRecorder {
	return t.requestJSON(http.MethodPatch, url, body, result)
}

func (t *TestServer) Post(url string, body interface{}, result interface{}) *httptest.ResponseRecorder {
	return t.requestJSON(http.MethodPost, url, body, result)
}

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

// TestMain is test entryPoint
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
	if err = service.InitDB(c.DB); err != nil {
		log.Panicf("failed to connect DB. err: %v", err)
	}

	testServer = &TestServer{}
	testServer.Init()
}

func teardown() {
}
