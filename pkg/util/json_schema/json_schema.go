// Copyright (C) 2025 wangyusong
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program. If not, see <https://www.gnu.org/licenses/>.

package jsonschema

import (
	"maps"
	"reflect"
	"strings"
	"time"

	"github.com/pkg/errors"
)

// ForType generates a JSON Schema for the given reflect.Type.
// It supports struct fields with json tags and desc tags for metadata.
func ForType(t reflect.Type) (map[string]any, error) {
	definitions := make(map[string]any)
	schema, err := forTypeInternal(t, "", make(map[reflect.Type]string), definitions)
	if err != nil {
		return nil, err
	}

	if len(definitions) == 0 {
		return schema, nil
	}

	result := map[string]any{
		"$schema":     "http://json-schema.org/draft-07/schema#",
		"definitions": definitions,
	}
	maps.Copy(result, schema)

	return result, nil
}

func forTypeInternal(
	t reflect.Type,
	fieldName string,
	visited map[reflect.Type]string,
	definitions map[string]any,
) (map[string]any, error) {
	if t == nil {
		return nil, errors.New("type cannot be nil")
	}

	// Dereference pointer types
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	// Handle previously visited types
	if refName, ok := visited[t]; ok {
		return map[string]any{"$ref": "#/definitions/" + refName}, nil
	}

	switch t.Kind() {
	case reflect.Struct:
		return handleStructType(t, fieldName, visited, definitions)

	case reflect.Slice, reflect.Array:
		return handleArrayType(t, visited, definitions)

	case reflect.Map:
		return handleMapType(t, visited, definitions)

	default:
		return handlePrimitiveType(t)
	}
}

func handleStructType(
	t reflect.Type,
	fieldName string,
	visited map[reflect.Type]string,
	definitions map[string]any,
) (map[string]any, error) {
	// Handle special types.
	if t == reflect.TypeOf(time.Time{}) {
		return map[string]any{
			"type":   "string",
			"format": "date-time",
		}, nil
	}

	if t == reflect.TypeOf(time.Duration(0)) {
		return map[string]any{
			"type":    "string",
			"format":  "duration",
			"pattern": "^([0-9]+(s|m|h))+$",
		}, nil
	}

	// Generate type name.
	typeName := t.Name()
	if typeName == "" {
		typeName = "Anonymous" + fieldName
	}
	visited[t] = typeName

	// Process schema.
	schema := map[string]any{"type": "object"}

	properties, err := handleStructFields(t, visited, definitions)
	if err != nil {
		return nil, errors.Wrap(err, "handle struct fields")
	}
	if len(properties) > 0 {
		schema["properties"] = properties
	}

	definitions[typeName] = schema

	return map[string]any{"$ref": "#/definitions/" + typeName}, nil
}

func handleStructFields(
	t reflect.Type,
	visited map[reflect.Type]string,
	definitions map[string]any,
) (properties map[string]any, err error) {
	properties = make(map[string]any, t.NumField())

	for i := range t.NumField() {
		field := t.Field(i)
		if !field.IsExported() {
			continue
		}

		propName := getPropertyName(field)
		if propName == "" {
			continue
		}

		if field.Anonymous {
			if err := handleEmbeddedStruct(field, visited, definitions, properties); err != nil {
				return nil, err
			}

			continue
		}

		fieldSchema, err := forTypeInternal(field.Type, field.Name, visited, definitions)
		if err != nil {
			return nil, errors.Wrapf(err, "generating schema for field %s", field.Name)
		}

		if desc := field.Tag.Get("desc"); desc != "" {
			fieldSchema["description"] = desc
		}

		properties[propName] = fieldSchema
	}

	return properties, nil
}

func handleArrayType(
	t reflect.Type,
	visited map[reflect.Type]string,
	definitions map[string]any,
) (map[string]any, error) {
	itemSchema, err := forTypeInternal(t.Elem(), "", visited, definitions)
	if err != nil {
		return nil, errors.Wrap(err, "generating array item schema")
	}

	return map[string]any{
		"type":  "array",
		"items": itemSchema,
	}, nil
}

func handleMapType(
	t reflect.Type,
	visited map[reflect.Type]string,
	definitions map[string]any,
) (map[string]any, error) {
	if t.Key().Kind() != reflect.String {
		return nil, errors.Errorf("unsupported map key type: %s (must be string)", t.Key().Kind())
	}

	valueSchema, err := forTypeInternal(t.Elem(), "", visited, definitions)
	if err != nil {
		return nil, errors.Wrap(err, "generating map value schema")
	}

	return map[string]any{
		"type":                 "object",
		"additionalProperties": valueSchema,
	}, nil
}

func handlePrimitiveType(t reflect.Type) (map[string]any, error) {
	schema := make(map[string]any)

	switch t.Kind() {
	case reflect.String:
		schema["type"] = "string"

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if t == reflect.TypeOf(time.Duration(0)) {
			schema["type"] = "string"
			schema["format"] = "duration"
			schema["pattern"] = "^([0-9]+(s|m|h))+$"
		} else {
			schema["type"] = "integer"
		}

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		schema["type"] = "integer"
		schema["minimum"] = 0

	case reflect.Float32, reflect.Float64:
		schema["type"] = "number"

	case reflect.Bool:
		schema["type"] = "boolean"

	default:
		return nil, errors.Errorf("unsupported type: %s", t.Kind())
	}

	return schema, nil
}

func getPropertyName(field reflect.StructField) string {
	jsonTag := field.Tag.Get("json")
	if jsonTag == "-" {
		return ""
	}

	if jsonTag != "" {
		parts := strings.Split(jsonTag, ",")

		return parts[0]
	}

	return field.Name
}

func handleEmbeddedStruct(
	field reflect.StructField,
	visited map[reflect.Type]string,
	definitions map[string]any,
	properties map[string]any,
) error {
	embeddedSchema, err := forTypeInternal(field.Type, "", visited, definitions)
	if err != nil {
		return errors.Wrapf(err, "generating schema for embedded field %s", field.Name)
	}

	if embeddedType, ok := embeddedSchema["$ref"]; ok {
		refType := embeddedType.(string)
		key := strings.TrimPrefix(refType, "#/definitions/")
		if def, ok := definitions[key]; ok {
			if embeddedProps, ok := def.(map[string]any)["properties"].(map[string]any); ok {
				maps.Copy(properties, embeddedProps)
			}

			delete(definitions, key)
		}
	}

	return nil
}
