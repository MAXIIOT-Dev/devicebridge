//go:generate go-bindata -prefix ../migrate/ -pkg storage -o ./migrations_gen.go ../migrate/
package storage

import (
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/rubenv/sql-migrate"
	log "github.com/sirupsen/logrus"
)

var db *sqlx.DB

// Connect 连接数据库
func Connect(dsn string) {
	log.Info("storage: setting up connect postgres.")
	var err error
	for {
		db, err = sqlx.Open("postgres", dsn)
		if err != nil {
			log.Errorf("connect postgres error: %s, will retry 2 seconds\n", err)
			time.Sleep(time.Second * 2)
		}
		log.Info("connect postgres success.")

		break
	}
}

// Migrate migrate to postgres
func Migrate() error {
	m := &migrate.AssetMigrationSource{
		Asset:    Asset,
		AssetDir: AssetDir,
		Dir:      "",
	}

	n, err := migrate.Exec(db.DB, "postgres", m, migrate.Up)
	if err != nil {
		return err
	}

	log.Infof("migrate success: %d", n)
	return nil
}
