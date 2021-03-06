package api

import (
	"bakery-project/entities"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

type CategoriesRepoMySql struct {
	db *sql.DB
}

func (u CategoriesRepoMySql) FindAll(start, count int) ([]entities.Category, error) {
	statement := "SELECT * FROM categories LIMIT ? OFFSET ?"
	rows, err := u.db.Query(statement, count, start)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	categories := []entities.Category{}
	for rows.Next() {
		var category entities.Category
		err := rows.Scan(&category.ID, &category.CategoryName)
		if err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}
	rows.Close()
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return categories, nil
}

//FindById return orders by order ID or error otherwise
func (u *CategoriesRepoMySql) FindByID(id int) (*entities.Category, error) {
	category := &entities.Category{}
	statement := "SELECT * from categories WHERE id = ?"
	err := u.db.QueryRow(statement, id).Scan(&category.ID, &category.CategoryName)
	if err != nil {
		return nil, err
	}
	return category, nil
}

//Create creates and returns new order with autogenerated ID
func (u *CategoriesRepoMySql) Create(category *entities.Category) (*entities.Category, error) {
	statement := "INSERT INTO categories(category_name) VALUES(?)"
	result, err := u.db.Exec(statement, category.CategoryName)
	if err != nil {
		return nil, err
	}
	id, err := result.LastInsertId()
	category.ID = int(id)
	//err = r.db.QueryRow("SELECT LAST_INSERT_ID()").Scan(&user.ID)
	if err != nil {
		return nil, err
	}
	return category, nil
}

//Update updates existing order data
func (u *CategoriesRepoMySql) Update(category *entities.Category) (*entities.Category, error) {
	statement := "UPDATE categories SET category_name=? WHERE id=? "
	_, err := u.db.Exec(statement, category.CategoryName, category.ID)
	if err != nil {
		return nil, err
	}
	return u.FindByID(int(category.ID))
}

//DeleteById removes and returns order with specified ID or error otherwise
func (u *CategoriesRepoMySql) DeleteByID(id int) (*entities.Category, error) {
	category, err := u.FindByID(id)
	if err != nil {
		return nil, err
	}
	_, err = u.db.Exec("DELETE FROM categories WHERE id=?", id)
	return category, err
}

func NewCategoriesRepoMysql(user, password, dbname string) *CategoriesRepoMySql {
	connectionString := fmt.Sprintf("%s:%s@/%s?parseTime=true", user, password, dbname)
	repo := &CategoriesRepoMySql{}
	var err error
	repo.db, err = sql.Open("mysql", connectionString)
	if err != nil {
		log.Fatal(err)
	}
	return repo
}
