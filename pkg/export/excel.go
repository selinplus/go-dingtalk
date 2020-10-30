package export

import (
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/selinplus/go-dingtalk/pkg/setting"
	"github.com/selinplus/go-dingtalk/pkg/upload"
	"strconv"
	"strings"
	"time"
)

const EXT = ".xlsx"

// GetExcelFullUrl get the full access path of the Excel file
func GetExcelFullUrl(name string) string {
	return setting.AppSetting.PrefixUrl + "/" + GetExcelPath() + name
}

// GetExcelPath get the relative save path of the Excel file
func GetExcelPath() string {
	return setting.AppSetting.ExportSavePath
}

// GetExcelFullPath Get the full save path of the Excel file
func GetExcelFullPath() string {
	return setting.AppSetting.RuntimeRootPath + GetExcelPath()
}

// Write into excel
func WriteIntoExecel(fileName string, tbHeader []string, records []map[string]string) (string, error) {
	var sheetName = "Sheet1"
	xlsx := excelize.NewFile()
	index := xlsx.NewSheet(sheetName)
	// Set table head
	var A = 'A'
	for i, field := range tbHeader {
		var cell string
		if i < 26 {
			cell = fmt.Sprintf("%c", A)
			A++
		} else if i == 26 {
			A = 'A'
			cell = fmt.Sprintf("A%c", A)
			A++
		} else {
			cell = fmt.Sprintf("A%c", A)
			A++
		}
		_ = xlsx.SetCellValue(sheetName, cell+"1", field)
		_ = xlsx.SetColWidth(sheetName, cell, cell, countWidth(field))
	}
	// Set cell value
	for row, record := range records {
		var A = 'A'
		for i, field := range tbHeader {
			var cell string
			if i < 26 {
				cell = fmt.Sprintf("%c", A) + strconv.Itoa(row+2)
				A++
			} else if i == 26 {
				A = 'A'
				cell = fmt.Sprintf("A%c", A) + strconv.Itoa(row+2)
				A++
			} else {
				cell = fmt.Sprintf("A%c", A) + strconv.Itoa(row+2)
				A++
			}
			_ = xlsx.SetCellValue(sheetName, cell, record[field])
		}
	}
	// Set active sheet of the workbook
	xlsx.SetActiveSheet(index)
	// Save xlsx file by the given path
	savePath := GetExcelFullPath()
	if err := upload.CheckImage(savePath); err != nil {
		//log.Println(err)
		return "", err
	}
	saveName := fileName + strconv.Itoa(int(time.Now().Unix())) + EXT
	scr := savePath + saveName
	if err := xlsx.SaveAs(scr); err != nil {
		//log.Println(err)
		return "", err
	}
	return GetExcelFullUrl(saveName), nil
}

// according to cell value, count colwidth
func countWidth(value string) float64 {
	letters := "abcdefghijklmnopqrstuvwxyz"
	letters = letters + strings.ToUpper(letters)
	nums := "0123456789"
	chars := "()#$￥%+-*/="

	numCnt := 0
	letterCnt := 0
	otherCnt := 0
	charsCnt := 0

	for _, i := range value {
		switch {
		case strings.ContainsRune(letters, i) == true:
			letterCnt += 1
		case strings.ContainsRune(nums, i) == true:
			numCnt += 1
		case strings.ContainsRune(chars, i) == true:
			charsCnt += 1
		default:
			otherCnt += 1
		}
	}
	return float64((numCnt + letterCnt + charsCnt + otherCnt*2) * 4)
}
