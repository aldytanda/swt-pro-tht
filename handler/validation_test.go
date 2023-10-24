package handler

import (
	"context"
	"testing"

	"github.com/aldytanda/swt-pro-tht/generated"
)

func TestValidationErrorResp_Error(t *testing.T) {
	type fields struct {
		Errors []ValidationError
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "contains 1 error",
			fields: fields{
				Errors: []ValidationError{
					{
						FieldName: "test_field",
						ErrRules:  []string{"test_err"},
					},
				},
			},
			want: "test_field: [test_err]",
		},
		{
			name: "contains 2 error",
			fields: fields{
				Errors: []ValidationError{
					{
						FieldName: "test_field",
						ErrRules:  []string{"test_err_1", "test_err_2"},
					},
				},
			},
			want: "test_field: [test_err_1 test_err_2]",
		},
		{
			name: "contains no error",
			fields: fields{
				Errors: []ValidationError{
					{
						FieldName: "test_field",
						ErrRules:  []string{},
					},
				},
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := ValidationErrorResp{
				Errors: tt.fields.Errors,
			}
			if got := e.Error(); got != tt.want {
				t.Errorf("ValidationErrorResp.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidationErrorResp_FromMap(t *testing.T) {
	type fields struct {
		Errors []ValidationError
	}
	type args struct {
		v map[string][]string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "contains 1 field and 1 error",
			fields: fields{
				Errors: []ValidationError{},
			},
			args: args{
				map[string][]string{
					"field": {"err"},
				},
			},
			want: true,
		},
		{
			name: "contains 2 field and 1 error",
			fields: fields{
				Errors: []ValidationError{},
			},
			args: args{
				map[string][]string{
					"field":  {"err"},
					"field2": {},
				},
			},
			want: true,
		},
		{
			name: "contains 1 field with no error",
			fields: fields{
				Errors: []ValidationError{},
			},
			args: args{
				map[string][]string{
					"field": {},
				},
			},
			want: false,
		},
		{
			name: "contains no field with no error",
			fields: fields{
				Errors: []ValidationError{},
			},
			args: args{
				map[string][]string{
					"field": {},
				},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &ValidationErrorResp{
				Errors: tt.fields.Errors,
			}
			if got := e.FromMap(tt.args.v); got != tt.want {
				t.Errorf("ValidationErrorResp.FromMap() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_validateRegister(t *testing.T) {
	type args struct {
		ctx context.Context
		v   generated.RegisterRequest
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "phone min length",
			args: args{
				ctx: context.Background(),
				v: generated.RegisterRequest{
					Name:     "Test Name",
					Password: "Abc123!",
					Phone:    "+621",
				},
			},
			wantErr: true,
		},
		{
			name: "phone max length",
			args: args{
				ctx: context.Background(),
				v: generated.RegisterRequest{
					Name:     "Test Name",
					Password: "Abc123!",
					Phone:    "+621234567890123123121",
				},
			},
			wantErr: true,
		},
		{
			name: "phone prefix +62",
			args: args{
				ctx: context.Background(),
				v: generated.RegisterRequest{
					Name:     "Test Name",
					Password: "Abc123!",
					Phone:    "+6112345678",
				},
			},
			wantErr: true,
		},
		{
			name: "name min length",
			args: args{
				ctx: context.Background(),
				v: generated.RegisterRequest{
					Name:     "Ab",
					Password: "Abc123!",
					Phone:    "+628131746432",
				},
			},
			wantErr: true,
		},
		{
			name: "name max length",
			args: args{
				ctx: context.Background(),
				v: generated.RegisterRequest{
					Name:     "AbcdefghijklmnopqrstuvwxyzAbcdefghijklmnopqrstuvwxyzAbcdefghijklmnopqrstuvwxyzAbcdefghijklmnopqrstuvwxyz",
					Password: "Abc123!",
					Phone:    "+628131746432",
				},
			},
			wantErr: true,
		},
		{
			name: "password min length",
			args: args{
				ctx: context.Background(),
				v: generated.RegisterRequest{
					Name:     "Test Name",
					Password: "Ab12!",
					Phone:    "+628131746432",
				},
			},
			wantErr: true,
		},
		{
			name: "password max length",
			args: args{
				ctx: context.Background(),
				v: generated.RegisterRequest{
					Name:     "Test Name",
					Password: "Abc123!AbcdefghijklmnopqrstuvwxyzAbcdefghijklmnopqrstuvwxyzAbcdefghijklmnopqrstuvwxyzAbcdefghijklmnopqrstuvwxyz",
					Phone:    "+628131746432",
				},
			},
			wantErr: true,
		},
		{
			name: "password uppercase",
			args: args{
				ctx: context.Background(),
				v: generated.RegisterRequest{
					Name:     "Test Name",
					Password: "abc123!",
					Phone:    "+628131746432",
				},
			},
			wantErr: true,
		},
		{
			name: "password numeric digit",
			args: args{
				ctx: context.Background(),
				v: generated.RegisterRequest{
					Name:     "Test Name",
					Password: "Abcdefg!",
					Phone:    "+628131746432",
				},
			},
			wantErr: true,
		},
		{
			name: "password special chars",
			args: args{
				ctx: context.Background(),
				v: generated.RegisterRequest{
					Name:     "Test Name",
					Password: "Abc1234",
					Phone:    "+628131746432",
				},
			},
			wantErr: true,
		},
		{
			name: "valid payload",
			args: args{
				ctx: context.Background(),
				v: generated.RegisterRequest{
					Name:     "Test Name",
					Password: "Abc123!",
					Phone:    "+628131746432",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateRegister(tt.args.ctx, tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf("validateRegister() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_validateUpdateUser(t *testing.T) {
	newUpdateProfileRequest := func(isName bool, name string, isPhone bool, phone string) generated.UpdateProfileRequest {
		p := generated.UpdateProfileRequest{}

		if isName {
			p.Name = &name
		}

		if isPhone {
			p.Phone = &phone
		}

		return generated.UpdateProfileRequest{}
	}

	type args struct {
		ctx context.Context
		v   generated.UpdateProfileRequest
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "phone min length",
			args: args{
				ctx: context.Background(),
				v:   newUpdateProfileRequest(true, "Test Name", true, "+621"),
			},
			wantErr: true,
		},
		{
			name: "phone max length",
			args: args{
				ctx: context.Background(),
				v:   newUpdateProfileRequest(true, "Test Name", true, "+621234567890123123121"),
			},
			wantErr: true,
		},
		{
			name: "phone prefix +62",
			args: args{
				ctx: context.Background(),
				v:   newUpdateProfileRequest(true, "Test Name", true, "+6112345678"),
			},
			wantErr: true,
		},
		{
			name: "name min length",
			args: args{
				ctx: context.Background(),
				v:   newUpdateProfileRequest(true, "Ab", true, "+628131746432"),
			},
			wantErr: true,
		},
		{
			name: "name max length",
			args: args{
				ctx: context.Background(),
				v:   newUpdateProfileRequest(true, "AbcdefghijklmnopqrstuvwxyzAbcdefghijklmnopqrstuvwxyzAbcdefghijklmnopqrstuvwxyzAbcdefghijklmnopqrstuvwxyz", true, "+628131746432"),
			},
			wantErr: true,
		},
		{
			name: "nil name and phone",
			args: args{
				ctx: context.Background(),
				v:   newUpdateProfileRequest(false, "", false, ""),
			},
			wantErr: true,
		},
		{
			name: "valid payload",
			args: args{
				ctx: context.Background(),
				v:   newUpdateProfileRequest(true, "Test Update", true, "+628136751232"),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateUpdateUser(tt.args.ctx, tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf("validateUpdateUser() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
