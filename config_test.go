package jsonnext

import (
	"strings"
	"testing"

	"foxygo.at/s/test"
	"github.com/google/go-jsonnet"
	"github.com/stretchr/testify/require"
)

var (
	extvarSnippet = "std.extVar('var')"
	tlavarSnippet = "function(var) var"
)

func testVars(t *testing.T, c *Config, snippet string) {
	vm := jsonnet.MakeVM()
	c.ConfigureVM(vm)
	s, err := vm.EvaluateSnippet("", snippet)
	require.NoError(t, err)
	require.Equal(t, `"val"`, strings.TrimSuffix(s, "\n"))
}

func TestAddVarsExtStr(t *testing.T) {
	c := NewConfig()
	c.ExtVars["var"] = NewExtStr("val")
	testVars(t, c, extvarSnippet)
}

func TestAddVarsExtCode(t *testing.T) {
	c := NewConfig()
	c.ExtVars["var"] = NewExtCode(`"val"`)
	testVars(t, c, extvarSnippet)
}

func TestAddVarsTLAStr(t *testing.T) {
	c := NewConfig()
	c.TLAVars["var"] = NewTLAStr("val")
	testVars(t, c, tlavarSnippet)
}

func TestAddVarsTLACode(t *testing.T) {
	c := NewConfig()
	c.TLAVars["var"] = NewTLACode(`"val"`)
	testVars(t, c, tlavarSnippet)
}

func TestAddVarsExtStrFile(t *testing.T) {
	c := NewConfig()
	c.ExtVars["var"] = NewExtStrFile("testdata/config/str")
	testVars(t, c, extvarSnippet)
}

func TestAddVarsExtCodeFile(t *testing.T) {
	c := NewConfig()
	c.ExtVars["var"] = NewExtCodeFile("testdata/config/code")
	testVars(t, c, extvarSnippet)
}

func TestAddVarsTLAStrFile(t *testing.T) {
	c := NewConfig()
	c.TLAVars["var"] = NewTLAStrFile("testdata/config/str")
	testVars(t, c, tlavarSnippet)
}

func TestAddVarsTLACodeFile(t *testing.T) {
	c := NewConfig()
	c.TLAVars["var"] = NewTLACodeFile("testdata/config/code")
	testVars(t, c, tlavarSnippet)
}

func TestConfigureImporter(t *testing.T) {
	i := Importer{}
	c := NewConfig()
	c.ImportPath = []string{"a", "b"}
	c.ConfigureImporter(&i, "")
	require.Equal(t, []string{"a", "b"}, i.SearchPath)
}

func TestConfigureImporterFromEnv(t *testing.T) {
	test.Env.Set("JPATH", "c:d")
	defer test.Env.Restore()
	i := Importer{}
	c := NewConfig()
	c.ConfigureImporter(&i, "JPATH")
	require.Equal(t, []string{"c", "d"}, i.SearchPath)
}

func TestConfigureImporterFromConfigAndEnv(t *testing.T) {
	test.Env.Set("JPATH", "c:d")
	defer test.Env.Restore()
	i := Importer{}
	c := NewConfig()
	c.ImportPath = []string{"a", "b"}
	c.ConfigureImporter(&i, "JPATH")
	require.Equal(t, []string{"a", "b", "c", "d"}, i.SearchPath)
}
