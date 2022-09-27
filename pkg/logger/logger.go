package logger

/*
type Logger struct{}

func (logger *Logger) WarnLogger(message, path, body string) {
	logg := fmt.Sprintf("[ WARN! ] Path: %s Message: %s Body: %s, %s ", path, message, body, time.Now())
	err := writeLogInFile(warnLogFile, logg)

	if err != nil {
		errMessage := fmt.Sprintf("can't write warn in file - %s", err.Error())
		log.Println(errMessage)
	}

	log.Println(logg)
}

func (logger *Logger) InfoLogger(message, path, body string) {
	logg := fmt.Sprintf("[ INFO! ] Path: %s Message: %s Body: %s, %s ", path, message, body, time.Now())
	err := writeLogInFile(infoLogFile, logg)

	if err != nil {
		errMessage := fmt.Sprintf("can't write info in file - %s", err.Error())
		log.Println(errMessage)
	}

	log.Println(logg)
}

func (logger *Logger) ErrorLogger(message, path, function, body string) {
	logg := fmt.Sprintf("[ ERROR ] Path: %s Function: %s Message: %s Body: %s, %s ", path, function, message, body, time.Now())
	err := writeLogInFile(errLogFile, logg)

	if err != nil {
		errMessage := fmt.Sprintf("can't write error in file - %s", err.Error())
		log.Println(errMessage)
	}

	log.Println(fmt.Errorf(logg))
}

func (logger *Logger) PanicLogger(message, path, body string) {
	logg := fmt.Sprintf("[ PANIC ] Path: %s Message: %s Body: %s, %s ", path, message, body, time.Now())
	err := writeLogInFile(panicLogFile, logg)

	if err != nil {
		errMessage := fmt.Sprintf("can't write patic in file - %s", err.Error())
		log.Println(errMessage)
	}

	log.Fatalln(logg)
}

func (logger *Logger) FatalLogger(message, path, body string) {
	logg := fmt.Sprintf("[ FATAL ] Path: %s Message: %s Body: %s, %s ", path, message, body, time.Now())
	err := writeLogInFile(fatalLogFile, logg)

	if err != nil {
		errMessage := fmt.Sprintf("can't write fatal in file - %s", err.Error())
		log.Println(errMessage)
	}

	log.Fatalln(logg)
}

func GetLoggerPackage() *Logger {
	return &Logger{}
}
*/
