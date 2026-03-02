package handler

import (
	"errors"
	"strconv"

	"github.com/gofiber/fiber/v3"
	"github.com/jhseoeo/fiber-skeleton/src/dto/errorcode"
	"github.com/jhseoeo/fiber-skeleton/src/dto/req"
	"github.com/jhseoeo/fiber-skeleton/src/dto/resp"
	"github.com/jhseoeo/fiber-skeleton/src/model"
	"github.com/jhseoeo/fiber-skeleton/src/pkg/typeerr"
	repositoryerror "github.com/jhseoeo/fiber-skeleton/src/repository/error"
	"github.com/jhseoeo/fiber-skeleton/src/service/serviceport"
)

type ExampleHandler struct {
	exampleService serviceport.ExampleServicePort
}

func NewExampleHandler(exampleService serviceport.ExampleServicePort) *ExampleHandler {
	return &ExampleHandler{
		exampleService: exampleService,
	}
}

func (h *ExampleHandler) RegisterRoutes(router fiber.Router) {
	router.Get("/example/:id", h.get)
	router.Post("/example", h.create)
	router.Put("/example/:id", h.update)
	router.Delete("/example/:id", h.delete)
}

func (h *ExampleHandler) get(c fiber.Ctx) error {
	idUint, err := parseID(c)
	if err != nil {
		return typeerr.NewErrorResp(err, errorcode.ErrInvalidID, "invalid id")
	}
	example, err := h.exampleService.GetExample(c.Context(), idUint)
	if err != nil {
		if errors.Is(err, repositoryerror.ErrNotFound) {
			return typeerr.NewErrorResp(err, errorcode.ErrNotFound, "example not found")
		}
		return typeerr.NewErrorResp(err, errorcode.ErrInternalServer, "failed to get example")
	}
	return c.JSON(resp.CommonResp{
		Code: errorcode.Success,
		Data: resp.GetExampleResp{
			ID:      example.ID,
			Content: example.Content,
		},
	})
}

func (h *ExampleHandler) create(c fiber.Ctx) error {
	var body req.CreateExampleReq
	if err := c.Bind().JSON(&body); err != nil {
		return typeerr.NewErrorResp(err, errorcode.ErrInvalidBody, "invalid request body")
	}
	example := &model.Example{Content: body.Content}
	if err := h.exampleService.CreateExample(c.Context(), example); err != nil {
		if errors.Is(err, repositoryerror.ErrAlreadyExists) {
			return typeerr.NewErrorResp(err, errorcode.ErrConflict, "example already exists")
		}
		return typeerr.NewErrorResp(err, errorcode.ErrInternalServer, "failed to create example")
	}
	return c.Status(fiber.StatusCreated).JSON(resp.CommonResp{
		Code: errorcode.Success,
		Data: resp.CreateExampleResp{
			ID: example.ID,
		},
	})
}

func (h *ExampleHandler) update(c fiber.Ctx) error {
	idUint, err := parseID(c)
	if err != nil {
		return typeerr.NewErrorResp(err, errorcode.ErrInvalidID, "invalid id")
	}
	var body req.UpdateExampleReq
	if err := c.Bind().JSON(&body); err != nil {
		return typeerr.NewErrorResp(err, errorcode.ErrInvalidBody, "invalid request body")
	}
	example := &model.Example{ID: idUint, Content: body.Content}
	if err := h.exampleService.UpdateExample(c.Context(), example); err != nil {
		if errors.Is(err, repositoryerror.ErrNotFound) {
			return typeerr.NewErrorResp(err, errorcode.ErrNotFound, "example not found")
		}
		return typeerr.NewErrorResp(err, errorcode.ErrInternalServer, "failed to update example")
	}
	return c.JSON(resp.CommonResp{
		Code: errorcode.Success,
		Data: resp.GetExampleResp{
			ID:      example.ID,
			Content: example.Content,
		},
	})
}

func (h *ExampleHandler) delete(c fiber.Ctx) error {
	idUint, err := parseID(c)
	if err != nil {
		return typeerr.NewErrorResp(err, errorcode.ErrInvalidID, "invalid id")
	}
	if err := h.exampleService.DeleteExample(c.Context(), idUint); err != nil {
		if errors.Is(err, repositoryerror.ErrNotFound) {
			return typeerr.NewErrorResp(err, errorcode.ErrNotFound, "example not found")
		}
		return typeerr.NewErrorResp(err, errorcode.ErrInternalServer, "failed to delete example")
	}
	return c.SendStatus(fiber.StatusNoContent)
}

func parseID(c fiber.Ctx) (uint, error) {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	return uint(id), err
}
