package main

import (
	"bakery-project/entities"
	"context"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
	"log"
	"time"
)

const (
	dbname = "bakery"
)

var (
	ctx context.Context
	db  *sql.DB
)

func main() {
	db, err := sql.Open("mysql", "root:admin123@tcp(127.0.0.1:3306)/"+dbname)
	if err != nil {
		panic(err.Error()) // Just for example purpose. You should use proper error handling instead of panic
	}
	defer db.Close()
	db.SetConnMaxLifetime(time.Minute * 5)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
	db.SetConnMaxIdleTime(time.Minute * 3)

	dropDbs(db)
	createDbs(db)

	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()
	res, err := db.ExecContext(ctx, "CREATE DATABASE IF NOT EXISTS "+dbname)
	if err != nil {
		log.Printf("Error %s when creating DB\n", err)
		return
	}
	no, err := res.RowsAffected()
	if err != nil {
		log.Printf("Error %s when fetching rows", err)
		return
	}
	log.Printf("rows affected %d\n", no)

	db.Close()
	db, err = sql.Open("mysql", "root:admin123@tcp(localhost:3306)/bakery")
	if err != nil {
		log.Printf("Error %s when opening DB", err)
		return
	}
	defer db.Close()

	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(20)
	db.SetConnMaxLifetime(time.Minute * 5)

	ctx, cancelfunc = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()
	err = db.PingContext(ctx)
	if err != nil {
		log.Printf("Errors %s pinging DB", err)
		return
	}
	log.Printf("Connected to DB %s successfully\n", dbname)

	// Insert baked goods

	bakedGoods := []entities.BakedGood{
		{Name: "Chocolate cake"},
		{Name: "Red velvet cake"},
		{Name: "Cheesecake"},
		{Name: "Profiteroles"},
	}
	stmt, err := db.Prepare("INSERT INTO bakedGoods(name) VALUES( ? )")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close() // Prepared statements take up server resources and should be closed after use.

	for i, c := range bakedGoods {
		res, err := stmt.Exec(c.Name)
		if err != nil {
			log.Fatal(err)
		}
		numRows, err := res.RowsAffected()
		if err != nil || numRows != 1 {
			log.Fatal("Error inserting new Baked Good", err)
		}
		insId, err := res.LastInsertId()
		if err != nil {
			log.Fatal(err)
		}
		bakedGoods[i].ID = insId
	}

	// Insert users
	users := []entities.User{
		{Email: "linus@linux.com", FirstName: "Linus", LastName: "Torvalds", Username: "linus", Password: "linus", DeliveryAddress: "Sofia,Reduta", IsAdmin: true},
		{Email: "gosling@java.com", FirstName: "James", LastName: "Gosling", Username: "james", Password: "james", DeliveryAddress: "Oborishte", IsAdmin: false},
		{Email: "pike@golang.com", FirstName: "Rob", LastName: "Pike", Username: "rob", Password: "rob", DeliveryAddress: "Varna", IsAdmin: false},
		{Email: "kamel@docker.com", FirstName: "Kamel", LastName: "Founadi", Username: "kamel", Password: "kamel", DeliveryAddress: "Lulin", IsAdmin: false},
	}

	stmt, err = db.Prepare(
		`INSERT INTO users(email, first_name, last_name, username, password, delivery_address, isAdmin, created, modified) 
		VALUES( ?, ?, ?, ?, ?, ?, ?, ?,?)`)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close() // Prepared statements take up server resources and should be closed after use.

	for i := range users {
		users[i].Created = time.Now()
		users[i].Modified = time.Now()
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(users[i].Password), bcrypt.DefaultCost)
		if err != nil {
			panic(err)
		}
		users[i].Password = "{bcrypt}" + string(hashedPassword)
		result, err := stmt.Exec(users[i].Email, users[i].FirstName, users[i].LastName, users[i].Username,
			users[i].Password, users[i].DeliveryAddress, users[i].IsAdmin, users[i].Created, users[i].Modified)
		if err != nil {
			log.Fatal(err)
		}
		numRows, err := result.RowsAffected()
		if err != nil || numRows != 1 {
			log.Fatal("Error inserting new User", err)
		}
		insId, err := result.LastInsertId()
		if err != nil {
			log.Fatal(err)
		}
		users[i].ID = insId
	}
}

func dropDbs(db *sql.DB) {
	res, err := db.Exec("DROP TABLE IF EXISTS `review`")
	if err != nil {
		log.Fatal(err)
	}
	rowsAffected, err := res.RowsAffected()
	log.Printf("'review' - Rows Affected: %d %v", rowsAffected, err)

	res, err = db.Exec("DROP TABLE IF EXISTS `orderedGoods`")
	if err != nil {
		log.Fatal(err)
	}
	rowsAffected, err = res.RowsAffected()
	log.Printf("'ordered goods' - Rows Affected: %d %v", rowsAffected, err)

	res, err = db.Exec("DROP TABLE IF EXISTS `bakedGoodsTags`")
	if err != nil {
		log.Fatal(err)
	}
	rowsAffected, err = res.RowsAffected()
	log.Printf("'bakedGoodstags' - Rows Affected: %d %v", rowsAffected, err)

	res, err = db.Exec("DROP TABLE IF EXISTS `tags`")
	if err != nil {
		log.Fatal(err)
	}
	rowsAffected, err = res.RowsAffected()
	log.Printf("tags - Rows Affected: %d %v", rowsAffected, err)

	res, err = db.Exec("DROP TABLE IF EXISTS `orders`")
	if err != nil {
		log.Fatal(err)
	}
	rowsAffected, err = res.RowsAffected()
	log.Printf("orders - Rows Affected: %d %v", rowsAffected, err)

	res, err = db.Exec("DROP TABLE IF EXISTS `bakedGoods`")
	if err != nil {
		log.Fatal(err)
	}
	rowsAffected, err = res.RowsAffected()
	log.Printf("bakedgoods - Rows Affected: %d %v", rowsAffected, err)

	res, err = db.Exec("DROP TABLE IF EXISTS `users`")
	if err != nil {
		log.Fatal(err)
	}
	rowsAffected, err = res.RowsAffected()
	log.Printf("users - Rows Affected: %d %v", rowsAffected, err)

}

