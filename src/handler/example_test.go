package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/jhseoeo/fiber-skeleton/src/dto/errorcode"
	"github.com/jhseoeo/fiber-skeleton/src/dto/resp"
	"github.com/jhseoeo/fiber-skeleton/src/handler"
	"github.com/jhseoeo/fiber-skeleton/src/middleware"
	"github.com/jhseoeo/fiber-skeleton/src/model"
	repositoryerror "github.com/jhseoeo/fiber-skeleton/src/repository/error"
	servicemock "github.com/jhseoeo/fiber-skeleton/src/service/mock"
	"github.com/jhseoeo/fiber-skeleton/src/testutil"
)

const testTimeout = 5 * time.Second

func newTestApp(svc *servicemock.MockExampleService) *fiber.App {
	app := fiber.New(fiber.Config{ErrorHandler: middleware.NewErrorHandler()})
	h := handler.NewExampleHandler(svc)
	h.RegisterRoutes(app)
	return app
}

func TestListExample(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		app := newTestApp(&servicemock.MockExampleService{
			ListExamplesFn: func(_ context.Context, page, limit int) ([]*model.Example, int, error) {
				return []*model.Example{
					{ID: 1, Content: "first"},
					{ID: 2, Content: "second"},
				}, 5, nil
			},
		})

		req := httptest.NewRequest(http.MethodGet, "/example?page=1&limit=2", nil)
		res, err := app.Test(req, fiber.TestConfig{Timeout: testTimeout})
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, res.StatusCode)

		r := testutil.DecodeResp[resp.CommonResp](t, res.Body)
		assert.Equal(t, errorcode.Success, r.Code)
		data, ok := r.Data.(map[string]any)
		require.True(t, ok)
		assert.Equal(t, float64(5), data["total"])
		assert.Equal(t, float64(1), data["page"])
		assert.Equal(t, float64(2), data["limit"])
		items, ok := data["data"].([]any)
		require.True(t, ok)
		assert.Len(t, items, 2)
	})

	t.Run("invalid query", func(t *testing.T) {
		app := newTestApp(&servicemock.MockExampleService{})

		req := httptest.NewRequest(http.MethodGet, "/example?page=0&limit=10", nil)
		res, err := app.Test(req, fiber.TestConfig{Timeout: testTimeout})
		require.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, res.StatusCode)
	})

	t.Run("service error", func(t *testing.T) {
		app := newTestApp(&servicemock.MockExampleService{
			ListExamplesFn: func(_ context.Context, page, limit int) ([]*model.Example, int, error) {
				return nil, 0, errors.New("db error")
			},
		})

		req := httptest.NewRequest(http.MethodGet, "/example?page=1&limit=10", nil)
		res, err := app.Test(req, fiber.TestConfig{Timeout: testTimeout})
		require.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
	})
}

func TestGetExample(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		app := newTestApp(&servicemock.MockExampleService{
			GetExampleFn: func(_ context.Context, id uint) (*model.Example, error) {
				return &model.Example{ID: id, Content: "hello"}, nil
			},
		})

		req := httptest.NewRequest(http.MethodGet, "/example/1", nil)
		res, err := app.Test(req, fiber.TestConfig{Timeout: testTimeout})
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, res.StatusCode)
		r := testutil.DecodeResp[resp.CommonResp](t, res.Body)
		assert.Equal(t, errorcode.Success, r.Code)
	})

	t.Run("invalid id", func(t *testing.T) {
		app := newTestApp(&servicemock.MockExampleService{})

		req := httptest.NewRequest(http.MethodGet, "/example/abc", nil)
		res, err := app.Test(req, fiber.TestConfig{Timeout: testTimeout})
		require.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, res.StatusCode)
	})

	t.Run("not found", func(t *testing.T) {
		app := newTestApp(&servicemock.MockExampleService{
			GetExampleFn: func(_ context.Context, id uint) (*model.Example, error) {
				return nil, repositoryerror.ErrNotFound.New("example 1")
			},
		})

		req := httptest.NewRequest(http.MethodGet, "/example/99", nil)
		res, err := app.Test(req, fiber.TestConfig{Timeout: testTimeout})
		require.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, res.StatusCode)
	})
}

func TestCreateExample(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		app := newTestApp(&servicemock.MockExampleService{
			CreateExampleFn: func(_ context.Context, example *model.Example) error {
				example.ID = 1
				return nil
			},
		})

		body, _ := json.Marshal(map[string]string{"content": "hello"})
		req := httptest.NewRequest(http.MethodPost, "/example", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		res, err := app.Test(req, fiber.TestConfig{Timeout: testTimeout})
		require.NoError(t, err)
		require.Equal(t, http.StatusCreated, res.StatusCode)

		r := testutil.DecodeResp[resp.CommonResp](t, res.Body)
		assert.Equal(t, errorcode.Success, r.Code)
		data, ok := r.Data.(map[string]any)
		require.True(t, ok)
		assert.Equal(t, float64(1), data["id"])
		assert.Equal(t, "hello", data["content"])
	})

	t.Run("invalid body", func(t *testing.T) {
		app := newTestApp(&servicemock.MockExampleService{})

		req := httptest.NewRequest(http.MethodPost, "/example", bytes.NewBufferString(`{"content":""}`))
		req.Header.Set("Content-Type", "application/json")
		res, err := app.Test(req, fiber.TestConfig{Timeout: testTimeout})
		require.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, res.StatusCode)
	})
}

func TestUpdateExample(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		app := newTestApp(&servicemock.MockExampleService{
			UpdateExampleFn: func(_ context.Context, example *model.Example) error {
				return nil
			},
		})

		body, _ := json.Marshal(map[string]string{"content": "updated"})
		req := httptest.NewRequest(http.MethodPut, "/example/1", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		res, err := app.Test(req, fiber.TestConfig{Timeout: testTimeout})
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, res.StatusCode)
	})

	t.Run("not found", func(t *testing.T) {
		app := newTestApp(&servicemock.MockExampleService{
			UpdateExampleFn: func(_ context.Context, example *model.Example) error {
				return repositoryerror.ErrNotFound.New("example 1")
			},
		})

		body, _ := json.Marshal(map[string]string{"content": "updated"})
		req := httptest.NewRequest(http.MethodPut, "/example/1", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		res, err := app.Test(req, fiber.TestConfig{Timeout: testTimeout})
		require.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, res.StatusCode)
	})
}

func TestDeleteExample(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		app := newTestApp(&servicemock.MockExampleService{
			DeleteExampleFn: func(_ context.Context, id uint) error {
				return nil
			},
		})

		req := httptest.NewRequest(http.MethodDelete, "/example/1", nil)
		res, err := app.Test(req, fiber.TestConfig{Timeout: testTimeout})
		require.NoError(t, err)
		assert.Equal(t, http.StatusNoContent, res.StatusCode)
	})

	t.Run("not found", func(t *testing.T) {
		app := newTestApp(&servicemock.MockExampleService{
			DeleteExampleFn: func(_ context.Context, id uint) error {
				return repositoryerror.ErrNotFound.New("example 1")
			},
		})

		req := httptest.NewRequest(http.MethodDelete, "/example/99", nil)
		res, err := app.Test(req, fiber.TestConfig{Timeout: testTimeout})
		require.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, res.StatusCode)
	})
}
