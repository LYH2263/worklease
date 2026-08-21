package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"

	"example.com/worklease"
	"example.com/worklease/internal/clock"
)

func main() {
	addr := flag.String("addr", ":8107", "listen")
	web := flag.String("web", "web", "web")
	data := flag.String("data", "data", "data")
	flag.Parse()
	_ = os.MkdirAll(*data, 0o755)
	q, err := worklease.New(worklease.Options{
		Clock: clock.Real{}, PersistPath: filepath.Join(*data, "queue.json"), LeaseTTL: 5 * time.Second,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer q.Close()
	q.StartReaper()

	var mu sync.Mutex
	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir(*web)))
	mux.HandleFunc("/api/jobs", func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		defer mu.Unlock()
		writeJSON(w, q.List())
	})
	mux.HandleFunc("/api/stats", func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		defer mu.Unlock()
		writeJSON(w, q.Stats())
	})
	mux.HandleFunc("/api/enqueue", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "POST", 405)
			return
		}
		var spec worklease.EnqueueSpec
		if err := json.NewDecoder(r.Body).Decode(&spec); err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		mu.Lock()
		err := q.Enqueue(spec)
		mu.Unlock()
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		writeJSON(w, map[string]string{"ok": "1"})
	})
	mux.HandleFunc("/api/reap", func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		err := q.ReapExpired()
		mu.Unlock()
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		writeJSON(w, map[string]string{"ok": "1"})
	})
	log.Printf("workd on %s", *addr)
	log.Fatal(http.ListenAndServe(*addr, mux))
}

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(v)
}
