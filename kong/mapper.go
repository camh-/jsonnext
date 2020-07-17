package kong

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/alecthomas/kong"
)

func init() {
	kong.DefaultRegistry.RegisterName("envmap", kong.MapperFunc(envMapDecoder))
}

func EnvMapMapper() kong.Option {
	return kong.NamedMapper("envmap", kong.MapperFunc(envMapDecoder))
}

func envMapDecoder(ctx *kong.DecodeContext, target reflect.Value) error {
	// Check that target is of type map[string]string, or compatible
	if target.Kind() != reflect.Map || target.Type().Key().Kind() != reflect.String || target.Type().Elem().Kind() != reflect.String {
		return fmt.Errorf("\"envmap\" must be applied to a map[string]string, not %s", target.Type())
	}
	if target.IsNil() {
		target.Set(reflect.MakeMap(target.Type()))
	}

	t := ctx.Scan.Pop()
	if t.IsEOL() {
		return errors.New("missing argument")
	}

	v, ok := t.Value.(string)
	if !ok {
		return fmt.Errorf("expected string, got %#v", t)
	}

	parts := strings.SplitN(v, "=", 2)
	if len(parts) == 1 {
		val, ok := os.LookupEnv(parts[0])
		if !ok {
			return fmt.Errorf("environment variable %s is not defined", parts[0])
		}
		parts = append(parts, val)
	}

	target.SetMapIndex(reflect.ValueOf(parts[0]), reflect.ValueOf(parts[1]))
	return nil
}
