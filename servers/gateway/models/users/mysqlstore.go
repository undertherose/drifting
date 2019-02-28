package users

import (
	"database/sql"
	"fmt"
)

const sqlInsertTask = "insert into users (email, passHash, userName, firstName, lastName, photoURL) values (?,?,?,?,?,?)"
const sqlSelectAll = "select * from users"
const sqlSelectTrie = "select id, userName, firstName, lastName from users"
const sqlSelectByID = sqlSelectAll + " where id=?"
const sqlSelectByEmail = sqlSelectAll + " where email=?"
const sqlSelectByUsername = sqlSelectAll + " where userName=?"
const sqlUpdate = "update users set firstName=?, lastName=? where id=?"
const sqlDel = "delete from users where id=?"

//MySQLStore represents a user.Store backed by mySQL
type MySQLStore struct {
	//DB pointer to be used to talk to SQL store
	db *sql.DB
}

//NewMySQLStore constructs a new MySQLStore
func NewMySQLStore(db *sql.DB) *MySQLStore {
	//initialize and return a new MySQLStore struct
	if db == nil {
		panic("nil database pointer")
	}
	return &MySQLStore{db}
}

//Store implementation

//Insert inserts the `user` into the store
func (ms *MySQLStore) Insert(user *User) (*User, error) {

	res, err := ms.db.Exec(sqlInsertTask, user.PassHash, user.UserName)
	if err != nil {
		fmt.Printf("error inserting new row: %v\n", err)
		return nil, err
	}
	//get the auto-assigned ID for the new row
	id, err := res.LastInsertId()
	if err != nil {
		fmt.Printf("error getting new ID: %v\n", id)
		return nil, err
	}

	user.ID = id
	return user, nil
}

//Update updates the user with `id` with values in `updates`
/* func (ms *MySQLStore) Update(id int64, updates *Updates) (*User, error) {
	results, err := ms.db.Exec(sqlUpdate, updates.FirstName, updates.LastName, id)
	if err != nil {
		return nil, fmt.Errorf("updating: %v", err)
	}
	affected, err := results.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("getting rows affected: %v", err)
	}
	//if no rows were affected, then the requested
	//ID was not in the database
	if affected == 0 {
		return nil, ErrUserNotFound
	}
	return ms.GetByID(id)
} */

//GetByID gets the user details with the specified `id`
func (ms *MySQLStore) GetByID(id int64) (*User, error) {
	row := ms.db.QueryRow(sqlSelectByID, id)
	user := &User{}

	if err := row.Scan(&user.ID, &user.PassHash, &user.UserName); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("scanning: %v", err)
	}
	return user, nil
}

//GetByEmail gets the user details with the specified `id`
func (ms *MySQLStore) GetByEmail(email string) (*User, error) {
	row := ms.db.QueryRow(sqlSelectByEmail, email)
	user := &User{}

	if err := row.Scan(&user.ID, &user.PassHash, &user.UserName); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("scanning: %v", err)
	}
	return user, nil
}

//GetByUserName gets the user details with the specified `id`
func (ms *MySQLStore) GetByUserName(username string) (*User, error) {
	row := ms.db.QueryRow(sqlSelectByUsername, username)
	user := &User{}

	if err := row.Scan(&user.ID, &user.PassHash, &user.UserName); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("scanning: %v", err)
	}
	return user, nil
}

//Delete deletes the user with the give `id`
func (ms *MySQLStore) Delete(id int64) error {
	_, err := ms.db.Exec(sqlDel, id)
	if err != nil {
		return err
	}
	return nil
}

//GetAll gets all the users from the db
func (ms *MySQLStore) GetAll() ([]*User, error) {
	rows, err := ms.db.Query(sqlSelectAll)

	users := []*User{}

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		user := &User{}
		if err := rows.Scan(&user.ID, &user.PassHash, &user.UserName); err != nil {
			if err == sql.ErrNoRows {
				return nil, ErrUserNotFound
			}
			return nil, fmt.Errorf("scanning: %v", err)
		}

		users = append(users, user)
	}

	return users, nil

}

const sqlSelectUserType = "select type from users where username=?"

//find user with specific URL
func (ms *MySQLStore) GetUserTypeByUsername(username string) (string, error) {
	row := ms.db.QueryRow(sqlSelectByUsername, username)
	user := &User{}

	if err := row.Scan(&user.ID, &user.PassHash, &user.UserName, &user.Type); err != nil {
		if err == sql.ErrNoRows {
			return "", ErrUserNotFound
		}
		return "", fmt.Errorf("scanning: %v", err)
	}
	return user.Type, nil

}
