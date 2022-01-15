package rest

import (
	"bakery-project/api"
	"bakery-project/entities"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"strconv"
	"time"
)

type App struct {
	Router *mux.Router

	Users      api.UserRepo
	BakedGoods api.BakedGoodRepo
	Orders     api.OrderRepo
	Reviews    api.ReviewRepo
	Tags       api.TagRepo
	Categories api.CategoryRepo
	Validator  *validator.Validate
	Translator ut.Translator
	Server     *http.Server
}

func (a *App) Init(user, password, dbname string) {
	a.Users = api.NewUserRepoMysql(user, password, dbname)
	a.BakedGoods = api.NewBakedGoodsRepoMysql(user, password, dbname)
	a.Orders = api.NewOrdersRepoMysql(user, password, dbname)
	a.Reviews = api.NewReviewsRepoMysql(user, password, dbname)
	a.Tags = api.NewTagsRepoMysql(user, password, dbname)
	a.Categories = api.NewCategoriesRepoMysql(user, password, dbname)

	// Create and configure validator and translator
	a.Validator = validator.New()
	eng := en.New()
	var uni *ut.UniversalTranslator
	uni = ut.New(eng, eng)
	// this is usually know or extracted from http 'Accept-Language' header
	// also see uni.FindTranslator(...)
	var found bool
	a.Translator, found = uni.GetTranslator("en")
	if !found {
		log.Fatal("translator not found")
	}
	if err := en_translations.RegisterDefaultTranslations(a.Validator, a.Translator); err != nil {
		log.Fatal(err)
	}

	a.Router = mux.NewRouter()
	a.initializeRoutes()
}

func (a *App) Run(addr string) {
	a.Server = &http.Server{
		Addr:    addr,
		Handler: a.Router,
	}
	go func() {
		log.Println("Starting HTTP server on", addr)
		log.Println("Press Ctrl+C to stop the server ...")
		if err := a.Server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Error server listen: %v\n", err)
		}
	}()
}

func (a *App) Shutdown(ctx context.Context) error {
	return a.Server.Shutdown(ctx)
}

