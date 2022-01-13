package api

import "bakery-project/entities"

type UserRepo interface {
	FindAll(start, count int) ([]entities.User, error)
	FindByID(id int) (*entities.User, error)
	FindByEmail(email string) (*entities.User, error)
	Create(user *entities.User) (*entities.User, error)
	Update(user *entities.User) (*entities.User, error)
	DeleteByID(id int) (*entities.User, error)
	//Count() (int, error)
}

type BakedGoodRepo interface {
	FindAll(start, count int) ([]entities.BakedGood, error)
	FindByID(id int) (*entities.BakedGood, error)
	FindByName(name string) (*entities.BakedGood, error)
	Create(bakedGood *entities.BakedGood) (*entities.BakedGood, error)
	Update(bakedGood *entities.BakedGood) (*entities.BakedGood, error)
	DeleteByID(id int) (*entities.BakedGood, error)
	Count() (int, error)
}

type OrderRepo interface {
	FindAll(start, count int) ([]entities.Order, error)
	FindByID(id int) (*entities.Order, error)
	FindByUserId(userId int) (*entities.Order, error)
	Create(order *entities.Order) (*entities.Order, error)
	Update(order *entities.Order) (*entities.Order, error)
	DeleteByID(id int) (*entities.Order, error)
}

type ReviewRepo interface {
	FindAll(start, count int) ([]entities.Review, error)
	FindByID(id int) (*entities.Review, error)
	FindByUserId(userId int) (*entities.Review, error)
	Create(review *entities.Review) (*entities.Review, error)
	Update(review *entities.Review) (*entities.Review, error)
	DeleteByID(id int) (*entities.Review, error)
}

type TagRepo interface {
	FindAll(start, count int) ([]entities.Tag, error)
	FindByID(id int) (*entities.Tag, error)
	Create(tag *entities.Tag) (*entities.Tag, error)
	Update(tag *entities.Tag) (*entities.Tag, error)
	DeleteByID(id int) (*entities.Tag, error)
}
