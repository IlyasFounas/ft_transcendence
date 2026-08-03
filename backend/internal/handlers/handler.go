package handlers

type Handlers struct {
	Location *LocationHandler
	// User *UserHandler
}

// NewHandlers instanciates and assembles all the handlers
func NewHandlers(locationRepo LocationGetter) *Handlers {
	return &Handlers{
		Location: NewLocationHandler(locationRepo),
		// User: NewUserHandler(userRepo),
	}
}
