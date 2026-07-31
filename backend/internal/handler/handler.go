package handler

type Handler struct {
	Location *LocationHandler
	// User *UserHandler
}

// NewHandlers instanciates and assembles all the handlers
func NewHandler(locationRepo LocationGetter) *Handler {
	return &Handler{
		Location: NewLocationHandler(locationRepo),
		// User: NewUserHandler(userRepo),
	}
}
