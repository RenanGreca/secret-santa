package main

import (
	"encoding/json"
	"log"
	"math/rand/v2"
	"net/http"
	"os"
	"secret-santa/file"
	"strings"

	"github.com/gosimple/slug"
)

type Players map[string]string

func main() {
	if len(os.Args) >= 2 {
		command := os.Args[1]
		if command == "shuffle" {
			shufflePlayers()
			os.Exit(0)
		}
		if command == "show" {
			getPlayers()
			os.Exit(0)
		}
		log.Fatalln("Please specify a command: shuffle or show")
	}

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatalln("Please set the PORT environment variable")
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /{player}", playerHandler)

	err := http.ListenAndServe(":"+port, mux)
	if err != nil {
		log.Fatalln("Error starting server:", err)
	}
}

func playerHandler(w http.ResponseWriter, req *http.Request) {
	player := req.PathValue("player")

	players := getPlayers()
	secret, ok := players[player]
	if !ok {
		log.Println("Player not found:", player)
	}

	w.Write([]byte(secret))
}

func getPlayers() Players {
	players := file.ReadFile("players.json")

	dict := map[string]string{}
	err := json.Unmarshal([]byte(players), &dict)
	if err != nil {
		log.Fatalln(err)
	}

	// for k, v := range dict {
	// 	log.Printf("%s: %s\n", k, v)
	// }

	return dict
}

func shufflePlayers() {
	fileExists := file.FileExists("players.txt")
	if !fileExists {
		log.Fatalln("players.txt does not exist")
	}

	players := strings.Split(file.ReadFile("players.txt"), "\n")
	players = players[:len(players)-1]

	dict := map[string]string{}
	for {
		no_conflicts := true
		shuffled_indexes := rand.Perm(len(players))
		for i := 0; i < len(players); i++ {
			if players[i] == players[shuffled_indexes[i]] {
				no_conflicts = false
			}
			slug := slug.Make(players[i])
			dict[slug] = players[shuffled_indexes[i]]
		}
		if no_conflicts {
			break
		}
		log.Println("No conflicts:", no_conflicts)
	}
	log.Println("Conflicts resolved!")

	json, err := json.Marshal(dict)
	if err != nil {
		log.Fatalln(err)
	}

	file.Create("players.json", string(json))
}
