package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const RoleClaimKey = "http://schemas.microsoft.com/ws/2008/06/identity/claims/role"

type UserDoc struct {
	MongoID      primitive.ObjectID `bson:"_id,omitempty" json:"-"`
	ID           int64              `bson:"id" json:"id"`
	Username     string             `bson:"username" json:"username"`
	Email        string             `bson:"email" json:"email"`
	PasswordHash string             `bson:"passwordHash" json:"-"`
	Role         string             `bson:"role" json:"role"`
	CreatedAt    time.Time          `bson:"createdAt" json:"createdAt"`
}

type RegisterReq struct {
	Email    string `json:"email" binding:"required,email"`
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
	Role     string `json:"role" binding:"required"` // GUIDE / TOURIST
}

type LoginReq struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type AuthResponse struct {
	ID          int64  `json:"id"`
	AccessToken string `json:"accessToken"`
}

type MeResponse struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Role     string `json:"role"`
}

type CounterDoc struct {
	ID    string `bson:"_id"`
	Value int64  `bson:"value"`
}
