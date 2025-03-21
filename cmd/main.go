package main

import (
	"context"
	"fmt"
	"log"

	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

const (
	spreadsheetID = "13UcQUHFmCa_oK3IXGth5BBtT8lpfUYQF7SvZM3Yp6iM"
	rangeData     = "Diary.xlsx"
)

func main() {
	// Загружаем учетные данные сервисного аккаунта
	ctx := context.Background()
	srv, err := sheets.NewService(ctx, option.WithCredentialsFile("tokens.json"))
	if err != nil {
		log.Fatalf("Не удалось создать клиент Sheets: %v", err)
	}

	// Новые данные для добавления (например, в строках A1, B1, C1)
	values := []interface{}{"Новая запись 1", 100, "Комментарий"}

	// Подготовка данных в формате, который можно записать в таблицу
	var vr sheets.ValueRange
	vr.Values = append(vr.Values, values)

	// Диапазон (пустая строка указывает на добавление в конец таблицы)
	// Убедись, что имя листа правильно

	// Добавление данных в таблицу
	_, err = srv.Spreadsheets.Values.Append(spreadsheetID, rangeData, &vr).ValueInputOption("RAW").Do()
	if err != nil {
		log.Fatalf("Ошибка добавления данных: %v", err)
	}

	fmt.Println("Данные успешно добавлены!")
}
