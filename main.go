package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

// Ошибки для парсера CSV
var (
	ErrQuote      = errors.New("excess or missing \" in quoted-field")
	ErrFieldCount = errors.New("wrong number of fields")
)

// Интерфейс для парсера CSV
type CSVParser interface {
	ReadLine(r io.Reader) (string, error)
	GetField(n int) (string, error)
	GetNumberOfFields() int
}

// Структура, которая реализует интерфейс CSVParser
type CSVParserImpl struct {
	currentLine string
	fields      []string
}

// Метод для чтения одной строки
func (p *CSVParserImpl) ReadLine(r io.Reader) (string, error) {
	scanner := bufio.NewScanner(r)
	if scanner.Scan() {
		line := scanner.Text()

		// Проверка на корректность кавычек
		if containsMismatchedQuotes(line) {
			return "", ErrQuote
		}

		// Убираем символы конца строки
		line = strings.TrimSpace(line)

		// Сохраняем текущую строку и разбиваем ее на поля
		p.currentLine = line
		p.fields = parseFields(line)

		return line, nil
	}

	// Ошибка при сканировании
	if scanner.Err() != nil {
		return "", scanner.Err()
	}

	return "", nil // EOF
}

// Проверка на некорректные кавычки в строке
func containsMismatchedQuotes(line string) bool {
	quoteCount := 0
	for _, ch := range line {
		if ch == '"' {
			quoteCount++
		}
	}
	return quoteCount%2 != 0
}

// Разделение строки на поля
func parseFields(line string) []string {
	var fields []string
	var field strings.Builder
	inQuotes := false

	// Проходим по строке символ за символом
	for _, ch := range line {
		if ch == '"' {
			inQuotes = !inQuotes // Меняем состояние (внутри или вне кавычек)
		} else if ch == ',' && !inQuotes {
			fields = append(fields, field.String())
			field.Reset()
		} else {
			field.WriteRune(ch)
		}
	}

	// Добавляем последнее поле, если оно не пустое
	if field.Len() > 0 {
		fields = append(fields, field.String())
	}

	return fields
}

// Метод для получения поля по индексу
func (p *CSVParserImpl) GetField(n int) (string, error) {
	if n < 0 || n >= len(p.fields) {
		return "", ErrFieldCount
	}
	return p.fields[n], nil
}

// Метод для получения количества полей в строке
func (p *CSVParserImpl) GetNumberOfFields() int {
	return len(p.fields)
}

func main() {
	// Открываем файл
	file, err := os.Open("example.csv")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	// Создаем экземпляр парсера CSV
	var csvparser CSVParser = &CSVParserImpl{}
	var count int
	// Чтение строк из файла
	for count == 0 {
		line, err := csvparser.ReadLine(file)
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println("Error reading line:", err)
			return
		}

		// Выводим строку
		fmt.Println("Line:", line)

		// Пример: получение первого поля
		field, err := csvparser.GetField(0)
		if err != nil {
			fmt.Println("Error getting field:", err)
			return
		}
		fmt.Println("First field:", field)

		// Пример: получение количества полей
		fmt.Println("Number of fields:", csvparser.GetNumberOfFields())
		count++
	}
}
