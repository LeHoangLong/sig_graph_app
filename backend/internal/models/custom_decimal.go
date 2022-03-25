package models

import (
	"database/sql/driver"
	"fmt"

	"github.com/shopspring/decimal"
)

type CustomDecimal struct {
	decimal.Decimal
}

func (c CustomDecimal) Value() (driver.Value, error) {
	return c.String(), nil
}

func (c *CustomDecimal) Scan(value interface{}) error {
	if value == nil {
		return fmt.Errorf("value is nil")
	}

	if decimalString, ok := value.([]byte); !ok {
		return fmt.Errorf("could not parse decimal")
	} else {
		newDec, err := decimal.NewFromString(string(decimalString))
		if err != nil {
			return err
		}

		c.Decimal = newDec
		return nil
	}
}

func NewCustomDecimalFromString(
	iData string,
) (CustomDecimal, error) {
	dec, err := decimal.NewFromString(iData)
	if err != nil {
		return CustomDecimal{}, err
	}

	return CustomDecimal{
		Decimal: dec,
	}, nil
}

func (d CustomDecimal) MarshalJSON() ([]byte, error) {
	//do your serializing here
	str := fmt.Sprintf("\"%s\"", d.String())
	return []byte(str), nil
}

func (d *CustomDecimal) UnmarshalJSON(iData []byte) error {
	strData := string(iData)
	if len(strData) <= 2 {
		return fmt.Errorf("invalid data %s", strData)
	}
	parsed, err := decimal.NewFromString(strData[1 : len(strData)-1])
	if err != nil {
		return err
	}

	*d = CustomDecimal{
		Decimal: parsed,
	}
	return nil
}
