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
	dbname   = "bakery"
	username = "root"
	password = "admin123"
)

var (
	ctx context.Context
	db  *sql.DB
)

func main() {
	db, err := sql.Open("mysql", username+":"+password+"@tcp(127.0.0.1:3306)/"+dbname)
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
	db, err = sql.Open("mysql", username+":"+password+"@tcp(localhost:3306)/"+dbname)
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

	//Insert categories
	categories := []entities.Category{
		{CategoryName: "Cakes"},
		{CategoryName: "Pies"},
		{CategoryName: "Tarts"},
		{CategoryName: "Cheesecakes"},
		{CategoryName: "Pastries"},
		{CategoryName: "Other"},
	}

	stmt, err := db.Prepare("INSERT INTO categories(category_name) VALUES( ? )")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close() // Prepared statements take up server resources and should be closed after use.

	for i, c := range categories {
		res, err := stmt.Exec(c.CategoryName)
		if err != nil {
			log.Fatal(err)
		}
		numRows, err := res.RowsAffected()
		if err != nil || numRows != 1 {
			log.Fatal("Error inserting new Category", err)
		}
		insId, err := res.LastInsertId()
		if err != nil {
			log.Fatal(err)
		}
		categories[i].ID = int(insId)
	}

	// Insert baked goods
	bakedGoods := []entities.BakedGood{
		{Name: "Chocolate cake", Price: 12.4,
			PhotoUrl: "https://img.taste.com.au/P9x2Yltr/taste/2016/11/homemade-chocolate-cake-85524-1.jpeg", CategoryId: 4},
		{Name: "Lemon cake", Price: 19.6,
			PhotoUrl: "https://www.recipetineats.com/wp-content/uploads/2021/09/Lemon-Cake-with-Lemon-Frosting_85-SQ.jpg", CategoryId: 4},
		{Name: "Red velvet cake", Price: 20,
			PhotoUrl: "https://www.cookingclassy.com/wp-content/uploads/2014/11/red-velvet-cake-5.jpg", CategoryId: 4},
		{Name: "Cheesecake", Price: 15.3,
			PhotoUrl: "https://natashaskitchen.com/wp-content/uploads/2020/05/Pefect-Cheesecake-7.jpg", CategoryId: 9},
		{Name: "Strawberry clafoutis", Price: 9.1,
			PhotoUrl: "https://thenewbaguette.com/wp-content/uploads/2019/06/strawberry-clafoutis-16.jpg", CategoryId: 9},
		{Name: "Croissants", Price: 7.1,
			PhotoUrl: "https://static01.nyt.com/images/2021/04/07/dining/06croissantsrex1/merlin_184841898_ccc8fb62-ee41-44e8-9ddf-b95b198b88db-articleLarge.jpg", CategoryId: 9},
		{Name: "Tart au Chocolat", Price: 22,
			PhotoUrl: "https://images.ricardocuisine.com/services/recipes/8284-portrait.jpg", CategoryId: 5},
		{Name: "Tart au Citron", Price: 12.3,
			PhotoUrl: "https://clemfoodie.com/wp-content/uploads/2021/03/tarte-citron-chantilly-1-scaled.jpg", CategoryId: 5},
		{Name: "Lavender Cake", Price: 32.3,
			PhotoUrl: "https://allysabakes.files.wordpress.com/2019/05/whitagram-image-4.jpg", CategoryId: 4},
		{Name: "Apple Pie", Price: 9.3,
			PhotoUrl: "https://cdn3.tmbi.com/toh/GoogleImagesPostCard/exps6086_HB133235C07_19_4b_WEB.jpg", CategoryId: 5},
		{Name: "Blueberry Pie", Price: 15.1,
			PhotoUrl: "https://www.tasteofhome.com/wp-content/uploads/2018/01/Contest-Winning-Fresh-Blueberry-Pie_exps10457_BS3149327B02_26_1bC_RMS-4.jpg", CategoryId: 5},
		{Name: "Mille Feuille", Price: 9.3,
			PhotoUrl: "http://oogio.net/wp-content/uploads/2016/11/tropical_mille_feuille-s.jpg", CategoryId: 8},
		{Name: "Saint Honore", Price: 14.8,
			PhotoUrl: "https://www.elle-et-vire.com/uploads/cache/930w/uploads/recip/recipe/870/230009.png", CategoryId: 8},
		{Name: "Vanilla Cake", Price: 53.8,
			PhotoUrl: "https://livforcake.com/wp-content/uploads/2017/06/vanilla-cake-4.jpg", CategoryId: 4},
		{Name: "Lemon Ricotta Cheesecake", Price: 15.3,
			PhotoUrl: "https://img.delicious.com.au/1H84Dz4k/del/2017/08/italian-style-ricotta-cheesecake-51070-2.jpg", CategoryId: 9},
		{Name: "Profiteroles", Price: 7.2,
			PhotoUrl: "https://images.immediate.co.uk/production/volatile/sites/30/2020/11/profiteroles-0dde0bb.jpg?quality=45&resize=504,458?quality=90&webp=true&resize=504,458", CategoryId: 8},
	}
	stmt, err = db.Prepare("INSERT INTO bakedGoods(name, price, photo_url, category_id) VALUES( ?,?,?,?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close() // Prepared statements take up server resources and should be closed after use.

	for i, c := range bakedGoods {
		res, err := stmt.Exec(c.Name, c.Price, c.PhotoUrl, c.CategoryId)
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
		{Email: "admin@admin.com", FirstName: "Admin", LastName: "Admin", Password: "admin", DeliveryAddress: "Sofia,Reduta", IsAdmin: true},
		{Email: "user@user.com", FirstName: "Default", LastName: "User", Password: "password", DeliveryAddress: "Oborishte"},
		{Email: "pike@golang.com", FirstName: "Rob", LastName: "Pike", Password: "rob", DeliveryAddress: "Varna"},
		{Email: "kamel@docker.com", FirstName: "Kamel", LastName: "Founadi", Password: "kamel", DeliveryAddress: "Lulin"},
	}

	stmt, err = db.Prepare(
		`INSERT INTO users(email, first_name, last_name, password, delivery_address, isAdmin, created, modified) 
		VALUES( ?, ?, ?, ?, ?, ?, ?,?)`)
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
		result, err := stmt.Exec(users[i].Email, users[i].FirstName, users[i].LastName, users[i].Password,
			users[i].DeliveryAddress, users[i].IsAdmin, users[i].Created, users[i].Modified)
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
	res, err := db.Exec("DROP TABLE IF EXISTS `reviews`")
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

	res, err = db.Exec("DROP TABLE IF EXISTS `categories`")
	if err != nil {
		log.Fatal(err)
	}
	rowsAffected, err = res.RowsAffected()
	log.Printf("'categories' - Rows Affected: %d %v", rowsAffected, err)

	res, err = db.Exec("DROP TABLE IF EXISTS `users`")
	if err != nil {
		log.Fatal(err)
	}
	rowsAffected, err = res.RowsAffected()
	log.Printf("users - Rows Affected: %d %v", rowsAffected, err)

}

func createDbs(db *sql.DB) {
	res, err := db.Exec("CREATE TABLE `users` (`id` bigint(20) NOT NULL AUTO_INCREMENT, `email` varchar(255) DEFAULT NULL, `first_name` varchar(255) DEFAULT NULL, `last_name` varchar(255) DEFAULT NULL, `password` varchar(255) DEFAULT NULL, `delivery_address` varchar(255) DEFAULT NULL, `isAdmin` tinyint(1) DEFAULT 0, `created` datetime(6) DEFAULT NULL, `modified` datetime(6) DEFAULT NULL, PRIMARY KEY (`id`), UNIQUE KEY `UK_6dotkott2kjsp8vw4d0m25fb7` (`email`)) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;")
	if err != nil {
		log.Fatal(err)
	}
	rowsAffected, err := res.RowsAffected()
	log.Printf("'USERS' - Rows Affected: %d %v", rowsAffected, err)

	res, err = db.Exec("CREATE TABLE `categories` (`id` bigint(20) NOT NULL AUTO_INCREMENT, `category_name` varchar(255) DEFAULT NULL, PRIMARY KEY (`id`)) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;")
	if err != nil {
		log.Fatal(err)
	}
	rowsAffected, err = res.RowsAffected()
	log.Printf("'category' - Rows Affected: %d %v", rowsAffected, err)

	res, err = db.Exec("CREATE TABLE `bakedGoods` (`id` bigint(20) NOT NULL AUTO_INCREMENT, `name` varchar(255) DEFAULT NULL, `photo_url` varchar(255), `price` double DEFAULT NULL, `category_id` bigint(20), PRIMARY KEY (`id`), FOREIGN KEY (`category_id`) REFERENCES `categories` (`id`)) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;")
	if err != nil {
		log.Fatal(err)
	}
	rowsAffected, err = res.RowsAffected()
	log.Printf("'BAKED GOODS' - Rows Affected: %d %v", rowsAffected, err)

	res, err = db.Exec("CREATE TABLE `tags` (`id` bigint(20) NOT NULL AUTO_INCREMENT, `tagName` varchar(255) DEFAULT NULL, PRIMARY KEY (`id`)) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;")
	if err != nil {
		log.Fatal(err)
	}
	rowsAffected, err = res.RowsAffected()
	log.Printf("'TAGS' - Rows Affected: %d %v", rowsAffected, err)

	res, err = db.Exec("CREATE TABLE `bakedGoodsTags` (`id` bigint(20) NOT NULL AUTO_INCREMENT, `tagId` bigint(20) DEFAULT NULL, `bakedGoodId` bigint(20) DEFAULT NULL, PRIMARY KEY (`id`), FOREIGN KEY (`tagId`) REFERENCES `tags` (`id`), FOREIGN KEY (`bakedGoodId`) REFERENCES `bakedGoods` (`id`)) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;")
	if err != nil {
		log.Fatal(err)
	}
	rowsAffected, err = res.RowsAffected()
	log.Printf("BAKED GOODS TAGS - Rows Affected: %d %v", rowsAffected, err)

	res, err = db.Exec("CREATE TABLE `orders` (`id` bigint(20) NOT NULL AUTO_INCREMENT,`userId` bigint(20),`status` varchar(50) DEFAULT NULL,`deliveryAddress` varchar(255) DEFAULT NULL,PRIMARY KEY (`id`),FOREIGN KEY (`userId`) REFERENCES `users` (`id`)) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;")
	if err != nil {
		log.Fatal(err)
	}
	rowsAffected, err = res.RowsAffected()
	log.Printf("ORDERS - Rows Affected: %d %v", rowsAffected, err)

	res, err = db.Exec("CREATE TABLE `orderedGoods` (`id` bigint(20) NOT NULL AUTO_INCREMENT,`bakedGoodId` bigint(20) DEFAULT NULL,`orderId` bigint(20) DEFAULT NULL, `quantity` bigint(20) DEFAULT NULL, PRIMARY KEY (`id`),FOREIGN KEY (`bakedGoodId`) REFERENCES `bakedGoods` (`id`),FOREIGN KEY (`orderId`) REFERENCES `orders` (`id`)) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;")
	if err != nil {
		log.Fatal(err)
	}
	rowsAffected, err = res.RowsAffected()
	log.Printf("ORDERED GOODS - Rows Affected: %d %v", rowsAffected, err)

	res, err = db.Exec("CREATE TABLE `reviews` (`id` bigint(20) NOT NULL AUTO_INCREMENT,`userId` bigint(20) DEFAULT NULL,`reviewText` varchar(255) DEFAULT NULL,PRIMARY KEY (`id`),FOREIGN KEY (`userId`) REFERENCES `users` (`id`)) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;")
	if err != nil {
		log.Fatal(err)
	}
	rowsAffected, err = res.RowsAffected()
	log.Printf("REVIEW - Rows Affected: %d %v", rowsAffected, err)
}
