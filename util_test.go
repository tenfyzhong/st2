package st2

import (
	"reflect"
	"testing"
)

func TestCamel(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args func(t *testing.T) args

		want1 string
	}{
		{
			name: "empty",
			args: func(t *testing.T) args {
				return args{
					s: "",
				}
			},
			want1: "",
		},
		{
			name: "_",
			args: func(t *testing.T) args {
				return args{
					s: "_",
				}
			},
			want1: "",
		},
		{
			name: "a_b",
			args: func(t *testing.T) args {
				return args{
					s: "a_b",
				}
			},
			want1: "AB",
		},
		{
			name: "_b",
			args: func(t *testing.T) args {
				return args{
					s: "_b",
				}
			},
			want1: "B",
		},
		{
			name: "a_",
			args: func(t *testing.T) args {
				return args{
					s: "a_",
				}
			},
			want1: "A",
		},
		{
			name: "hello_world",
			args: func(t *testing.T) args {
				return args{
					s: "hello_world",
				}
			},
			want1: "HelloWorld",
		},
		{
			name: "hEllo_World",
			args: func(t *testing.T) args {
				return args{
					s: "hEllo_World",
				}
			},
			want1: "HElloWorld",
		},
		{
			name: "hElloWorld",
			args: func(t *testing.T) args {
				return args{
					s: "hElloWorld",
				}
			},
			want1: "HElloWorld",
		},
		{
			name: "helloWorld",
			args: func(t *testing.T) args {
				return args{
					s: "helloWorld",
				}
			},
			want1: "HelloWorld",
		},
		{
			name: "HelloWorld",
			args: func(t *testing.T) args {
				return args{
					s: "HelloWorld",
				}
			},
			want1: "HelloWorld",
		},
		{
			name: "id_name",
			args: func(t *testing.T) args {
				return args{
					s: "id_name",
				}
			},
			want1: "IdName",
		},
		{
			name: "id_url",
			args: func(t *testing.T) args {
				return args{
					s: "id_url",
				}
			},
			want1: "IdUrl",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tArgs := tt.args(t)

			got1 := Camel(tArgs.s)

			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("Camel got1 = %v, want1: %v", got1, tt.want1)
			}
		})
	}
}

func TestSnake(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args func(t *testing.T) args

		want1 string
	}{
		{
			name: "A",
			args: func(t *testing.T) args {
				return args{
					s: "A",
				}
			},
			want1: "a",
		},
		{
			name: "ID",
			args: func(t *testing.T) args {
				return args{
					s: "ID",
				}
			},
			want1: "id",
		},
		{
			name: "IDName",
			args: func(t *testing.T) args {
				return args{
					s: "IDName",
				}
			},
			want1: "id_name",
		},
		{
			name: "IDURL",
			args: func(t *testing.T) args {
				return args{
					s: "IDURL",
				}
			},
			want1: "idurl",
		},
		{
			name: "id_name",
			args: func(t *testing.T) args {
				return args{
					s: "id_name",
				}
			},
			want1: "id_name",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tArgs := tt.args(t)

			got1 := Snake(tArgs.s)

			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("Snake got1 = %v, want1: %v", got1, tt.want1)
			}
		})
	}
}
