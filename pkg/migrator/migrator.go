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
	entemailtmpl "github.com/NpoolPlatform/notif-middleware/pkg/db/ent/emailtemplate"
	entfrontendtmpl "github.com/NpoolPlatform/notif-middleware/pkg/db/ent/frontendtemplate"
	entnotif "github.com/NpoolPlatform/notif-middleware/pkg/db/ent/notif"
	entchannel "github.com/NpoolPlatform/notif-middleware/pkg/db/ent/notifchannel"
	entuser "github.com/NpoolPlatform/notif-middleware/pkg/db/ent/notifuser"
	entsmstmpl "github.com/NpoolPlatform/notif-middleware/pkg/db/ent/smstemplate"

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
	return db.WithTx(ctx, func(_ctx context.Context, tx *ent.Tx) error {
		oldEventType := "GoodBenefit"
		_, err := tx.
			Notif.
			Update().
			Where(
				entnotif.EventType(oldEventType),
			).
			SetEventType(basetypes.UsedFor_GoodBenefit1.String()).
			Save(_ctx)
		if err != nil {
			return err
		}

		return nil
	})
}

func migrateNotifChannel(ctx context.Context) error {
	return db.WithTx(ctx, func(_ctx context.Context, tx *ent.Tx) error {
		oldEventType := "GoodBenefit"
		_, err := tx.
			NotifChannel.
			Update().
			Where(
				entchannel.EventType(oldEventType),
			).
			SetEventType(basetypes.UsedFor_GoodBenefit1.String()).
			Save(_ctx)
		if err != nil {
			return err
		}

		return nil
	})
}

func migrateEmailTmpl(ctx context.Context) error {
	return db.WithTx(ctx, func(_ctx context.Context, tx *ent.Tx) error {
		oldEventType := "GoodBenefit"
		_, err := tx.
			EmailTemplate.
			Update().
			Where(
				entemailtmpl.UsedFor(oldEventType),
			).
			SetUsedFor(basetypes.UsedFor_GoodBenefit1.String()).
			Save(_ctx)
		if err != nil {
			return err
		}

		return nil
	})
}

func migrateFrontendTmpl(ctx context.Context) error {
	return db.WithTx(ctx, func(_ctx context.Context, tx *ent.Tx) error {
		oldEventType := "GoodBenefit"
		_, err := tx.
			FrontendTemplate.
			Update().
			Where(
				entfrontendtmpl.UsedFor(oldEventType),
			).
			SetUsedFor(basetypes.UsedFor_GoodBenefit1.String()).
			Save(_ctx)
		if err != nil {
			return err
		}

		return nil
	})
}

func migrateSMSTmpl(ctx context.Context) error {
	return db.WithTx(ctx, func(_ctx context.Context, tx *ent.Tx) error {
		oldEventType := "GoodBenefit"
		_, err := tx.
			SMSTemplate.
			Update().
			Where(
				entsmstmpl.UsedFor(oldEventType),
			).
			SetUsedFor(basetypes.UsedFor_GoodBenefit1.String()).
			Save(_ctx)
		if err != nil {
			return err
		}

		return nil
	})
}

func migrateNotifUser(ctx context.Context) error {
	return db.WithTx(ctx, func(_ctx context.Context, tx *ent.Tx) error {
		oldEventType := "GoodBenefit"
		_, err := tx.
			NotifUser.
			Update().
			Where(
				entuser.EventType(oldEventType),
			).
			SetEventType(basetypes.UsedFor_GoodBenefit1.String()).
			Save(_ctx)
		if err != nil {
			return err
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
	if err := migrateNotifChannel(ctx); err != nil {
		logger.Sugar().Errorw("Migrate", "error", err)
		return err
	}
	if err := migrateEmailTmpl(ctx); err != nil {
		logger.Sugar().Errorw("Migrate", "error", err)
		return err
	}
	if err := migrateFrontendTmpl(ctx); err != nil {
		logger.Sugar().Errorw("Migrate", "error", err)
		return err
	}
	if err := migrateSMSTmpl(ctx); err != nil {
		logger.Sugar().Errorw("Migrate", "error", err)
		return err
	}

	logger.Sugar().Infow("Migrate", "Done", "success")

	return nil
}

func Abort(ctx context.Context) {
	_ = redis2.Unlock(lockKey())
}
