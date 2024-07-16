package postgres

import "database/sql"

type repositoryImpl struct {
	conn *sql.DB
}

func (r *repositoryImpl) GetByOriginalLink(link []byte) (string, error) {
	rows := r.conn.QueryRow("SELECT token FROM link where original=$1", link)
	if err := rows.Err(); err != nil {
		return "", err
	}

	var token string
	err := rows.Scan(&token)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil
		}
		return "", err
	}

	return token, nil
}

func (r *repositoryImpl) Create(token string, original []byte) error {
	_, err := r.conn.Exec("INSERT INTO link  (token, original) values ($1, $2)", token, original)
	if err != nil {
		return err
	}
	return nil
}

func (r *repositoryImpl) GetByToken(token string) (string, error) {
	rows := r.conn.QueryRow("SELECT original FROM link where token=$1", token)
	var link string
	err := rows.Scan(&link)
	if err != nil {
		return "", err
	}
	return link, nil
}
