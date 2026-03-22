package postgres

import (
	"context"
	"database/sql"
	"errors"
	"github.com/lib/pq"
	"transaction_routine/internal/domain"
)

const PostgresDuplicatedValue = "23505"

type AccountRepository struct {
	db *sql.DB
}

func NewAccountRepository(db *sql.DB) *AccountRepository {
	return &AccountRepository{db: db}
}

func (ar *AccountRepository) Insert(ctx context.Context, documentHash string, documentEncrypted string) (*domain.Account, error) {
	var account domain.Account
	const query = `
					INSERT INTO accounts(document_hash, document_encrypted) 
					VALUES ($1, $2) 
					RETURNING account_id, document_hash, document_encrypted
					`
	err := ar.db.QueryRowContext(ctx, query, documentHash, documentEncrypted).Scan(&account.AccountId, &account.DocumentHash, &account.DocumentEncrypted)

	if err != nil {
		var pgErr *pq.Error

		if errors.As(err, &pgErr) {
			if pgErr.Code == PostgresDuplicatedValue { // document_hash is unique
				return nil, domain.ErrAccountAlreadyExists
			}
		}
		return nil, err

	}
	return &account, nil
}

func (ar *AccountRepository) GetById(ctx context.Context, accountId int64) (*domain.Account, error) {
	var account domain.Account
	const query = `
		SELECT account_id, document_hash, document_encrypted
		FROM accounts
		WHERE account_id = $1
	`

	err := ar.db.QueryRowContext(ctx, query, accountId).Scan(
		&account.AccountId,
		&account.DocumentHash,
		&account.DocumentEncrypted,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrAccountNotFound
		}

		return nil, err
	}

	return &account, nil
}

func (ar *AccountRepository) ExistsById(ctx context.Context, accountId int64) (bool, error) {
	var exists bool
	const query = `
        SELECT EXISTS (
            SELECT 1
            FROM accounts
            WHERE account_id = $1
        )
    `

	err := ar.db.QueryRowContext(ctx, query, accountId).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}
