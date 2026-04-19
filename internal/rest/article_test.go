package rest_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"

	faker "github.com/go-faker/faker/v4"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/josemolinai-nuamx/go-clean-arch/domain"
	"github.com/josemolinai-nuamx/go-clean-arch/internal/rest"
	"github.com/josemolinai-nuamx/go-clean-arch/internal/rest/mocks"
)

func TestFetch(t *testing.T) {
	var mockArticle domain.Article
	err := faker.FakeData(&mockArticle)
	assert.NoError(t, err)
	mockUCase := new(mocks.ArticleService)
	mockListArticle := make([]domain.Article, 0)
	mockListArticle = append(mockListArticle, mockArticle)
	num := 1
	cursor := "2"
	mockUCase.On("Fetch", mock.Anything, cursor, int64(num)).Return(mockListArticle, "10", nil)

	e := echo.New()
	req, err := http.NewRequestWithContext(context.TODO(),
		echo.GET, "/article?num=1&cursor="+cursor, strings.NewReader(""))
	assert.NoError(t, err)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	handler := rest.ArticleHandler{
		Service: mockUCase,
	}
	err = handler.FetchArticle(c)
	require.NoError(t, err)

	responseCursor := rec.Header().Get("X-Cursor")
	assert.Equal(t, "10", responseCursor)
	assert.Equal(t, http.StatusOK, rec.Code)
	mockUCase.AssertExpectations(t)
}

func TestFetchError(t *testing.T) {
	mockUCase := new(mocks.ArticleService)
	num := 1
	cursor := "2"
	mockUCase.On("Fetch", mock.Anything, cursor, int64(num)).Return(nil, "", domain.ErrInternalServerError)

	e := echo.New()
	req, err := http.NewRequestWithContext(context.TODO(), echo.GET, "/article?num=1&cursor="+cursor, strings.NewReader(""))
	assert.NoError(t, err)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	handler := rest.ArticleHandler{
		Service: mockUCase,
	}
	err = handler.FetchArticle(c)
	require.NoError(t, err)

	responseCursor := rec.Header().Get("X-Cursor")
	assert.Equal(t, "", responseCursor)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	mockUCase.AssertExpectations(t)
}

func TestGetByID(t *testing.T) {
	var mockArticle domain.Article
	err := faker.FakeData(&mockArticle)
	assert.NoError(t, err)

	mockUCase := new(mocks.ArticleService)

	num := int(mockArticle.ID)

	mockUCase.On("GetByID", mock.Anything, int64(num)).Return(mockArticle, nil)

	e := echo.New()
	req, err := http.NewRequestWithContext(context.TODO(), echo.GET, "/article/"+strconv.Itoa(num), strings.NewReader(""))
	assert.NoError(t, err)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("article/:id")
	c.SetParamNames("id")
	c.SetParamValues(strconv.Itoa(num))
	handler := rest.ArticleHandler{
		Service: mockUCase,
	}
	err = handler.GetByID(c)
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, rec.Code)
	mockUCase.AssertExpectations(t)
}

