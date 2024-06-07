package db

import (
	"cmp"
	"context"
	"fmt"
	"github.com/jaswdr/faker"
	"reflect"
	"strconv"
	"strings"
)

type InMemoryUserRepo struct {
	data map[uint64]User
}

func NewInMemoryUserRepo() *InMemoryUserRepo {
	r := &InMemoryUserRepo{
		data: make(map[uint64]User, 100),
	}
	r.init()
	return r
}

func (r *InMemoryUserRepo) init() {
	fake := faker.New()
	for i := 0; i < 100; i++ {
		r.data[uint64(i)] = newFakeUser(fake, uint64(i))
	}
}

func newFakeUser(fake faker.Faker, id uint64) User {
	return User{
		Id:          id,
		FirstName:   fake.Person().FirstName(),
		City:        fake.Address().City(),
		PhoneNumber: fake.Phone().Number(),
		Height:      fake.Float32(2, 4, 7),
		Married:     fake.Bool(),
	}
}

func (r *InMemoryUserRepo) GetUserById(_ context.Context, id uint64) (User, error) {
	u, ok := r.data[id]
	if !ok {
		return User{}, ErrUserNotExists
	}
	return u, nil
}

func (r *InMemoryUserRepo) GetUsersByIds(_ context.Context, ids []uint64) ([]User, error) {
	users := make([]User, 0, len(ids))
	for _, id := range ids {
		if u, ok := r.data[id]; ok {
			users = append(users, u)
		}
	}
	return users, nil
}

func (r *InMemoryUserRepo) SearchUsers(_ context.Context, criteria SearchCriteria) ([]User, error) {
	users := make([]User, 0, 20)
	searchFn := func(u User) (bool, error) {
		switch criteria.Field {
		case "fname":
			return matchString(criteria, u.FirstName)
		case "city":
			return matchString(criteria, u.City)
		case "phone":
			return matchNumber(criteria, u.PhoneNumber)
		case "height":
			return matchNumber(criteria, u.Height)
		case "Married":
			return matchBool(criteria, u.Married)
		}
		return false, ErrInvalidField
	}
	for _, u := range r.data {
		ok, err := searchFn(u)
		if err != nil {
			return nil, err
		}
		if ok {
			users = append(users, u)
		}
	}
	return users, nil
}

func matchString(criteria SearchCriteria, val string) (bool, error) {
	var matchMode MatchMode
	if criteria.MatchMode != nil {
		matchMode = MatchMode(*criteria.MatchMode)
	} else {
		matchMode = MatchModeContains
	}
	switch matchMode {
	case MatchModeContains:
		return strings.Contains(strings.ToLower(val), strings.ToLower(criteria.Value)), nil
	case MatchModeStartsWith:
		return strings.HasPrefix(strings.ToLower(val), strings.ToLower(criteria.Value)), nil
	case MatchModeEndsWith:
		return strings.HasSuffix(strings.ToLower(val), strings.ToLower(criteria.Value)), nil
	case MatchModeExact:
		return strings.ToLower(val) == strings.ToLower(criteria.Value), nil
	default:
		return false, ErrInvalidMatchMode
	}
}

func matchNumber(criteria SearchCriteria, val any) (bool, error) {
	var matchMode MatchMode
	if criteria.MatchMode != nil {
		matchMode = MatchMode(*criteria.MatchMode)
	} else {
		matchMode = MatchModeEquals
	}
	switch val.(type) {
	case int, int8, int16, int32, int64:
		cv, err := strconv.ParseInt(criteria.Value, 10, 64)
		if err != nil {
			return false, ErrInvalidCriteriaValue
		}
		intVal, err := getInt64(val)
		if err != nil {
			return false, err
		}
		return compare(matchMode, cv, intVal)
	case uint, uint8, uint16, uint32, uint64:
		cv, err := strconv.ParseUint(criteria.Value, 10, 64)
		if err != nil {
			return false, ErrInvalidCriteriaValue
		}
		uintVal, err := getUint64(val)
		if err != nil {
			return false, err
		}
		return compare(matchMode, cv, uintVal)
	case float32, float64:
		cv, err := strconv.ParseFloat(criteria.Value, 10)
		if err != nil {
			return false, ErrInvalidCriteriaValue
		}
		floatVal, err := getFloat64(val)
		if err != nil {
			return false, err
		}
		return compare(matchMode, cv, floatVal)
	default:
		return false, ErrInvalidField
	}
}

func matchBool(criteria SearchCriteria, val bool) (bool, error) {
	var matchMode MatchMode
	if criteria.MatchMode != nil {
		matchMode = MatchMode(*criteria.MatchMode)
	} else {
		matchMode = MatchModeEquals
	}
	boolVal, err := strconv.ParseBool(criteria.Value)
	if err != nil {
		return false, ErrInvalidCriteriaValue
	}
	switch matchMode {
	case MatchModeEquals:
		return boolVal == val, nil
	case MatchModeNot:
		return boolVal != val, nil
	default:
		return false, ErrInvalidMatchMode
	}
}

func compare[N cmp.Ordered](matchMode MatchMode, criteria N, model N) (bool, error) {
	switch matchMode {
	case MatchModeEquals:
		return criteria == model, nil
	case MatchModeNotEquals:
		return criteria != model, nil
	case MatchModeGreaterThan:
		return criteria < model, nil
	case MatchModeGreaterThanOrEqual:
		return criteria <= model, nil
	case MatchModeLessThan:
		return criteria > model, nil
	case MatchModeLessThanOrEqual:
		return criteria >= model, nil
	default:
		return false, ErrInvalidMatchMode
	}
}

func getInt64(v interface{}) (int64, error) {
	switch reflect.TypeOf(v).Kind() {
	case reflect.Int8:
		d, _ := v.(int8)
		return int64(d), nil
	case reflect.Int16:
		d, _ := v.(int16)
		return int64(d), nil
	case reflect.Int32:
		d, _ := v.(int32)
		return int64(d), nil
	case reflect.Int64:
		d, _ := v.(int64)
		return d, nil
	case reflect.Int:
		d, _ := v.(int)
		return int64(d), nil
	default:
		return 0, fmt.Errorf("not interger: %v", v)
	}

}

func getUint64(v interface{}) (uint64, error) {
	switch reflect.TypeOf(v).Kind() {
	case reflect.Uint8:
		d, _ := v.(uint8)
		return uint64(d), nil
	case reflect.Uint16:
		d, _ := v.(uint16)
		return uint64(d), nil
	case reflect.Uint32:
		d, _ := v.(uint32)
		return uint64(d), nil
	case reflect.Uint64:
		d, _ := v.(uint64)
		return d, nil
	case reflect.Uint:
		d, _ := v.(uint)
		return uint64(d), nil
	default:
		return 0, fmt.Errorf("not unsigned interger: %v", v)
	}
}

func getFloat64(v interface{}) (float64, error) {
	switch reflect.TypeOf(v).Kind() {
	case reflect.Float32:
		d, _ := v.(float32)
		return float64(d), nil
	case reflect.Float64:
		d, _ := v.(float64)
		return d, nil
	default:
		return 0, fmt.Errorf("not float: %v", v)
	}
}
