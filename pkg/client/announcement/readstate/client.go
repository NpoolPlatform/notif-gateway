package readstate

import (
	"context"
	"fmt"
	"time"

	mgrpb "github.com/NpoolPlatform/message/npool/notif/mgr/v1/announcement/readstate"

	grpc2 "github.com/NpoolPlatform/go-service-framework/pkg/grpc"

	"github.com/NpoolPlatform/libent-cruder/pkg/cruder"

	npool "github.com/NpoolPlatform/message/npool/notif/mw/v1/announcement/readstate"

	constant "github.com/NpoolPlatform/notif-gateway/pkg/message/const"
)

var timeout = 10 * time.Second

type handler func(context.Context, npool.MiddlewareClient) (cruder.Any, error)

func withCRUD(ctx context.Context, handler handler) (cruder.Any, error) {
	_ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	conn, err := grpc2.GetGRPCConn(constant.ServiceName, grpc2.GRPCTAG)
	if err != nil {
		return nil, fmt.Errorf("fail get readstate connection: %v", err)
	}

	defer conn.Close()

	cli := npool.NewMiddlewareClient(conn)

	return handler(_ctx, cli)
}

func GetReadState(ctx context.Context, announcementID, userID string) (*npool.ReadState, error) {
	info, err := withCRUD(ctx, func(_ctx context.Context, cli npool.MiddlewareClient) (cruder.Any, error) {
		resp, err := cli.GetReadState(ctx, &npool.GetReadStateRequest{
			AnnouncementID: announcementID,
			UserID:         userID,
		})
		if err != nil {
			return nil, fmt.Errorf("fail get readstate: %v", err)
		}
		return resp.GetInfo(), nil
	})
	if err != nil {
		return nil, fmt.Errorf("fail get readstate: %v", err)
	}
	return info.(*npool.ReadState), nil
}

func GetReadStates(ctx context.Context, conds *mgrpb.Conds, offset, limit int32) ([]*npool.ReadState, uint32, error) {
	var total uint32
	infos, err := withCRUD(ctx, func(_ctx context.Context, cli npool.MiddlewareClient) (cruder.Any, error) {
		resp, err := cli.GetReadStates(ctx, &npool.GetReadStatesRequest{
			Conds:  conds,
			Limit:  limit,
			Offset: offset,
		})
		if err != nil {
			return nil, fmt.Errorf("fail get readstates: %v", err)
		}
		total = resp.GetTotal()
		return resp.GetInfos(), nil
	})
	if err != nil {
		return nil, 0, fmt.Errorf("fail get readstates: %v", err)
	}
	return infos.([]*npool.ReadState), total, nil
}
