package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type InsertData struct {
	Feature  int32
	TagIds   []int32
	Content  JSONMap
	IsActive bool
}

type JSONMap map[string]interface{}

// Value - реализация интерфейса driver.Valuer
func (j *JSONMap) Value() (driver.Value, error) {
	return json.Marshal(j)
}

// Scan - реализация интерфейса sql.Scanner
func (j *JSONMap) Scan(value interface{}) error {
	// Проверяем тип значения
	if value == nil {
		*j = nil
		return nil
	}

	// Преобразуем значение в []byte
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("ошибка преобразования типа %T в []byte", value)
	}

	// Декодируем JSON в map[string]interface{}
	var data map[string]interface{}
	if err := json.Unmarshal(bytes, &data); err != nil {
		return err
	}

	// Устанавливаем значение в j
	*j = data
	return nil
}
