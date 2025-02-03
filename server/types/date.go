package types

import (
	"database/sql/driver"
	"fmt"
	"strings"
	"time"
)

// Date is a custom type that represents a date without time information.
// It implements json.Unmarshaler to handle date strings in the format "YYYY-MM-DD".
type Date time.Time

// Time converts Date to time.Time
func (d Date) Time() time.Time {
	return time.Time(d)
}

// UnmarshalJSON implements the json.Unmarshaler interface.
// It parses a JSON string containing a date in the format "YYYY-MM-DD".
// The method expects the input to be a quoted string like "2024-02-28".
// It removes the quotes and converts the string to a time.Time value
// which is then assigned to the Date.
func (d *Date) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"") // trim the quotes from the byte value
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return err
	}
	*d = Date(t)
	return nil
}

// Value implements the driver.Valuer interface
func (d Date) Value() (driver.Value, error) {
	return time.Time(d).Format("2006-01-02"), nil
}

// Scan implements the sql.Scanner interface
func (d *Date) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	t, ok := value.(time.Time)
	if !ok {
		return fmt.Errorf("expected time.Time got %T", value)
	}
	*d = Date(t)
	return nil
}
