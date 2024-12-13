package databaseHelper

import (
	"context"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

func TransactionHandler(ctx context.Context, txFunc func(tx *gorm.DB) error) error {
	ctx, span := tracer.Start(ctx, "HandleTransaction")
	defer span.End()

	dbh := GetDB(ctx)

	// Start the transaction
	tx := dbh.Begin()
	if tx.Error != nil {
		err := errors.Wrap(tx.Error, "Failed to start the transaction")
		return err
	}

	// Execute the transaction function
	if err := txFunc(tx); err != nil {
		// Rollback if the creation failed
		tx.Rollback()

		err = errors.Wrap(err, "Transaction failed, rolled back")
		return err
	}

	// Commit the transaction if everything went well
	if err := tx.Commit().Error; err != nil {
		err = errors.Wrap(err, "Failed to commit the transaction")
		return err
	}

	return nil
}
