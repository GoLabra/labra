package main

//go:generate rm -rf ./ent/schema
//go:generate mkdir -p ./ent/schema
//go:generate cp -a ./schema/_admin/. ./ent/schema
//go:generate mkdir -p ./schema/_user
//go:generate cp -a ./schema/_user/. ./ent/schema
//go:generate mkdir -p ./schema/_user_backup/
//go:generate mkdir -p ./schema/_admin_backup/

//go:generate go run -mod=mod entc.go
//go:generate go run -mod=mod github.com/99designs/gqlgen
