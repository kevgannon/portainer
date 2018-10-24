package schedules

import (
	"net/http"

	"github.com/gorilla/mux"
	httperror "github.com/portainer/libhttp/error"
	"github.com/portainer/portainer"
	"github.com/portainer/portainer/http/security"
)

var mockSchedules []*portainer.Schedule = []*portainer.Schedule{
	&portainer.Schedule{
		ID:         0,
		Schedule:   "every day",
		Endpoints:  []portainer.EndpointID{1},
		ScriptPath: "/tmp/scripts",
		Name:       "First schedule",
	},
	&portainer.Schedule{
		ID:         1,
		Schedule:   "every monday",
		Endpoints:  []portainer.EndpointID{0},
		ScriptPath: "/tmp/scripts",
		Name:       "Second schedule",
	},
}

type Handler struct {
	*mux.Router
	scheduleService portainer.ScheduleService
	fileService     portainer.FileService
	scheduler       portainer.JobScheduler
}

func NewHandler(bouncer *security.RequestBouncer, scheduleService portainer.ScheduleService, fileService portainer.FileService, scheduler portainer.JobScheduler) *Handler {
	h := &Handler{
		Router:          mux.NewRouter(),
		scheduleService: scheduleService,
		fileService:     fileService,
		scheduler:       scheduler,
	}

	h.Handle("/schedules",
		bouncer.AuthenticatedAccess(httperror.LoggerHandler(h.listSchedules))).Methods(http.MethodGet)

	h.Handle("/schedules",
		bouncer.AuthenticatedAccess(httperror.LoggerHandler(h.createScheduleHandler))).Methods(http.MethodPost)

	h.Handle("/schedules/{id}",
		bouncer.AuthenticatedAccess(httperror.LoggerHandler(h.inspectSchedule))).Methods(http.MethodGet)

	h.Handle("/schedules/{id}",
		bouncer.AdministratorAccess(httperror.LoggerHandler(h.updateSchedule))).Methods(http.MethodPut)

	h.Handle("/schedules/{id}",
		bouncer.AuthenticatedAccess(httperror.LoggerHandler(h.deleteSchedule))).Methods(http.MethodDelete)

	return h
}
