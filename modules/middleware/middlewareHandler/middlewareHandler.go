package middlewareHandler

import "github.com/Rayato159/hello-sekai-shop-tutorial/modules/middleware/middlewareRepository"

type (
	MiddlewareHandlerService interface{}

	middlewareHandler struct {
		middlewareRepository middlewareRepository.MiddlewareRepositoryService
	}
)

func NewMiddlewareHandler(middlewareRepository middlewareRepository.MiddlewareRepositoryService) MiddlewareHandlerService {
	return &middlewareHandler{middlewareRepository}
}