func TestStore(t *testing.T) {
	mockArticle := domain.Article{
		Title:     "Title",
		Content:   "Content",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	tempMockArticle := mockArticle
	tempMockArticle.ID = 0
	mockUCase := new(mocks.ArticleService)

	j, err := json.Marshal(tempMockArticle)
	assert.NoError(t, err)

	mockUCase.On("Store", mock.Anything, mock.AnythingOfType("*domain.Article")).Return(nil)

	e := echo.New()
	req, err := http.NewRequestWithContext(context.TODO(), echo.POST, "/article", strings.NewReader(string(j)))
	assert.NoError(t, err)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/article")

	handler := rest.ArticleHandler{
		Service: mockUCase,
	}
	err = handler.Store(c)
	require.NoError(t, err)

	assert.Equal(t, http.StatusCreated, rec.Code)
	mockUCase.AssertExpectations(t)
}

func TestDelete(t *testing.T) {
	var mockArticle domain.Article
	err := faker.FakeData(&mockArticle)
	assert.NoError(t, err)

	mockUCase := new(mocks.ArticleService)

	num := int(mockArticle.ID)

	mockUCase.On("Delete", mock.Anything, int64(num)).Return(nil)

	e := echo.New()
	req, err := http.NewRequestWithContext(context.TODO(), echo.DELETE, "/article/"+strconv.Itoa(num), strings.NewReader(""))
	assert.NoError(t, err)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("article/:id")
	c.SetParamNames("id")
	c.SetParamValues(strconv.Itoa(num))
	handler := rest.ArticleHandler{
		Service: mockUCase,
	}
	err = handler.Delete(c)
	require.NoError(t, err)

	assert.Equal(t, http.StatusNoContent, rec.Code)
	mockUCase.AssertExpectations(t)
}

func TestNewArticleHandler(t *testing.T) {
	e := echo.New()
	mockUCase := new(mocks.ArticleService)

	rest.NewArticleHandler(e, mockUCase)

	routes := e.Routes()
	assert.Len(t, routes, 4)

	got := map[string]bool{}
	for _, route := range routes {
		got[route.Method+" "+route.Path] = true
	}

	assert.True(t, got[http.MethodGet+" /articles"])
	assert.True(t, got[http.MethodPost+" /articles"])
	assert.True(t, got[http.MethodGet+" /articles/:id"])
	assert.True(t, got[http.MethodDelete+" /articles/:id"])
}

func TestFetchWithInvalidNumUsesDefault(t *testing.T) {
	mockUCase := new(mocks.ArticleService)
	mockUCase.On("Fetch", mock.Anything, "2", int64(10)).Return([]domain.Article{}, "", nil)

	e := echo.New()
	req, err := http.NewRequestWithContext(context.TODO(), echo.GET, "/article?num=invalid&cursor=2", strings.NewReader(""))
	assert.NoError(t, err)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	handler := rest.ArticleHandler{Service: mockUCase}

	err = handler.FetchArticle(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	mockUCase.AssertExpectations(t)
}

func TestFetchErrorUnknownMappedToInternalServerError(t *testing.T) {
	mockUCase := new(mocks.ArticleService)
	mockUCase.On("Fetch", mock.Anything, "", int64(10)).Return(nil, "", errors.New("boom"))

	e := echo.New()
	req, err := http.NewRequestWithContext(context.TODO(), echo.GET, "/article", strings.NewReader(""))
	assert.NoError(t, err)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	handler := rest.ArticleHandler{Service: mockUCase}

	err = handler.FetchArticle(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	mockUCase.AssertExpectations(t)
}

func TestGetByIDInvalidID(t *testing.T) {
	mockUCase := new(mocks.ArticleService)
	e := echo.New()
	req, err := http.NewRequestWithContext(context.TODO(), echo.GET, "/article/not-a-number", strings.NewReader(""))
	assert.NoError(t, err)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("article/:id")
	c.SetParamNames("id")
	c.SetParamValues("not-a-number")
	handler := rest.ArticleHandler{Service: mockUCase}

	err = handler.GetByID(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, rec.Code)
	mockUCase.AssertNotCalled(t, "GetByID", mock.Anything, mock.Anything)
}

func TestGetByIDNotFound(t *testing.T) {
	mockUCase := new(mocks.ArticleService)
	mockUCase.On("GetByID", mock.Anything, int64(99)).Return(domain.Article{}, domain.ErrNotFound)

	e := echo.New()
	req, err := http.NewRequestWithContext(context.TODO(), echo.GET, "/article/99", strings.NewReader(""))
	assert.NoError(t, err)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("article/:id")
	c.SetParamNames("id")
	c.SetParamValues("99")
	handler := rest.ArticleHandler{Service: mockUCase}

	err = handler.GetByID(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, rec.Code)
	mockUCase.AssertExpectations(t)
}

func TestStoreInvalidJSON(t *testing.T) {
	mockUCase := new(mocks.ArticleService)
	e := echo.New()
	req, err := http.NewRequestWithContext(context.TODO(), echo.POST, "/article", strings.NewReader("{invalid-json"))
	assert.NoError(t, err)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	handler := rest.ArticleHandler{Service: mockUCase}

	err = handler.Store(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
	mockUCase.AssertNotCalled(t, "Store", mock.Anything, mock.Anything)
}

func TestStoreValidationError(t *testing.T) {
	mockUCase := new(mocks.ArticleService)
	e := echo.New()
	req, err := http.NewRequestWithContext(context.TODO(), echo.POST, "/article", strings.NewReader(`{"title":""}`))
	assert.NoError(t, err)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	handler := rest.ArticleHandler{Service: mockUCase}

	err = handler.Store(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	mockUCase.AssertNotCalled(t, "Store", mock.Anything, mock.Anything)
}

func TestStoreConflict(t *testing.T) {
	mockUCase := new(mocks.ArticleService)
	mockUCase.On("Store", mock.Anything, mock.AnythingOfType("*domain.Article")).Return(domain.ErrConflict)

	e := echo.New()
	req, err := http.NewRequestWithContext(context.TODO(), echo.POST, "/article", strings.NewReader(`{"title":"t","content":"c"}`))
	assert.NoError(t, err)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	handler := rest.ArticleHandler{Service: mockUCase}

	err = handler.Store(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusConflict, rec.Code)
	mockUCase.AssertExpectations(t)
}

func TestDeleteInvalidID(t *testing.T) {
	mockUCase := new(mocks.ArticleService)
	e := echo.New()
	req, err := http.NewRequestWithContext(context.TODO(), echo.DELETE, "/article/not-a-number", strings.NewReader(""))
	assert.NoError(t, err)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("article/:id")
	c.SetParamNames("id")
	c.SetParamValues("not-a-number")
	handler := rest.ArticleHandler{Service: mockUCase}

	err = handler.Delete(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, rec.Code)
	mockUCase.AssertNotCalled(t, "Delete", mock.Anything, mock.Anything)
}

func TestDeleteNotFound(t *testing.T) {
	mockUCase := new(mocks.ArticleService)
	mockUCase.On("Delete", mock.Anything, int64(77)).Return(domain.ErrNotFound)

	e := echo.New()
	req, err := http.NewRequestWithContext(context.TODO(), echo.DELETE, "/article/77", strings.NewReader(""))
	assert.NoError(t, err)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("article/:id")
	c.SetParamNames("id")
	c.SetParamValues("77")
	handler := rest.ArticleHandler{Service: mockUCase}

	err = handler.Delete(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, rec.Code)
	mockUCase.AssertExpectations(t)
}

func TestGetByIDInternalServerError(t *testing.T) {
	mockUCase := new(mocks.ArticleService)
	mockUCase.On("GetByID", mock.Anything, int64(13)).Return(domain.Article{}, domain.ErrInternalServerError)

	e := echo.New()
	req, err := http.NewRequestWithContext(context.TODO(), echo.GET, "/article/13", strings.NewReader(""))
	assert.NoError(t, err)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("article/:id")
	c.SetParamNames("id")
	c.SetParamValues("13")
	handler := rest.ArticleHandler{Service: mockUCase}

	err = handler.GetByID(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	mockUCase.AssertExpectations(t)
}

func TestDeleteInternalServerError(t *testing.T) {
	mockUCase := new(mocks.ArticleService)
	mockUCase.On("Delete", mock.Anything, int64(14)).Return(domain.ErrInternalServerError)

	e := echo.New()
	req, err := http.NewRequestWithContext(context.TODO(), echo.DELETE, "/article/14", strings.NewReader(""))
	assert.NoError(t, err)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("article/:id")
	c.SetParamNames("id")
	c.SetParamValues("14")
	handler := rest.ArticleHandler{Service: mockUCase}

	err = handler.Delete(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	mockUCase.AssertExpectations(t)
}
