package main

import (
	"context"
	"fmt"
	"html/template"
	"net/http"

	"github.com/jackc/pgx/v4/pgxpool"
)

type News struct {
	Url           string `json:"url"`           //Url новости
	UrlSmallImage string `json:"urlsmallimage"` //получаем маленькое изображение
	UrlImage      string `json:"urlimage"`      //получаем полноразмерное изображение
	Category      string `json:"category"`      //категория новости
	Time          string `json:"time"`          //дата публикации
	Text          string `json:"text"`          //текст новости
	Title         string `json:"title"`         //Заголовок новости
}

func main() {
	handleRequest()
	//selectNews()
}

func selectNews() []News {
	url_database := "postgres://nwuser:nwpassword@supportdev.ru:5432/news"
	dbPool, err := pgxpool.Connect(context.Background(), url_database)
	var accidents []News //массив новостей
	var accident News    //буферная переменная для получения записи из бд
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Подключение открыто успешно!" + "\n" + "Получаем список новостей из базы для отображения во View!")

	rows, err := dbPool.Query(context.Background(), "select * from news.news")
	if err != nil {
		fmt.Println(err)
	}

	defer dbPool.Close()
	for rows.Next() {
		values, err := rows.Values()
		if err != nil {
			fmt.Println(err)
		}
		/*первая колонка в бд автоинкремент записи поэтому ее не учитываем парсим значения со второй колонки*/
		accident.Category = values[1].(string)
		accident.Time = values[2].(string)
		accident.Text = values[3].(string)
		accident.Title = values[4].(string)
		accident.UrlImage = values[5].(string)
		accident.UrlSmallImage = values[6].(string)
		accident.Url = values[7].(string)

		accidents = append(accidents, accident)
	}
	fmt.Println(accidents)
	return accidents
}

// //WEB logic
func viewPage(w http.ResponseWriter, r *http.Request) {
	tmpl, _ := template.ParseFiles("templates/viewpage.html")
	tmpl.Execute(w, selectNews())
}

func handleRequest() {
	http.HandleFunc("/view", viewPage)
	http.ListenAndServe(":8081", nil)
}
