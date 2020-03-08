package controller

import (
	"bytes"
	"encoding/json"
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
)

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
}

func teardown() {

}

// makeGet makes a GET request
func makeGet(url string) (req *http.Request, write *httptest.ResponseRecorder) {
	return makeRequest(http.MethodGet, url)
}

// makePostJSON makes a POST request with JSON content-type
func makePostJSON(url string, body interface{}) (req *http.Request, write *httptest.ResponseRecorder, err error) {
	return makeRequestJSON(http.MethodPost, url, body)
}

// makeDelete makes a request
func makeRequest(method, url string) (req *http.Request, write *httptest.ResponseRecorder) {
	r := httptest.NewRequest(method, url, nil)
	w := httptest.NewRecorder()
	return r, w
}

// makeRequestJSON makes a request with JSON content-type
func makeRequestJSON(method, url string, body interface{}) (req *http.Request, write *httptest.ResponseRecorder, err error) {
	var reader io.Reader

	switch body.(type) {
	case string:
		reader = strings.NewReader(body.(string))
		break
	default:
		jsonValue, e := json.Marshal(body)
		if e != nil {
			return nil, nil, e
		}
		reader = bytes.NewBuffer(jsonValue)
		break
	}

	r := httptest.NewRequest(method, url, reader)
	r.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	return r, w, nil
}

func newRouter() *gin.Engine {
	r := gin.New()

	return r
}

func parseJSON(resp *httptest.ResponseRecorder, v interface{}) error {
	return json.NewDecoder(resp.Body).Decode(v)
}

func isIntegrationTest() bool {
	return os.Getenv("INTEGRATION_TEST") == "true"
}
