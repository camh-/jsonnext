package jsonnext

import (
	"strings"

	"github.com/google/go-jsonnet"
)

// Config holds configuration for a jsonnet VM and other types defined in this
// jsonnext package. Command line parsing into this struct
type Config struct {
	ImportPath []string         `name:"jpath" sep:"none" short:"J" placeholder:"dir" help:"Add a library search dir"`
	ExtVars    map[string]VMVar `kong:"-"`
	TLAVars    map[string]VMVar `kong:"-"`
}

// VMVar is a variable that can be pre-set in a jsonnet VM, either as an extVar
// or a top-level arg (TLA) as a string or code and either a literal or from a
// file. It maps to the --ext-* and --tla-* command line arguments to the
// standard jsonnet binary.
type VMVar interface {
	Set(key string, vm *jsonnet.VM)
}

// NewConfig returns a new Config struct with the field initialised.
func NewConfig() *Config {
	return &Config{
		ExtVars: map[string]VMVar{},
		TLAVars: map[string]VMVar{},
	}
}

// ConfigureImporter sets up a jsonnext.Importer with the import path from
// the config and from a PATH-style environment variable. If ennvar is the
// empty string, no paths are taken from the environment.
func (c *Config) ConfigureImporter(i *Importer, envvar string) {
	i.SearchPath = c.ImportPath
	if envvar != "" {
		i.AppendSearchFromEnv(envvar)
	}
}

// ConfigureVM sets the VMVars in the config in the jsonnet VM.
func (c *Config) ConfigureVM(vm *jsonnet.VM) {
	for k, v := range c.ExtVars {
		v.Set(k, vm)
	}
	for k, v := range c.TLAVars {
		v.Set(k, vm)
	}
}

type extStr string
type extCode string
type extStrFile string
type extCodeFile string

func NewExtStr(s string) VMVar                       { return extStr(s) }
func (v extStr) Set(key string, vm *jsonnet.VM)      { vm.ExtVar(key, string(v)) }
func NewExtCode(s string) VMVar                      { return extCode(s) }
func (v extCode) Set(key string, vm *jsonnet.VM)     { vm.ExtCode(key, string(v)) }
func NewExtStrFile(s string) VMVar                   { return extStrFile(s) }
func (v extStrFile) Set(key string, vm *jsonnet.VM)  { vm.ExtCode(key, mkImportStr(string(v))) }
func NewExtCodeFile(s string) VMVar                  { return extCodeFile(s) }
func (v extCodeFile) Set(key string, vm *jsonnet.VM) { vm.ExtCode(key, mkImport(string(v))) }

type tlaStr string
type tlaCode string
type tlaStrFile string
type tlaCodeFile string

func NewTLAStr(s string) VMVar                       { return tlaStr(s) }
func (v tlaStr) Set(key string, vm *jsonnet.VM)      { vm.TLAVar(key, string(v)) }
func NewTLACode(s string) VMVar                      { return tlaCode(s) }
func (v tlaCode) Set(key string, vm *jsonnet.VM)     { vm.TLACode(key, string(v)) }
func NewTLAStrFile(s string) VMVar                   { return tlaStrFile(s) }
func (v tlaStrFile) Set(key string, vm *jsonnet.VM)  { vm.TLACode(key, mkImportStr(string(v))) }
func NewTLACodeFile(s string) VMVar                  { return tlaCodeFile(s) }
func (v tlaCodeFile) Set(key string, vm *jsonnet.VM) { vm.TLACode(key, mkImport(string(v))) }

// Quote string using verbatim string: @'...'
func quoteStr(s string) string    { return "@'" + strings.ReplaceAll(s, "'", "''") + "'" }
func mkImport(f string) string    { return "import " + quoteStr(f) }
func mkImportStr(f string) string { return "importstr" + quoteStr(f) }
