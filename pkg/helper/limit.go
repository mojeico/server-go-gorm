package helper

func CheckLimit(offSet, limit string) (string, string) {

	if offSet == "" {
		offSet = "0"
	}

	if limit == "" {
		limit = "20"
	}

	return offSet, limit

}
