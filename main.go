package main

import (
	"check_vacancy_status/db"
	"check_vacancy_status/models"
	"check_vacancy_status/utils"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/joho/godotenv"
	"github.com/tidwall/gjson"
)

var userAgent = "Mozilla/5.0 (iPad; CPU OS 7_2_1 like Mac OS X; en-US) AppleWebKit/533.14.6 (KHTML, like Gecko) Version/3.0.5 Mobile/8B116 Safari/6533.14.6"
var HeadersHH = map[string]string{
	"User-Agent": userAgent,
	"Authorization": "Bearer " + os.Getenv("HEADHUNTER_TOKEN"),
}
var HeadersSuperjob = map[string]string{
	"User-Agent": userAgent,
	"Authorization": "Bearer " + os.Getenv("SUPERJOB_TOKEN"),
	"X-Api-App-Id": os.Getenv("SUPERJOB_SECRET"),
}

var updatedItems []models.Vacancy

func main() {
	s := time.Now().Unix()
	err := godotenv.Load(".env")
	utils.CheckErr(err)
	database := db.Database{
		Port:     os.Getenv("MYSQL_PORT"),
		Host:     os.Getenv("MYSQL_HOST"),
		User:     os.Getenv("MYSQL_USER"),
		Password: os.Getenv("MYSQL_PASSWORD"),
		Name:     os.Getenv("MYSQL_DATABASE"),
	}
	database.Connect()
	var last_id = 0
	for {
		vacancies := database.GetVacancies(100, last_id)
		check(vacancies)
		database.UpdateVacanciesStatus(updatedItems)
		log.Printf("Updated %d vacancies", len(updatedItems))
		if vacancies[len(vacancies)-1].Id != last_id {
			last_id = vacancies[len(vacancies)-1].Id
		} else {
			break
		}
	} 
	
	database.Close()
	fmt.Println("Time: ", time.Now().Unix() - s)
}


func check(vacancies []models.Vacancy) {
	updatedItems = []models.Vacancy{}
	var wg sync.WaitGroup
	wg.Add(len(vacancies))
	for _, vacancy := range vacancies {
		go checkVacancyStatus(vacancy, &wg)
	}
	wg.Wait()
}

func checkVacancyStatus(vacancy models.Vacancy, wg *sync.WaitGroup) {
	var updated models.Vacancy
	defer wg.Done()

	switch vacancy.Platform {
		case "hh": updated = checkVacancyStatusInHeadHunter(vacancy)
		case "superjob": updated = checkVacancyStatusInSuperjob(vacancy)
	}
	if updated.Status != vacancy.Status {
		updatedItems = append(updatedItems, updated)
	}
}

func checkVacancyStatusInHeadHunter(vacancy models.Vacancy) models.Vacancy {
	url := fmt.Sprintf("https://api.hh.ru/vacancies/%d", vacancy.Id)
	json, err := utils.GetJson(url, HeadersHH)
	utils.CheckErr(err)

	if gjson.Get(json, "archived").Bool() {
		vacancy.Status = false
		return vacancy
	}

	status := gjson.Get(json, "type.id").String()
	if status != "open" {
		vacancy.Status = false
	} else {
		vacancy.Status = true
	}
	return vacancy
}

func checkVacancyStatusInSuperjob(vacancy models.Vacancy) models.Vacancy {
	url := fmt.Sprintf("https://api.superjob.ru/2.0/vacancies/%d", vacancy.Id)
	json, err := utils.GetJson(url, HeadersSuperjob)
	utils.CheckErr(err)

	if gjson.Get(json, "is_archive").Bool() {
		vacancy.Status = false
	} else if gjson.Get(json, "is_closed").Bool() {
		vacancy.Status = false
	}
	return vacancy
}