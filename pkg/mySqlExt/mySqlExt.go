package mySqlExt

import (
	"boilerplate-service/constant"
	"boilerplate-service/pkg/newRelicExt"
	"context"

	"database/sql"
	"fmt"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/newrelic/go-agent/v3/newrelic"
)

type IMySqlExt interface {
	Close() error
	QueryContext(
		ctx context.Context,
		query string,
		args ...interface{},
	) (*sql.Rows, error)
	ExecContext(
		ctx context.Context,
		query string,
		args ...interface{},
	) (bool, error)
	NamedExecContext(
		ctx context.Context,
		query string,
		args interface{},
	) (bool, error)
	GetContext(
		ctx context.Context,
		dest interface{},
		query string,
		args ...interface{},
	) error
	Ping() error
}

type Config struct {
	Host         string
	Port         string
	Username     string
	Password     string
	DBName       string
	MaxIdleConns int
	MaxIdleTime  int
	MaxLifeTime  int
	MaxOpenConns int
}

type mySqlExt struct {
	db *sqlx.DB
}

func New(config Config) (IMySqlExt, error) {
	db, err := sqlx.Connect(
		"mysql",
		fmt.Sprintf(
			"%s:%s@tcp(%s:%s)/%s?parseTime=true",
			config.Username,
			config.Password,
			config.Host,
			config.Port,
			config.DBName,
		),
	)
	if err != nil {
		return nil, err
	}

	setDBConfig(db, config)

	return &mySqlExt{db}, nil
}

func (m *mySqlExt) Close() error {
	return m.db.Close()
}

func setDBConfig(db *sqlx.DB, config Config) {
	if config.MaxIdleConns == 0 {
		config.MaxIdleConns = 15
	}

	if config.MaxOpenConns == 0 {
		config.MaxOpenConns = 25
	}

	if config.MaxIdleTime == 0 {
		config.MaxIdleTime = 300 // 5 Mins
	}

	if config.MaxLifeTime == 0 {
		config.MaxLifeTime = 300 // 5 Mins
	}

	db.SetMaxIdleConns(config.MaxIdleConns)
	db.SetMaxOpenConns(config.MaxOpenConns)
	db.SetConnMaxIdleTime(time.Duration(config.MaxIdleTime) * time.Second)
	db.SetConnMaxLifetime(time.Duration(config.MaxLifeTime) * time.Second)
}

func (m *mySqlExt) getTableName(ctx context.Context) string {
	if ctx.Value(constant.CtxSQLTableNameKey) != nil {
		return ctx.Value(constant.CtxSQLTableNameKey).(string)
	}
	return ""
}

func (m *mySqlExt) QueryContext(
	ctx context.Context,
	query string,
	args ...interface{},
) (*sql.Rows, error) {
	txn := newRelicExt.GetTxnFromCtx(ctx)
	operationsQuery := strings.Split(query, " ")
	dbSegment := newrelic.DatastoreSegment{
		StartTime:  txn.StartSegmentNow(),
		Product:    newrelic.DatastoreMySQL,
		Collection: m.getTableName(ctx),
		Operation:  operationsQuery[0],
		RawQuery:   query,
	}
	defer dbSegment.End()

	ctx = newrelic.NewContext(ctx, txn)

	return m.db.QueryContext(ctx, query, args...)
}

func (m *mySqlExt) ExecContext(
	ctx context.Context,
	query string,
	args ...interface{},
) (bool, error) {
	txn := newRelicExt.GetTxnFromCtx(ctx)
	operationsQuery := strings.Split(query, " ")
	dbSegment := newrelic.DatastoreSegment{
		StartTime:  txn.StartSegmentNow(),
		Product:    newrelic.DatastoreMySQL,
		Collection: m.getTableName(ctx),
		Operation:  operationsQuery[0],
		RawQuery:   query,
	}
	defer dbSegment.End()

	ctx = newrelic.NewContext(ctx, txn)

	sqlResults, err := m.db.ExecContext(ctx, query, args...)
	if err != nil {
		return false, err
	}

	affected, err := sqlResults.RowsAffected()
	return affected != 0, err
}

func (m *mySqlExt) NamedExecContext(
	ctx context.Context,
	query string,
	args interface{},
) (bool, error) {
	txn := newRelicExt.GetTxnFromCtx(ctx)
	operationsQuery := strings.Split(query, " ")
	dbSegment := newrelic.DatastoreSegment{
		StartTime:  txn.StartSegmentNow(),
		Product:    newrelic.DatastoreMySQL,
		Collection: m.getTableName(ctx),
		Operation:  operationsQuery[0],
		RawQuery:   query,
	}
	defer dbSegment.End()

	ctx = newrelic.NewContext(ctx, txn)

	sqlResults, err := m.db.NamedExecContext(ctx, query, args)
	if err != nil {
		return false, err
	}

	affected, err := sqlResults.RowsAffected()
	return affected != 0, err
}

func (m *mySqlExt) GetContext(
	ctx context.Context,
	dest interface{},
	query string,
	args ...interface{},
) error {
	txn := newRelicExt.GetTxnFromCtx(ctx)
	dbSegment := newrelic.DatastoreSegment{
		StartTime:  txn.StartSegmentNow(),
		Product:    newrelic.DatastoreMySQL,
		Collection: m.getTableName(ctx),
		Operation:  "SELECT",
		RawQuery:   query,
	}
	defer dbSegment.End()

	ctx = newrelic.NewContext(ctx, txn)

	return m.db.GetContext(ctx, dest, query, args...)
}

func (m *mySqlExt) Ping() error {
	return m.db.Ping()
}
