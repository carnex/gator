package main

import (
	gatorConfig "github.com/carnex/gator/internal/config"
)

func main() {
	cfg, err := gatorConfig.Read()
	if err != nil {
		println(err)
	}
	cfg.SetUser("Erik")
	post, err := gatorConfig.Read()
	if err != nil {
		println(err)
	}
	println(post.DbURL, post.CurrentUserName)

}
