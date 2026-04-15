package resource

import (
	"aIBuildService/aPI/service"
	"net/http"
)

type WelcomeResource struct {
	WelcomeService service.WelcomeService
}

func NewWelcomeResource(welcomeService service.WelcomeService) *WelcomeResource {
	return &WelcomeResource{
		WelcomeService: welcomeService,
	}
}

func (app *WelcomeResource) WelcomeGetHandler(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("Welcome Resource"))
	if err != nil {
		return
	}
}
