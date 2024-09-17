package repository

import (
	"context"

	models "github.com/idontknowtoobrother/stripe-go-lang/Models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repo interface {
	GetAll() (*[]models.Product, error)
	GetByUuid(string) (*models.Product, error)
	Create(*models.Product) error
}

type repository struct {
	ctx context.Context
	db  *mongo.Database
}

func NewRepo(ctx context.Context, db *mongo.Database) Repo {
	return &repository{
		ctx: ctx,
		db:  db,
	}
}

func (r *repository) GetAll() (*[]models.Product, error) {
	var products []models.Product
	cursor, err := r.db.Collection("products").Find(r.ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(r.ctx)
	if err = cursor.All(r.ctx, &products); err != nil {
		return nil, err
	}
	return &products, nil
}

func (r *repository) GetByUuid(uuid string) (*models.Product, error) {
	var product models.Product
	err := r.db.Collection("products").FindOne(r.ctx, bson.M{
		"uuid": uuid,
	}).Decode(&product)
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *repository) Create(product *models.Product) error {
	_, err := r.db.Collection("products").InsertOne(r.ctx, product)
	if err != nil {
		return err
	}
	return nil
}
