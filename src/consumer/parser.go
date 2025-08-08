package consumer

import (
	"fmt"
	"net/http"
	"university_bot/src/db"
	"university_bot/src/service"
	"university_bot/src/types"
)

func ParserScore(user *types.User) (http.Response, error) {
	return http.Response{}, nil
	// TODO: сделать
}

func Parser() chan interface{} {
	c := make(chan interface{})
	go func() {
		for {
			ParserChangedScore(c)
		}
	}()
	return c
}

func ParserChangedScore(c chan interface{}) {
	var users, err = db.GetAllUsers()
	if err != nil {
		fmt.Println("Возникли трудности с получением слайса всех юзеров", err)
	}
	for _, user := range users {
		var response, err = ParserScore(&user)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		sliceScores, err := service.ComparisonScore(user, response)
		if err != nil {
			fmt.Println(err.Error())
		}
		if sliceScores != nil {
			oldScores, err := db.GetUserScores(user.Login)
			if err != nil {
				fmt.Println(err.Error())
			}
			var sliceNotifications = service.ConvertToNotifications(oldScores, sliceScores)
			c <- sliceNotifications
			if err := db.UpdateUserScores(sliceScores); err != nil {
				fmt.Println(err.Error())
			}
		}
	}
}
