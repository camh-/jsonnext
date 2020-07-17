package jsonnext

import (
	"flag"
	"testing"

	"foxygo.at/s/test"
	"github.com/stretchr/testify/require"
)

func TestConfigFlags(t *testing.T) {
	fs := &flag.FlagSet{}
	c := ConfigFlags(fs)

	args := []string{
		"-jpath", "a", "-J", "b",
		"-ext-str", "world=123", "-V", "extra=789",
		"-ext-code", "bar=21 + 21", "-ext-code", "extra=false",
		"-ext-str-file", "yaml=hello.yaml",
		"-ext-code-file", "jsonnet=hello.jsonnet",
		"-tla-str", "world=456", "-A", "extra=123",
		"-tla-code", "bar=20 + 22", "-tla-code", "extra=true",
		"-tla-str-file", "yaml=goodbye.yaml",
		"-tla-code-file", "jsonnet=goodbye.jsonnet",
	}
	expectedConfig := &Config{
		ImportPath: []string{"a", "b"},
		ExtVars: map[string]VMVar{
			"world":   extStr("123"),
			"bar":     extCode("21 + 21"),
			"extra":   extCode("false"),
			"yaml":    extStrFile("hello.yaml"),
			"jsonnet": extCodeFile("hello.jsonnet"),
		},
		TLAVars: map[string]VMVar{
			"world":   tlaStr("456"),
			"bar":     tlaCode("20 + 22"),
			"extra":   tlaCode("true"),
			"yaml":    tlaStrFile("goodbye.yaml"),
			"jsonnet": tlaCodeFile("goodbye.jsonnet"),
		},
	}

	err := fs.Parse(args)
	require.NoError(t, err)
	require.Equal(t, expectedConfig, c)
}

func TestStringSlice(t *testing.T) {
	fs := &flag.FlagSet{}
	ss := StringSlice(fs, "food", "food to eat")
	err := fs.Parse([]string{"-food", "banana", "-food", "pie"})
	require.NoError(t, err)
	require.Equal(t, &[]string{"banana", "pie"}, ss)

	// Test Get() method
	f := fs.Lookup("food")
	require.NotNil(t, f)
	require.Equal(t, *ss, f.Value.(flag.Getter).Get())
}

func TestExtVars(t *testing.T) {
	test.Env.Set("hello", "abc").Set("foo", "42")
	defer test.Env.Restore()
	fs := &flag.FlagSet{}
	m := map[string]VMVar{}
	ExtStrMapVar(fs, m, "ext-str", "ext var string")
	ExtCodeMapVar(fs, m, "ext-code", "ext var code")
	ExtStrFileMapVar(fs, m, "ext-str-file", "ext var string file")
	ExtCodeFileMapVar(fs, m, "ext-code-file", "ext var code file")

	args := []string{
		"-ext-str", "hello", "-ext-str", "world=123", "-ext-str", "extra=789",
		"-ext-code", "foo", "-ext-code", "bar=21 + 21", "-ext-code", "extra=false",
		"-ext-str-file", "yaml=hello.yaml",
		"-ext-code-file", "jsonnet=hello.jsonnet",
	}
	expected := map[string]VMVar{
		"hello":   extStr("abc"),
		"world":   extStr("123"),
		"foo":     extCode("42"),
		"bar":     extCode("21 + 21"),
		"extra":   extCode("false"),
		"yaml":    extStrFile("hello.yaml"),
		"jsonnet": extCodeFile("hello.jsonnet"),
	}

	err := fs.Parse(args)
	require.NoError(t, err)
	require.Equal(t, expected, m)
}

func TestTLAVars(t *testing.T) {
	test.Env.Set("hello", "abc").Set("foo", "42")
	defer test.Env.Restore()
	fs := &flag.FlagSet{}
	m := map[string]VMVar{}
	TLAStrMapVar(fs, m, "tla-str", "tla var string")
	TLACodeMapVar(fs, m, "tla-code", "tla var code")
	TLAStrFileMapVar(fs, m, "tla-str-file", "tla var string file")
	TLACodeFileMapVar(fs, m, "tla-code-file", "tla var code file")

	args := []string{
		"-tla-str", "hello", "-tla-str", "world=123", "-tla-str", "extra=789",
		"-tla-code", "foo", "-tla-code", "bar=21 + 21", "-tla-code", "extra=false",
		"-tla-str-file", "yaml=hello.yaml",
		"-tla-code-file", "jsonnet=hello.jsonnet",
	}
	expected := map[string]VMVar{
		"hello":   tlaStr("abc"),
		"world":   tlaStr("123"),
		"foo":     tlaCode("42"),
		"bar":     tlaCode("21 + 21"),
		"extra":   tlaCode("false"),
		"yaml":    tlaStrFile("hello.yaml"),
		"jsonnet": tlaCodeFile("hello.jsonnet"),
	}

	err := fs.Parse(args)
	require.NoError(t, err)
	require.Equal(t, expected, m)
}

func TestMapValueMissing(t *testing.T) {
	fs := &flag.FlagSet{}
	m := map[string]VMVar{}
	TLAStrFileMapVar(fs, m, "tla-str-file", "tla var string file")

	err := fs.Parse([]string{"-tla-str-file", "hello"})

	require.Error(t, err)
}

func TestMapValueMissingEnv(t *testing.T) {
	test.Env.Unset("foo")
	defer test.Env.Restore()
	fs := &flag.FlagSet{}
	m := map[string]VMVar{}
	TLAStrMapVar(fs, m, "tla-str", "tla var string")

	err := fs.Parse([]string{"-tla-str", "foo"})

	require.Error(t, err)
}

func TestVMVarMapGet(t *testing.T) {
	fs := &flag.FlagSet{}
	m := map[string]VMVar{}
	fs.Var(vmVarMapValue{m, NewExtStr, true}, "ext-str", "ext var string")
	err := fs.Parse([]string{"-ext-str", "foo=bar"})
	require.NoError(t, err)
	f := fs.Lookup("ext-str")
	require.NotNil(t, f)
	require.Equal(t, m, f.Value.(flag.Getter).Get())
}
