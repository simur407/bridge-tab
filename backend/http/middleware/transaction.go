package middleware

import (
	"context"
	"database/sql"
	"time"

	"github.com/gofiber/fiber/v2"
)

func Transaction(db *sql.DB, txOptions *sql.TxOptions) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
		defer cancel()

		c.SetUserContext(ctx)

		tx, txErr := db.BeginTx(ctx, txOptions)
		if txErr != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		c.Locals("tx", tx)
		err := c.Next()

		if err != nil {
			tx.Rollback()
			return err
		}

		return tx.Commit()
	}
}

func GetTransaction(c *fiber.Ctx) *sql.Tx {
	tx, ok := c.Locals("tx").(*sql.Tx)

	if !ok {
		panic("transaction not hooked")
	}

	return tx
}
