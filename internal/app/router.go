package app

import "net/http"

func NewRouter(c *Container) {
	http.HandleFunc("/users/register", c.UserHandler.Register)
}
