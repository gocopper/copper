package phone

import (
	"github.com/asaskevich/govalidator"
	"github.com/tusharsoni/copper/clogger"
)

func AddPhoneNumberValidator(logger clogger.Logger) {
	logger.Info("Adding auth.PhoneNumber validator..")
	govalidator.TagMap["auth.PhoneNumber"] = govalidator.Validator(func(str string) bool {
		if len(str) != 12 {
			return false
		}

		return govalidator.IsNumeric(str[1:])
	})
}
