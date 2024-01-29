package util

import (
	"time"
)

type LocationLoader func(name string) (*time.Location, error)

func GetJakartaTimeWithLoader(loader LocationLoader) (time.Time, error) {
	t, err := loader("Asia/Jakarta")
	if err != nil {
		return time.Time{}, err
	}

	return time.Now().In(t), nil
}

func GetJakartaTime() (time.Time, error) {
	return GetJakartaTimeWithLoader(time.LoadLocation)
}
