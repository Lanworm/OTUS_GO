package internalhttp

import (
	"github.com/Lanworm/OTUS_GO/hw12_13_14_15_calendar/internal/server/http/httphandler"
)

const baseContentType = "application/json"

func (s *Server) RegisterRoutes(handler *httphandler.Handler) {
	s.AddRoute("/event/create", ContentType(baseContentType, Method("POST", handler.CreateEvent)))
	s.AddRoute("/event", ContentType(baseContentType, Method("GET", handler.GetEvent)))
	s.AddRoute("/event/update", ContentType(baseContentType, Method("PUT", handler.UpdateEvent)))
	s.AddRoute("/event/delete", ContentType(baseContentType, Method("DELETE", handler.DeleteEvent)))
	s.AddRoute("/event/list", ContentType(baseContentType, Method("GET", handler.ListEvent)))

	s.AddRoute("/ready", ContentType(baseContentType, Method("GET", handler.Ready)))
	s.AddRoute("/health", ContentType(baseContentType, Method("GET", handler.Health)))
}
