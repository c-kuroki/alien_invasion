package http

import (
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"

	"github.com/c-kuroki/alien_invasion/pkg/app"
)

type HTTPService struct {
	invasion       *app.AlienInvasionApp
	serviceAddress string
}

func NewHTTPService(invasion *app.AlienInvasionApp, serviceAddress string) *HTTPService {
	return &HTTPService{
		invasion:       invasion,
		serviceAddress: serviceAddress,
	}
}

func (srv *HTTPService) Start() {
	r := chi.NewRouter()

	//r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	// public endpoint
	r.Get("/", srv.GetIndex)
	r.Get("/map", srv.GetMap)

	log.Printf("Starting http server at %s\n", srv.serviceAddress)
	httpServer := &http.Server{Addr: srv.serviceAddress, Handler: r}
	err := httpServer.ListenAndServe()
	if err != nil {
		panic(err)
	}
}

func (srv *HTTPService) GetIndex(w http.ResponseWriter, r *http.Request) {
	render.Status(r, http.StatusOK)
	render.HTML(w, r, `
<html>
<head>
<meta http-equiv="refresh" content="1" />
</head>
<body>
<img src="/map" />
</body>
</html>
`)
}

func (srv *HTTPService) GetMap(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "image/svg+xml")
	err := srv.invasion.RenderMap(r.Context(), w)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, err.Error())
		return
	}
	render.Status(r, http.StatusOK)
}
