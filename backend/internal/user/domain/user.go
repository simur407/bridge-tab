package user

import "slices"

type UserId string

type UserState struct {
	Id   UserId
	Name string
}

type UserRegistered struct {
	Id   UserId
	Name string
}

type User struct {
	State UserState
	Events []any
}

func RegisterUser(id *UserId, name string) *User {
	return &User{
		State: UserState{Id: *id, Name: name},
		Events: []any{UserRegistered{Id: *id, Name: name}},
	}
}

func (u *User) GetEvents() []any {
	return u.Events
}

func (u *User) Commit() {
	u.Events = slices.Delete(u.Events, 0, len(u.Events))
}
