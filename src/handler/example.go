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
	"github.com/jhseoeo/fiber-skeleton/src/pkg/validate"
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
	router.Get("/example", h.list)
	router.Get("/example/:id", h.get)
	router.Post("/example", h.create)
	router.Put("/example/:id", h.update)
	router.Delete("/example/:id", h.delete)
}

// list godoc
//
//	@Summary		List examples
//	@Description	Returns a paginated list of example items
//	@Tags			example
//	@Produce		json
//	@Param			page	query		int					true	"Page number (>=1)"
//	@Param			limit	query		int					true	"Items per page (1-100)"
//	@Success		200		{object}	resp.PaginatedResp	"ok"
//	@Failure		400		{object}	resp.CommonResp		"invalid query"
//	@Failure		500		{object}	resp.CommonResp		"internal error"
//	@Router			/example [get]
func (h *ExampleHandler) list(c fiber.Ctx) error {
	var query req.PaginationReq
	if err := bindQuery(c, &query); err != nil {
		return err
	}
	examples, total, err := h.exampleService.ListExamples(c.Context(), query.Page, query.Limit)
	if err != nil {
		return typeerr.NewErrorResp(err, errorcode.ErrInternalServer, "failed to list examples")
	}

	items := make([]resp.GetExampleResp, 0, len(examples))
	for _, e := range examples {
		items = append(items, resp.GetExampleResp{ID: e.ID, Content: e.Content})
	}
	return c.JSON(resp.CommonResp{
		Code: errorcode.Success,
		Data: resp.PaginatedResp{
			Total: total,
			Page:  query.Page,
			Limit: query.Limit,
			Data:  items,
		},
	})
}

// get godoc
//
//	@Summary		Get example
//	@Description	Get an example item by ID
//	@Tags			example
//	@Produce		json
//	@Param			id	path		int				true	"Example ID"
//	@Success		200	{object}	resp.CommonResp	"ok"
//	@Failure		400	{object}	resp.CommonResp	"invalid id"
//	@Failure		404	{object}	resp.CommonResp	"not found"
//	@Failure		500	{object}	resp.CommonResp	"internal error"
//	@Router			/example/{id} [get]
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

// create godoc
//
//	@Summary		Create example
//	@Description	Create a new example item
//	@Tags			example
//	@Accept			json
//	@Produce		json
//	@Param			body	body		req.CreateExampleReq	true	"Request body"
//	@Success		201		{object}	resp.CommonResp			"created"
//	@Failure		400		{object}	resp.CommonResp			"invalid body"
//	@Failure		409		{object}	resp.CommonResp			"already exists"
//	@Failure		500		{object}	resp.CommonResp			"internal error"
//	@Router			/example [post]
func (h *ExampleHandler) create(c fiber.Ctx) error {
	var body req.CreateExampleReq
	if err := bindJSON(c, &body); err != nil {
		return err
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
		Data: resp.GetExampleResp{
			ID:      example.ID,
			Content: example.Content,
		},
	})
}

// update godoc
//
//	@Summary		Update example
//	@Description	Update an existing example item by ID
//	@Tags			example
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int						true	"Example ID"
//	@Param			body	body		req.UpdateExampleReq	true	"Request body"
//	@Success		200		{object}	resp.CommonResp			"ok"
//	@Failure		400		{object}	resp.CommonResp			"invalid id or body"
//	@Failure		404		{object}	resp.CommonResp			"not found"
//	@Failure		500		{object}	resp.CommonResp			"internal error"
//	@Router			/example/{id} [put]
func (h *ExampleHandler) update(c fiber.Ctx) error {
	idUint, err := parseID(c)
	if err != nil {
		return typeerr.NewErrorResp(err, errorcode.ErrInvalidID, "invalid id")
	}
	var body req.UpdateExampleReq
	if err := bindJSON(c, &body); err != nil {
		return err
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

// delete godoc
//
//	@Summary		Delete example
//	@Description	Delete an example item by ID
//	@Tags			example
//	@Produce		json
//	@Param			id	path		int				true	"Example ID"
//	@Success		204	"no content"
//	@Failure		400	{object}	resp.CommonResp	"invalid id"
//	@Failure		404	{object}	resp.CommonResp	"not found"
//	@Failure		500	{object}	resp.CommonResp	"internal error"
//	@Router			/example/{id} [delete]
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

// bindJSON binds a JSON request body into dst and validates it.
// Returns a ready-to-return ErrorResp on failure, or nil on success.
func bindJSON(c fiber.Ctx, dst any) error {
	if err := c.Bind().JSON(dst); err != nil {
		return typeerr.NewErrorResp(err, errorcode.ErrInvalidBody, "invalid request body")
	}
	if err := validate.Struct(dst); err != nil {
		return typeerr.NewErrorRespWithData(err, errorcode.ErrInvalidBody, "validation failed", err)
	}
	return nil
}

// bindQuery binds query parameters into dst and validates them.
// Returns a ready-to-return ErrorResp on failure, or nil on success.
func bindQuery(c fiber.Ctx, dst any) error {
	if err := c.Bind().Query(dst); err != nil {
		return typeerr.NewErrorResp(err, errorcode.ErrBadRequest, "invalid query parameters")
	}
	if err := validate.Struct(dst); err != nil {
		return typeerr.NewErrorRespWithData(err, errorcode.ErrBadRequest, "validation failed", err)
	}
	return nil
}
