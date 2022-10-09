package config

type Args struct {
	Format string `short:"f" long:"format" description:"The desired output format" choice:"excel" choice:"json" default:"excel"`

	Verbose bool `short:"v" long:"verbose" description:"Show verbose debug information"`

	StdOut bool `short:"o" long:"out" description:"Print to stdout instead of file"`

	EnableMissingOperatorsWarning bool `long:"missing-operator-warning" description:"Enable warnings for missing operators"`

	PositionalArgs struct {
		InputFile string `positional-arg-name:"input-file" required:"yes"`
		Rest      []string
	} `positional-args:"yes"`
}
