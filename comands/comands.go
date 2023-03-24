package comands

import (
	"errors"
	"os"
	"time"
)

func CompareTimes(first_arg, second_arg time.Time, operator string) (bool, error) {
	switch operator {
	case "<":
		return first_arg.After(second_arg), nil
	}
	return false, errors.New("unknown parameter")
}

func TimeFromString(in string) (time.Time, error) {
	return time.Parse("", in)
}

func CreateFile(fileName string) error {
	f, err := os.Create(fileName)
	if err != nil {
		return errors.New("failed to create a file")
	}
	defer f.Close()
	return err
}

func RenameFile(fileName, rename string) error {
	return os.Rename(fileName, rename)
}

func GetCreationTime(fileName string) (string, error) {
	file, err := os.Stat(fileName)
	if err != nil {
		return "", errors.New("failed to open a file")
	}
	return file.ModTime().String(), nil
}

func AppendFile(fileName, text string) error {
	err := os.WriteFile(fileName, []byte(text), 0)
	if err != nil {
		return errors.New("cannot write a file")
	}
	return nil
}
