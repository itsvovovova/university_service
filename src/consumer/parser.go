package consumer

import (
	"fmt"
	"net/http"
	"time"
	"university_bot/src/db"
	"university_bot/src/types"
)

func parserScore(user *types.User) (http.Response, error) {
	return http.Response{}, nil
	// TODO: сделать
}

func Parser() chan interface{} {
	c := make(chan interface{})
	go func() {
		for {
			// расписание лк, оценки лк, дедлайны еду [range(chat_id) -> (user, password) -> parserScore, ...]
			var users, err = db.GetAllUsers()
			if err != nil {
				fmt.Println("Возникли трудности с получением мапы всех юзеров", err)
			}
			for _, user := range users {
				var response, err = parserScore(&user)
				if err != nil {
					fmt.Println(err.Error())
					return
				}
				ComparisonScore(user, response)
				c <- data                          // отправка данных в канал
				time.Sleep(100 * time.Millisecond) // пауза между запросами
			}
		}
	}()
	return c
}

func ProcessData(data chan interface{}) {
	for v := range data {
		switch v.(type) {
		case types.Score:
			// pass
		case types.Deadline:
			// pass
		case types.Schedule:
			// pass
		}
	}
}