func (a *App) initializeRoutes() {
	// User routes
	a.Router.HandleFunc("/users", a.getUsers).Methods("GET")
	a.Router.HandleFunc("/signUp", a.createUser).Methods("POST")
	a.Router.HandleFunc("/users/{id:[0-9]+}", a.getUser).Methods("GET")
	a.Router.HandleFunc("/users/{id:[0-9]+}", a.updateUser).Methods("PUT")
	a.Router.HandleFunc("/users/{id:[0-9]+}", a.deleteUser).Methods("DELETE")
	a.Router.HandleFunc("/login", a.login).Methods("POST")

	// Baked goods routes
	a.Router.HandleFunc("/baked-goods", a.getBakedGoods).Methods("GET")
	a.Router.HandleFunc("/baked-goods", a.createBakedGood).Methods("POST")
	a.Router.HandleFunc("/baked-goods/{id:[0-9]+}", a.getBakedGood).Methods("GET")
	a.Router.HandleFunc("/baked-goods/{id:[0-9]+}", a.updateBakedGood).Methods("PUT")
	a.Router.HandleFunc("/baked-goods/{id:[0-9]+}", a.deleteBakedGood).Methods("DELETE")

	// Orders routes
	a.Router.HandleFunc("/orders", a.getOrders).Methods("GET")
	a.Router.HandleFunc("/orders", a.createOrder).Methods("POST")
	a.Router.HandleFunc("/orders/{id:[0-9]+}", a.getOrder).Methods("GET")
	a.Router.HandleFunc("/orders/{id:[0-9]+}", a.updateOrder).Methods("PUT")
	a.Router.HandleFunc("/orders/{id:[0-9]+}", a.deleteOrder).Methods("DELETE")

	// Reviews routes
	a.Router.HandleFunc("/reviews", a.getReviews).Methods("GET")
	a.Router.HandleFunc("/reviews", a.createReview).Methods("POST")
	a.Router.HandleFunc("/reviews/{id:[0-9]+}", a.getReview).Methods("GET")
	a.Router.HandleFunc("/reviews/{id:[0-9]+}", a.updateReview).Methods("PUT")
	a.Router.HandleFunc("/reviews/{id:[0-9]+}", a.deleteReview).Methods("DELETE")

	// Tags routes
	a.Router.HandleFunc("/tags", a.getTags).Methods("GET")
	a.Router.HandleFunc("/tags", a.createTag).Methods("POST")
	a.Router.HandleFunc("/tags/{id:[0-9]+}", a.getTag).Methods("GET")
	a.Router.HandleFunc("/tags/{id:[0-9]+}", a.updateTag).Methods("PUT")
	a.Router.HandleFunc("/tags/{id:[0-9]+}", a.deleteTag).Methods("DELETE")

	// Categories routes
	a.Router.HandleFunc("/categories", a.getCategories).Methods("GET")
	a.Router.HandleFunc("/categories", a.createCategory).Methods("POST")
	a.Router.HandleFunc("/categories/{id:[0-9]+}", a.getCategory).Methods("GET")
	a.Router.HandleFunc("/categories/{id:[0-9]+}", a.updateCategory).Methods("PUT")
	a.Router.HandleFunc("/categories/{id:[0-9]+}", a.deleteCategory).Methods("DELETE")

	// Auth route
	s := a.Router.PathPrefix("/auth").Subrouter()
	s.Use(JwtVerify)
	s.HandleFunc("/users", a.getUsers).Methods(http.MethodGet)
	s.HandleFunc("/users", a.createUser).Methods("POST")
	s.HandleFunc("/users/{id:[0-9]+}", a.getUser).Methods("GET")
	s.HandleFunc("/users/{id:[0-9]+}", a.updateUser).Methods("PUT")
	s.HandleFunc("/users/{id:[0-9]+}", a.deleteUser).Methods("DELETE")

}

func (a *App) login(w http.ResponseWriter, r *http.Request) {
	userCredentials := &entities.UserLogin{}
	err := json.NewDecoder(r.Body).Decode(userCredentials)
	if err != nil {
		fmt.Printf("Error logging user %v: %v", userCredentials, err)
		var resp = map[string]interface{}{"status": false, "message": "Invalid request"}
		json.NewEncoder(w).Encode(resp)
		return
	}
	resp, err := a.checkEmailPassword(w, userCredentials.Email, userCredentials.Password)
	if err == nil {
		json.NewEncoder(w).Encode(resp)
	}
}

