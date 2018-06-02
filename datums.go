// Datums of XPipe
//

package xpipe

import (
    "fmt"
    "github.com/moovweb/gokogiri/xml"
)


// Base datum of an XProc pipeline.
type Datum interface {

    // Returns the standard string representation of the datum
    String() string
}

// Base config argument to a process.
type ConfigArg interface {

    // Returns the string representation of the config argument.
    String() string

    // Returns true if the config argument is a constant
    IsConst() bool
}


// A string datum
type StringDatum        string

func (s StringDatum) String() string {
    return string(s)
}

func (s StringDatum) IsConst() bool {
    return true
}


// A boolean datum
type BoolDatum          bool

func (b BoolDatum) String() string {
    if bool(b) {
        return "true"
    } else {
        return "false"
    }
}

func (n BoolDatum) IsConst() bool {
    return true
}


// A number datum
type NumberDatum        float64

func (n NumberDatum) String() string {
    return fmt.Sprintf("%g", float64(n))
}

func (n NumberDatum) IsConst() bool {
    return true
}


// A node datum
type NodeDatum          struct {
    Node        xml.Node
}

func (n NodeDatum) String() string {
    return n.Node.String()
}

func (n NodeDatum) IsConst() bool {
    return false
}


// A document datum
type DocDatum           struct {
    Doc         xml.Document
}

func (d DocDatum) String() string {
    return d.Doc.String()
}

func (d DocDatum) IsConst() bool {
    return false
}

// --------------------------------------------------------------------------
//

// Execution context for a pipeline.
type ProcessContext    struct {
    Runtime     *Runtime
    Filename    string
}
