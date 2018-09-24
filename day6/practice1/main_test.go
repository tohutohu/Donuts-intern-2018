package main

import (
	"reflect"
	"testing"
	"time"

	"github.com/go-redis/redis"
)

func TestC_Get(t *testing.T) {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6378",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	client.FlushAll()
	client.Set("po", "po", time.Hour)

	type fields struct {
		client *redis.Client
	}
	type args struct {
		key string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
		err    bool
	}{
		{
			name:   "success",
			fields: fields{client},
			args:   args{"po"},
			want:   "po",
			err:    false,
		},
		{
			name:   "error",
			fields: fields{client},
			args:   args{"pi"},
			want:   "",
			err:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &C{
				client: tt.fields.client,
			}
			got := c.Get(tt.args.key)
			if !reflect.DeepEqual(got.Val(), tt.want) {
				t.Errorf("C.Get() = %v, want %v", got.Val(), tt.want)
			}
			if tt.err {
				if got.Err() == nil {
					t.Errorf("use invalid key but not return error")
				}
			}
		})
	}
}

func TestC_Set(t *testing.T) {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6378",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	client.FlushAll()
	client.Set("po", "po", time.Hour)

	type fields struct {
		client *redis.Client
	}
	type args struct {
		key      string
		value    string
		duration time.Duration
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name:   "success",
			fields: fields{client},
			args: args{
				key:      "pi",
				value:    "pi",
				duration: 1 * time.Hour,
			},
			want: "OK",
		},
		{
			name:   "error",
			fields: fields{client},
			args: args{
				key:      "po",
				value:    "pi",
				duration: 1 * time.Hour,
			},
			want: "NG",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &C{
				client: tt.fields.client,
			}
			if got := c.Set(tt.args.key, tt.args.value, tt.args.duration); !reflect.DeepEqual(got.Val(), tt.want) {
				t.Errorf("C.Set() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestC_HGet(t *testing.T) {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6378",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	client.FlushAll()
	client.Set("po", "po", time.Hour)
	type fields struct {
		client *redis.Client
	}
	type args struct {
		key   string
		field string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *redis.StringCmd
	}{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &C{
				client: tt.fields.client,
			}
			if got := c.HGet(tt.args.key, tt.args.field); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("C.HGet() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestC_HSet(t *testing.T) {
	type fields struct {
		client *redis.Client
	}
	type args struct {
		key   string
		field string
		value interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *redis.BoolCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &C{
				client: tt.fields.client,
			}
			if got := c.HSet(tt.args.key, tt.args.field, tt.args.value); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("C.HSet() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestC_LPush(t *testing.T) {
	type fields struct {
		client *redis.Client
	}
	type args struct {
		key   string
		value string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *redis.IntCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &C{
				client: tt.fields.client,
			}
			if got := c.LPush(tt.args.key, tt.args.value); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("C.LPush() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestC_RPush(t *testing.T) {
	type fields struct {
		client *redis.Client
	}
	type args struct {
		key   string
		value string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *redis.IntCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &C{
				client: tt.fields.client,
			}
			if got := c.RPush(tt.args.key, tt.args.value); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("C.RPush() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestC_LRange(t *testing.T) {
	type fields struct {
		client *redis.Client
	}
	type args struct {
		key   string
		start int64
		stop  int64
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *redis.StringSliceCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &C{
				client: tt.fields.client,
			}
			if got := c.LRange(tt.args.key, tt.args.start, tt.args.stop); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("C.LRange() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestC_LLen(t *testing.T) {
	type fields struct {
		client *redis.Client
	}
	type args struct {
		key string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *redis.IntCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &C{
				client: tt.fields.client,
			}
			if got := c.LLen(tt.args.key); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("C.LLen() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestC_ZAdd(t *testing.T) {
	type fields struct {
		client *redis.Client
	}
	type args struct {
		key   string
		value redis.Z
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *redis.IntCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &C{
				client: tt.fields.client,
			}
			if got := c.ZAdd(tt.args.key, tt.args.value); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("C.ZAdd() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestC_ZRange(t *testing.T) {
	type fields struct {
		client *redis.Client
	}
	type args struct {
		key   string
		start int64
		stop  int64
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *redis.StringSliceCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &C{
				client: tt.fields.client,
			}
			if got := c.ZRange(tt.args.key, tt.args.start, tt.args.stop); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("C.ZRange() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestC_ZRemRangeByRank(t *testing.T) {
	type fields struct {
		client *redis.Client
	}
	type args struct {
		key   string
		start int64
		stop  int64
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *redis.IntCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &C{
				client: tt.fields.client,
			}
			if got := c.ZRemRangeByRank(tt.args.key, tt.args.start, tt.args.stop); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("C.ZRemRangeByRank() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestC_ZRevRange(t *testing.T) {
	type fields struct {
		client *redis.Client
	}
	type args struct {
		key   string
		start int64
		stop  int64
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *redis.StringSliceCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &C{
				client: tt.fields.client,
			}
			if got := c.ZRevRange(tt.args.key, tt.args.start, tt.args.stop); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("C.ZRevRange() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestC_ZIncrBy(t *testing.T) {
	type fields struct {
		client *redis.Client
	}
	type args struct {
		key   string
		score float64
		value string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *redis.FloatCmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &C{
				client: tt.fields.client,
			}
			if got := c.ZIncrBy(tt.args.key, tt.args.score, tt.args.value); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("C.ZIncrBy() = %v, want %v", got, tt.want)
			}
		})
	}
}
