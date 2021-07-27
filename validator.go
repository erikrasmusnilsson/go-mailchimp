package mailchimp

import "reflect"

func validate(value interface{}) ([]string, bool) {
	v := reflect.Indirect(reflect.ValueOf(value))
	invalidParameters := make([]string, 0)
	isValid := true
	for currentField := 0; currentField < v.NumField(); currentField++ {
		fieldValue := v.Field(currentField)
		field := v.Type().Field(currentField)
		if isRequired(&field) && !isSet(&fieldValue, &field) {
			invalidParameters = append(invalidParameters, field.Name)
			isValid = false
		}
	}
	return invalidParameters, isValid
}

func isSet(fieldValue *reflect.Value, field *reflect.StructField) bool {
	switch fieldValue.Interface().(type) {
	case string:
		return isStringSet(fieldValue)
	}
	return true
}

func isStringSet(fieldValue *reflect.Value) bool {
	return fieldValue.String() != ""
}

func isRequired(field *reflect.StructField) bool {
	return validatorTag(field) == "required"
}

func validatorTag(field *reflect.StructField) string {
	return field.Tag.Get("mc_validator")
}
