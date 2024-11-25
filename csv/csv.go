package csv

import (
	"bufio"
	"errors"
	"io"
	"strings"
)

// Ошибки, которые могут быть возвращены парсером
var (
	ErrQuote      = errors.New("excess or missing \" in quoted-field")
	ErrFieldCount = errors.New("wrong number of fields")
)

// Интерфейс для парсера CSV
type CSVParser interface {
	ReadLine(r io.Reader) (string, error) // Чтение строки из файла
	GetField(n int) (string, error)       // Получение поля по индексу
	GetNumberOfFields() int               // Количество полей в строке
}

// Структура парсера CSV
type StructCSVParser struct {
	currentLine string
	fields      []string
}

// Чтение одной строки из источника
func (p *StructCSVParser) ReadLine(r io.Reader) (string, error) {
	scanner := bufio.NewScanner(r)
	if scanner.Scan() {
		line := scanner.Text()
		// Проверяем на корректность кавычек
		if containsMismatchedQuotes(line) {
			return "", ErrQuote
		}

		// Убираем символы окончания строки
		line = strings.TrimSpace(line)
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

// Проверка на наличие несоответствующих кавычек
func containsMismatchedQuotes(line string) bool {
	quoteCount := 0
	for _, ch := range line {
		if ch == '"' {
			quoteCount++
		}
	}
	// Если количество кавычек нечетное, значит, есть ошибка
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

// Получение поля по индексу
func (p *StructCSVParser) GetField(n int) (string, error) {
	if n < 0 || n >= len(p.fields) {
		return "", ErrFieldCount
	}
	return p.fields[n], nil
}

// Получение количества полей в строке
func (p *StructCSVParser) GetNumberOfFields() int {
	return len(p.fields)
}
