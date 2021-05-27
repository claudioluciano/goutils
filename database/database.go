package database

import (
	"fmt"

	"github.com/lithammer/shortuuid"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"

	"github.com/claudioluciano/goutils/logger"
)

type DB struct {
	table  string
	logger *logger.Logger
	gormDB *gorm.DB
}

type NewPostgresOpts struct {
	Table    string
	Host     string
	Port     string
	DBName   string
	User     string
	Password string
	Logger   *logger.Logger
}

func NewPostgres(opts *NewPostgresOpts) (*DB, error) {
	dsn := fmt.Sprintf("host=%v port=%v user=%v password=%v dbname=%v", opts.Host, opts.Port, opts.User, opts.Password, opts.DBName)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		opts.Logger.Error("db error when initialize the database", err)

		return nil, err
	}

	return &DB{
		table:  opts.Table,
		logger: opts.Logger,
		gormDB: db,
	}, nil
}

type NewSqliteOpts struct {
	Table  string
	DBName string
	Logger *logger.Logger
}

func NewSqlite(opts *NewSqliteOpts) (*DB, error) {
	db, err := gorm.Open(sqlite.Open(opts.DBName), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		opts.Logger.Error("db error when initialize the database", err)

		return nil, err
	}

	return &DB{
		table:  opts.Table,
		logger: opts.Logger,
		gormDB: db,
	}, nil

}

func (db *DB) GormDB() *gorm.DB {
	return db.gormDB
}

func (db *DB) DropTable() error {
	return db.gormDB.Migrator().DropTable(db.table)
}

func (db *DB) AutoMigrate(models ...interface{}) error {
	return db.gormDB.Migrator().AutoMigrate(models...)
}

func (db *DB) Create(target interface{}) error {
	db.gormDB.Transaction(func(tx *gorm.DB) error {
		// do some database operations in the transaction (use 'tx' from this point, not 'db')
		if err := tx.Table(db.table).Create(target).Error; err != nil {
			db.logger.Error("db error when create entity", err)
			// return any error will rollback
			return err
		}

		// return nil will commit the whole transaction
		return nil
	})

	return nil
}

func (db *DB) Update(target interface{}, newValues interface{}) error {
	db.gormDB.Transaction(func(tx *gorm.DB) error {
		// do some database operations in the transaction (use 'tx' from this point, not 'db')
		if err := tx.Table(db.table).Model(target).Updates(newValues).Error; err != nil {
			db.logger.Error("db error when update entity", err)
			// return any error will rollback
			return err
		}

		// return nil will commit the whole transaction
		return nil
	})

	return nil
}

func (db *DB) Delete(target interface{}) error {
	db.gormDB.Transaction(func(tx *gorm.DB) error {
		// do some database operations in the transaction (use 'tx' from this point, not 'db')
		if err := tx.Table(db.table).Delete(target).Error; err != nil {
			db.logger.Error("db error when delete entity", err)
			// return any error will rollback
			return err
		}

		// return nil will commit the whole transaction
		return nil
	})

	return nil
}

func (db *DB) FindByID(target interface{}, id string) error {
	db.gormDB.Transaction(func(tx *gorm.DB) error {
		// do some database operations in the transaction (use 'tx' from this point, not 'db')
		if err := tx.Table(db.table).First(target, "id = ?", id).Error; err != nil {
			db.logger.Error("db error when find by id entity", err)
			// return any error will rollback
			return err
		}

		// return nil will commit the whole transaction
		return nil
	})

	return nil
}

func (db *DB) Query(target interface{}, query string, orderBy string, args ...interface{}) error {
	db.gormDB.Transaction(func(tx *gorm.DB) error {
		// do some database operations in the transaction (use 'tx' from this point, not 'db')
		q := tx.Table(db.table).Where(query, args...)
		if orderBy != "" {
			q.Order(orderBy)
		}

		if err := q.Find(target).Error; err != nil {
			db.logger.Error("db error when query first entity", err)
			// return any error will rollback
			return err
		}

		// return nil will commit the whole transaction
		return nil
	})

	return nil
}

func (db *DB) Exec(raw string, args ...interface{}) error {
	db.gormDB.Transaction(func(tx *gorm.DB) error {
		// do some database operations in the transaction (use 'tx' from this point, not 'db')
		if err := tx.Exec(raw, args...).Error; err != nil {
			db.logger.Error("db error when create entity", err)
			// return any error will rollback
			return err
		}

		// return nil will commit the whole transaction
		return nil
	})

	return nil
}

func (db *DB) NewID(prefix string) string {
	var newID string

	if prefix != "" {
		prefix += "_"
	}

	newID = shortuuid.New()

	return prefix + newID
}
