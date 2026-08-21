package repository

import "database/sql"

// DB is alias for *sql.DB to allow future interface mocking.
type DB = sql.DB

// Tx is alias for *sql.Tx
type Tx = sql.Tx