func (a *App) checkEmailPassword(w http.ResponseWriter, email, password string) (map[string]interface{}, error) {
	user, err := a.Users.FindByEmail(email)
	fmt.Println(email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Email address not found")
		return nil, err
	}
	expiresAt := time.Now().Add(time.Minute * 10).Unix()

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword { //Password does not match!
		respondWithError(w, http.StatusUnauthorized, "Invalid login credentials. Please try again")
		return nil, err
	}

	claims := &entities.UserToken{
		UserID:          int(user.ID),
		FirstName:       user.FirstName,
		LastName:        user.LastName,
		Email:           user.Email,
		DeliveryAddress: user.DeliveryAddress,
		IsAdmin:         user.IsAdmin,
		Created:         user.Created,
		Modified:        user.Modified,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expiresAt,
			Issuer:    "test",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, error := token.SignedString([]byte("secret"))
	if error != nil {
		fmt.Println(error)
	}

	var resp = map[string]interface{}{"status": false, "message": "logged in"}
	resp["token"] = tokenString //Store the token in the response
	// remove user password
	user.Password = ""

	resp["user"] = user
	return resp, nil
}

//func enableCors(w *http.ResponseWriter) {
//	(*w).Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
//	(*w).Header().Set("Access-Control-Allow-Headers", "Content-Type")
//}
//
//func setupResponse(w *http.ResponseWriter) {
//	(*w).Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
//	(*w).Header().Set("Content-Type", "text/html; charset=utf-8")
//	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
//	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
//}

// User handlers
func (a *App) getUsers(w http.ResponseWriter, r *http.Request) {
	count, err := strconv.Atoi(r.FormValue("count"))
	if err != nil && r.FormValue("count") != "" {
		respondWithError(w, http.StatusBadRequest, "Invalid request count parameter")
		return
	}
	start, err := strconv.Atoi(r.FormValue("start"))
	if err != nil && r.FormValue("start") != "" {
		respondWithError(w, http.StatusBadRequest, "Invalid request start parameter")
		return
	}
	start--
	if count > 20 || count < 1 {
		count = 20
	}
	if start < 0 {
		start = 0
	}
	users, err := a.Users.FindAll(start, count)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	// remove user passwords
	for i := range users {
		users[i].Password = ""
	}
	respondWithJSON(w, http.StatusOK, users)
}

func (a *App) createUser(w http.ResponseWriter, r *http.Request) {
	user := &entities.User{}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(user); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Validate User struct
	err := a.Validator.Struct(user)
	if err != nil {
		// translate all error at once
		errs := err.(validator.ValidationErrors)
		respondWithValidationError(errs.Translate(a.Translator), w)
		return
	}

	// Hash the pasword with bcrypt
	pass, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println(err)
		respondWithError(w, http.StatusInternalServerError, "Password Encryption  failed")
		return
	}
	user.Password = string(pass)
	user.Created = time.Now()
	user.Modified = time.Now()
	if user, err = a.Users.Create(user); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	//// remove user password
	user.Password = ""

	respondWithJSON(w, http.StatusCreated, user)
}

func (a *App) getUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	var user *entities.User
	if user, err = a.Users.FindByID(id); err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithError(w, http.StatusNotFound, "User not found")
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}
	// remove user password
	user.Password = ""

	respondWithJSON(w, http.StatusOK, user)
}

func (a *App) updateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	user := &entities.User{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(user); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid resquest payload")
		return
	}
	// Validate User struct
	err = a.Validator.Struct(user)
	if err != nil {
		// translate all error at once
		errs := err.(validator.ValidationErrors)
		respondWithValidationError(errs.Translate(a.Translator), w)
		return
	}

	//if int(user.ID) != id {
	//	respondWithError(w, http.StatusBadRequest, fmt.Sprintf("ID %d in URL path is different from ID %d in request payload", user.ID, id))
	//	return
	//}

	// Find if user exists in DB
	oldUser, err := a.Users.FindByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			respondWithError(w, http.StatusNotFound, fmt.Sprintf("user with ID='%d' does not exist", id))
		} else {
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}
	user.ID = int64(id)
	// Encrypt password if sent otherwise use old password
	if user.Password != "" {
		// Hash the pasword with bcrypt
		pass, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			fmt.Println(err)
			respondWithError(w, http.StatusInternalServerError, "Password Encryption  failed")
			return
		}
		user.Password = string(pass)
	} else {
		user.Password = oldUser.Password
	}

	// Do update user
	user.Modified = time.Now()
	if user, err = a.Users.Update(user); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	//// remove user password
	//user.Password = ""

	respondWithJSON(w, http.StatusOK, user)
}

func (a *App) deleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid User ID")
		return
	}
	// Do delete user in DB
	user, err := a.Users.DeleteByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			respondWithError(w, http.StatusNotFound, fmt.Sprintf("user with ID='%d' does not exist", id))
		} else {
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}
	// remove user password
	user.Password = ""

	respondWithJSON(w, http.StatusOK, user)
}

//Baked Goods handlers
func (a *App) getBakedGoods(w http.ResponseWriter, r *http.Request) {
	count, err := strconv.Atoi(r.FormValue("count"))
	if err != nil && r.FormValue("count") != "" {
		respondWithError(w, http.StatusBadRequest, "Invalid request count parameter")
		return
	}
	start, err := strconv.Atoi(r.FormValue("start"))
	if err != nil && r.FormValue("start") != "" {
		respondWithError(w, http.StatusBadRequest, "Invalid request start parameter")
		return
	}
	start--
	if count > 20 || count < 1 {
		count = 20
	}
	if start < 0 {
		start = 0
	}
	bakedGoods, err := a.BakedGoods.FindAll(start, count)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, bakedGoods)
}

