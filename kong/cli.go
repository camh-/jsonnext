package kong

type Config struct {
	ExtStr      map[string]string `type:"envmap" short:"V" placeholder:"var[=str]" help:"Add extVar string (from environment if <str> is omitted)"`
	ExtCode     map[string]string `type:"envmap" placeholder:"var[=code]" help:"Add extVar code (from environment if <code> is omitted)"`
	TLAStr      map[string]string `type:"envmap" short:"A" placeholder:"var[=str]" help:" Add top-level arg string (from environment if <str> is omitted)"`
	TLACode     map[string]string `type:"envmap" placeholder:"var[=code]" help:"Add top-level arg code (from environment if <code> is omitted)"`
	ExtStrFile  map[string]string `mapsep:"none" placeholder:"var=file" help:"Add extVar string from a file"`
	ExtCodeFile map[string]string `mapsep:"none" placeholder:"var=file" help:"Add extVar code from a file"`
	TLAStrFile  map[string]string `mapsep:"none" placeholder:"var=file" help:"Add top-level arg string from a file"`
	TLACodeFile map[string]string `mapsep:"none" placeholder:"var=file" help:"Add top-level arg code from a file"`
}
