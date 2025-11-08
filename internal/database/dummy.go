package database

func (db *DB) InsertDummy() error {
	query := `CREATE TABLE IF NOT EXISTS user (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	name TEXT NOT NULL
);

INSERT INTO user (name) VALUES
	('Kevin'), ('Mike'), ('Brandon');
`

	if _, err := db.Exec(query); err != nil {
		return err
	}

	return nil
}
