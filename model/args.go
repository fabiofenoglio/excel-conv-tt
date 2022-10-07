package model

type Args struct {
	Verbose bool `short:"v" long:"verbose" description:"Show verbose debug information"`

	PositionalArgs struct {
		InputFile string `positional-arg-name:"input-file" required:"yes"`
		Rest      []string
	} `positional-args:"yes"`
}
