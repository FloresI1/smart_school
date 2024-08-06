package main

import (
    "log"
    "net/http"
    "github.com/FloresI1/smart_school/db"
    "github.com/FloresI1/smart_school/handlers"
)

func main() {
    db.InitDB()

    http.HandleFunc("/create-material", handlers.CreateMaterialHandler)
    http.HandleFunc("/get-material", handlers.GetMaterialHandler)
    http.HandleFunc("/update-material", handlers.UpdateMaterialHandler)
    http.HandleFunc("/get-all-materials", handlers.GetAllMaterialsHandler)

    log.Println("Server is running on port 8080...")
    log.Fatal(http.ListenAndServe(":8080", nil))
}