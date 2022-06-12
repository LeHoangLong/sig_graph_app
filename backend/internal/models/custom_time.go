package models

import (
	"database/sql/driver"
	"fmt"
	"time"
)

type CustomTime time.Time

func (t CustomTime) MarshalJSON() ([]byte, error) {
	//do your serializing here

	stamp := fmt.Sprintf("\"%s\"", time.Time(t).Format(time.RFC3339))
	return []byte(stamp), nil
}

func (t CustomTime) String() string {
	return time.Time(t).Format(time.RFC3339)
}

func (t *CustomTime) UnmarshalJSON(iData []byte) error {
	strData := string(iData)
	if len(strData) <= 2 {
		return fmt.Errorf("invalid data %s", strData)
	}
	parsed, err := time.Parse(time.RFC3339, strData[1:len(strData)-1])
	if err != nil {
		return err
	}

	*t = CustomTime(parsed)
	return nil
}

func NewCustomTimeFromString(iData string) (CustomTime, error) {
	parsed, err := time.Parse(time.RFC3339, iData)
	if err != nil {
		return CustomTime{}, err
	}
	return CustomTime(parsed), nil
}

func (c CustomTime) Value() (driver.Value, error) {
	return time.Time(c), nil
}

func (c *CustomTime) Scan(value interface{}) error {
	if value == nil {
		return fmt.Errorf("value is nil")
	}

	if time, ok := value.(time.Time); !ok {
		return fmt.Errorf("could not parse time")
	} else {
		*c = CustomTime(time)
		return nil
	}
}
