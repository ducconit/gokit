package pagination

import (
	"fmt"
	"net/url"
	"reflect"
	"strconv"
	"time"
)

type Normalizer interface {
	Normalize()
}

func BindQuery(values url.Values, dst any) error {
	rv := reflect.ValueOf(dst)
	if rv.Kind() != reflect.Pointer || rv.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("pagination: dst must be pointer to struct")
	}

	st := rv.Elem().Type()
	sv := rv.Elem()

	for i := 0; i < st.NumField(); i++ {
		f := st.Field(i)
		if f.PkgPath != "" {
			continue
		}
		fv := sv.Field(i)
		if !fv.CanSet() {
			continue
		}

		name := pickName(f)
		raw := values.Get(name)
		if raw == "" {
			raw = f.Tag.Get("default")
		}
		if raw == "" {
			continue
		}

		if err := setValue(fv, raw); err != nil {
			return fmt.Errorf("pagination: field %s: %w", f.Name, err)
		}
	}

	Normalize(dst)
	return nil
}

func Normalize(v any) {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Pointer || rv.Elem().Kind() != reflect.Struct {
		return
	}

	rv = rv.Elem()
	rt := rv.Type()

	for i := 0; i < rt.NumField(); i++ {
		f := rt.Field(i)
		fv := rv.Field(i)

		// Handle embedded structs
		if f.Anonymous && fv.Kind() == reflect.Struct {
			Normalize(fv.Addr().Interface())
			continue
		}

		if !fv.CanSet() {
			continue
		}

		// Apply default if zero value
		if fv.IsZero() {
			if def := f.Tag.Get("default"); def != "" {
				_ = setValue(fv, def)
			}
		}

		// Apply max limit for Limit field
		if (f.Name == "Limit" || f.Tag.Get("query") == "limit" || f.Tag.Get("form") == "limit") && fv.Kind() == reflect.Int {
			if maxStr := f.Tag.Get("max"); maxStr != "" {
				if maxVal, err := strconv.Atoi(maxStr); err == nil && maxVal > 0 {
					if fv.Int() > int64(maxVal) {
						fv.SetInt(int64(maxVal))
					}
				}
			}
		}
	}
}

func pickName(f reflect.StructField) string {
	if v := f.Tag.Get("query"); v != "" && v != "-" {
		return v
	}
	if v := f.Tag.Get("form"); v != "" && v != "-" {
		return v
	}
	if v := f.Tag.Get("json"); v != "" && v != "-" {
		return v
	}
	return f.Name
}

func setValue(v reflect.Value, raw string) error {
	switch v.Kind() {
	case reflect.String:
		v.SetString(raw)
		return nil
	case reflect.Bool:
		b, err := strconv.ParseBool(raw)
		if err != nil {
			return err
		}
		v.SetBool(b)
		return nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if v.Type() == reflect.TypeOf(time.Duration(0)) {
			d, err := time.ParseDuration(raw)
			if err != nil {
				return err
			}
			v.SetInt(int64(d))
			return nil
		}
		n, err := strconv.ParseInt(raw, 10, v.Type().Bits())
		if err != nil {
			return err
		}
		v.SetInt(n)
		return nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		n, err := strconv.ParseUint(raw, 10, v.Type().Bits())
		if err != nil {
			return err
		}
		v.SetUint(n)
		return nil
	default:
		return fmt.Errorf("unsupported kind: %s", v.Kind())
	}
}
