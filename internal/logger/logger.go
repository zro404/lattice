package logger

import (
	"io"
	"log"
	"os"
)

var l *log.Logger

func Init() error {
	file, err := os.OpenFile(
		"lattice.log",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0666,
	)
	if err != nil {
		panic(err)
	}

	multi := io.MultiWriter(os.Stdout, file)

	l = log.New(multi, "", log.Ldate|log.Ltime)

	return nil
}

func Writer() io.Writer {
	return l.Writer()
}

func Println(v ...any) {
	l.Println(v...)
}

func Printf(format string, v ...any) {
	l.Printf(format, v...)
}

func Fatal(v ...any) {
	l.Fatal(v...)
}

func Fatalf(format string, v ...any) {
	l.Fatalf(format, v...)
}
