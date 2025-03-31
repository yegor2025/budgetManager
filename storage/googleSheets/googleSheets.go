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

	_, err := s.Client.Spreadsheets.Values.Append(spreadsheetID, sheetName, &vr).ValueInputOption("RAW").Do()
	if err != nil {
		log.Printf("Storage.Append: %v", err)
	}

	return err
}

func (s Storage) InsertBeforeLast(sheetName, message string, isCalculation bool) error {
	// Получаем количество строк в листе
	rangeData := fmt.Sprintf("%s!A:A", sheetName) // Берем первую колонку, чтобы определить длину таблицы
	resp, err := s.Client.Spreadsheets.Values.Get(spreadsheetID, rangeData).Do()
	if err != nil {
		log.Printf("Ошибка при получении количества строк: %v", err)
		return err
	}

	rowCount := len(resp.Values) // Количество заполненных строк
	if rowCount < 2 {
		log.Println("Недостаточно строк для вставки")
		return fmt.Errorf("недостаточно строк для вставки")
	}

	insertRow := rowCount // Последняя строка
	insertRow--           // Берем предпоследнюю строку

	// Разбиваем сообщение на части
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

	// Определяем диапазон вставки
	insertRange := fmt.Sprintf("%s!A%d", sheetName, insertRow)

	_, err = s.Client.Spreadsheets.Values.Update(spreadsheetID, insertRange, &vr).ValueInputOption("RAW").Do()
	if err != nil {
		log.Printf("Ошибка при вставке данных: %v", err)
	}

	return err
}

func generateDate() string {
	now := time.Now()

	formattedDate := fmt.Sprintf("%d %s", now.Day(), months[now.Month()])
	return formattedDate
}
