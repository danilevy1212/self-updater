package server

func (a *Application) RegisterRoutes() {
	r := a.Router

	r.GET("/health", a.HealthCheck)
}
