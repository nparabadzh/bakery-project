package api

import (
	"bakery-project/entities"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

type TagsRepoMySql struct {
	db *sql.DB
}

func (u TagsRepoMySql) FindAll(start, count int) ([]entities.Tag, error) {
	statement := "SELECT * FROM tags LIMIT ? OFFSET ?"
	rows, err := u.db.Query(statement, count, start)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	tags := []entities.Tag{}
	for rows.Next() {
		var tag entities.Tag
		err := rows.Scan(&tag.ID, &tag.TagName)
		if err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}
	rows.Close()
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return tags, nil
}

//FindById return orders by order ID or error otherwise
func (u *TagsRepoMySql) FindByID(id int) (*entities.Tag, error) {
	tag := &entities.Tag{}
	statement := "SELECT * from tags WHERE id = ?"
	err := u.db.QueryRow(statement, id).Scan(&tag.ID, &tag.TagName)
	if err != nil {
		return nil, err
	}
	return tag, nil
}

//Create creates and returns new order with autogenerated ID
func (u *TagsRepoMySql) Create(tag *entities.Tag) (*entities.Tag, error) {
	statement := "INSERT INTO tags(tagName) VALUES(?)"
	result, err := u.db.Exec(statement, tag.TagName)
	if err != nil {
		return nil, err
	}
	id, err := result.LastInsertId()
	tag.ID = int(id)
	//err = r.db.QueryRow("SELECT LAST_INSERT_ID()").Scan(&user.ID)
	if err != nil {
		return nil, err
	}
	return tag, nil
}

//Update updates existing order data
func (u *TagsRepoMySql) Update(tag *entities.Tag) (*entities.Tag, error) {
	statement := "UPDATE tags SET tagName=? WHERE id=? "
	_, err := u.db.Exec(statement, tag.TagName, tag.ID)
	if err != nil {
		return nil, err
	}
	return u.FindByID(int(tag.ID))
}

//DeleteById removes and returns order with specified ID or error otherwise
func (u *TagsRepoMySql) DeleteByID(id int) (*entities.Tag, error) {
	tag, err := u.FindByID(id)
	if err != nil {
		return nil, err
	}
	_, err = u.db.Exec("DELETE FROM tags WHERE id=?", id)
	return tag, err
}

func NewTagsRepoMysql(user, password, dbname string) *TagsRepoMySql {
	connectionString := fmt.Sprintf("%s:%s@/%s?parseTime=true", user, password, dbname)
	repo := &TagsRepoMySql{}
	var err error
	repo.db, err = sql.Open("mysql", connectionString)
	if err != nil {
		log.Fatal(err)
	}
	return repo
}