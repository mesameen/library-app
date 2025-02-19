package handlertest

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/test/library-app/internal/config"
	"github.com/test/library-app/internal/handler"
	"github.com/test/library-app/internal/model"
	"github.com/test/library-app/internal/store"
)

var reqHandler *handler.Handler

func TestMain(m *testing.M) {
	// loading configuration
	config.LoadConfig()
	// updating store type to local for unit testing
	config.CommonConfig.StoreType = "local"
	// initializing the store
	store, err := store.NewStore()
	assert.Nil(&testing.T{}, err)
	// initializing the handler
	reqHandler = handler.NewHandler(store)
	m.Run()
}

// initiazes the gin context
func GetTestGinContext(w *httptest.ResponseRecorder) *gin.Context {
	gin.SetMode(gin.TestMode)
	// creating gin test context by passing the response writer
	c, _ := gin.CreateTestContext(w)
	// adding headers if any
	c.Request = &http.Request{
		Header: make(http.Header),
	}
	return c
}

func TestGetBook(t *testing.T) {
	// creating response writer by calling httptest recorder
	w := httptest.NewRecorder()
	c := GetTestGinContext(w)
	// adding path params to the request
	params := []gin.Param{
		{
			Key:   "title",
			Value: "alchemist",
		},
	}
	c.Params = params

	reqHandler.GetBook(c)
	assert.EqualValues(t, http.StatusOK, w.Code)

	// failure case
	w = httptest.NewRecorder()
	c = GetTestGinContext(w)
	params = []gin.Param{
		{
			Key:   "title",
			Value: "book_200",
		},
	}
	c.Params = params

	reqHandler.GetBook(c)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestBorrowBook(t *testing.T) {
	// creating response writer by calling httptest recorder
	w := httptest.NewRecorder()
	c := GetTestGinContext(w)
	req := model.LoanRequest{
		NameOfBorrower: "test_user",
		Title:          "alchemist",
	}
	reqBytes, _ := json.Marshal(&req)
	c.Request.Body = io.NopCloser(bytes.NewBuffer(reqBytes))
	defer c.Request.Body.Close()

	reqHandler.LoanBook(c)
	assert.EqualValues(t, http.StatusCreated, w.Code)

	// failure case
	w = httptest.NewRecorder()
	c = GetTestGinContext(w)
	req = model.LoanRequest{
		NameOfBorrower: "test_user",
		Title:          "book_100",
	}
	reqBytes, _ = json.Marshal(&req)
	c.Request.Body = io.NopCloser(bytes.NewBuffer(reqBytes))
	defer c.Request.Body.Close()
	reqHandler.LoanBook(c)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestExtendLoan(t *testing.T) {
	w := httptest.NewRecorder()
	c := GetTestGinContext(w)
	req := model.LoanRequest{
		NameOfBorrower: "test_user",
		Title:          "alchemist",
	}
	reqBytes, _ := json.Marshal(&req)
	c.Request.Body = io.NopCloser(bytes.NewBuffer(reqBytes))
	defer c.Request.Body.Close()

	reqHandler.LoanBook(c)
	assert.EqualValues(t, http.StatusCreated, w.Code)

	// creating response writer by calling httptest recorder
	w = httptest.NewRecorder()
	c = GetTestGinContext(w)
	params := []gin.Param{
		{
			Key:   "id",
			Value: "1",
		},
	}
	c.Params = params

	reqHandler.ExtendLoan(c)
	assert.EqualValues(t, http.StatusAccepted, w.Code)

	// failure case
	// creating response writer by calling httptest recorder
	w = httptest.NewRecorder()
	c = GetTestGinContext(w)
	params = []gin.Param{
		{
			Key:   "id",
			Value: "100",
		},
	}
	c.Params = params

	reqHandler.ExtendLoan(c)
	assert.EqualValues(t, http.StatusNotFound, w.Code)
}

func TestReturnBook(t *testing.T) {
	w := httptest.NewRecorder()
	c := GetTestGinContext(w)
	req := model.LoanRequest{
		NameOfBorrower: "test_user",
		Title:          "alchemist",
	}
	reqBytes, _ := json.Marshal(&req)
	c.Request.Body = io.NopCloser(bytes.NewBuffer(reqBytes))
	defer c.Request.Body.Close()

	reqHandler.LoanBook(c)
	assert.EqualValues(t, http.StatusCreated, w.Code)

	// creating response writer by calling httptest recorder
	w = httptest.NewRecorder()
	c = GetTestGinContext(w)
	params := []gin.Param{
		{
			Key:   "id",
			Value: "1",
		},
	}
	c.Params = params

	reqHandler.ReturnBook(c)
	assert.EqualValues(t, http.StatusAccepted, w.Code)

	// failure case
	// creating response writer by calling httptest recorder
	w = httptest.NewRecorder()
	c = GetTestGinContext(w)
	params = []gin.Param{
		{
			Key:   "id",
			Value: "100",
		},
	}
	c.Params = params

	reqHandler.ReturnBook(c)
	assert.EqualValues(t, http.StatusNotFound, w.Code)
}
