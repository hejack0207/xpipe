// Main process

package main

import (
	"flag"
	"fmt"
	"os"
	"os/user"
	"path/filepath"

	"github.com/lmika/xpipe"
)

// Flag for describing the expression
var flagExpression *string = flag.String("e", "", "Process expression")

// Flag for starting with an XPath expression.  This doesn't require an -e
var flagXPathExpr *string = flag.String("x", "", "Select XPath")

// List files with at-least one datum
var flagListWithDatum *bool = flag.Bool("l", false, "List file with results")

// List files with no datums
var flagListWithoutDatum *bool = flag.Bool("L", false, "List files without results")

// Die with an error message
func die(err error) {
	fmt.Fprintf(os.Stderr, "%s: %s\n", os.Args[0], err.Error())
	os.Exit(1)
}

func main() {
	var err error
	var usageOk bool = false
	var addEmptyPipelineIfNoExpr bool = false

	flag.Parse()

	rt := xpipe.NewRuntime()

	// Setup a common end depending on what the user wants to see
	if *flagListWithDatum {
		rt.CommonEnd.Append(rt.Registry.MustNewProcess("printfile", nil))
	} else if *flagListWithoutDatum {
		rt.CommonEnd.Append(rt.Registry.MustNewProcess("printemptyfile", nil))
	} else {
		rt.CommonEnd.Append(rt.Registry.MustNewProcess("print", nil))
	}

	// Read the users RC file
	if user, err := user.Current(); err == nil {
		rcFilename := filepath.Join(user.HomeDir, ".xpiperc")
		err := rt.EvalFile(rcFilename)
		if err != nil && !os.IsNotExist(err) {
			die(err)
		}
	}

	if *flagXPathExpr != "" {
		rt.CommonStart.Append(rt.Registry.MustNewProcess("xpath", []xpipe.ConfigArg{xpipe.StringDatum(*flagXPathExpr)}))
		usageOk = true
		addEmptyPipelineIfNoExpr = true
	}

	// Parse the expression
	if *flagExpression != "" {
		err = rt.EvalString(*flagExpression, "-e")
		if err != nil {
			die(err)
		}
		usageOk = true
	} else if addEmptyPipelineIfNoExpr {
		// If there is a need to add a default pipeline (because someone has specified -x), add an empty
		// pipeline.
		rt.AddPipeline(xpipe.NewPipeline())
	}

	if !usageOk {
		flag.Usage()
		os.Exit(2)
	}

	// Execute for files
	if flag.NArg() == 0 {
		err = rt.ExecuteForFile("-")
		if err != nil {
			die(err)
		}
	} else {
		for _, filename := range flag.Args() {
			err = rt.ExecuteForFile(filename)
			if err != nil {
				die(err)
			}
		}
	}
}
