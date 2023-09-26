package main

import (
	"context"
	"fmt"
	"html/template"
	"net/http"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

// структура чтобы хранить состояние погоды время и курсы валют которые выводятся в header
type Header struct {
	Weather    string `json:"weather"`    //текущая погода воронеж
	Time       string `json:"time"`       //текущие точное время воронеж
	DollarRate string `json:"dollarrate"` //курс доллара
	EuroRate   string `json:"eurorate"`   //курс евро
	UanRate    string `json:"uanrate"`    //курсюаня
}

type News struct {
	Url           string `json:"url"`           //Url новости
	UrlSmallImage string `json:"urlsmallimage"` //получаем маленькое изображение
	UrlImage      string `json:"urlimage"`      //получаем полноразмерное изображение
	Category      string `json:"category"`      //категория новости
	Time          string `json:"time"`          //дата публикации
	Text          string `json:"text"`          //текст новости
	Title         string `json:"title"`         //Заголовок новости
}

// функция возвращает последнюю новость по времени
func (n News) LastforNews() News {
	url_database := "postgres://nwuser:nwpassword@supportdev.ru:5432/news"
	dbPool, err := pgxpool.Connect(context.Background(), url_database)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Подключение открыто успешно!" + "\n" + "Получаем последнюю новость по времени")
	/*так как у нас в базе используется поле даты и времени в типе стринг получаем значение всех новостей, а затем сортируем, пильнуть тикет на изменение типа поля*/
	rows, err := dbPool.Query(context.Background(), "select * from news.news order by timenews desc limit 1;")
	if err != nil {
		fmt.Println(err)
	}
	defer dbPool.Close()

	var accident News //буферная переменная для получения записи из бд
	for rows.Next() {
		values, err := rows.Values()
		if err != nil {
			fmt.Println(err)
		}
		accident.Category = values[1].(string)
		accident.Time = values[2].(string)
		accident.Text = values[3].(string)
		accident.Title = values[4].(string)
		accident.UrlImage = values[5].(string)
		accident.UrlSmallImage = values[6].(string)
		accident.Url = values[7].(string)
	}
	return accident
	/*можно будет удалить после изменения типа поля*/
	/*var accidents []News //массив новостей
	var accident News    //буферная переменная для получения записи из бд
	for rows.Next() {
		values, err := rows.Values()
		if err != nil {
			fmt.Println(err)
		}
		/*первая колонка в бд автоинкремент записи поэтому ее не учитываем парсим значения со второй колонки*/
	/*accident.Category = values[1].(string)
		accident.Time = values[2].(string)
		accident.Text = values[3].(string)
		accident.Title = values[4].(string)
		accident.UrlImage = values[5].(string)
		accident.UrlSmallImage = values[6].(string)
		accident.Url = values[7].(string)
		fmt.Println("----------------------------------------------------------------------------------------------------------------------------------------------------------------")
		fmt.Println(accident)
		fmt.Println("----------------------------------------------------------------------------------------------------------------------------------------------------------------")
		accidents = append(accidents, accident)
	}
	SortUPTime(accidents)
	fmt.Println()*/
}

// функция сортировки по времени чтобы вернуть значение последней по времени
// функция сортировки работает не правильно
func SortUPTime(accidents []News) News {
	for i := 0; i < len(accidents)-1; i++ {
		//4minIndex := i
		for j := 0; j < len(accidents)-i-1; j++ {
			date1, _ := time.Parse("20060102", accidents[i].Time)
			date2, _ := time.Parse("20060102", accidents[j].Time)
			if date1.Unix() < date2.Unix() {
				//minIndex = j
				buff := accidents[i]
				accidents[i] = accidents[j]
				accidents[j] = buff
			}
		}
		//accidents[i], accidents[minIndex] = accidents[minIndex], accidents[i]
	}
	return accidents[len(accidents)-1]
}

// Функция сортировки по категориям
func SortCategory(accidents []News) []News {
	return accidents
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
		fmt.Println("----------------------------------------------------------------------------------------------------------------------------------------------------------------")
		fmt.Println(accident)
		fmt.Println("----------------------------------------------------------------------------------------------------------------------------------------------------------------")
		accidents = append(accidents, accident)
	}
	//Вывод всех новостей
	//fmt.Println(accidents)
	fmt.Println("\n\n\n\n Последняя новость была в: " + accident.LastforNews().Time + "\n\n\n\n")

	return accidents
}

// //WEB logic
func viewPage(w http.ResponseWriter, r *http.Request) {
	tmpl, _ := template.ParseFiles("templates/viewpage.html")
	tmpl.Execute(w, selectNews())
}

func handleRequest() {
	http.HandleFunc("/view", viewPage)
	http.Handle("/templates/images/", http.StripPrefix("/templates/images/", http.FileServer(http.Dir("templates/images"))))
	http.Handle("/templates/templates/assets/", http.StripPrefix("/templates/templates/assets/", http.FileServer(http.Dir("templates/templates/assets"))))
	http.Handle("/templates/templates/assets/js/", http.StripPrefix("/templates/templates/assets/js/", http.FileServer(http.Dir("templates/templates/assets/js"))))
	http.Handle("/templates/templates/assets/css/", http.StripPrefix("/templates/templates/assets/css/", http.FileServer(http.Dir("templates/templates/assets/css"))))
	http.Handle("/templates/templates/assets/img/", http.StripPrefix("/templates/templates/assets/img/", http.FileServer(http.Dir("templates/templates/assets/img"))))
	http.Handle("/templates/templates/assets/img/logo/", http.StripPrefix("/templates/templates/assets/img/logo/", http.FileServer(http.Dir("templates/templates/assets/img/logo"))))
	http.ListenAndServe(":8082", nil)
}
