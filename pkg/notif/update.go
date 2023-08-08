//nolint:nolintlint,dupl
package notif

import (
	"context"
	"fmt"

	appmwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/app"
	usermwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/user"
	applangmwcli "github.com/NpoolPlatform/g11n-middleware/pkg/client/applang"

	appmwpb "github.com/NpoolPlatform/message/npool/appuser/mw/v1/app"
	usermwpb "github.com/NpoolPlatform/message/npool/appuser/mw/v1/user"
	applangmwpb "github.com/NpoolPlatform/message/npool/g11n/mw/v1/applang"

	mwcli "github.com/NpoolPlatform/notif-middleware/pkg/client/notif"

	npool "github.com/NpoolPlatform/message/npool/notif/gw/v1/notif"
	notifmwpb "github.com/NpoolPlatform/message/npool/notif/mw/v1/notif"

	cruder "github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"
	constant "github.com/NpoolPlatform/notif-gateway/pkg/const"
)

type updateHandler struct {
	*Handler
}

func (h *updateHandler) createNotifsResp(ctx context.Context, notifs []*notifmwpb.Notif) ([]*npool.Notif, error) {
	appIDs := []string{}
	userIDs := []string{}
	langIDs := []string{}

	for _, val := range notifs {
		appIDs = append(appIDs, val.AppID)
		userIDs = append(userIDs, val.UserID)
		langIDs = append(langIDs, val.LangID)
	}
	appInfos, _, err := appmwcli.GetApps(ctx, &appmwpb.Conds{
		IDs: &basetypes.StringSliceVal{Op: cruder.IN, Value: appIDs},
	}, 0, int32(len(appIDs)))
	if err != nil {
		return nil, err
	}
	appMap := map[string]*appmwpb.App{}
	langMap := map[string]*applangmwpb.Lang{}

	for _, val := range appInfos {
		appMap[val.ID] = val
		langs, _, err := applangmwcli.GetLangs(ctx, &applangmwpb.Conds{
			AppID:   &basetypes.StringVal{Op: cruder.EQ, Value: val.ID},
			LangIDs: &basetypes.StringSliceVal{Op: cruder.IN, Value: langIDs},
		}, 0, int32(len(langIDs)))
		if err != nil {
			return nil, err
		}
		for _, lang := range langs {
			langMap[lang.AppID+"-"+lang.LangID] = lang
		}
	}

	userInfos, _, err := usermwcli.GetUsers(ctx, &usermwpb.Conds{
		IDs: &basetypes.StringSliceVal{Op: cruder.IN, Value: userIDs},
	}, 0, int32(len(userIDs)))
	if err != nil {
		return nil, err
	}
	userMap := map[string]*usermwpb.User{}
	for _, val := range userInfos {
		userMap[val.ID] = val
	}

	infos := []*npool.Notif{}
	for _, val := range notifs {
		app, ok := appMap[val.AppID]
		if !ok {
			continue
		}
		user, ok := userMap[val.UserID]
		if !ok {
			continue
		}
		lang, ok := langMap[val.AppID+"-"+val.LangID]
		if !ok {
			continue
		}

		infos = append(infos, &npool.Notif{
			ID:           val.ID,
			AppID:        val.AppID,
			AppName:      app.Name,
			UserID:       val.UserID,
			EmailAddress: user.EmailAddress,
			PhoneNO:      user.PhoneNO,
			Username:     user.Username,
			EventID:      val.EventID,
			EventType:    val.EventType,
			UseTemplate:  val.UseTemplate,
			Title:        val.Title,
			Content:      val.Content,
			Channel:      val.Channel,
			LangID:       lang.LangID,
			Lang:         lang.Lang,
			NotifType:    val.NotifType,
			Notified:     val.Notified,
			CreatedAt:    val.CreatedAt,
			UpdatedAt:    val.UpdatedAt,
		})
	}
	return infos, nil
}

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
	}
	rows, err := mwcli.UpdateNotifs(ctx, reqs)
	if err != nil {
		return nil, err
	}

	if len(rows) == 0 {
		return nil, nil
	}

	handler := &updateHandler{
		Handler: h,
	}
	infos, err := handler.createNotifsResp(ctx, rows)
	if err != nil {
		return nil, err
	}
	return infos, nil
}
