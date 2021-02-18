package filter

import "testing"

func TestCalLuaFilter(t *testing.T) {
	luaCode1 := `
					function Filter(metadata)
					  userId = metadata["userid"]
					  if( tonumber(userId) < 1000 )
					  then
						 return true
					  else
						 return false
					  end
					end
`
	type args struct {
		code      string
		metatdata map[string]string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "test-01",
			args: args{
				code: luaCode1,
				metatdata: map[string]string{
					"userid": "1000",
				}},
			want: false,
		},
		{
			name: "test-02",
			args: args{
				code: luaCode1,
				metatdata: map[string]string{
					"userid": "10",
				}},
			want: true,
		},
		{
			name: "test-03",
			args: args{
				code: luaCode1,
				metatdata: map[string]string{
					"userid": "10000",
				}},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CalLuaFilter(tt.args.code, tt.args.metatdata); got != tt.want {
				t.Errorf("CalLuaFilter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCalMepFilter(t *testing.T) {
	type args struct {
		code      string
		metatdata map[string]string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CalMepFilter(tt.args.code, tt.args.metatdata); got != tt.want {
				t.Errorf("CalMepFilter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCalSimpleFilter(t *testing.T) {
	type args struct {
		code      string
		metatdata map[string]string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CalSimpleFilter(tt.args.code, tt.args.metatdata); got != tt.want {
				t.Errorf("CalSimpleFilter() = %v, want %v", got, tt.want)
			}
		})
	}
}
