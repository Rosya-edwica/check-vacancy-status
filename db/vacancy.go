package db

import (
	"check_vacancy_status/models"
	"fmt"
)

func (d *Database) GetVacancies(limit int, lastId int) (vacancies []models.Vacancy) {
	query := fmt.Sprintf("SELECT DISTINCT id, is_open, platform FROM h_vacancy WHERE id > %d ORDER BY id LIMIT %d", lastId, limit)
	rows, err := d.Connection.Query(query)
	checkErr(err)
	defer rows.Close()

	for rows.Next() {
		var id int
		var status bool
		var platform string
		err = rows.Scan(&id, &status, &platform)
		checkErr(err)
		vacancies = append(vacancies, models.Vacancy{
			Id: id,
			Status: status,
			Platform: platform,
		})
	}
	return
}


func (d *Database) UpdateVacanciesStatus(vacancies []models.Vacancy) {
	for _, v := range vacancies {
		query := fmt.Sprintf(`UPDATE h_vacancy SET is_open=%t WHERE id=%d AND platform = '%s';`, v.Status, v.Id, v.Platform)
		tx, _ := d.Connection.Begin()
		_, err := d.Connection.Exec(query)
		checkErr(err)
		tx.Commit()
	}
}