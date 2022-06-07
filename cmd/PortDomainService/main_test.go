package main_test

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	p "github.com/bjornaer/sailor/cmd/PortDomainService"
	"github.com/bjornaer/sailor/internal/db"
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
	s.dbClient = db.InitDBClient()
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
		expectedResponse interface{}
	}

	testCases := []testCase{
		{
			name:             "Hello Endpoint",
			endpoint:         "/",
			expectedResponse: "Hello Sailor! Welcome to the Port Domain Service!",
		},
		{
			name:             "Port Processing Endpoint",
			endpoint:         "/process",
			expectedResponse: "Finished updating DB with ports data!",
		},
	}

	for _, testCase := range testCases {

		s.Run(testCase.name, func() {
			router := s.router
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, testCase.endpoint, nil)
			router.ServeHTTP(w, req)

			assert.Equal(s.T(), http.StatusOK, w.Code)
			assert.Equal(s.T(), testCase.expectedResponse, w.Body.String())
		})
	}
}

func TestUnitTestSuite(t *testing.T) {
	suite.Run(t, new(UnitTestSuite))
}
