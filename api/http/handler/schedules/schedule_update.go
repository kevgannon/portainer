package schedules

import (
	"net/http"

	httperror "github.com/portainer/libhttp/error"
	"github.com/portainer/libhttp/request"
	"github.com/portainer/libhttp/response"
	"github.com/portainer/portainer"
)

type payload struct {
	Name      string
	Endpoints []portainer.EndpointID
	Schedule  string
}

func (payload *payload) Validate(r *http.Request) error {
	return nil
}

func (handler *Handler) updateSchedule(w http.ResponseWriter, r *http.Request) *httperror.HandlerError {
	id, err := request.RetrieveNumericRouteVariableValue(r, "id")
	if err != nil {
		return &httperror.HandlerError{http.StatusBadRequest, "Invalid schedule identifier route variable", err}
	}

	var payload payload
	err = request.DecodeAndValidateJSONPayload(r, &payload)
	if err != nil {
		return &httperror.HandlerError{http.StatusBadRequest, "Invalid request payload", err}
	}

	schedule, err := handler.scheduleService.Schedule(portainer.ScheduleID(id))

	if err != nil {
		return &httperror.HandlerError{http.StatusInternalServerError, "Can't find schedule", err}
	}

	if payload.Endpoints != nil {
		schedule.Endpoints = payload.Endpoints
	}

	if payload.Name != "" {
		schedule.Name = payload.Name
	}

	if payload.Schedule != "" {
		schedule.Schedule = payload.Schedule
		handler.scheduler.UpdateScriptJob(schedule.ID, schedule.Schedule)
	}

	handler.scheduleService.UpdateSchedule(portainer.ScheduleID(id), schedule)

	return response.JSON(w, schedule)
}
