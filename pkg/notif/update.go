//nolint:nolintlint,dupl
package notif

import (
	"context"
	"fmt"

	usermwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/user"
	mwcli "github.com/NpoolPlatform/notif-middleware/pkg/client/notif"

	npool "github.com/NpoolPlatform/message/npool/notif/gw/v1/notif"
	notifmwpb "github.com/NpoolPlatform/message/npool/notif/mw/v1/notif"

	cruder "github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"
	constant "github.com/NpoolPlatform/notif-gateway/pkg/const"
)

type updateHandler struct {
	*Handler
	notifs   []*notifmwpb.Notif
	reqIDMap map[uint32]*uint32
}

func (h *updateHandler) validateNotifs(ctx context.Context) error {
	exist, err := usermwcli.ExistUser(ctx, *h.AppID, *h.UserID)
	if err != nil {
		return err
	}
	if !exist {
		return fmt.Errorf("invalid user")
	}
	idMap := make(map[uint32]*uint32)
	for _, row := range h.Reqs {
		h.IDs = append(h.IDs, *row.ID)
		idMap[*row.ID] = row.ID
	}
	h.reqIDMap = idMap

	limit := int32(len(h.IDs))
	notifs, _, err := mwcli.GetNotifs(ctx, &notifmwpb.Conds{
		IDs: &basetypes.Uint32SliceVal{Op: cruder.IN, Value: h.IDs},
	}, 0, limit)
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

func (h *updateHandler) updateNotifs(ctx context.Context) error {
	eventIDs := []string{}
	for _, row := range h.notifs {
		if row.Channel == basetypes.NotifChannel_ChannelFrontend {
			eventIDs = append(eventIDs, row.EventID)
		}
	}

	offset := int32(0)
	for {
		notifs, _, err := mwcli.GetNotifs(ctx, &notifmwpb.Conds{
			AppID:    &basetypes.StringVal{Op: cruder.EQ, Value: *h.AppID},
			UserID:   &basetypes.StringVal{Op: cruder.EQ, Value: *h.UserID},
			EventIDs: &basetypes.StringSliceVal{Op: cruder.IN, Value: eventIDs},
			Channel:  &basetypes.Uint32Val{Op: cruder.EQ, Value: uint32(basetypes.NotifChannel_ChannelFrontend)},
		}, offset, constant.DefaultRowLimit)
		if err != nil {
			return err
		}

		reqs := []*notifmwpb.NotifReq{}
		notified := true
		for _, _notif := range notifs {
			reqNotifID := h.reqIDMap[_notif.ID]
			if reqNotifID == nil {
				reqs = append(reqs, &notifmwpb.NotifReq{
					ID:       &_notif.ID,
					Notified: &notified,
				})
			}
		}
		h.Reqs = append(h.Reqs, reqs...)
		if int32(len(notifs)) < constant.DefaultRowLimit {
			break
		}
		offset += constant.DefaultRowLimit
	}

	_, err := mwcli.UpdateNotifs(ctx, h.Reqs)
	if err != nil {
		return err
	}

	return nil
}

//nolint:gocyclo
func (h *Handler) UpdateNotifs(ctx context.Context) ([]*npool.Notif, error) {
	handler := &updateHandler{
		Handler: h,
	}

	if err := handler.validateNotifs(ctx); err != nil {
		return nil, err
	}
	if err := handler.updateNotifs(ctx); err != nil {
		return nil, err
	}

	infos, _, err := h.GetNotifs(ctx)
	if err != nil {
		return nil, err
	}

	return infos, nil
}
