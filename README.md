# Secret Santa distributor service made in Go

## Usage
Create a file `players/players.txt` containing the list of players (one per line).
Create a file `.env` specifying the `SECRET_SANTA_URL` variable.

Feel free to customize the `templates/template.html` file to your liking before building the image.

Build the docker image:
```sh
docker build --tag 'secret-santa' .
```

During build, the program will shuffle the list of players and store the player/friend pairs internally in the image, in a `friends.json` file.
This is done to ensure the pairs remain stable for each image and also to obfuscate the results to the person building it.

Run the docker image:
```sh
docker compose up -d
```

The links can be found in the `players/links.txt` file.
They can be


