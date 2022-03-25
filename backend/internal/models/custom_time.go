package models

import (
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
