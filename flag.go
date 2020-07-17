package jsonnext

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

// ConfigFlags defines a set of flags for a Config struct to populate the
// fields from the command line. The return value is the address of the
// Config struct that stores the values of the flags.
// The fields are populated with the following flags:
//  Config.ImportPath:
//   -J, -jpath
//  Config.ExtVars:
//   -V, -ext-str: ext var as string literal
//   -ext-code: ext var as code literal
//   -ext-str-file: ext var as string from file
//   -ext-code-file: ext var as code from file
//  Config.TLAVars:
//   -A, -tla-str: top-level arg as string literal
//   -tla-code: top-level arg as code literal
//   -tla-str-file: top-level arg as string from file
//   -tla-code-file: top-level arg as code from file
func ConfigFlags(fs *flag.FlagSet) *Config {
	c := NewConfig()
	ConfigFlagsVar(fs, c)
	return c
}

// ConfigFlagsVar defines a set of flags for a Config struct to populate
// the fields from the command line. The argument c points to the Config
// struct to populate. The set of flags defined is described in the
// ConfigFlags function description.
func ConfigFlagsVar(fs *flag.FlagSet, c *Config) {
	StringSliceVar(fs, &c.ImportPath, "jpath", "Add a library search `dir`")
	ExtStrMapVar(fs, c.ExtVars, "ext-str", "Add extVar `var[=str]` (from environment if <str> is omitted)")
	ExtCodeMapVar(fs, c.ExtVars, "ext-code", "Add extVar `var[=code]` (from environment if <code> is omitted)")
	ExtStrFileMapVar(fs, c.ExtVars, "ext-str-file", "Add extVar `var=file` string from a file")
	ExtCodeFileMapVar(fs, c.ExtVars, "ext-code-file", "Add extVar `var=file` code from a file")

	TLAStrMapVar(fs, c.TLAVars, "tla-str", "Add top-level arg `var=[=str]` (from environment if <str> is omitted)")
	TLACodeMapVar(fs, c.TLAVars, "tla-code", "Add top-level arg `var[=code]` (from environment if <code> is omitted)")
	TLAStrFileMapVar(fs, c.TLAVars, "tla-str-file", "Add top-level arg `var=file` string from a file")
	TLACodeFileMapVar(fs, c.TLAVars, "tla-code-file", "Add top-level arg `var=file` code from a file")

	// Add short flags. TODO(camh): consider making this optional
	StringSliceVar(fs, &c.ImportPath, "J", "Add a library search `dir`")
	ExtStrMapVar(fs, c.ExtVars, "V", "Add extVar `var[=str]` (from environment if <str> is omitted)")
	TLAStrMapVar(fs, c.TLAVars, "A", "Add top-level arg `var[=str]` (from environment if <str> is omitted)")
}

// StringSliceVar defines a flag in the specified FlagSet with the specified
// name and usage string. The argument p is a pointer to a []string variable in
// which to store the value of the flag. The value given to each instance of
// the flag is appended to the slice.
func StringSliceVar(fs *flag.FlagSet, p *[]string, name, usage string) {
	fs.Var((*stringSliceValue)(p), name, usage)
}

// StringSlice defines a flag in the specified FlagSet with the specified name
// and usage string. The return value is the address of a []string variable
// that stores the value of the flag. The value given to each instance of the
// flag is appended to the slice.
func StringSlice(fs *flag.FlagSet, name, usage string) *[]string {
	p := &[]string{}
	StringSliceVar(fs, p, name, usage)
	return p
}

type stringSliceValue []string

func (s *stringSliceValue) Set(v string) error {
	*s = append(*s, v)
	return nil
}

func (s *stringSliceValue) Get() interface{} { return []string(*s) }
func (s *stringSliceValue) String() string   { return "" }

// ExtStrMapVar defines a flag in the specified FlagSet with the specified
// name and usage string. The argument m is the map that
func ExtStrMapVar(fs *flag.FlagSet, m map[string]VMVar, name, usage string) {
	fs.Var(vmVarMapValue{m, NewExtStr, true}, name, usage)
}

func ExtCodeMapVar(fs *flag.FlagSet, m map[string]VMVar, name, usage string) {
	fs.Var(vmVarMapValue{m, NewExtCode, true}, name, usage)
}

func ExtStrFileMapVar(fs *flag.FlagSet, m map[string]VMVar, name, usage string) {
	fs.Var(vmVarMapValue{m, NewExtStrFile, false}, name, usage)
}

func ExtCodeFileMapVar(fs *flag.FlagSet, m map[string]VMVar, name, usage string) {
	fs.Var(vmVarMapValue{m, NewExtCodeFile, false}, name, usage)
}

func TLAStrMapVar(fs *flag.FlagSet, m map[string]VMVar, name, usage string) {
	fs.Var(vmVarMapValue{m, NewTLAStr, true}, name, usage)
}

func TLACodeMapVar(fs *flag.FlagSet, m map[string]VMVar, name, usage string) {
	fs.Var(vmVarMapValue{m, NewTLACode, true}, name, usage)
}

func TLAStrFileMapVar(fs *flag.FlagSet, m map[string]VMVar, name, usage string) {
	fs.Var(vmVarMapValue{m, NewTLAStrFile, false}, name, usage)
}

func TLACodeFileMapVar(fs *flag.FlagSet, m map[string]VMVar, name, usage string) {
	fs.Var(vmVarMapValue{m, NewTLACodeFile, false}, name, usage)
}

type vmVarMapValue struct {
	m       map[string]VMVar
	makevar func(string) VMVar
	fromenv bool
}

func (mv vmVarMapValue) Set(v string) error {
	parts := strings.SplitN(v, "=", 2)
	if len(parts) == 1 {
		if !mv.fromenv {
			return fmt.Errorf("expected <key>=<value> but got \"%s\"", v)
		}
		val, ok := os.LookupEnv(parts[0])
		if !ok {
			return fmt.Errorf("environment variable %s is not defined", parts[0])
		}
		parts = append(parts, val)
	}
	mv.m[parts[0]] = mv.makevar(parts[1])
	return nil
}
func (mv vmVarMapValue) Get() interface{} { return map[string]VMVar(mv.m) }
func (mv vmVarMapValue) String() string   { return "" }
