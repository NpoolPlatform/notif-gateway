package notif

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	valuedef "github.com/NpoolPlatform/message/npool"
	"github.com/NpoolPlatform/message/npool/notif/mgr/v1/channel"

	"github.com/NpoolPlatform/go-service-framework/pkg/config"

	"bou.ke/monkey"
	grpc2 "github.com/NpoolPlatform/go-service-framework/pkg/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/NpoolPlatform/notif-gateway/pkg/testinit"

	mgrpb "github.com/NpoolPlatform/message/npool/notif/mgr/v1/notif"
	npool "github.com/NpoolPlatform/message/npool/notif/mw/v1/notif"
	"github.com/stretchr/testify/assert"

	"github.com/google/uuid"
)

func init() {
	if runByGithubAction, err := strconv.ParseBool(os.Getenv("RUN_BY_GITHUB_ACTION")); err == nil && runByGithubAction {
		return
	}
	if err := testinit.Init(); err != nil {
		fmt.Printf("cannot init test stub: %v\n", err)
	}
}

var data = &npool.Notif{
	ID:          uuid.NewString(),
	AppID:       uuid.NewString(),
	UserID:      uuid.NewString(),
	AlreadyRead: true,
	LangID:      uuid.NewString(),
	EventType:   mgrpb.EventType_KycReviewApproved,
	UseTemplate: true,
	Title:       uuid.NewString(),
	Content:     uuid.NewString(),
	Channels:    []channel.NotifChannel{channel.NotifChannel_ChannelSMS, channel.NotifChannel_ChannelEmail},
	EmailSend:   true,
}

var dataReq = &mgrpb.NotifReq{
	ID:          &data.ID,
	AppID:       &data.AppID,
	UserID:      &data.UserID,
	AlreadyRead: &data.AlreadyRead,
	LangID:      &data.LangID,
	EventType:   &data.EventType,
	UseTemplate: &data.UseTemplate,
	Title:       &data.Title,
	Content:     &data.Content,
	Channels:    data.Channels,
	EmailSend:   &data.EmailSend,
}

func createNotif(t *testing.T) {
	info, err := CreateNotif(context.Background(), dataReq)
	if assert.Nil(t, err) {
		data.CreatedAt = info.CreatedAt
		data.UpdatedAt = info.UpdatedAt
		assert.Equal(t, data, info)
	}
}

func getNotif(t *testing.T) {
	info, err := GetNotif(context.Background(), data.ID)
	if assert.Nil(t, err) {
		assert.Equal(t, data, info)
	}
}

func getNotifs(t *testing.T) {
	infos, total, err := GetNotifs(context.Background(), &mgrpb.Conds{
		ID: &valuedef.StringVal{
			Op:    cruder.EQ,
			Value: data.ID,
		},
	}, 0, 1)
	if assert.Nil(t, err) {
		assert.Equal(t, total, uint32(1))
		assert.Equal(t, infos[0], data)
	}
}

func getNotifOnly(t *testing.T) {
	info, err := GetNotifOnly(context.Background(), &mgrpb.Conds{
		ID: &valuedef.StringVal{
			Op:    cruder.EQ,
			Value: data.ID,
		},
	})
	if assert.Nil(t, err) {
		assert.Equal(t, info, data)
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

	t.Run("createNotif", createNotif)
	t.Run("getNotif", getNotif)
	t.Run("getNotifs", getNotifs)
	t.Run("getNotifOnly", getNotifOnly)
}