func (a *App) createBakedGood(w http.ResponseWriter, r *http.Request) {
	bakedGood := &entities.BakedGood{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(bakedGood); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Validate User struct
	err := a.Validator.Struct(bakedGood)
	if err != nil {
		// translate all error at once
		errs := err.(validator.ValidationErrors)
		respondWithValidationError(errs.Translate(a.Translator), w)
		return
	}

	if bakedGood, err = a.BakedGoods.Create(bakedGood); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, bakedGood)
}

func (a *App) getBakedGood(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid baked good ID")
		return
	}

	var bakedGood *entities.BakedGood
	if bakedGood, err = a.BakedGoods.FindByID(id); err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithError(w, http.StatusNotFound, "Baked Good not found")
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	respondWithJSON(w, http.StatusOK, bakedGood)
}

func (a *App) updateBakedGood(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid baked good ID")
		return
	}

	bakedGood := &entities.BakedGood{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(bakedGood); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid resquest payload")
		return
	}
	// Validate User struct
	err = a.Validator.Struct(bakedGood)
	if err != nil {
		// translate all error at once
		errs := err.(validator.ValidationErrors)
		respondWithValidationError(errs.Translate(a.Translator), w)
		return
	}

	// Find if baked good exists in DB
	_, err = a.BakedGoods.FindByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			respondWithError(w, http.StatusNotFound, fmt.Sprintf("baked good with ID='%d' does not exist", id))
		} else {
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}
	bakedGood.ID = int64(id)

	// Do update baked good
	if bakedGood, err = a.BakedGoods.Update(bakedGood); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, bakedGood)
}

func (a *App) deleteBakedGood(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Baked Good ID")
		return
	}
	// Do delete baked good in DB
	bakedGood, err := a.BakedGoods.DeleteByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			respondWithError(w, http.StatusNotFound, fmt.Sprintf("baked good with ID='%d' does not exist", id))
		} else {
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	respondWithJSON(w, http.StatusOK, bakedGood)
}

//Orders handlers
func (a *App) getOrders(w http.ResponseWriter, r *http.Request) {
	count, err := strconv.Atoi(r.FormValue("count"))
	if err != nil && r.FormValue("count") != "" {
		respondWithError(w, http.StatusBadRequest, "Invalid request count parameter")
		return
	}
	start, err := strconv.Atoi(r.FormValue("start"))
	if err != nil && r.FormValue("start") != "" {
		respondWithError(w, http.StatusBadRequest, "Invalid request start parameter")
		return
	}
	start--
	if count > 20 || count < 1 {
		count = 20
	}
	if start < 0 {
		start = 0
	}
	orders, err := a.Orders.FindAll(start, count)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, orders)
}

func (a *App) createOrder(w http.ResponseWriter, r *http.Request) {
	order := &entities.Order{}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(order); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Validate User struct
	err := a.Validator.Struct(order)
	if err != nil {
		// translate all error at once
		errs := err.(validator.ValidationErrors)
		respondWithValidationError(errs.Translate(a.Translator), w)
		return
	}

	if order, err = a.Orders.Create(order); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, order)
}

func (a *App) getOrder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid order ID")
		return
	}

	var order *entities.Order
	if order, err = a.Orders.FindByID(id); err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithError(w, http.StatusNotFound, "Order not found")
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	respondWithJSON(w, http.StatusOK, order)
}

