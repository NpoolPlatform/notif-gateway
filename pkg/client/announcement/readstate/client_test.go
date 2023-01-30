package readstate

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/google/uuid"

	"github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	valuedef "github.com/NpoolPlatform/message/npool"

	"github.com/NpoolPlatform/go-service-framework/pkg/config"

	"bou.ke/monkey"
	grpc2 "github.com/NpoolPlatform/go-service-framework/pkg/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/NpoolPlatform/notif-gateway/pkg/testinit"

	announcementmgrpb "github.com/NpoolPlatform/message/npool/notif/mgr/v1/announcement"
	readstatemgrpb "github.com/NpoolPlatform/message/npool/notif/mgr/v1/announcement/readstate"
	channelpb "github.com/NpoolPlatform/message/npool/notif/mgr/v1/channel"
	npool "github.com/NpoolPlatform/message/npool/notif/mw/v1/announcement/readstate"
	"github.com/stretchr/testify/assert"

	announcementcrud "github.com/NpoolPlatform/notif-manager/pkg/crud/announcement"
	readstatecrud "github.com/NpoolPlatform/notif-manager/pkg/crud/announcement/readstate"
)

func init() {
	if runByGithubAction, err := strconv.ParseBool(os.Getenv("RUN_BY_GITHUB_ACTION")); err == nil && runByGithubAction {
		return
	}
	if err := testinit.Init(); err != nil {
		fmt.Printf("cannot init test stub: %v\n", err)
	}
}

var data = npool.ReadState{
	AnnouncementID: uuid.NewString(),
	AppID:          uuid.NewString(),
	UserID:         uuid.NewString(),
	Title:          uuid.NewString(),
	Content:        uuid.NewString(),
	ChannelsStr:    `["ChannelEmail", "ChannelSMS"]`,
	Channels:       []channelpb.NotifChannel{channelpb.NotifChannel_ChannelEmail, channelpb.NotifChannel_ChannelSMS},
	EmailSend:      true,
	AlreadyRead:    true,
}

func getReadState(t *testing.T) {
	_, err := announcementcrud.Create(context.Background(), &announcementmgrpb.AnnouncementReq{
		ID:        &data.AnnouncementID,
		AppID:     &data.AppID,
		Title:     &data.Title,
		Content:   &data.Content,
		Channels:  data.Channels,
		EmailSend: &data.EmailSend,
	})
	assert.Nil(t, err)

	_, err = readstatecrud.Create(context.Background(), &readstatemgrpb.ReadStateReq{
		AppID:          &data.AppID,
		UserID:         &data.UserID,
		AnnouncementID: &data.AnnouncementID,
	})
	assert.Nil(t, err)

	info, err := GetReadState(context.Background(), data.AnnouncementID, data.UserID)
	if assert.Nil(t, err) {
		data.CreatedAt = info.CreatedAt
		data.UpdatedAt = info.UpdatedAt
		assert.Equal(t, &data, info)
	}
}

func getReadStates(t *testing.T) {
	infos, total, err := GetReadStates(context.Background(), &readstatemgrpb.Conds{
		AppID: &valuedef.StringVal{
			Op:    cruder.EQ,
			Value: data.AppID,
		},
		UserID: &valuedef.StringVal{
			Op:    cruder.EQ,
			Value: data.UserID,
		},
		AnnouncementID: &valuedef.StringVal{
			Op:    cruder.EQ,
			Value: data.AnnouncementID,
		},
	}, 0, 1)
	if assert.Nil(t, err) {
		assert.Equal(t, total, uint32(1))
		data.CreatedAt = infos[0].CreatedAt
		data.UpdatedAt = infos[0].UpdatedAt
		assert.Equal(t, infos[0], &data)
	}
}

func TestClient(t *testing.T) {
	if runByGithubAction, err := strconv.ParseBool(os.Getenv("RUN_BY_GITHUB_ACTION")); err == nil && runByGithubAction {
		return
	}

	gport := config.GetIntValueWithNameSpace("", config.KeyGRPCPort)

	monkey.Patch(grpc2.GetGRPCConn, func(service string, tags ...string) (*grpc.ClientConn, error) {
		return grpc.Dial(fmt.Sprintf("localhost:%v", gport), grpc.WithTransportCredentials(insecure.NewCredentials()))
	})
	t.Run("getReadState", getReadState)
	t.Run("getReadStates", getReadStates)
}
