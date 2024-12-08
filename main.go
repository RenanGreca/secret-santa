package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand/v2"
	"net/http"
	"os"
	"secret-santa/file"
	"strings"
	"text/template"

	"github.com/gosimple/slug"
)

type Players map[string]map[string]string
type Player struct {
	Slug   string
	Name   string
	Token  string
	Friend string
}

const (
	playersFile  = "players/players.txt"
	friendsFile  = "friends.json"
	templateFile = "templates/template.html"
	linksFile    = "players/links.txt"
)

func main() {

	// Command-line arguments for the initial generation of player-friend pairs
	// Usage: go run main.go shuffle
	if len(os.Args) >= 2 {
		command := os.Args[1]
		if command == "shuffle" {
			shufflePlayers()
			os.Exit(0)
		}
		if command == "show" {
			players := getPlayers()
			fmt.Println(players)
			os.Exit(0)
		}
		log.Fatalln("Please specify a command: shuffle or show")
	}

	generateLinks()

	mux := http.NewServeMux()
	mux.HandleFunc("GET /{player}", playerHandler)

	err := http.ListenAndServe(":3000", mux)
	if err != nil {
		log.Fatalln("Error starting server:", err)
	}
}

// playerHandler handles GET requests to /{player}.
// It displays the player's secret friend.
func playerHandler(w http.ResponseWriter, req *http.Request) {
	slug := req.PathValue("player")

	players := getPlayers()
	player := parsePlayer(players, slug)
	query := req.URL.Query()
	token := query.Get("token")
	if token != player.Token {
		log.Println("Wrong token")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	tmpl := template.Must(template.ParseFiles(templateFile))
	tmpl.Execute(w, player)
}

func parsePlayer(players Players, slug string) Player {
	playerdict, ok := players[slug]
	if !ok {
		log.Println("Player not found:", slug)
	}

	return Player{
		Slug:   slug,
		Name:   playerdict["name"],
		Token:  playerdict["token"],
		Friend: playerdict["friend"],
	}
}

func getPlayers() Players {
	players := file.ReadFile(friendsFile)

	dict := Players{}
	err := json.Unmarshal([]byte(players), &dict)
	if err != nil {
		log.Fatalln(err)
	}

	return dict
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randomString(n int) string {

	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.IntN(len(letters))]
	}
	return string(b)
}

func shufflePlayers() {
	fileExists := file.FileExists(playersFile)
	if !fileExists {
		log.Fatalln(playersFile + " does not exist")
	}

	players := strings.Split(file.ReadFile(playersFile), "\n")
	players = players[:len(players)-1]

	dict := Players{}
	for {
		no_conflicts := true
		shuffled_indexes := rand.Perm(len(players))
		for i := 0; i < len(players); i++ {
			if players[i] == players[shuffled_indexes[i]] {
				no_conflicts = false
				break
			}
			slug := slug.Make(players[i])
			dict[slug] = map[string]string{}
			dict[slug]["slug"] = slug
			dict[slug]["name"] = players[i]
			dict[slug]["token"] = randomString(10)
			dict[slug]["friend"] = players[shuffled_indexes[i]]
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

	file.Create(friendsFile, string(json))
}

func generateLinks() {
	dict := getPlayers()

	links := ""
	for k, v := range dict {
		link := fmt.Sprintf("%s: https://amigosecreto.renangreca.com/%s?token=%s\n", v["name"], k, v["token"])
		links = links + link
	}

	file.Create(linksFile, links)
}
