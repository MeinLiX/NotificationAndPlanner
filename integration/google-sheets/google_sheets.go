package google_sheets

import (
	"fmt"
	"strings"
	"time"

	"gopkg.in/Iwark/spreadsheet.v2"

	utils "meinlix.inc/NotificationAndPlanner.v1/utils"
)

type SheetContext struct {
	Spreadsheet       *spreadsheet.Spreadsheet
	IsInitialized bool

	PlanActionArgs        SheetArgs
	PlanActionDetailsArgs SheetArgs
	PlanActionUsersArgs   SheetArgs

	AvaliableActions []PlanAction
	AvaliableWorks   []PlanActionDetails
	AvaliableUsers   []PlanActionUsers

	PlanActionSheet        *spreadsheet.Sheet
	PlanActionDetailsSheet *spreadsheet.Sheet
	PlanActionUsersSheet   *spreadsheet.Sheet
}

//моделі

// modified (required sheet)
type PlanAction struct {
	Uuid           string        `json:"PA_id"`
	TgUserName     string        `json:"PA_tg_user_name"`
	SelectedWorkId uint          `json:"PA_selected_work_id"`
	StartWork      time.Time     `json:"PA_start_work_time"`
	DurationWork   time.Duration `json:"PA_duration_work"`
	CompletedWork  uint          `json:"PA_completed_work"`
	rowNumber      int           //local field
}

type PlanActionDetails struct {
	WorkId               uint          `json:"PAD_work_id"`
	DescriptionWork      string        `json:"PAD_description_work"`
	WorkPeriodicity      uint          `json:"PAD_work_periodicity"`
	PlannedDurationWork  time.Duration `json:"PAD_planned_duration_work"`
	TgUsernamesPossibles []string      `json:"PAD_tg_usernames_possibles"`
	rowNumber         	 int           //local field
}

type PlanActionUsers struct {
	TgUserName  string `json:"PAU_tg_user_name"`
	Name        string `json:"PAU_name"`
	rowNumber 	int    //local fields
}

type SheetArgs struct {
	RowStart    uint
	ColumnStart uint
	CountItem   uint
	SheetID     uint
}

type ISheetContext interface {
	GetPlanActionDetails(id int, toSync bool) (PlanActionDetails, error)        // GetPlanActionDetails not modified from bot
	GetPlanActionUsers(tgUserName string, toSync bool) (PlanActionUsers, error) // PlanActionUsers not modified from bot
	GetPlanAction(guid string, toSync bool) (PlanAction, error)

	GetPlanActionsByUser(tgUserName string) ([]PlanAction, error)
	InsertOrUpdatePlanAction(planAction PlanAction) (bool, error)
	RemovePlanAction(guid string) (bool, error)

	SyncDatas(syncCheets bool) (*SheetContext, error)
	SyncSheets() (*SheetContext, error)
	GetDescriptionsSheets() string

	RegeneratePlanActionDetails(days int) 
}

// ready
func (ctx *SheetContext) SyncDatas(syncSheets bool) (*SheetContext, error) {
	if syncSheets || !ctx.IsInitialized {
		if ctxnew, err := ctx.SyncSheets(); err != nil {
			return nil, err
		} else {
			ctx = ctxnew
		}
	}

	//PlanActionUsers
	firstRow := ctx.PlanActionUsersArgs.RowStart + 1
	firstColumn := ctx.PlanActionUsersArgs.ColumnStart
	currentSheet := ctx.PlanActionUsersSheet
	ctx.AvaliableUsers = make([]PlanActionUsers, ctx.PlanActionUsersArgs.CountItem)
	for i := uint(0); i < ctx.PlanActionUsersArgs.CountItem; i++ {
		currentRow := firstRow + i
		ctx.AvaliableUsers[i] = PlanActionUsers{
			TgUserName:  currentSheet.Rows[currentRow][firstColumn].Value,
			Name:        currentSheet.Rows[currentRow][firstColumn+1].Value,
			rowNumber: int(currentRow),
		}
	}

	//PlanActionDetails
	firstRow = ctx.PlanActionDetailsArgs.RowStart + 1
	firstColumn = ctx.PlanActionDetailsArgs.ColumnStart
	currentSheet = ctx.PlanActionDetailsSheet
	ctx.AvaliableWorks = make([]PlanActionDetails, ctx.PlanActionDetailsArgs.CountItem)
	for i := uint(0); i < ctx.PlanActionDetailsArgs.CountItem; i++ {
		currentRow := firstRow + i
		workId, err := utils.StringToUint(currentSheet.Rows[currentRow][firstColumn].Value)
		if err != nil {
			return nil, err
		}
		workPeriodicity, err := utils.StringToUint(currentSheet.Rows[currentRow][firstColumn+2].Value)
		if err != nil {
			return nil, err
		}
		avaliableUsers:=strings.Split(currentSheet.Rows[currentRow][firstColumn+4].Value, ";")
		if(utils.Contains(avaliableUsers, "")){
			avaliableUsers=nil;
		}
		ctx.AvaliableWorks[i] = PlanActionDetails{
			WorkId:              workId,
			DescriptionWork:     currentSheet.Rows[currentRow][firstColumn+1].Value,
			WorkPeriodicity:     workPeriodicity,
			PlannedDurationWork: *utils.StringToDuration(currentSheet.Rows[currentRow][firstColumn+3].Value),
			TgUsernamesPossibles: avaliableUsers,
			rowNumber:         int(currentRow),
		}
	}

	//PlanAction
	firstRow = ctx.PlanActionArgs.RowStart + 1
	firstColumn = ctx.PlanActionArgs.ColumnStart
	currentSheet = ctx.PlanActionSheet
	ctx.AvaliableActions = make([]PlanAction, ctx.PlanActionArgs.CountItem)
	for i := uint(0); i < ctx.PlanActionArgs.CountItem; i++ {
		currentRow := firstRow + i
		selectedWorkId, err := utils.StringToUint(currentSheet.Rows[currentRow][firstColumn+2].Value)
		if err != nil {
			return nil, err
		}
		startWork, err := time.Parse(time.DateOnly, currentSheet.Rows[currentRow][firstColumn+3].Value)
		if err != nil {
			return nil, err
		}
		completedWork, err := utils.StringToUint(currentSheet.Rows[currentRow][firstColumn+5].Value)
		if err != nil {
			return nil, err
		}
		ctx.AvaliableActions[i] = PlanAction{
			Uuid:           currentSheet.Rows[currentRow][firstColumn].Value,
			TgUserName:     currentSheet.Rows[currentRow][firstColumn+1].Value,
			SelectedWorkId: selectedWorkId,
			StartWork:      startWork,
			DurationWork:   *utils.StringToDuration(currentSheet.Rows[currentRow][firstColumn+4].Value),
			CompletedWork:  completedWork,
			rowNumber:    int(currentRow),
		}
	}

	return ctx, nil
}

