package converters

import (
	"database/sql"
	"fmt"
)

func NullStringToString(nullString sql.NullString) string {
	if !nullString.Valid {
		return ""
	}
	return nullString.String
}

func StringToNullString(str string) sql.NullString {
	var sqlString sql.NullString
	if str != "" {
		sqlString.Valid = true
		sqlString.String = str
	}
	return sqlString
}

func GenderStringToUint8(gender string) (uint8, error) {
	switch gender {
	case "female":
		return 0, nil
	case "male":
		return 1, nil
	default:
		return 2, fmt.Errorf("unsupported gender type: %s", gender)
	}
}

func GenderUint8ToString(gender uint8) string {
	switch gender {
	case 0:
		return "female"
	case 1:
		return "male"
	default:
		return "unknown"
	}
}

func StringToStringPointer(str string) *string {
	if str == "" {
		return nil
	}
	return &str
}
