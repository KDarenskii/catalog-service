package rcpostgres

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"
	"strconv"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/migrate"

	"github.com/KDarenskii/catalog-service/internal/app/config/section"
	"github.com/KDarenskii/catalog-service/migration"
)

type (
	Client struct {
		_bunDB
		rawBunDB *bun.DB

		cfg section.RepositoryPostgres
	}

	_bunDB = bun.IDB
)

func (c *Client) GetRawBunDB() *bun.DB {
	return c.rawBunDB
}

func NewClient(ctx context.Context, cfg section.RepositoryPostgres) (*Client, error) {
	var u url.URL

	u.Scheme = "postgres"
	u.Host = cfg.Address
	u.User = url.UserPassword(cfg.Username, cfg.Password)
	u.Path = cfg.Name

	args := make(url.Values)
	args.Set("sslmode", "disable")
	u.RawQuery = args.Encode()

	dsn := u.String()

	log.Info().Str("address", cfg.Address).
		Str("read_timeout", cfg.ReadTimeout.String()).
		Str("write_timeout", cfg.WriteTimeout.String()).
		Msg("Initializing PostgreSQL connection")

	sqlDB := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn),
		pgdriver.WithReadTimeout(cfg.ReadTimeout),
		pgdriver.WithWriteTimeout(cfg.WriteTimeout)))

	sqlDB.SetMaxOpenConns(10)

	bunDB := bun.NewDB(sqlDB, pgdialect.New(), bun.WithDiscardUnknownColumns())

	pingCtx, cancel := context.WithTimeout(ctx, 2*time.Second)

	defer cancel()

	if err := sqlDB.PingContext(pingCtx); err != nil {
		_ = bunDB.Close()

		return nil, fmt.Errorf("ping postgres: %w", err)
	}

	log.Info().Msg("PostgreSQL connection established")

	return &Client{
		rawBunDB: bunDB,
		cfg:      cfg,
		_bunDB:   newTxInjector(bunDB),
	}, nil
}

func (c *Client) Migrate(ctx context.Context) (oldVer, newVer int64, err error) {
	migrations := migrate.NewMigrations()

	if err := migrations.Discover(migration.Postgres); err != nil {
		return 0, 0, fmt.Errorf("failed to discover migrations: %w", err)
	}

	migratorOpts := []migrate.MigratorOption{
		migrate.WithTableName(c.cfg.MigrationTable),
		migrate.WithLocksTableName(c.cfg.MigrationTable + "_lock"),
		migrate.WithMarkAppliedOnSuccess(true),
	}

	migrator := migrate.NewMigrator(c.rawBunDB, migrations, migratorOpts...)

	if err := migrator.Init(ctx); err != nil {
		return 0, 0, fmt.Errorf("init migrations: %w", err)
	}

	applied, err := migrator.AppliedMigrations(ctx)
	if err != nil {
		return 0, 0, fmt.Errorf("get applied migrations: %w", err)
	}

	oldVer = getMaxMigrationVersion(&applied)

	appliedGroup, err := migrator.Migrate(ctx)
	if err != nil {
		return oldVer, oldVer, fmt.Errorf("migrate: %w", err)
	}

	if len(appliedGroup.Migrations) == 0 {
		newVer = oldVer
	} else {
		newVer = getMaxMigrationVersion(&appliedGroup.Migrations)
	}

	return oldVer, newVer, nil
}

func getMaxMigrationVersion(migrations *migrate.MigrationSlice) (maxVersion int64) {
	for _, mg := range *migrations {
		v, _ := strconv.ParseInt(mg.Name, 10, 64)
		if v > maxVersion {
			maxVersion = v
		}
	}

	return maxVersion
}

func (c *Client) InsideTx(ctx context.Context, fn func(ctx context.Context) error) error {
	if tx := getTxFromContext(ctx); tx.Tx != nil {
		return fn(ctx)
	}

	tx, err := c.rawBunDB.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}

	done := false

	defer func() {
		if !done {
			_ = tx.Rollback()
		}
	}()

	ctx = setTxToContext(ctx, tx)

	if err := fn(ctx); err != nil {
		return fmt.Errorf("exec tx: %w", err)
	}

	done = true

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit tx: %w", err)
	}

	return nil
}