func createDbs(db *sql.DB) {
	res, err := db.Exec("CREATE TABLE `users` (`id` bigint(20) NOT NULL AUTO_INCREMENT, `email` varchar(255) DEFAULT NULL, `first_name` varchar(255) DEFAULT NULL, `last_name` varchar(255) DEFAULT NULL, `username` varchar(30) DEFAULT NULL, `password` varchar(255) DEFAULT NULL, `delivery_address` varchar(255) DEFAULT NULL, `isAdmin` tinyint(1) DEFAULT 0, `created` datetime(6) DEFAULT NULL, `modified` datetime(6) DEFAULT NULL, PRIMARY KEY (`id`), UNIQUE KEY `UK_6dotkott2kjsp8vw4d0m25fb7` (`email`), UNIQUE KEY `UK_r43af9ap4edm43mmtq01oddj6` (`username`)) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;")
	if err != nil {
		log.Fatal(err)
	}
	rowsAffected, err := res.RowsAffected()
	log.Printf("'companies' - Rows Affected: %d %v", rowsAffected, err)

	res, err = db.Exec("CREATE TABLE `bakedGoods` (`id` bigint(20) NOT NULL AUTO_INCREMENT, `name` varchar(255) DEFAULT NULL, `photoUrl` varchar(255), `price` double DEFAULT NULL, PRIMARY KEY (`id`)) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;")
	if err != nil {
		log.Fatal(err)
	}
	rowsAffected, err = res.RowsAffected()
	log.Printf("'projects' - Rows Affected: %d %v", rowsAffected, err)

	res, err = db.Exec("CREATE TABLE `tags` (`id` bigint(20) NOT NULL AUTO_INCREMENT, `tagName` varchar(255) DEFAULT NULL, PRIMARY KEY (`id`)) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;")
	if err != nil {
		log.Fatal(err)
	}
	rowsAffected, err = res.RowsAffected()
	log.Printf("'users' - Rows Affected: %d %v", rowsAffected, err)

	res, err = db.Exec("CREATE TABLE `bakedGoodsTags` (`id` bigint(20) NOT NULL AUTO_INCREMENT, `tagId` bigint(20) DEFAULT NULL, `bakedGoodId` bigint(20) DEFAULT NULL, PRIMARY KEY (`id`), FOREIGN KEY (`tagId`) REFERENCES `tags` (`id`), FOREIGN KEY (`bakedGoodId`) REFERENCES `bakedGoods` (`id`)) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;")
	if err != nil {
		log.Fatal(err)
	}
	rowsAffected, err = res.RowsAffected()
	log.Printf("user_roles - Rows Affected: %d %v", rowsAffected, err)

	res, err = db.Exec("CREATE TABLE `orders` (`id` bigint(20) NOT NULL AUTO_INCREMENT,`userId` bigint(20) DEFAULT NULL,`status` varchar(50) DEFAULT NULL,`deliveryAddress` varchar(255) DEFAULT NULL,PRIMARY KEY (`id`),FOREIGN KEY (`userId`) REFERENCES `users` (`id`)) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;")
	if err != nil {
		log.Fatal(err)
	}
	rowsAffected, err = res.RowsAffected()
	log.Printf("user_roles - Rows Affected: %d %v", rowsAffected, err)

	res, err = db.Exec("CREATE TABLE `orderedGoods` (`id` bigint(20) NOT NULL AUTO_INCREMENT,`bakedGoodId` bigint(20) DEFAULT NULL,`orderId` bigint(20) DEFAULT NULL,PRIMARY KEY (`id`),FOREIGN KEY (`bakedGoodId`) REFERENCES `bakedGoods` (`id`),FOREIGN KEY (`orderId`) REFERENCES `orders` (`id`)) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;")
	if err != nil {
		log.Fatal(err)
	}
	rowsAffected, err = res.RowsAffected()
	log.Printf("user_roles - Rows Affected: %d %v", rowsAffected, err)

	res, err = db.Exec("CREATE TABLE `review` (`id` bigint(20) NOT NULL AUTO_INCREMENT,`userId` bigint(20) DEFAULT NULL,`reviewText` varchar(255) DEFAULT NULL,PRIMARY KEY (`id`),FOREIGN KEY (`userId`) REFERENCES `users` (`id`)) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;")
	if err != nil {
		log.Fatal(err)
	}
	rowsAffected, err = res.RowsAffected()
	log.Printf("user_roles - Rows Affected: %d %v", rowsAffected, err)
}
