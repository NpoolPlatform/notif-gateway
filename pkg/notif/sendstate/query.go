package sendstate

import (
	"context"
	"fmt"

	appusermwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/user"
	"github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	usermwpb "github.com/NpoolPlatform/message/npool/appuser/mw/v1/user"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"
	npool "github.com/NpoolPlatform/message/npool/notif/gw/v1/notif/sendstate"
	mwpb "github.com/NpoolPlatform/message/npool/notif/mw/v1/notif/sendstate"
	mwcli "github.com/NpoolPlatform/notif-middleware/pkg/client/notif/sendstate"
)

func (h *Handler) GetSendStates(ctx context.Context) ([]*npool.SendState, uint32, error) {
	if h.AppID == nil {
		return nil, 0, fmt.Errorf("invalid appid")
	}
	if h.UserID == nil {
		return nil, 0, fmt.Errorf("invalid userid")
	}

	user, err := appusermwcli.GetUser(ctx, *h.AppID, *h.UserID)
	if err != nil {
		return nil, 0, err
	}

	if user == nil {
		return nil, 0, fmt.Errorf("invalid user")
	}

	conds := &mwpb.Conds{
		AppID: &basetypes.StringVal{
			Op:    cruder.EQ,
			Value: *h.AppID,
		},
		UserID: &basetypes.StringVal{
			Op:    cruder.EQ,
			Value: *h.UserID,
		},
	}
	if h.EventID != nil {
		conds.EventID = &basetypes.StringVal{
			Op:    cruder.EQ,
			Value: *h.EventID,
		}
	}
	if h.Channel != nil {
		conds.Channel = &basetypes.Uint32Val{
			Op:    cruder.EQ,
			Value: uint32(*h.Channel),
		}
	}

	rows, total, err := mwcli.GetSendStates(ctx, conds, h.Offset, h.Limit)
	if err != nil {
		return nil, 0, err
	}
	if len(rows) == 0 {
		return nil, total, nil
	}

	infos := []*npool.SendState{}
	for _, val := range rows {
		infos = append(infos, &npool.SendState{
			ID:           val.ID,
			EventID:      val.EventID,
			AppID:        val.AppID,
			UserID:       val.UserID,
			EmailAddress: user.EmailAddress,
			PhoneNO:      user.PhoneNO,
			Username:     user.Username,
			Channel:      val.Channel,
			CreatedAt:    val.CreatedAt,
			UpdatedAt:    val.UpdatedAt,
		})
	}
	return infos, total, nil
}

func (h *Handler) GetAppSendStates(ctx context.Context) ([]*npool.SendState, uint32, error) {
	if h.AppID == nil {
		return nil, 0, fmt.Errorf("invalid appid")
	}
	conds := &mwpb.Conds{
		AppID: &basetypes.StringVal{
			Op:    cruder.EQ,
			Value: *h.AppID,
		},
	}
	if h.EventID != nil {
		conds.EventID = &basetypes.StringVal{
			Op:    cruder.EQ,
			Value: *h.EventID,
		}
	}
	if h.Channel != nil {
		conds.Channel = &basetypes.Uint32Val{
			Op:    cruder.EQ,
			Value: uint32(*h.Channel),
		}
	}

	rows, total, err := mwcli.GetSendStates(ctx, conds, h.Offset, h.Limit)
	if err != nil {
		return nil, 0, err
	}
	if len(rows) == 0 {
		return nil, total, nil
	}

	userIDs := []string{}

	for _, val := range rows {
		if val.UserID != "" {
			userIDs = append(userIDs, val.UserID)
		}
	}

	userMap := map[string]*usermwpb.User{}
	if len(userIDs) > 0 {
		userInfos, _, err := appusermwcli.GetUsers(ctx, &usermwpb.Conds{
			IDs: &basetypes.StringSliceVal{
				Op: cruder.IN, Value: userIDs,
			},
		}, 0, int32(len(userIDs)))
		if err != nil {
			return nil, 0, err
		}

		for _, val := range userInfos {
			userMap[val.ID] = val
		}
	}

	infos := []*npool.SendState{}
	for _, val := range rows {
		user, ok := userMap[val.UserID]
		if !ok {
			continue
		}
		infos = append(infos, &npool.SendState{
			ID:           val.ID,
			EventID:      val.EventID,
			AppID:        val.AppID,
			UserID:       val.UserID,
			EmailAddress: user.EmailAddress,
			PhoneNO:      user.PhoneNO,
			Username:     user.Username,
			Channel:      val.Channel,
			CreatedAt:    val.CreatedAt,
			UpdatedAt:    val.UpdatedAt,
		})
	}
	return infos, total, nil
}
