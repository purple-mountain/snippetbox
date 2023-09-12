package main

import (
	"log"
	"net/http"
)

func helloWorld(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	w.Write([]byte("Hello from Snake Den"))
}

func createSnippet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		// w.Header().Set("Allow", "POST")
		// w.WriteHeader(405)
		// w.Write([]byte("Only POST method allowed\n"))
		// Same thing
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Write([]byte("Create a new snippet\n"))
}

func deleteSnippet(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Delete a snippet"))
}

func viewSnippet(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("View snippet"))
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", helloWorld)
	mux.HandleFunc("/snippet/view", viewSnippet)
	mux.HandleFunc("/snippet/create", createSnippet)
	mux.HandleFunc("/snippet/delete", deleteSnippet)

	log.Println("The server is running on http://localhost:4000")
	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}
