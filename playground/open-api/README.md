# commands
validate

```sh
$ swagger validate ./swagger.yml
```

generate

```sh
$ swagger generate server -A todo-list -f ./swagger.yml
```

serve
```sh
$ go run cmd/todo-list-server/main.go --port 3003
```

server: http://127.0.0.1:3003
swagger ui: http://127.0.0.1:3003/docs
