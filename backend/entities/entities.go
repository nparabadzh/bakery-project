package entities

import (
	"github.com/dgrijalva/jwt-go"
	"time"
)

type User struct {
	ID              int64     `json:"id"`
	Email           string    `json:"email"`
	FirstName       string    `json:"first_name"`
	LastName        string    `json:"last_name"`
	Username        string    `json:"username"`
	Password        string    `json:"password"`
	DeliveryAddress string    `json:"delivery_address"`
	IsAdmin         bool      `json:"is_admin"`
	Created         time.Time `json:"created_at"`
	Modified        time.Time `json:"modified_at"`
}

type UserLogin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type BakedGood struct {
	ID       int64   `json:"id"`
	Name     string  `json:"name"`
	PhotoUrl []byte  `json:"photo"`
	Price    float64 `json:"price"`
	TagsIds  []int64 `json:"tags_ids"`
}

type Order struct {
	ID              int    `json:"id"`
	UserId          int64  `json:"user_id"`
	Status          string `json:"status"`
	DeliveryAddress string `json:"delivery_address"`
}

type Review struct {
	ID         int64  `json:"id"`
	UserId     int64  `json:"user_id"`
	ReviewText string `json:"review_text"`
}

type Tag struct {
	ID      int    `json:"id"`
	TagName string `json:"tag_name"`
}

type BakedGoodTag struct {
	ID          int64 `json:"id"`
	TagID       int64 `json:"tag_id"`
	BakedGoodID int64 `json:"baked_good_id"`
}

type OrderedGoods struct {
	ID      int64 `json:"id"`
	GoodID  int64 `json:"good_id"`
	OrderID int64 `json:"order_id"`
}

type UserToken struct {
	UserID          int       `json:"id"`
	FirstName       string    `json:"first-name"`
	LastName        string    `json:"last-name"`
	Email           string    `json:"email"`
	DeliveryAddress string    `json:"delivery_address"`
	IsAdmin         bool      `json:"is_admin"`
	Created         time.Time `json:"created_at"`
	Modified        time.Time `json:"modified_at"`
	jwt.StandardClaims
}
