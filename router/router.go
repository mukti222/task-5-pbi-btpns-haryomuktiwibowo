// project/router/router.go
package router

import (
    "net/http"
    "task-5-pbi-btpns-haryomuktiwibowo/controllers"
    "task-5-pbi-btpns-haryomuktiwibowo/middlewares"
)

func InitRouter() {
    // Mengatur route untuk endpoint POST /users/register
    http.HandleFunc("/users/register", controllers.RegisterUser)

    // Mengatur route untuk endpoint POST /users/login
    http.HandleFunc("/users/login", controllers.LoginUser)

    // Mengatur route untuk endpoint DELETE /users/:userId
    http.HandleFunc("/users/", middlewares.Authenticate(controllers.DeleteUser))
}
