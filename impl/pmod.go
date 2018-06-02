// Pipelines that performs modifications to datums.

package xpipe

import (
    "fmt"
)

//// settext <newvalue>
//
//  Changes the contents of nodes and document datums to *newvalue*.  All other datums
//  are sent to the sink unchanged.
//
type SetTextProcess struct {
    NewValue    ConfigArg
}

// Configures the process using the arguments from pipeline definition
func (p *SetTextProcess) Config(args []ConfigArg) error {
    if len(args) != 1 {
        return fmt.Errorf("settext expects 1 argument")
    }

    p.NewValue = args[0]
    return nil
}

func (p *SetTextProcess) Open(ctx *ProcessContext, sink ProcessSink) error {
    return SendOpen(sink, ctx)
}

func (p *SetTextProcess) Close(ctx *ProcessContext, sink ProcessSink) error {
    return SendClose(sink, ctx)
}

// Applies the process with the specific datum.
func (p *SetTextProcess) Apply(ctx *ProcessContext, in Datum, sink ProcessSink) error {
    switch x := in.(type) {
    case NodeDatum:
        x.Node.SetContent(p.NewValue.String())
    case DocDatum:
        x.Doc.Root().SetContent(p.NewValue.String())
    }

    return SendToSink(sink, ctx, in)
}
