package helpers

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"sync"

	utils "github.com/flabio/safe_constants"
	"github.com/safe_msvc_city/usecase/dto"
)

var fieldCache sync.Map // Caché para almacenar los índices de los campos
func MapToStruct(dataDto *dto.CityDTO, dataMap map[string]interface{}) {
	city := dto.CityDTO{
		Name:   dataMap[utils.NAME].(string),
		Active: dataMap[utils.ACTIVE].(bool),
	}
	*dataDto = city
}

func ValidateFieldCity(value map[string]interface{}) string {
	var msg string = utils.EMPTY
	if value[utils.NAME] == nil {
		msg = utils.NAME_FIELD_IS_REQUIRED
	}
	if value[utils.ACTIVE] == nil {
		msg = utils.ACTIVE_FIELD_IS_REQUIRED
	}
	return msg
}

func ValidateRequiredCity(field dto.CityDTO) string {
	var msg string = utils.EMPTY
	if field.Name == utils.EMPTY {
		msg = utils.NAME_IS_REQUIRED
	}
	return msg
}

// MapToStruct es una función genérica que mapea un map[string]interface{} a cualquier estructura proporcionada
func MapToStructState(dataMap map[string]interface{}, structPointer interface{}) error {
	// Verificamos que structPointer sea un puntero a una estructura
	value := reflect.ValueOf(structPointer)
	if value.Kind() != reflect.Ptr || value.Elem().Kind() != reflect.Struct {
		return errors.New("structPointer debe ser un puntero a una estructura")
	}

	structValue := value.Elem()
	structType := structValue.Type()

	for key, mapValue := range dataMap {
		// Intentamos obtener el índice del campo desde el caché
		fieldIdx, found := getFieldIndex(structType, key)
		if !found {
			// Si no encontramos el campo, continuamos con la siguiente clave del map
			continue
		}

		fieldValue := structValue.Field(fieldIdx)
		if !fieldValue.CanSet() {
			// Si el campo no se puede asignar (no exportado), continuamos
			continue
		}

		mapValueReflect := reflect.ValueOf(mapValue)

		// Verificamos si el valor es asignable al tipo del campo
		if mapValueReflect.Type().AssignableTo(fieldValue.Type()) {
			fieldValue.Set(mapValueReflect)
		} else {
			// Intentamos convertir el valor si es posible
			convertedValue, err := tryConvert(mapValueReflect, fieldValue.Type())
			if err == nil {
				fieldValue.Set(convertedValue)
			} else {
				return fmt.Errorf("error asignando el campo '%s': %v", key, err)
			}
		}
	}

	return nil
}

// getFieldIndex intenta obtener el índice del campo desde el caché, o lo calcula y almacena si no está en el caché
func getFieldIndex(structType reflect.Type, fieldName string) (int, bool) {
	cacheKey := structType.String() + "." + fieldName

	if index, ok := fieldCache.Load(cacheKey); ok {
		return index.(int), true
	}
	/*
		field, found := structType.FieldByName(fieldName)
		if found {
			fieldCache.Store(cacheKey, field.Index[0])
			return field.Index[0], true
		}*/
	// Buscar el campo correspondiente a la etiqueta json
	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)
		tag := field.Tag.Get("json")
		if tag == "" {
			tag = field.Name
		} else {
			tag = strings.Split(tag, ",")[0] // Obtenemos el nombre antes de cualquier opción (como omitir)
		}

		if tag == fieldName {
			fieldCache.Store(cacheKey, i)
			return i, true
		}
	}
	return -1, false
}

// tryConvert intenta convertir un valor a un tipo objetivo compatible
func tryConvert(value reflect.Value, targetType reflect.Type) (reflect.Value, error) {
	if value.Type().ConvertibleTo(targetType) {
		return value.Convert(targetType), nil
	}
	return reflect.Value{}, fmt.Errorf("no se puede convertir el valor '%v' al tipo '%v'", value, targetType)
}
