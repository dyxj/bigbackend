package faker

import (
	"context"
	"database/sql"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/mock"
)

type TransactionManagerMock struct {
	mock.Mock
	db      *sql.DB
	sqlMock sqlmock.Sqlmock
}

func NewTransactionManagerMock() (*TransactionManagerMock, error) {
	db, sqlMock, err := sqlmock.New()
	if err != nil {
		return nil, err
	}

	return &TransactionManagerMock{
		db:      db,
		sqlMock: sqlMock,
	}, nil
}

func (m *TransactionManagerMock) BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error) {
	returnArgs := m.Called(ctx, opts)
	return returnArgs.Get(0).(*sql.Tx), returnArgs.Error(1)
}

// ReturnCommitedTx sets up the mock to return a committed transaction
// Can be used to simulate error, as performing Tx.Commit returns sql.ErrTxDone
func (m *TransactionManagerMock) ReturnCommitedTx(args mock.Arguments) {
	tx, err := m.db.Begin()
	if err != nil {
		m.ExpectedCalls[0].ReturnArguments = mock.Arguments{nil, err}
	}
	_ = tx.Commit()
	m.ExpectedCalls[0].ReturnArguments = mock.Arguments{tx, err}
}

func (m *TransactionManagerMock) ReturnTx(args mock.Arguments) {
	tx, err := m.db.Begin()
	if err != nil {
		m.ExpectedCalls[0].ReturnArguments = mock.Arguments{nil, err}
	}
	m.ExpectedCalls[0].ReturnArguments = mock.Arguments{tx, err}
}

func (m *TransactionManagerMock) Close() error {
	return m.db.Close()
}

func (m *TransactionManagerMock) SqlMock() sqlmock.Sqlmock {
	return m.sqlMock
}
