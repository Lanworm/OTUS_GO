package grpchandler

import (
	"context"
	"errors"
	"fmt"

	"github.com/Lanworm/OTUS_GO/hw12_13_14_15_calendar/internal/enum"
	"github.com/Lanworm/OTUS_GO/hw12_13_14_15_calendar/internal/logger"
	"github.com/Lanworm/OTUS_GO/hw12_13_14_15_calendar/internal/service"
	"github.com/Lanworm/OTUS_GO/hw12_13_14_15_calendar/internal/storage"
	grpc_calendar "github.com/Lanworm/OTUS_GO/hw12_13_14_15_calendar/pb/grpc"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Handler struct {
	grpc_calendar.UnimplementedCalendarServer
	logger  *logger.Logger
	service *service.Event
}

func NewHandler(
	logger *logger.Logger,
	service *service.Event,
) *Handler {
	return &Handler{
		logger:  logger,
		service: service,
	}
}

func (h *Handler) CreateEvent(
	_ context.Context,
	req *grpc_calendar.CreateEvent_Request,
) (*grpc_calendar.CreateEvent_Reply, error) {
	err := req.GetStartDatetime().CheckValid()
	if err != nil {
		return nil, err
	}

	err = req.GetEndDatetime().CheckValid()
	if err != nil {
		return nil, err
	}

	res, err := h.service.CreateEvent(&storage.Event{
		Title:         req.GetTitle(),
		StartDatetime: req.GetStartDatetime().AsTime(),
		EndDatetime:   req.GetEndDatetime().AsTime(),
		Description:   req.GetDescription(),
		UserID:        uuid.NewString(), // TODO !!!
		RemindBefore:  req.GetRemindBefore(),
	})
	if err != nil {
		return nil, err
	}

	return &grpc_calendar.CreateEvent_Reply{
		Id: res,
	}, nil
}

func (h *Handler) GetEvent(
	_ context.Context,
	req *grpc_calendar.GetEvent_Request,
) (*grpc_calendar.GetEvent_Reply, error) {
	id := req.GetId()
	if id == "" {
		return nil, errors.New("parameter id: is required")
	}

	evt, err := h.service.GetEvent(id)
	if err != nil {
		return nil, err
	}

	return &grpc_calendar.GetEvent_Reply{
		Result: &grpc_calendar.Domain_Event{
			Id:            evt.ID,
			Title:         evt.Title,
			Description:   evt.Description,
			UserId:        evt.UserID,
			RemindBefore:  evt.RemindBefore,
			StartDatetime: timestamppb.New(evt.StartDatetime),
			EndDatetime:   timestamppb.New(evt.EndDatetime),
		},
	}, nil
}

func (h *Handler) UpdateEvent(
	_ context.Context,
	req *grpc_calendar.UpdateEvent_Request,
) (*grpc_calendar.UpdateEvent_Reply, error) {
	id := req.GetId()
	if id == "" {
		return nil, errors.New("parameter id: is required")
	}

	origEvt, err := h.service.GetEvent(id)
	if err != nil {
		return nil, fmt.Errorf("find item (%s) for update: %w", id, err)
	}

	err = req.GetStartDatetime().CheckValid()
	if err != nil {
		return nil, err
	}

	err = req.GetEndDatetime().CheckValid()
	if err != nil {
		return nil, err
	}

	origEvt.ID = req.Id
	origEvt.Title = req.Title
	origEvt.Description = req.Description
	origEvt.RemindBefore = req.RemindBefore
	origEvt.StartDatetime = req.StartDatetime.AsTime()
	origEvt.EndDatetime = req.EndDatetime.AsTime()

	err = h.service.UpdateEvent(origEvt)
	if err != nil {
		return nil, fmt.Errorf("update event (%s): %w", id, err)
	}

	return &grpc_calendar.UpdateEvent_Reply{
		Result: &grpc_calendar.Result{
			Code:    0,
			Message: "success",
		},
	}, nil
}

func (h *Handler) DeleteEvent(
	_ context.Context,
	req *grpc_calendar.DeleteEvent_Request,
) (*grpc_calendar.DeleteEvent_Reply, error) {
	id := req.GetId()
	if id == "" {
		return nil, errors.New("parameter event_id is required")
	}

	err := h.service.DeleteEvent(id)
	if err != nil {
		return nil, err
	}

	return &grpc_calendar.DeleteEvent_Reply{
		Result: &grpc_calendar.Result{
			Code:    0,
			Message: "success",
		},
	}, nil
}

func (h *Handler) ListEvent(
	_ context.Context,
	req *grpc_calendar.ListEvent_Request,
) (*grpc_calendar.ListEvent_Reply, error) {
	mode := req.GetMode()
	if mode == grpc_calendar.ListEvent_UNSPECIFIED {
		return nil, errors.New("specify for which period to upload the data")
	}

	res, err := h.service.ListEvent(enum.RangeDuration(mode))
	if err != nil {
		return nil, err
	}

	result := make([]*grpc_calendar.Domain_Event, 0, len(res))
	for _, evt := range res {
		result = append(result, &grpc_calendar.Domain_Event{
			Id:            evt.ID,
			Title:         evt.Title,
			Description:   evt.Description,
			UserId:        evt.UserID,
			RemindBefore:  evt.RemindBefore,
			StartDatetime: timestamppb.New(evt.StartDatetime),
			EndDatetime:   timestamppb.New(evt.EndDatetime),
		})
	}

	return &grpc_calendar.ListEvent_Reply{Result: result}, nil
}
