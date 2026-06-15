package http

type Handler struct{}

type Dependencies struct{}

func New(deps Dependencies) *Handler {
	return &Handler{}
}
