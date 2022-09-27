package helper

import (
	"fmt"
	"github.com/trucktrace/pkg/logger"

	"github.com/trucktrace/internal/models"
)

func CheckChargesRate(charges models.Charges) error {

	if charges.TypeDeductions == "deduction" && charges.Rate > 0 {
		logger.ErrorLogger("CheckChargesRate", "Rate is more than 0").Error("Error - " + fmt.Sprintf("type Charge is deduction but rate is more than 0 - %v", charges.Rate))
		return fmt.Errorf("type Charge is deduction but rate is more than 0 - %v", charges.Rate)
	}

	if charges.TypeDeductions == "earning" && charges.Rate < 0 {
		logger.ErrorLogger("CheckChargesRate", "Rate is less than 0").Error("Error - " + fmt.Sprintf("type Charge is earning but Rate is less than 0 - %v", charges.Rate))
		return fmt.Errorf("type Charge is earning but rate is less than 0 - %v", charges.Rate)
	}

	return nil
}