func (a *App) updateOrder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid baked good ID")
		return
	}

	order := &entities.Order{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(order); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid resquest payload")
		return
	}
	// Validate User struct
	err = a.Validator.Struct(order)
	if err != nil {
		// translate all error at once
		errs := err.(validator.ValidationErrors)
		respondWithValidationError(errs.Translate(a.Translator), w)
		return
	}

	// Find if baked good exists in DB
	_, err = a.Orders.FindByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			respondWithError(w, http.StatusNotFound, fmt.Sprintf("order with ID='%d' does not exist", id))
		} else {
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}
	order.ID = id

	// Do update baked good
	if order, err = a.Orders.Update(order); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, order)
}

func (a *App) deleteOrder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Order ID")
		return
	}
	// Do delete baked good in DB
	order, err := a.Orders.DeleteByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			respondWithError(w, http.StatusNotFound, fmt.Sprintf("order with ID='%d' does not exist", id))
		} else {
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	respondWithJSON(w, http.StatusOK, order)
}

//Review handlers
func (a *App) getReviews(w http.ResponseWriter, r *http.Request) {
	count, err := strconv.Atoi(r.FormValue("count"))
	if err != nil && r.FormValue("count") != "" {
		respondWithError(w, http.StatusBadRequest, "Invalid request count parameter")
		return
	}
	start, err := strconv.Atoi(r.FormValue("start"))
	if err != nil && r.FormValue("start") != "" {
		respondWithError(w, http.StatusBadRequest, "Invalid request start parameter")
		return
	}
	start--
	if count > 20 || count < 1 {
		count = 20
	}
	if start < 0 {
		start = 0
	}
	reviews, err := a.Reviews.FindAll(start, count)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, reviews)
}

func (a *App) createReview(w http.ResponseWriter, r *http.Request) {
	review := &entities.Review{}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(review); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Validate User struct
	err := a.Validator.Struct(review)
	if err != nil {
		// translate all error at once
		errs := err.(validator.ValidationErrors)
		respondWithValidationError(errs.Translate(a.Translator), w)
		return
	}

	if review, err = a.Reviews.Create(review); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, review)
}

func (a *App) getReview(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid review ID")
		return
	}

	var review *entities.Review
	if review, err = a.Reviews.FindByID(id); err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithError(w, http.StatusNotFound, "review not found")
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	respondWithJSON(w, http.StatusOK, review)
}

func (a *App) updateReview(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid  review ID")
		return
	}

	review := &entities.Review{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(review); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid resquest payload")
		return
	}
	// Validate User struct
	err = a.Validator.Struct(review)
	if err != nil {
		// translate all error at once
		errs := err.(validator.ValidationErrors)
		respondWithValidationError(errs.Translate(a.Translator), w)
		return
	}

	// Find if baked good exists in DB
	_, err = a.Reviews.FindByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			respondWithError(w, http.StatusNotFound, fmt.Sprintf("review with ID='%d' does not exist", id))
		} else {
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}
	review.ID = int64(id)

	// Do update baked good
	if review, err = a.Reviews.Update(review); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, review)
}

func (a *App) deleteReview(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid review ID")
		return
	}
	// Do delete baked good in DB
	review, err := a.Reviews.DeleteByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			respondWithError(w, http.StatusNotFound, fmt.Sprintf("order with ID='%d' does not exist", id))
		} else {
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	respondWithJSON(w, http.StatusOK, review)
}

//Tags handlers
func (a *App) getTags(w http.ResponseWriter, r *http.Request) {
	count, err := strconv.Atoi(r.FormValue("count"))
	if err != nil && r.FormValue("count") != "" {
		respondWithError(w, http.StatusBadRequest, "Invalid request count parameter")
		return
	}
	start, err := strconv.Atoi(r.FormValue("start"))
	if err != nil && r.FormValue("start") != "" {
		respondWithError(w, http.StatusBadRequest, "Invalid request start parameter")
		return
	}
	start--
	if count > 20 || count < 1 {
		count = 20
	}
	if start < 0 {
		start = 0
	}
	tags, err := a.Tags.FindAll(start, count)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, tags)
}

