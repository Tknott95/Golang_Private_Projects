package sqlModel

import "database/sql"

type SQLStore struct {
	DB *sql.DB
}
