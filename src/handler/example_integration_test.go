package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v3"
	"github.com/jhseoeo/fiber-skeleton/src/handler"
	"github.com/jhseoeo/fiber-skeleton/src/middleware"
	"github.com/jhseoeo/fiber-skeleton/src/repository"
	"github.com/jhseoeo/fiber-skeleton/src/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// newIntegrationApp wires the real in-memory repository and service so the
// full handler→service→repository call chain is exercised.
func newIntegrationApp() *fiber.App {
	repo := repository.NewExampleRepository()
	svc := service.NewExampleService(repo)
	h := handler.NewExampleHandler(svc)

	app := fiber.New(fiber.Config{ErrorHandler: middleware.NewErrorHandler()})
	h.RegisterRoutes(app)
	return app
}

func TestIntegration_CreateAndGet(t *testing.T) {
	app := newIntegrationApp()

	// Create
	body, _ := json.Marshal(map[string]string{"content": "integration"})
	req := httptest.NewRequest(http.MethodPost, "/example", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	res, err := app.Test(req, fiber.TestConfig{Timeout: testTimeout})
	require.NoError(t, err)
	assert.Equal(t, http.StatusCreated, res.StatusCode)

	// Get the created item (ID=1 for fresh repository)
	req = httptest.NewRequest(http.MethodGet, "/example/1", nil)
	res, err = app.Test(req, fiber.TestConfig{Timeout: testTimeout})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func TestIntegration_UpdateAndGet(t *testing.T) {
	app := newIntegrationApp()

	// Create
	body, _ := json.Marshal(map[string]string{"content": "original"})
	req := httptest.NewRequest(http.MethodPost, "/example", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	res, err := app.Test(req, fiber.TestConfig{Timeout: testTimeout})
	require.NoError(t, err)
	require.Equal(t, http.StatusCreated, res.StatusCode)

	// Update
	body, _ = json.Marshal(map[string]string{"content": "updated"})
	req = httptest.NewRequest(http.MethodPut, "/example/1", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	res, err = app.Test(req, fiber.TestConfig{Timeout: testTimeout})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)

	// Verify content changed
	req = httptest.NewRequest(http.MethodGet, "/example/1", nil)
	res, err = app.Test(req, fiber.TestConfig{Timeout: testTimeout})
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, res.StatusCode)

	var buf bytes.Buffer
	_, _ = buf.ReadFrom(res.Body)
	var r map[string]any
	require.NoError(t, json.NewDecoder(&buf).Decode(&r))
	data, _ := r["data"].(map[string]any)
	assert.Equal(t, "updated", data["content"])
}

func TestIntegration_DeleteAndNotFound(t *testing.T) {
	app := newIntegrationApp()

	// Create
	body, _ := json.Marshal(map[string]string{"content": "to delete"})
	req := httptest.NewRequest(http.MethodPost, "/example", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	res, err := app.Test(req, fiber.TestConfig{Timeout: testTimeout})
	require.NoError(t, err)
	require.Equal(t, http.StatusCreated, res.StatusCode)

	// Delete
	req = httptest.NewRequest(http.MethodDelete, "/example/1", nil)
	res, err = app.Test(req, fiber.TestConfig{Timeout: testTimeout})
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, res.StatusCode)

	// Should be gone
	req = httptest.NewRequest(http.MethodGet, "/example/1", nil)
	res, err = app.Test(req, fiber.TestConfig{Timeout: testTimeout})
	require.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, res.StatusCode)
}
