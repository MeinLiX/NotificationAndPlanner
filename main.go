package main

import (
	"encoding/json"
	"gopkg.in/Iwark/spreadsheet.v2"
	google_sheets "meinlix.inc/NotificationAndPlanner.v1/integration/google-sheets"
	utils "meinlix.inc/NotificationAndPlanner.v1/utils"
)

// main
func main() {
	
		// bot, err := tgbotapi.NewBotAPI("6247398053:AAFxrmydgxG7kVxooS0BDOSoHpIzTOlVcBA")
		// if err != nil {

		// }

		// bot.Debug = true

		// u := tgbotapi.NewUpdate(0)
		// u.Timeout = 60

		// updates := bot.GetUpdatesChan(u)

		// for update := range updates {
		// 	if update.Message != nil {
		// 		msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
		// 		msg.ReplyToMessageID = update.Message.MessageID

		// 		bot.Send(msg)
		// 	}
		// }
	
	config := utils.GetConfig()
	service, err := spreadsheet.NewService()
	utils.PanicIfError(err)
	spreadsheet, err := service.FetchSpreadsheet(config.SheetId)
	//https://docs.google.com/spreadsheets/d/1p78c-6nlL2HKfKpJpuLr_kYEumI_ZooracEd_c7e7hQ
	utils.PanicIfError(err)

	sheetsContext := google_sheets.SheetContext{
		Spreadsheet:       &spreadsheet,
		IsInitialized: false,

		PlanActionArgs: google_sheets.SheetArgs{
			SheetID: 1508707366,
		},
		PlanActionDetailsArgs: google_sheets.SheetArgs{
			SheetID: 0,
		},
		PlanActionUsersArgs: google_sheets.SheetArgs{
			SheetID: 0,
		},
	}

	success, err := sheetsContext.SyncData(true)
	utils.PanicIfError(err)
	if success == nil {
		panic("sync sheets failed")
	}
	//TEST
	f1, err := json.Marshal(success.AvailableUsers)
	utils.PanicIfError(err)
	f2, err := json.Marshal(success.AvailableWorks)
	utils.PanicIfError(err)
	f3, err := json.Marshal(success.AvailableActions)
	utils.PanicIfError(err)
	println(string(f1))
	println(string(f2))
	println(string(f3))
	/*
	   sheetsContext.PlanActionSheet.Update(0, 0, "DB")
	*/
}

//CONTEXT Sheets
//CONTEXT Sheets
//CONTEXT Sheets
//CONTEXT Sheets
//CONTEXT Sheets
//CONTEXT Sheets
//CONTEXT Sheets
//CONTEXT Sheets