func (a *App) createTag(w http.ResponseWriter, r *http.Request) {
	tag := &entities.Tag{}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(tag); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Validate User struct
	err := a.Validator.Struct(tag)
	if err != nil {
		// translate all error at once
		errs := err.(validator.ValidationErrors)
		respondWithValidationError(errs.Translate(a.Translator), w)
		return
	}

	if tag, err = a.Tags.Create(tag); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, tag)
}

func (a *App) getTag(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid tag ID")
		return
	}

	var tag *entities.Tag
	if tag, err = a.Tags.FindByID(id); err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithError(w, http.StatusNotFound, "tag not found")
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	respondWithJSON(w, http.StatusOK, tag)
}

func (a *App) updateTag(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid  tag ID")
		return
	}

	tag := &entities.Tag{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(tag); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid resquest payload")
		return
	}
	// Validate User struct
	err = a.Validator.Struct(tag)
	if err != nil {
		// translate all error at once
		errs := err.(validator.ValidationErrors)
		respondWithValidationError(errs.Translate(a.Translator), w)
		return
	}

	// Find if baked good exists in DB
	_, err = a.Tags.FindByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			respondWithError(w, http.StatusNotFound, fmt.Sprintf("tag with ID='%d' does not exist", id))
		} else {
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}
	tag.ID = id

	// Do update baked good
	if tag, err = a.Tags.Update(tag); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, tag)
}

func (a *App) deleteTag(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid tag ID")
		return
	}
	// Do delete baked good in DB
	tag, err := a.Tags.DeleteByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			respondWithError(w, http.StatusNotFound, fmt.Sprintf("tag with ID='%d' does not exist", id))
		} else {
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	respondWithJSON(w, http.StatusOK, tag)
}

//Category handlers
func (a *App) getCategories(w http.ResponseWriter, r *http.Request) {
	count, err := strconv.Atoi(r.FormValue("count"))
	if err != nil && r.FormValue("count") != "" {
		respondWithError(w, http.StatusBadRequest, "Invalid request count parameter")
		return
	}
	start, err := strconv.Atoi(r.FormValue("start"))
	if err != nil && r.FormValue("start") != "" {
		respondWithError(w, http.StatusBadRequest, "Invalid request start parameter")
		return
	}
	start--
	if count > 20 || count < 1 {
		count = 20
	}
	if start < 0 {
		start = 0
	}
	categories, err := a.Tags.FindAll(start, count)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, categories)
}

func (a *App) createCategory(w http.ResponseWriter, r *http.Request) {
	category := &entities.Category{}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(category); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Validate User struct
	err := a.Validator.Struct(category)
	if err != nil {
		// translate all error at once
		errs := err.(validator.ValidationErrors)
		respondWithValidationError(errs.Translate(a.Translator), w)
		return
	}

	if category, err = a.Categories.Create(category); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, category)
}

func (a *App) getCategory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid category ID")
		return
	}

	var category *entities.Category
	if category, err = a.Categories.FindByID(id); err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithError(w, http.StatusNotFound, "category not found")
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	respondWithJSON(w, http.StatusOK, category)
}

func (a *App) updateCategory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid category ID")
		return
	}

	category := &entities.Category{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(category); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid resquest payload")
		return
	}
	// Validate User struct
	err = a.Validator.Struct(category)
	if err != nil {
		// translate all error at once
		errs := err.(validator.ValidationErrors)
		respondWithValidationError(errs.Translate(a.Translator), w)
		return
	}

	// Find if baked good exists in DB
	_, err = a.Categories.FindByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			respondWithError(w, http.StatusNotFound, fmt.Sprintf("category with ID='%d' does not exist", id))
		} else {
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}
	category.ID = id

	// Do update baked good
	if category, err = a.Categories.Update(category); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, category)
}

func (a *App) deleteCategory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid category ID")
		return
	}
	// Do delete baked good in DB
	category, err := a.Categories.DeleteByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			respondWithError(w, http.StatusNotFound, fmt.Sprintf("category with ID='%d' does not exist", id))
		} else {
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	respondWithJSON(w, http.StatusOK, category)
}
