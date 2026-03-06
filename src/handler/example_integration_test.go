package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/jhseoeo/fiber-skeleton/src/handler"
	"github.com/jhseoeo/fiber-skeleton/src/middleware"
	"github.com/jhseoeo/fiber-skeleton/src/repository"
	"github.com/jhseoeo/fiber-skeleton/src/service"
	"github.com/jhseoeo/fiber-skeleton/src/testutil"
)

func newIntegrationApp() *fiber.App {
	repo := repository.NewExampleRepository()
	svc := service.NewExampleService(repo)
	h := handler.NewExampleHandler(svc)

	app := fiber.New(fiber.Config{ErrorHandler: middleware.NewErrorHandler()})
	h.RegisterRoutes(app)
	return app
}

func TestIntegrationCreate(t *testing.T) {
	t.Run("creates item and returns it on get", func(t *testing.T) {
		app := newIntegrationApp()

		body, _ := json.Marshal(map[string]string{"content": "integration"})
		req := httptest.NewRequest(http.MethodPost, "/example", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		res, err := app.Test(req, fiber.TestConfig{Timeout: testTimeout})
		require.NoError(t, err)
		assert.Equal(t, http.StatusCreated, res.StatusCode)

		req = httptest.NewRequest(http.MethodGet, "/example/1", nil)
		res, err = app.Test(req, fiber.TestConfig{Timeout: testTimeout})
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, res.StatusCode)
	})
}

func TestIntegrationUpdate(t *testing.T) {
	t.Run("updates item content", func(t *testing.T) {
		app := newIntegrationApp()

		body, _ := json.Marshal(map[string]string{"content": "original"})
		req := httptest.NewRequest(http.MethodPost, "/example", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		res, err := app.Test(req, fiber.TestConfig{Timeout: testTimeout})
		require.NoError(t, err)
		require.Equal(t, http.StatusCreated, res.StatusCode)

		body, _ = json.Marshal(map[string]string{"content": "updated"})
		req = httptest.NewRequest(http.MethodPut, "/example/1", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		res, err = app.Test(req, fiber.TestConfig{Timeout: testTimeout})
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, res.StatusCode)

		req = httptest.NewRequest(http.MethodGet, "/example/1", nil)
		res, err = app.Test(req, fiber.TestConfig{Timeout: testTimeout})
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, res.StatusCode)
		r := testutil.DecodeResp[map[string]any](t, res.Body)
		data, _ := r["data"].(map[string]any)
		assert.Equal(t, "updated", data["content"])
	})
}

func TestIntegrationDelete(t *testing.T) {
	t.Run("deletes item and returns 404 on subsequent get", func(t *testing.T) {
		app := newIntegrationApp()

		body, _ := json.Marshal(map[string]string{"content": "to delete"})
		req := httptest.NewRequest(http.MethodPost, "/example", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		res, err := app.Test(req, fiber.TestConfig{Timeout: testTimeout})
		require.NoError(t, err)
		require.Equal(t, http.StatusCreated, res.StatusCode)

		req = httptest.NewRequest(http.MethodDelete, "/example/1", nil)
		res, err = app.Test(req, fiber.TestConfig{Timeout: testTimeout})
		require.NoError(t, err)
		assert.Equal(t, http.StatusNoContent, res.StatusCode)

		req = httptest.NewRequest(http.MethodGet, "/example/1", nil)
		res, err = app.Test(req, fiber.TestConfig{Timeout: testTimeout})
		require.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, res.StatusCode)
	})
}

func TestIntegrationList(t *testing.T) {
	t.Run("empty repository", func(t *testing.T) {
		app := newIntegrationApp()

		req := httptest.NewRequest(http.MethodGet, "/example?page=1&limit=10", nil)
		res, err := app.Test(req, fiber.TestConfig{Timeout: testTimeout})
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, res.StatusCode)

		r := testutil.DecodeResp[map[string]any](t, res.Body)
		data, _ := r["data"].(map[string]any)
		assert.Equal(t, float64(0), data["total"])
	})

	t.Run("pagination", func(t *testing.T) {
		app := newIntegrationApp()

		for _, content := range []string{"alpha", "beta", "gamma"} {
			body, _ := json.Marshal(map[string]string{"content": content})
			req := httptest.NewRequest(http.MethodPost, "/example", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			res, err := app.Test(req, fiber.TestConfig{Timeout: testTimeout})
			require.NoError(t, err)
			require.Equal(t, http.StatusCreated, res.StatusCode)
		}

		req := httptest.NewRequest(http.MethodGet, "/example?page=1&limit=2", nil)
		res, err := app.Test(req, fiber.TestConfig{Timeout: testTimeout})
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, res.StatusCode)
		r := testutil.DecodeResp[map[string]any](t, res.Body)
		data, _ := r["data"].(map[string]any)
		assert.Equal(t, float64(3), data["total"])
		assert.Equal(t, float64(1), data["page"])
		items, _ := data["data"].([]any)
		assert.Len(t, items, 2)

		req = httptest.NewRequest(http.MethodGet, "/example?page=2&limit=2", nil)
		res, err = app.Test(req, fiber.TestConfig{Timeout: testTimeout})
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, res.StatusCode)
		r = testutil.DecodeResp[map[string]any](t, res.Body)
		data, _ = r["data"].(map[string]any)
		items, _ = data["data"].([]any)
		assert.Len(t, items, 1)
	})

	t.Run("invalid query", func(t *testing.T) {
		app := newIntegrationApp()

		req := httptest.NewRequest(http.MethodGet, "/example?limit=10", nil)
		res, err := app.Test(req, fiber.TestConfig{Timeout: testTimeout})
		require.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, res.StatusCode)

		req = httptest.NewRequest(http.MethodGet, "/example?page=1&limit=0", nil)
		res, err = app.Test(req, fiber.TestConfig{Timeout: testTimeout})
		require.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, res.StatusCode)
	})

	t.Run("max page boundary", func(t *testing.T) {
		app := newIntegrationApp()

		req := httptest.NewRequest(http.MethodGet, "/example?page=10000&limit=100", nil)
		res, err := app.Test(req, fiber.TestConfig{Timeout: testTimeout})
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, res.StatusCode)
		r := testutil.DecodeResp[map[string]any](t, res.Body)
		data, _ := r["data"].(map[string]any)
		assert.Equal(t, float64(0), data["total"])
		items, _ := data["data"].([]any)
		assert.Empty(t, items)

		req = httptest.NewRequest(http.MethodGet, "/example?page=10001&limit=10", nil)
		res, err = app.Test(req, fiber.TestConfig{Timeout: testTimeout})
		require.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, res.StatusCode)
	})
}
