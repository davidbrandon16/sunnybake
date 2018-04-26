package Controller

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	)

func Connect() (*sqlx.DB,error) {
	return sqlx.Open("postgres","postgres://gfvnnmotqekepv:d40007172557c0bef9ede300751f9d2ad209464a3594d514bff6448fa5349dde@ec2-79-125-110-209.eu-west-1.compute.amazonaws.com:5432/d9r4or09435uqm")
}