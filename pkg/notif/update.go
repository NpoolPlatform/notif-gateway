//nolint:nolintlint,dupl
package notif

import (
	"context"
	"fmt"

	mwcli "github.com/NpoolPlatform/notif-middleware/pkg/client/notif"

	npool "github.com/NpoolPlatform/message/npool/notif/gw/v1/notif"
	notifmwpb "github.com/NpoolPlatform/message/npool/notif/mw/v1/notif"

	cruder "github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"
	constant "github.com/NpoolPlatform/notif-gateway/pkg/const"
)

type updateHandler struct {
	*Handler
	notifs     []*notifmwpb.Notif
	updateReqs []*notifmwpb.NotifReq
	idMap      map[string]*string
}

func (h *updateHandler) validateReqs() error {
	reqs := []*notifmwpb.NotifReq{}
	idMap := make(map[string]*string)
	for _, row := range h.Reqs {
		if row.ID == nil {
			return fmt.Errorf("invalid id")
		}
		if row.Notified == nil {
			return fmt.Errorf("invalid notified")
		}
		if !*row.Notified {
			return fmt.Errorf("invalid notified %v", *row.Notified)
		}

		_req := &notifmwpb.NotifReq{
			ID:       row.ID,
			Notified: row.Notified,
		}
		reqs = append(reqs, _req)
		h.IDs = append(h.IDs, *row.ID)
		idMap[*row.ID] = row.ID
	}
	h.idMap = idMap
	h.updateReqs = append(h.updateReqs, reqs...)
	return nil
}

func (h *updateHandler) getNotifs(ctx context.Context) error {
	notifs, _, err := mwcli.GetNotifs(ctx, &notifmwpb.Conds{
		IDs: &basetypes.StringSliceVal{
			Op: cruder.IN, Value: h.IDs,
		},
	}, 0, constant.DefaultRowLimit)
	if err != nil {
		return err
	}
	if len(notifs) != len(h.IDs) {
		return fmt.Errorf("notif not exist")
	}

	for _, row := range notifs {
		if row.AppID != *h.AppID || row.UserID != *h.UserID {
			return fmt.Errorf("permission denied")
		}
	}

	h.notifs = notifs
	return nil
}

func (h *updateHandler) getSameEventIDNotifReqs(ctx context.Context) error {
	eventIDs := []string{}
	for _, row := range h.notifs {
		if row.Channel == basetypes.NotifChannel_ChannelFrontend {
			eventIDs = append(eventIDs, row.EventID)
		}
	}
	offset := int32(0)
	for {
		notifs, _, err := mwcli.GetNotifs(ctx, &notifmwpb.Conds{
			AppID: &basetypes.StringVal{
				Op: cruder.EQ, Value: *h.AppID,
			},
			UserID: &basetypes.StringVal{
				Op: cruder.EQ, Value: *h.UserID,
			},
			EventIDs: &basetypes.StringSliceVal{
				Op: cruder.IN, Value: eventIDs,
			},
			Channel: &basetypes.Uint32Val{
				Op: cruder.EQ, Value: uint32(basetypes.NotifChannel_ChannelFrontend),
			},
		}, offset, constant.DefaultRowLimit)
		if err != nil {
			return err
		}

		reqs := []*notifmwpb.NotifReq{}
		notified := true
		for _, _notif := range notifs {
			reqNotifID := h.idMap[_notif.ID]
			if reqNotifID == nil {
				reqs = append(reqs, &notifmwpb.NotifReq{
					ID:       &_notif.ID,
					Notified: &notified,
				})
			}
		}
		h.updateReqs = append(h.updateReqs, reqs...)
		if int32(len(notifs)) < constant.DefaultRowLimit {
			break
		}
		offset += 1
	}

	return nil
}

//nolint:gocyclo
func (h *Handler) UpdateNotifs(ctx context.Context) ([]*npool.Notif, error) {
	if h.AppID == nil || *h.AppID == "" {
		return nil, fmt.Errorf("invalid appid")
	}
	if h.UserID == nil || *h.UserID == "" {
		return nil, fmt.Errorf("invalid userid")
	}

	handler := &updateHandler{
		Handler: h,
	}

	if err := handler.validateReqs(); err != nil {
		return nil, err
	}
	if err := handler.getNotifs(ctx); err != nil {
		return nil, err
	}
	if err := handler.getSameEventIDNotifReqs(ctx); err != nil {
		return nil, err
	}

	_, err := mwcli.UpdateNotifs(ctx, handler.updateReqs)
	if err != nil {
		return nil, err
	}

	infos, _, err := h.GetNotifs(ctx)
	if err != nil {
		return nil, err
	}

	return infos, nil
}
