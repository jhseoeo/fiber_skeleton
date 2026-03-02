package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v3"
	"github.com/jhseoeo/fiber-skeleton/src/dto/errorcode"
	"github.com/jhseoeo/fiber-skeleton/src/dto/resp"
	"github.com/jhseoeo/fiber-skeleton/src/handler"
	"github.com/jhseoeo/fiber-skeleton/src/middleware"
	"github.com/jhseoeo/fiber-skeleton/src/model"
	repositoryerror "github.com/jhseoeo/fiber-skeleton/src/repository/error"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockExampleService is a manual mock of ExampleServicePort.
type mockExampleService struct {
	getExample    func(ctx context.Context, id uint) (*model.Example, error)
	createExample func(ctx context.Context, example *model.Example) error
	updateExample func(ctx context.Context, example *model.Example) error
	deleteExample func(ctx context.Context, id uint) error
}

func (m *mockExampleService) GetExample(ctx context.Context, id uint) (*model.Example, error) {
	return m.getExample(ctx, id)
}
func (m *mockExampleService) CreateExample(ctx context.Context, example *model.Example) error {
	return m.createExample(ctx, example)
}
func (m *mockExampleService) UpdateExample(ctx context.Context, example *model.Example) error {
	return m.updateExample(ctx, example)
}
func (m *mockExampleService) DeleteExample(ctx context.Context, id uint) error {
	return m.deleteExample(ctx, id)
}

func newTestApp(svc *mockExampleService) *fiber.App {
	app := fiber.New(fiber.Config{ErrorHandler: middleware.NewErrorHandler()})
	h := handler.NewExampleHandler(svc)
	h.RegisterRoutes(app)
	return app
}

func decodeResp(t *testing.T, body *bytes.Buffer) resp.CommonResp {
	t.Helper()
	var r resp.CommonResp
	require.NoError(t, json.NewDecoder(body).Decode(&r))
	return r
}

// --- GET /example/:id ---

func TestGetExample_Success(t *testing.T) {
	app := newTestApp(&mockExampleService{
		getExample: func(_ context.Context, id uint) (*model.Example, error) {
			return &model.Example{ID: id, Content: "hello"}, nil
		},
	})

	req := httptest.NewRequest(http.MethodGet, "/example/1", nil)
	res, err := app.Test(req)
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, res.StatusCode)

	var buf bytes.Buffer
	_, _ = buf.ReadFrom(res.Body)
	r := decodeResp(t, &buf)
	assert.Equal(t, errorcode.Success, r.Code)
}

func TestGetExample_InvalidID(t *testing.T) {
	app := newTestApp(&mockExampleService{})

	req := httptest.NewRequest(http.MethodGet, "/example/abc", nil)
	res, err := app.Test(req)
	require.NoError(t, err)

	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
}

func TestGetExample_NotFound(t *testing.T) {
	app := newTestApp(&mockExampleService{
		getExample: func(_ context.Context, id uint) (*model.Example, error) {
			return nil, repositoryerror.ErrNotFound.New("example 1")
		},
	})

	req := httptest.NewRequest(http.MethodGet, "/example/99", nil)
	res, err := app.Test(req)
	require.NoError(t, err)

	assert.Equal(t, http.StatusNotFound, res.StatusCode)
}

// --- POST /example ---

func TestCreateExample_Success(t *testing.T) {
	app := newTestApp(&mockExampleService{
		createExample: func(_ context.Context, example *model.Example) error {
			example.ID = 1
			return nil
		},
	})

	body, _ := json.Marshal(map[string]string{"content": "hello"})
	req := httptest.NewRequest(http.MethodPost, "/example", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	res, err := app.Test(req)
	require.NoError(t, err)

	assert.Equal(t, http.StatusCreated, res.StatusCode)
}

func TestCreateExample_InvalidBody(t *testing.T) {
	app := newTestApp(&mockExampleService{})

	req := httptest.NewRequest(http.MethodPost, "/example", bytes.NewBufferString(`{"content":""}`))
	req.Header.Set("Content-Type", "application/json")
	res, err := app.Test(req)
	require.NoError(t, err)

	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
}

// --- PUT /example/:id ---

func TestUpdateExample_Success(t *testing.T) {
	app := newTestApp(&mockExampleService{
		updateExample: func(_ context.Context, example *model.Example) error {
			return nil
		},
	})

	body, _ := json.Marshal(map[string]string{"content": "updated"})
	req := httptest.NewRequest(http.MethodPut, "/example/1", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	res, err := app.Test(req)
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func TestUpdateExample_NotFound(t *testing.T) {
	app := newTestApp(&mockExampleService{
		updateExample: func(_ context.Context, example *model.Example) error {
			return repositoryerror.ErrNotFound.New("example 1")
		},
	})

	body, _ := json.Marshal(map[string]string{"content": "updated"})
	req := httptest.NewRequest(http.MethodPut, "/example/1", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	res, err := app.Test(req)
	require.NoError(t, err)

	assert.Equal(t, http.StatusNotFound, res.StatusCode)
}

// --- DELETE /example/:id ---

func TestDeleteExample_Success(t *testing.T) {
	app := newTestApp(&mockExampleService{
		deleteExample: func(_ context.Context, id uint) error {
			return nil
		},
	})

	req := httptest.NewRequest(http.MethodDelete, "/example/1", nil)
	res, err := app.Test(req)
	require.NoError(t, err)

	assert.Equal(t, http.StatusNoContent, res.StatusCode)
}

func TestDeleteExample_NotFound(t *testing.T) {
	app := newTestApp(&mockExampleService{
		deleteExample: func(_ context.Context, id uint) error {
			return repositoryerror.ErrNotFound.New("example 1")
		},
	})

	req := httptest.NewRequest(http.MethodDelete, "/example/99", nil)
	res, err := app.Test(req)
	require.NoError(t, err)

	assert.Equal(t, http.StatusNotFound, res.StatusCode)
}
