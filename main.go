package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

type Products struct {
	quantity       int64
	price          int64
	isRefrigirator bool
}

func main() {
	_ = RunWebPortal()
}

func RunWebPortal() error {
	http.HandleFunc("/", rootHandler)
	return http.ListenAndServe("localhost", nil)
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	if _, err := os.Stat("app.txt"); os.IsNotExist(err) {
		_, err := os.Create("app.txt")
		if err != nil {
			panic(err)
		}
	}

	file, _ := os.OpenFile("app.txt", os.O_RDWR|os.O_APPEND, 0660)

	productName := r.URL.Query().Get("name")
	quantity := r.URL.Query().Get("quantity")
	price := r.URL.Query().Get("price")
	selectedType := r.URL.Query().Get("type")

	success := true

	if selectedType == "normal" {
		_, _ = file.WriteString("Product is of type: normal \n")

		if idx := strings.Index(productName, "_"); idx != -1 {
			productName = productName + "_"
			result := productName[:idx]
			productName = strings.ReplaceAll(productName, result, "")
			productName = strings.TrimLeft(productName, "_")
			pos := strings.Index(productName, "_")
			id, _ := strconv.ParseInt(productName[:pos], 0, 64)

			db, err := sql.Open("mysql", "snappmarket:snappmarket@/snappmarket")
			if err != nil {
				log.Fatal(err)
			}

			defer db.Close()

			row := db.QueryRow("select * from products where id = ?", id)
			product := Products{}
			err = row.Scan(&product.price, &product.quantity)

			updateResult, err := db.Exec("update products set `quantity` = ? and price = ? where id = ?", quantity, price, id)
		}
	} else if selectedType == "refrigerator" {
		_, _ = file.WriteString("Product is of type: refrigerator \n")

		if idx := strings.Index(productName, "_"); idx != -1 {
			productName = productName + "_"
			result := productName[:idx]
			productName = strings.ReplaceAll(productName, result, "")
			productName = strings.TrimLeft(productName, "_")
			pos := strings.Index(productName, "_")
			id, _ := strconv.ParseInt(productName[:pos], 0, 64)

			db, err := sql.Open("mysql", "snappmarket:snappmarket@/snappmarket")
			if err != nil {
				log.Fatal(err)
			}

			defer db.Close()

			row := db.QueryRow("select * from products where id = ?", id)
			product := Products{}
			err = row.Scan(&product.price, &product.quantity)

			updateResult, err := db.Exec("update products set `quantity` = ? and price = ? and isRefrigirator = ? where id = ?", quantity, price, true, id)
		}
	}

	if err := recover(); err != nil {
		success = false
	}

	statusCode := http.StatusInternalServerError
	statusText := "error"

	if success {
		statusCode = http.StatusOK
		statusText = "ok"
	}

	data := make(map[interface{}]interface{})
	data["status"] = statusText

	w.WriteHeader(statusCode)

	_ = json.NewEncoder(w).Encode(data)
}
