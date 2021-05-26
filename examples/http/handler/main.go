package main

import (
	"context"
	"log"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/gorilla/mux"
)

// User is a user model.
type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// UserService is a user service.
type UserService struct{ users map[string]*User }

// Get get a user from memory.
func (u *UserService) Get(ctx context.Context, user *User) (*User, error) {
	return u.users[user.ID], nil
}

// Add add a user to memory.
func (u *UserService) Add(ctx context.Context, user *User) (*User, error) {
	u.users[user.ID] = user
	return user, nil
}

func main() {
	us := &UserService{
		users: make(map[string]*User),
	}
	router := mux.NewRouter()
	router.Handle("/users/add", http.NewHandler(us.Add)).Methods("POST")
	router.Handle("/users/detail", http.NewHandler(us.Get)).Methods("GET")

	httpSrv := http.NewServer(http.Address(":8000"))
	httpSrv.HandlePrefix("/", router)

	app := kratos.New(
		kratos.Name("handler"),
		kratos.Server(
			httpSrv,
		),
	)
	if err := app.Run(); err != nil {
		log.Println(err)
	}
}
