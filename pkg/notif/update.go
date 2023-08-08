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

//nolint:gocyclo
func (h *Handler) UpdateNotifs(ctx context.Context) ([]*npool.Notif, error) {
	if h.AppID == nil || *h.AppID == "" {
		return nil, fmt.Errorf("invalid appid")
	}
	if h.UserID == nil || *h.UserID == "" {
		return nil, fmt.Errorf("invalid userid")
	}
	reqs := []*notifmwpb.NotifReq{}
	for _, row := range h.Reqs {
		if row.ID == nil {
			return nil, fmt.Errorf("invalid id")
		}
		if row.Notified == nil {
			return nil, fmt.Errorf("invalid notified")
		}
		if !*row.Notified {
			return nil, fmt.Errorf("invalid notified %v", *row.Notified)
		}

		notifID := *row.ID
		notifInfo, err := mwcli.GetNotif(ctx, *row.ID)
		if err != nil {
			return nil, err
		}
		if notifInfo == nil {
			return nil, fmt.Errorf("notif not exist")
		}
		if notifInfo.AppID != *h.AppID || notifInfo.UserID != *h.UserID {
			return nil, fmt.Errorf("permission denied")
		}

		// notif state of frontend channel need to keep consistent in multi language
		if notifInfo.Channel == basetypes.NotifChannel_ChannelFrontend {
			notifs, _, err := mwcli.GetNotifs(ctx, &notifmwpb.Conds{
				AppID: &basetypes.StringVal{
					Op: cruder.EQ, Value: *h.AppID,
				},
				UserID: &basetypes.StringVal{
					Op: cruder.EQ, Value: *h.UserID,
				},
				EventID: &basetypes.StringVal{
					Op: cruder.EQ, Value: notifInfo.EventID,
				},
				Channel: &basetypes.Uint32Val{
					Op: cruder.EQ, Value: uint32(basetypes.NotifChannel_ChannelFrontend),
				},
			}, 0, constant.DefaultRowLimit)
			if err != nil {
				return nil, err
			}

			for _, _notif := range notifs {
				if _notif.ID != *row.ID {
					reqs = append(reqs, &notifmwpb.NotifReq{
						ID:       &_notif.ID,
						Notified: row.Notified,
					})
				}
			}
		}

		_req := &notifmwpb.NotifReq{
			ID:       &notifID,
			Notified: row.Notified,
		}
		reqs = append(reqs, _req)
		h.IDs = append(h.IDs, *row.ID)
	}

	_, err := mwcli.UpdateNotifs(ctx, reqs)
	if err != nil {
		return nil, err
	}

	infos, _, err := h.GetNotifs(ctx)
	if err != nil {
		return nil, err
	}

	return infos, nil
}
