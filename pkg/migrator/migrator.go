package migrator

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/NpoolPlatform/go-service-framework/pkg/config"
	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	redis2 "github.com/NpoolPlatform/go-service-framework/pkg/redis"

	constant "github.com/NpoolPlatform/go-service-framework/pkg/mysql/const"
	constant1 "github.com/NpoolPlatform/notif-gateway/pkg/message/const"

	"github.com/NpoolPlatform/notif-manager/pkg/db"
	"github.com/NpoolPlatform/notif-manager/pkg/db/ent"
	entemailtmpl "github.com/NpoolPlatform/notif-manager/pkg/db/ent/emailtemplate"

	_ "github.com/NpoolPlatform/notif-manager/pkg/db/ent/runtime"

	"github.com/google/uuid"
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

func migrateEmailTemplate(ctx context.Context) error {
	type tmpl struct {
		ID                uuid.UUID
		AppID             uuid.UUID
		LangID            uuid.UUID
		DefaultToUsername string
		UsedFor           string
		Sender            string
		ReplyTos          []uint8
		CCTos             []uint8
		Subject           string
		Body              string
		CreatedAt         uint32
		UpdatedAt         uint32
		DeletedAt         uint32
	}

	return db.WithClient(ctx, func(_ctx context.Context, cli *ent.Client) error {
		rows, err := cli.Debug().QueryContext(
			ctx,
			"select "+
				"id,"+
				"app_id,"+
				"lang_id,"+
				"default_to_username,"+
				"used_for,"+
				"sender,"+
				"reply_tos,"+
				"cc_tos,"+
				"subject,"+
				"body,"+
				"created_at,"+
				"updated_at,"+
				"deleted_at "+
				"from third_manager.email_templates",
		)
		if err != nil {
			return err
		}

		bulk := []*ent.EmailTemplateCreate{}

		for rows.Next() {
			tmpl := tmpl{}
			err := rows.Scan(
				&tmpl.ID,
				&tmpl.AppID,
				&tmpl.LangID,
				&tmpl.DefaultToUsername,
				&tmpl.UsedFor,
				&tmpl.Sender,
				&tmpl.ReplyTos,
				&tmpl.CCTos,
				&tmpl.Subject,
				&tmpl.Body,
				&tmpl.CreatedAt,
				&tmpl.UpdatedAt,
				&tmpl.DeletedAt,
			)
			if err != nil {
				return err
			}

			exist, err := cli.
				EmailTemplate.
				Query().
				Where(
					entemailtmpl.ID(tmpl.ID),
				).
				Exist(_ctx)
			if err != nil {
				return err
			}
			if exist {
				continue
			}

			replyTos := []string{}
			_ = json.Unmarshal(tmpl.ReplyTos, &replyTos)
			ccTos := []string{}
			_ = json.Unmarshal(tmpl.CCTos, &ccTos)

			bulk = append(
				bulk,
				cli.
					EmailTemplate.
					Create().
					SetID(tmpl.ID).
					SetAppID(tmpl.AppID).
					SetLangID(tmpl.LangID).
					SetDefaultToUsername(tmpl.DefaultToUsername).
					SetUsedFor(tmpl.UsedFor).
					SetSender(tmpl.Sender).
					SetReplyTos(replyTos).
					SetCcTos(ccTos).
					SetSubject(tmpl.Subject).
					SetBody(tmpl.Body).
					SetCreatedAt(tmpl.CreatedAt).
					SetUpdatedAt(tmpl.UpdatedAt).
					SetUpdatedAt(tmpl.DeletedAt),
			)
		}

		_, err = cli.
			EmailTemplate.
			CreateBulk(bulk...).
			Save(_ctx)
		return err

		return nil
	})
}

func migrateSMSTemplate(ctx context.Context) error {
	return nil
}

func migrateFrontendTemplate(ctx context.Context) error {
	return nil
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

	if err := migrateEmailTemplate(ctx); err != nil {
		logger.Sugar().Errorw("Migrate", "error", err)
		return err
	}

	if err := migrateSMSTemplate(ctx); err != nil {
		logger.Sugar().Errorw("Migrate", "error", err)
		return err
	}

	if err := migrateFrontendTemplate(ctx); err != nil {
		logger.Sugar().Errorw("Migrate", "error", err)
		return err
	}

	logger.Sugar().Infow("Migrate", "Done", "success")

	return nil
}

func Abort(ctx context.Context) {
	_ = redis2.Unlock(lockKey())
}
