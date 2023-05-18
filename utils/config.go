package utils

import (
	"encoding/json"
	"os"
)

type Config struct {
	SheetId string `json:"sheet_id"`
}

func GetConfig() Config {
	jsonFile, err := os.ReadFile("config.json")
	PanicIfError(err)
	var res Config
	json.Unmarshal(jsonFile, &res)
	return res
}