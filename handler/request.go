package handler

import (
	"net/http"
	"strconv"
)

//Request
type Request struct {
	HttpRequest *http.Request
}

func (r *Request) Reset(request *http.Request) {
	r.HttpRequest = request
}

// Query returns input data item string by a given string.
func (r *Request) Query(key string) string {
	//if val := r.Param(key); val != "" {
	//	return val
	//}
	if r.HttpRequest.Form == nil {
		r.HttpRequest.ParseForm()
	}
	return r.HttpRequest.Form.Get(key)
}

func (r *Request) GetString(key string, def ...string) string {
	if v := r.Query(key); v != "" {
		return v
	}
	if len(def) > 0 {
		return def[0]
	}
	return ""
}

// GetInt returns input as an int or the default value while it's present and input is blank
func (r *Request) GetInt(key string, def ...int) (int, error) {
	strv := r.Query(key)
	if len(strv) == 0 && len(def) > 0 {
		return def[0], nil
	}
	return strconv.Atoi(strv)
}

// GetInt8 return input as an int8 or the default value while it's present and input is blank
func (r *Request) GetInt8(key string, def ...int8) (int8, error) {
	strv := r.Query(key)
	if len(strv) == 0 && len(def) > 0 {
		return def[0], nil
	}
	i64, err := strconv.ParseInt(strv, 10, 8)
	return int8(i64), err
}

// GetUint8 return input as an uint8 or the default value while it's present and input is blank
func (r *Request) GetUint8(key string, def ...uint8) (uint8, error) {
	strv := r.Query(key)
	if len(strv) == 0 && len(def) > 0 {
		return def[0], nil
	}
	u64, err := strconv.ParseUint(strv, 10, 8)
	return uint8(u64), err
}

// GetInt16 returns input as an int16 or the default value while it's present and input is blank
func (r *Request) GetInt16(key string, def ...int16) (int16, error) {
	strv := r.Query(key)
	if len(strv) == 0 && len(def) > 0 {
		return def[0], nil
	}
	i64, err := strconv.ParseInt(strv, 10, 16)
	return int16(i64), err
}

// GetUint16 returns input as an uint16 or the default value while it's present and input is blank
func (r *Request) GetUint16(key string, def ...uint16) (uint16, error) {
	strv := r.Query(key)
	if len(strv) == 0 && len(def) > 0 {
		return def[0], nil
	}
	u64, err := strconv.ParseUint(strv, 10, 16)
	return uint16(u64), err
}

// GetInt32 returns input as an int32 or the default value while it's present and input is blank
func (r *Request) GetInt32(key string, def ...int32) (int32, error) {
	strv := r.Query(key)
	if len(strv) == 0 && len(def) > 0 {
		return def[0], nil
	}
	i64, err := strconv.ParseInt(strv, 10, 32)
	return int32(i64), err
}

// GetUint32 returns input as an uint32 or the default value while it's present and input is blank
func (r *Request) GetUint32(key string, def ...uint32) (uint32, error) {
	strv := r.Query(key)
	if len(strv) == 0 && len(def) > 0 {
		return def[0], nil
	}
	u64, err := strconv.ParseUint(strv, 10, 32)
	return uint32(u64), err
}

// GetInt64 returns input value as int64 or the default value while it's present and input is blank.
func (r *Request) GetInt64(key string, def ...int64) (int64, error) {
	strv := r.Query(key)
	if len(strv) == 0 && len(def) > 0 {
		return def[0], nil
	}
	return strconv.ParseInt(strv, 10, 64)
}

// GetUint64 returns input value as uint64 or the default value while it's present and input is blank.
func (r *Request) GetUint64(key string, def ...uint64) (uint64, error) {
	strv := r.Query(key)
	if len(strv) == 0 && len(def) > 0 {
		return def[0], nil
	}
	return strconv.ParseUint(strv, 10, 64)
}

// GetBool returns input value as bool or the default value while it's present and input is blank.
func (r *Request) GetBool(key string, def ...bool) (bool, error) {
	strv := r.Query(key)
	if len(strv) == 0 && len(def) > 0 {
		return def[0], nil
	}
	return strconv.ParseBool(strv)
}

// GetFloat returns input value as float64 or the default value while it's present and input is blank.
func (r *Request) GetFloat(key string, def ...float64) (float64, error) {
	strv := r.Query(key)
	if len(strv) == 0 && len(def) > 0 {
		return def[0], nil
	}
	return strconv.ParseFloat(strv, 64)
}
