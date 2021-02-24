package ext_test

import (
	"reflect"
	"testing"

	"github.com/chatto-extensions/weatherapi/internal/ext"
	"github.com/jaimeteb/chatto/extension"
)

func TestWeather(t *testing.T) {
	type args struct {
		req *extension.Request
	}
	tests := []struct {
		name    string
		args    args
		wantRes *extension.Response
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := ext.Weather(tt.args.req); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Weather() = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}
