package Controller

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	)

func Connect() (*sqlx.DB,error) {
	return sqlx.Open("postgres","postgres://avqthqbjsuqktk:5ad82047c4e29a180a9d32aa9f7a127ac21a6baef70deaa0f2de21cd11b54c4f@ec2-54-247-118-139.eu-west-1.compute.amazonaws.com:5432/dark17f0m8bla5")
}