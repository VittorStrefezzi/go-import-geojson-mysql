package main

import (
	"database/sql"
	"encoding/json"
	"io/ioutil"

	_ "github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
)

type GeoJSON struct {
	Type     string `json:"type"`
	Features []struct {
		Type       string `json:"type"`
		Properties struct {
			ID          string `json:"id"`
			Name        string `json:"name"`
			Description string `json:"description"`
		} `json:"properties"`
		Geometry struct {
			Type        string        `json:"type"`
			Coordinates [][][]float64 `json:"coordinates"`
		} `json:"geometry"`
	} `json:"features"`
}

func main() {
	arq, err := ioutil.ReadFile("geo.json")
	if err != nil {
		logrus.Errorf("Error read file => %s", err)
	}

	var geojson GeoJSON
	json.Unmarshal(arq, &geojson)

	userDB := "ysk_dbu"
	passDB := "ysk_pass"
	schema := "yii2-starter-kit"

	db, err := sql.Open("mysql", userDB+":"+passDB+"@/"+schema)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	tx, _ := db.Begin()
	stmt, _ := tx.Prepare("insert into city(ibge_id, name, description, json) values(?,?,?,?)")

	for i, item := range geojson.Features {
		data, _ := json.Marshal(geojson.Features[i])
		_, err = stmt.Exec(item.Properties.ID, item.Properties.Name, item.Properties.Description, data)

		if err != nil {
			tx.Rollback()
			logrus.Fatal("error in transaction => %s", err)
		}
	}

	tx.Commit()

}