func (ctx SheetContext) GetDescriptionsSheets() string {
	return "TOTO INSTRUCTIONS HOW TO FILL GOOGLE SHEETS"
}

// ready
func (ctx *SheetContext) SyncSheets() (*SheetContext, error) {
	ctx.IsInitialized = false
	if ctx.Spreadsheet == nil {
		return nil, fmt.Errorf("spreadsheet not initialized")
	}

	if planActionSheet, err := ctx.Spreadsheet.SheetByID(ctx.PlanActionArgs.SheetID); err != nil {
		return nil, err
	} else {
		planActionSheet.Synchronize()
		ctx.PlanActionSheet = planActionSheet

		if args, err := getSheetArgs(ctx.PlanActionSheet, "PA_id", false); err != nil {
			return nil, err
		} else {
			ctx.PlanActionArgs = *args
		}
	}
	if planActionDetailsSheet, err := ctx.Spreadsheet.SheetByID(ctx.PlanActionDetailsArgs.SheetID); err != nil {
		return nil, err
	} else {
		planActionDetailsSheet.Synchronize()
		ctx.PlanActionDetailsSheet = planActionDetailsSheet

		if args, err := getSheetArgs(ctx.PlanActionDetailsSheet, "PAD_work_id", false); err != nil {
			return nil, err
		} else {
			ctx.PlanActionDetailsArgs = *args
		}
	}
	if planActionUsersSheet, err := ctx.Spreadsheet.SheetByID(ctx.PlanActionUsersArgs.SheetID); err != nil {
		return nil, err
	} else {
		planActionUsersSheet.Synchronize()
		ctx.PlanActionUsersSheet = planActionUsersSheet

		if args, err := getSheetArgs(ctx.PlanActionUsersSheet, "PAU_tg_user_name", false); err != nil {
			return nil, err
		} else {
			ctx.PlanActionUsersArgs = *args
		}
	}
	ctx.IsInitialized = true
	return ctx, nil
}

func (ctx SheetContext) InsertOrUpdateAction(planAction PlanAction) (bool, error) {
	if _, err := ctx.Spreadsheet.SheetByIndex(ctx.PlanActionArgs.SheetID); err != nil {
		return false, err
	}

	return false, fmt.Errorf("TODO")
}

func (ctx SheetContext) GetActions(tgUserName string) ([]PlanAction, error) {
	//planActions := list.New()

	return nil, fmt.Errorf("TODO")
}

func getCellByValue(sheet *spreadsheet.Sheet, valueToFound string) (spreadsheet.Cell, error) {
	for _, row := range sheet.Rows {
		for _, cell := range row {
			if cell.Value == valueToFound {
				return cell, nil
			}
		}
	}
	return spreadsheet.Cell{}, fmt.Errorf("cell with '%s' value not found", valueToFound)
}

func getSheetArgs(sheet *spreadsheet.Sheet, startTableColumnName string, syncSheet bool) (*SheetArgs, error) {
	if syncSheet {
		sheet.Synchronize()
	}
	if startCell, err := getCellByValue(sheet, startTableColumnName); err != nil {
		return nil, err
	} else {
		counterCellValue := (sheet.Rows[startCell.Row][startCell.Column-1]).Value
		counterFounded, err := utils.StringToUint(counterCellValue)
		if err != nil {
			return nil, err
		}

		sheetArgs := SheetArgs{
			RowStart:    startCell.Row,
			ColumnStart: startCell.Column,
			CountItem:   counterFounded,
			SheetID:     sheet.Properties.Index,
		}

		return &sheetArgs, nil
	}
}
