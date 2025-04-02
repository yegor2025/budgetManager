package googleSheets

import (
	"context"
	"fmt"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
	"log"
	"strings"
	"time"
)

const (
	spreadsheetID = "13UcQUHFmCa_oK3IXGth5BBtT8lpfUYQF7SvZM3Yp6iM"
	sheetName     = "Diary.xlsx"
)

var months = []string{
	"", "января", "февраля", "марта", "апреля", "мая", "июня",
	"июля", "августа", "сентября", "октября", "ноября", "декабря",
}

type Storage struct {
	Client *sheets.Service
}

func New(ctx context.Context, pathToken string) *Storage {
	// Загружаем учетные данные сервисного аккаунта
	srv, err := sheets.NewService(ctx, option.WithCredentialsFile(pathToken))
	if err != nil {
		log.Fatalf("Не удалось создать клиент Sheets: %v", err)
	}

	return &Storage{
		Client: srv,
	}
}

func (s Storage) Append(message string) error {
	parts := strings.Split(message, ",")
	for i := range parts {
		parts[i] = strings.TrimSpace(parts[i])
	}

	var values []interface{}
	values = append(values, generateDate())

	for _, part := range parts {
		if part == "0" {
			values = append(values, nil)
		} else {
			values = append(values, part)
		}
	}

	var vr sheets.ValueRange
	vr.Values = append(vr.Values, values)

	_, err := s.Client.Spreadsheets.Values.Append(spreadsheetID, sheetName+"!A:A", &vr).ValueInputOption("RAW").InsertDataOption("INSERT_ROWS").Do()
	if err != nil {
		log.Printf("Storage.Append: %v", err)
	}

	return err
}

func (s Storage) InsertBeforeLast(message string, isCalculation bool) error {
	// Получаем количество строк
	rangeData := fmt.Sprintf("%s!A:A", sheetName)
	resp, err := s.Client.Spreadsheets.Values.Get(spreadsheetID, rangeData).Do()
	if err != nil {
		log.Printf("Ошибка при получении количества строк: %v", err)
		return err
	}

	rowCount := len(resp.Values) // Количество заполненных строк
	if rowCount < 2 {
		return fmt.Errorf("недостаточно строк для вставки")
	}

	insertRow := rowCount // Определяем предпоследнюю строку (нумерация с 1)

	// Сдвигаем последнюю строку вниз
	requests := []*sheets.Request{
		{
			InsertDimension: &sheets.InsertDimensionRequest{
				Range: &sheets.DimensionRange{
					SheetId:    getSheetID(sheetName, &s), // Функция получения SheetId
					Dimension:  "ROWS",
					StartIndex: int64(insertRow), // Google Sheets использует 0-based index
					EndIndex:   int64(insertRow + 1),
				},
				InheritFromBefore: true,
			},
		},
	}

	_, err = s.Client.Spreadsheets.BatchUpdate(spreadsheetID, &sheets.BatchUpdateSpreadsheetRequest{
		Requests: requests,
	}).Do()
	if err != nil {
		log.Printf("Ошибка при сдвиге строк: %v", err)
		return err
	}

	// Заполняем новую строку
	parts := strings.Split(message, ",")
	for i := range parts {
		parts[i] = strings.TrimSpace(parts[i])
	}

	var values []interface{}
	if !isCalculation {
		values = append(values, generateDate()) // Добавляем дату, если это не расчет
	}
	for _, part := range parts {
		if part == "0" {
			values = append(values, nil)
		} else {
			values = append(values, part)
		}
	}

	var vr sheets.ValueRange
	vr.Values = append(vr.Values, values)

	// Записываем в предпоследнюю строку
	insertRange := fmt.Sprintf("%s!A%d", sheetName, insertRow)

	_, err = s.Client.Spreadsheets.Values.Update(spreadsheetID, insertRange, &vr).ValueInputOption("RAW").Do()
	if err != nil {
		log.Printf("Ошибка при вставке данных: %v", err)
	}

	return err
}

func getSheetID(sheetName string, s *Storage) int64 {
	resp, err := s.Client.Spreadsheets.Get(spreadsheetID).Do()
	if err != nil {
		log.Fatalf("Ошибка при получении таблицы: %v", err)
	}

	for _, sheet := range resp.Sheets {
		if sheet.Properties.Title == sheetName {
			return sheet.Properties.SheetId
		}
	}
	log.Fatalf("Лист %s не найден", sheetName)
	return -1
}

func generateDate() string {
	now := time.Now()

	formattedDate := fmt.Sprintf("%d %s", now.Day(), months[now.Month()])
	return formattedDate
}
