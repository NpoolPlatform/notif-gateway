//nolint
package migrator

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"

	"github.com/NpoolPlatform/go-service-framework/pkg/config"
	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	redis2 "github.com/NpoolPlatform/go-service-framework/pkg/redis"

	constant "github.com/NpoolPlatform/go-service-framework/pkg/mysql/const"
	constant1 "github.com/NpoolPlatform/notif-gateway/pkg/message/const"

	"github.com/NpoolPlatform/notif-middleware/pkg/db"
	"github.com/NpoolPlatform/notif-middleware/pkg/db/ent"
	entnotif "github.com/NpoolPlatform/notif-middleware/pkg/db/ent/notif"
	entuser "github.com/NpoolPlatform/notif-middleware/pkg/db/ent/notifuser"

	_ "github.com/NpoolPlatform/notif-middleware/pkg/db/ent/runtime"
)

const (
	keyUsername  = "username"
	keyPassword  = "password"
	keyServiceID = "serviceid"
	keyDBName    = "database_name"
	maxOpen      = 10
	maxIdle      = 10
	MaxLife      = 3
)

func dsn(hostname string) (string, error) {
	username := config.GetStringValueWithNameSpace(constant.MysqlServiceName, keyUsername)
	password := config.GetStringValueWithNameSpace(constant.MysqlServiceName, keyPassword)
	dbname := config.GetStringValueWithNameSpace(hostname, keyDBName)

	svc, err := config.PeekService(constant.MysqlServiceName)
	if err != nil {
		logger.Sugar().Warnw("dsb", "error", err)
		return "", err
	}

	return fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?parseTime=true&interpolateParams=true",
		username, password,
		svc.Address,
		svc.Port,
		dbname,
	), nil
}

func open(hostname string) (conn *sql.DB, err error) {
	hdsn, err := dsn(hostname)
	if err != nil {
		return nil, err
	}

	logger.Sugar().Infow("open", "hdsn", hdsn)

	conn, err = sql.Open("mysql", hdsn)
	if err != nil {
		return nil, err
	}

	// https://github.com/go-sql-driver/mysql
	// See "Important settings" section.

	conn.SetConnMaxLifetime(time.Minute * MaxLife)
	conn.SetMaxOpenConns(maxOpen)
	conn.SetMaxIdleConns(maxIdle)

	return conn, nil
}

func lockKey() string {
	serviceID := config.GetStringValueWithNameSpace(constant1.ServiceName, keyServiceID)
	return fmt.Sprintf("migrator:%v", serviceID)
}

func migrateNotif(ctx context.Context) error {
	return db.WithClient(ctx, func(_ctx context.Context, cli *ent.Client) error {
		_, err := cli.
			Notif.
			Update().
			Where(
				entnotif.CreatedAtLT(1685030400),
			).
			SetNotified(true).
			Save(_ctx)
		return err
	})
}

func migrateNotifUser(ctx context.Context) error {
	return db.WithTx(ctx, func(_ctx context.Context, tx *ent.Tx) error {
		oldEventType := "GoodBenefit"
		infos, err := tx.
			NotifUser.
			Query().
			Where(
				entuser.EventType(oldEventType),
			).
			All(_ctx)
		if err != nil {
			return err
		}

		for _, info := range infos {
			_, err = tx.
				NotifUser.
				UpdateOneID(info.ID).
				SetEventType(basetypes.UsedFor_GoodBenefit1.String()).
				Save(_ctx)
			if err != nil {
				return err
			}
		}

		return nil
	})
}

func Migrate(ctx context.Context) error {
	if err := redis2.TryLock(lockKey(), 0); err != nil {
		return err
	}
	defer func() {
		_ = redis2.Unlock(lockKey())
	}()

	logger.Sugar().Infow("Migrate", "Start", "...")

	if err := db.Init(); err != nil {
		logger.Sugar().Errorw("Migrate", "error", err)
		return err
	}

	if err := migrateNotif(ctx); err != nil {
		logger.Sugar().Errorw("Migrate", "error", err)
		return err
	}

	if err := migrateNotifUser(ctx); err != nil {
		logger.Sugar().Errorw("Migrate", "error", err)
		return err
	}

	logger.Sugar().Infow("Migrate", "Done", "success")

	return nil
}

func Abort(ctx context.Context) {
	_ = redis2.Unlock(lockKey())
}
