package mysql

import (
	"database/sql"

	"github.com/yunushamod/snippetbox/pkg/models"
)

// Define a SnippetModel type which wraps a sql.DB connection pool
type SnippetModel struct {
	DB *sql.DB
}

// This will insert a new snippet into the database.
func (m *SnippetModel) Insert(title, content, expires string) (int, error) {
	//Write the SQL statement we want to execute. I've split it over two lines for readability (which is why
	// it's surrounded with backquotes instead of normal double quotes),
	stmt := `INSERT INTO snippets (title, content, created, expires) VALUES (?,?,UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`
	// Use the Exec() method on the embedded connection pool to execute the
	// statement.
	//This method returns a sql.Result object, which contains some basic information about what happened when the statement
	// was executed.
	result, err := m.DB.Exec(stmt, title, content, expires)
	if err != nil {
		return 0, err
	}
	// Use the LastInsertId() method on the result object to get the ID of our newly inserted
	// record in the snippets table.
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	// The ID returned has type int64, so we convert it to an int type
	// before returning
	return int(id), nil
}

// This will return a specific snippet based on its id.
func (m *SnippetModel) Get(id int) (*models.Snippet, error) {
	// Write the SQL statement we want to execute.
	stmt := `SELECT id, title, content, created, expires FROM snippets where expires > UTC_TIMESTAMP() and id = ?`
	// Use the QueryRow() method on the connection pool to execute our SQL statement, passing
	// in the undertrusted id variable as the value for the placeholder parameter.
	// This returns a pointer to a sql.Row object which holds the result from the database.
	row := m.DB.QueryRow(stmt, id)
	// Initialize a pointer to a new zeroed Snippet struct
	s := &models.Snippet{}
	// Use row.Scan to copy tthe values from each field in sql.Row to the corresponding field in the
	// Snippet struct. Notice that the arguments
	// to row.Scan are *pointers* to the place you want to copy the data into,
	// and the number of arguments must be exactly the same as the number of
	// columns returned by your statement. If the query returns no rows, then
	// row.Scan() will return a sql.ErrNoRows error. We check for that and return
	err := row.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
	if err == sql.ErrNoRows {
		return nil, models.ErrNoRecord
	} else if err != nil {
		return nil, err
	}
	return s, nil
}

// This will return the 10 most recently created snippets.
func (m *SnippetModel) Latest() ([]*models.Snippet, error) {
	// Write the SQL statement we want to execute.
	stmt := `SELECT id, title, content, created, expires from snippets where expires > UTC_TIMESTAMP() ORDER BY created DESC LIMIT 10`
	// Use the Query() method on the connection pool to execute our
	// SQL statement. This returns a sql.Rows result set containing
	// the result of our query.
	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	// We defer rows.Close() to ensure resultset is always properly closed
	// after the Latest() method returns. This defer statement should *after* you
	// check for an error from the Query() method. Otherwise, if Query() returns an error,you'll
	// get a panic trying to close a nil resultset.
	defer rows.Close()
	// Initialize an empty slice to hold the models.Snippets objects.
	snippets := []*models.Snippet{}
	// Use rows.Next to iterate through the rows in the resultset.
	// This prepares the first (and then each subsequent) row to be acted on
	// rows.Can() method.
	for rows.Next() {
		// Create a pointer to a new zeroed Snippet struct.
		s := &models.Snippet{}
		err = rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
		if err != nil {
			return nil, err
		}
		snippets = append(snippets, s)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	// If everything went OK then return the Snippets slice.
	return snippets, nil
}
