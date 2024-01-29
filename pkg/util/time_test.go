package util_test

import (
	"boilerplate-service/pkg/util"
	"errors"
	"testing"
	"time"
)

func mockLocationLoader(name string) (*time.Location, error) {
	return nil, errors.New("mock location error")
}

func TestGetTimeWithLoader(t *testing.T) {
	tests := []struct {
		name       string
		loaderFunc util.LocationLoader
		wantErr    bool
	}{
		{
			name:       "Successful Timezone Retrieval",
			loaderFunc: time.LoadLocation,
			wantErr:    false,
		},
		{
			name:       "Error Handling",
			loaderFunc: mockLocationLoader,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := util.GetJakartaTimeWithLoader(tt.loaderFunc)
			if (err != nil) != tt.wantErr {
				t.Errorf("%s: GetJakartaTimeWithLoader() error = %v, wantErr %v", tt.name, err, tt.wantErr)
			}
		})
	}
}

func TestGetJakartaTime(t *testing.T) {
	jakartaTime, err := util.GetJakartaTime()
	if err != nil {
		t.Errorf("GetJakartaTime() returned an error: %v", err)
	}

	_, offset := jakartaTime.Zone()
	if offset != 7*60*60 {
		t.Errorf("GetJakartaTime() did not return time in 'Asia/Jakarta' timezone: got offset %d", offset)
	}
}
