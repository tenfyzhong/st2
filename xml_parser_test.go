package st2

import (
	"reflect"
	"testing"
)

func TestXMLUnmarshalTagFormat_Unmarshal(t *testing.T) {
	type fields struct {
		ContentTagPrefix   string
		AttributeTagPrefix string
	}
	type args struct {
		data []byte
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantErr    bool
		expectJson string
	}{
		{
			name: "succ",
			fields: fields{
				ContentTagPrefix:   "",
				AttributeTagPrefix: ",",
			},
			args: args{
				data: []byte(`<?xml version="1.0" encoding="UTF-8"?>
  <osm version="0.6" generator="CGImap 0.0.2">
   <bounds minlat="54.0889580" minlon="12.2487570" maxlat="54.0913900" maxlon="12.2524800"/>
   <foo>bar</foo>
  </osm>`),
			},
			wantErr:    false,
			expectJson: `{"osm":{",version":"0.6",",generator":"CGImap 0.0.2","bounds":{",minlon":"12.2487570",",maxlat":"54.0913900",",maxlon":"12.2524800",",minlat":"54.0889580"},"foo":"bar"}}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			x := XMLUnmarshalTagFormat{
				ContentTagPrefix:   tt.fields.ContentTagPrefix,
				AttributeTagPrefix: tt.fields.AttributeTagPrefix,
			}
			var v any
			if err := x.Unmarshal(tt.args.data, &v); (err != nil) != tt.wantErr {
				t.Errorf("XMLUnmarshalTagFormat.Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
			}

			var expectV any
			if err := jsonapi.UnmarshalFromString(tt.expectJson, &expectV); err != nil {
				t.Errorf("UnmarshalFromString error = %v", err)
			}

			if !reflect.DeepEqual(v, expectV) {
				t.Errorf("%+v should equal to %+v", v, expectV)
			}
		})
	}
}

func TestXMLUnmarshalTagFormat_TagFormat(t *testing.T) {
	type fields struct {
		ContentTagPrefix   string
		AttributeTagPrefix string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "test",
			want: `xml:"%s"`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			x := XMLUnmarshalTagFormat{
				ContentTagPrefix:   tt.fields.ContentTagPrefix,
				AttributeTagPrefix: tt.fields.AttributeTagPrefix,
			}
			if got := x.TagFormat(); got != tt.want {
				t.Errorf("XMLUnmarshalTagFormat.TagFormat() = %v, want %v", got, tt.want)
			}
		})
	}
}
