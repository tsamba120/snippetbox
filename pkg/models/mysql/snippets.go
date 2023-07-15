// this file contains code specifically for working with the snippets in our mysql database
// we'll define a SnippetModel type and implement models against it for CRUD

package mysql

import (
	"database/sql"
	"errors"

	// Import the models package that we just created. You need to prefix this with
	// whatever module path you set up
	"github.com/tsamba120/snippetbox/pkg/models"
)

// Define a SnippetModel type which wraps a sql.DB connection pool.
type SnippetModel struct {
	DB *sql.DB
}

func (m *SnippetModel) Insert(title, content, expires string) (int, error) {
	// use backticks for multiline strings
	// '?' character indicates placeholder params for non-validated user input
	// best practice to include these instead of interpolating data in the query
	stmt := `
		INSERT INTO snippets (title, content, created, expires)
		VALUES(?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY));
	`
	// insert values for placeholder params here
	// result is of type `sql.Result` which is an interface
	result, err := m.DB.Exec(stmt, title, content, expires)
	if err != nil {
		return 0, err
	}

	// LastInsertId() to get ID of newly inserted record in snippets table
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	// cast ID from type int64 to an general int type
	return int(id), nil
}

func (m *SnippetModel) Get(id int) (*models.Snippet, error) {

	stmt := `
		SELECT
			id
			, title
			, content
			, created
			, expires
		FROM snippets
		WHERE
			expires > UTC_TIMESTAMP()
			AND id = ?;
	`

	// initialize a pointer to a new zeroed Snippet struct object
	s := &models.Snippet{}

	// Execute query, returning a pointer to a sql.Row object
	// use row.Scan() to map row values to pointer. must pass in pointers to destination
	err := m.DB.QueryRow(stmt, id).Scan(
		&s.ID,
		&s.Title,
		&s.Content,
		&s.Created,
		&s.Expires,
	)

	if err != nil {
		// if row returned no rows then return models.ErrNoRecord
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}

	return s, nil
}

func (m *SnippetModel) Latest() ([]*models.Snippet, error) {

	stmt := `
		SELECT
			id
			, title
			, content
			, created
			, expires
		FROM snippets
		WHERE
			expires > UTC_TIMESTAMP()
		ORDER BY created DESC
		LIMIT 10;
	`

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}

	// ensure the resultset is closed after execution of function!!
	// important to do this *after* checking for errors from DB.Query()
	defer rows.Close()

	// initialize an empty slice to hold the models.Snippets objects
	snippets := []*models.Snippet{}

	// Use rows.Next to iterate through the rows in the resultset. This
	// prepares the first (and then each subsequent) row to be acted on by the
	// rows.Scan() method. If iteration over all the rows completes then the // resultset automatically closes itself and frees-up the underlying
	// database connection.
	for rows.Next() {
		s := &models.Snippet{}
		err := rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
		if err != nil {
			return nil, err
		}

		// append populated model to the slice
		snippets = append(snippets, s)
	}

	// When the rows.Next() loop has finished we call rows.Err() to retrieve any
	// error that was encountered during the iteration. It's important to
	// call this - don't assume that a successful iteration was completed
	// over the whole resultset
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return snippets, err
}
