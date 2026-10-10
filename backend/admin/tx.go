package main

import (
	"context"
	"database/sql"
	"time"

	"github.com/gofiber/fiber/v2"
)

func transaction(db *sql.DB, txOptions *sql.TxOptions) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()
		c.SetUserContext(ctx)

		tx, err := db.BeginTx(ctx, txOptions)
		if err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}
		c.Locals("tx", tx)

		err = c.Next()
		if c.Locals("txClosed") != nil {
			return err
		}
		if err != nil {
			_ = tx.Rollback()
			return err
		}
		return tx.Commit()
	}
}

func readTx(db *sql.DB) fiber.Handler {
	tx := transaction(db, &sql.TxOptions{ReadOnly: true})
	return func(c *fiber.Ctx) error {
		if c.Method() != fiber.MethodGet {
			return c.Next()
		}
		return tx(c)
	}
}

func writeTx(db *sql.DB) fiber.Handler {
	tx := transaction(db, nil)
	return func(c *fiber.Ctx) error {
		if c.Method() != fiber.MethodPost || c.Path() == "/logout" {
			return c.Next()
		}
		if err := checkCSRF(c); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}
		return tx(c)
	}
}
