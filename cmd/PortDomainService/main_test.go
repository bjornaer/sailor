package main_test

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	p "github.com/bjornaer/sailor/cmd/PortDomainService"
	"github.com/bjornaer/sailor/internal/db"
	"github.com/bjornaer/sailor/internal/port"
	sm "github.com/bjornaer/sailor/internal/sessionmanager"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type UnitTestSuite struct {
	suite.Suite
	router   *gin.Engine
	tmpfile  *os.File
	dbClient db.DBClient
}

func (s *UnitTestSuite) SetupTest() {
	content := []byte(`{
		"PORTKEY": {
			"name": "test",
			"city": "test",
			"country": "test",
			"alias": [],
			"regions": [],
			"coordinates": [
			  55.5136433,
			  25.4052165
			],
			"province": "test",
			"timezone": "test",
			"unlocs": [
			  "test"
			],
			"code": "00000"
		  }
	}`)
	tmpfile, err := ioutil.TempFile("", "test_file")
	if err != nil {
		log.Fatal(err)
	}

	// defer os.Remove(tmpfile.Name()) // clean up

	if _, err := tmpfile.Write(content); err != nil {
		log.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		log.Fatal(err)
	}
	s.dbClient, _ = db.InitDBClient()
	s.tmpfile = tmpfile
}

func (s *UnitTestSuite) BeforeTest(suiteName, testName string) {
	tmpfile := s.tmpfile
	dbClient := s.dbClient
	sm := &sm.SessionManager{UpdatesFile: tmpfile.Name(), DBClient: dbClient}
	s.router = p.SetupRouter(sm)
}

func (s *UnitTestSuite) AfterTest(suiteName, testName string) {
	os.Remove(s.tmpfile.Name())
}

func (s *UnitTestSuite) Test_TableTest() {

	type testCase struct {
		name             string
		endpoint         string
		requestMethod    string
		expectedResponse interface{}
		payload          string
	}

	testCases := []testCase{
		{
			name:             "Hello Endpoint",
			endpoint:         "/",
			requestMethod:    http.MethodGet,
			payload:          "",
			expectedResponse: "Hello Sailor! Welcome to the Port Domain Service!",
		},
		{
			name:             "Port Processing Endpoint",
			endpoint:         "/process",
			requestMethod:    http.MethodGet,
			payload:          "",
			expectedResponse: "Finished updating DB with ports data!",
		},
		{
			name:             "Port GET",
			endpoint:         "/port/PORTKEY",
			requestMethod:    http.MethodGet,
			payload:          "",
			expectedResponse: "00000",
		},
		{
			name:     "Port POST",
			endpoint: "/port",
			payload: `{
				"name": "test",
				"city": "test",
				"country": "test",
				"alias": [],
				"regions": [],
				"coordinates": [
				  55.5136433,
				  25.4052165
				],
				"province": "test",
				"timezone": "test",
				"unlocs": [
				  "PORTKEY"
				],
				"code": "00000"
			  }`,
			requestMethod:    http.MethodPost,
			expectedResponse: "Port Created/Updated successfully",
		},
	}

	for _, testCase := range testCases {

		s.Run(testCase.name, func() {
			router := s.router
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(testCase.requestMethod, testCase.endpoint, strings.NewReader(testCase.payload))
			router.ServeHTTP(w, req)

			assert.Equal(s.T(), http.StatusOK, w.Code)
			if testCase.name == "Port GET" {
				var p port.PortData
				json.NewDecoder(w.Body).Decode(&p)
				assert.Equal(s.T(), testCase.expectedResponse, p.Code)
			} else {
				assert.Equal(s.T(), testCase.expectedResponse, w.Body.String())
			}
		})
	}
}

func TestUnitTestSuite(t *testing.T) {
	suite.Run(t, new(UnitTestSuite))
}
