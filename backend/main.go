package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"
)

type Item struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func env(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func main() {
	dsn := "host=" + env("DB_HOST", "db") +
		" port=" + env("DB_PORT", "5432") +
		" user=" + env("DB_USER", "") +
		" password=" + env("DB_PASSWORD", "") +
		" dbname=" + env("DB_NAME", "") +
		" sslmode=disable"

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("open database: ", err)
	}
	defer db.Close()

	// PostgreSQL может запускаться дольше приложения: повторяем подключение.
	connected := false
	for attempt := 1; attempt <= 30; attempt++ {
		if err := db.Ping(); err == nil {
			connected = true
			break
		}
		log.Printf("database is not ready (attempt %d/30): %v", attempt, err)
		time.Sleep(2 * time.Second)
	}
	if !connected {
		log.Fatal("could not connect to PostgreSQL after 60 seconds")
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS items (
		id SERIAL PRIMARY KEY,
		name TEXT NOT NULL
	)`)
	if err != nil {
		log.Fatal("create items table: ", err)
	}

	var count int
	if err := db.QueryRow("SELECT COUNT(*) FROM items").Scan(&count); err != nil {
		log.Fatal("count items: ", err)
	}
	if count == 0 {
		_, err = db.Exec("INSERT INTO items (name) VALUES ($1), ($2)", "Первая запись", "Вторая запись")
		if err != nil {
			log.Fatal("insert demo items: ", err)
		}
	}

	staticDir := env("STATIC_DIR", "/app/static")
	fileServer := http.FileServer(http.Dir(staticDir))

	http.HandleFunc("/api/items", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		rows, err := db.Query("SELECT id, name FROM items ORDER BY id")
		if err != nil {
			http.Error(w, "database query failed", http.StatusInternalServerError)
			log.Println("query items:", err)
			return
		}
		defer rows.Close()

		items := make([]Item, 0)
		for rows.Next() {
			var item Item
			if err := rows.Scan(&item.ID, &item.Name); err != nil {
				http.Error(w, "database read failed", http.StatusInternalServerError)
				return
			}
			items = append(items, item)
		}
		if err := rows.Err(); err != nil {
			http.Error(w, "database read failed", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		_ = json.NewEncoder(w).Encode(items)
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			r.URL.Path = "/index.html"
		}
		fileServer.ServeHTTP(w, r)
	})

	addr := ":" + env("HTTP_PORT", "8080")
	server := &http.Server{Addr: addr, ReadHeaderTimeout: 5 * time.Second}
	log.Printf("web application listening on %s", addr)
	log.Fatal(server.ListenAndServe())
}
