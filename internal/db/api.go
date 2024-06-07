package db

import (
	"context"
	"errors"
)

var ErrUserNotExists = errors.New("user does not exist")
var ErrInvalidField = errors.New("invalid field")
var ErrInvalidCriteriaValue = errors.New("invalid criteria value")
var ErrInvalidMatchMode = errors.New("invalid match mode")

type UserRepository interface {
	GetUserById(context.Context, uint64) (User, error)
	GetUsersByIds(context.Context, []uint64) ([]User, error)
	SearchUsers(context.Context, SearchCriteria) ([]User, error)
}

type User struct {
	Id          uint64
	FirstName   string
	City        string
	PhoneNumber string
	Height      float32
	Married     bool
}

type SearchCriteria struct {
	Field     string
	Value     string
	MatchMode *MatchMode
}

type MatchMode string

const (
	MatchModeContains           MatchMode = "CONTAINS"
	MatchModeStartsWith         MatchMode = "STARTS_WITH"
	MatchModeEndsWith           MatchMode = "ENDS_WITH"
	MatchModeExact              MatchMode = "EXACT"
	MatchModeEquals             MatchMode = "EQUALS"
	MatchModeNotEquals          MatchMode = "NOT_EQUALS"
	MatchModeGreaterThan        MatchMode = "GREATER_THAN"
	MatchModeGreaterThanOrEqual MatchMode = "GREATER_THAN_OR_EQUAL"
	MatchModeLessThan           MatchMode = "LESS_THAN"
	MatchModeLessThanOrEqual    MatchMode = "LESS_THAN_OR_EQUAL"
	MatchModeNot                MatchMode = "NOT"
)
