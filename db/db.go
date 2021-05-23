package db

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/claudioluciano/goutils/logger"
)

type DB struct {
	logger *logger.Logger

	gormDB *gorm.DB
}

type NewPostgressOpts struct {
	Host     string
	Port     string
	DbName   string
	User     string
	Password string
	Logger   *logger.Logger
}

func NewPostgress(opts *NewPostgressOpts) (*DB, error) {
	// TODO change to environment variable

	dsn := fmt.Sprintf("host=%v port=5432 user=%v password=%v dbname=%v", opts.Host, opts.Port, opts.User, opts.Password)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		opts.Logger.ErrorWithError("db error when initialize the database", err)

		return nil, err
	}

	return &DB{
		logger: opts.Logger,
		gormDB: db,
	}, nil
}

func (db *DB) AutoMigrate(models ...interface{}) error {
	return db.gormDB.AutoMigrate(models)
}

func (db *DB) Create(target interface{}) error {
	db.gormDB.Transaction(func(tx *gorm.DB) error {
		// do some database operations in the transaction (use 'tx' from this point, not 'db')
		if err := tx.Create(target).Error; err != nil {
			db.logger.ErrorWithError("db error when create entity", err)

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
		if err := tx.Model(target).Updates(newValues).Error; err != nil {
			db.logger.ErrorWithError("db error when update entity", err)
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
		if err := tx.Delete(target).Error; err != nil {
			db.logger.ErrorWithError("db error when delete entity", err)
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
		if err := tx.First(target, id).Error; err != nil {
			db.logger.ErrorWithError("db error when find by id entity", err)
			// return any error will rollback
			return err
		}

		// return nil will commit the whole transaction
		return nil
	})

	return nil
}

func (db *DB) FindMany(targets []interface{}, limit int) error {
	db.gormDB.Transaction(func(tx *gorm.DB) error {
		// do some database operations in the transaction (use 'tx' from this point, not 'db')
		if err := tx.Limit(limit).Find(targets).Error; err != nil {
			db.logger.ErrorWithError("db error when find many entity", err)

			// return any error will rollback
			return err
		}

		// return nil will commit the whole transaction
		return nil
	})

	return nil
}

func (db *DB) QueryFirst(query string, target interface{}) error {
	db.gormDB.Transaction(func(tx *gorm.DB) error {
		// do some database operations in the transaction (use 'tx' from this point, not 'db')
		if err := tx.Where(query).Find(target).Error; err != nil {
			db.logger.ErrorWithError("db error when query first entity", err)
			// return any error will rollback
			return err
		}

		// return nil will commit the whole transaction
		return nil
	})

	return nil
}

func (db *DB) QueryMany(query string, targets []interface{}) error {
	db.gormDB.Transaction(func(tx *gorm.DB) error {
		// do some database operations in the transaction (use 'tx' from this point, not 'db')
		if err := tx.Where(query).Find(targets).Error; err != nil {
			db.logger.ErrorWithError("db error when query many entity", err)
			// return any error will rollback
			return err
		}

		// return nil will commit the whole transaction
		return nil
	})

	return nil
}
