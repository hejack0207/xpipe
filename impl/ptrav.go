// Pipelines that performs traversals over the nodes.  Non-nodes will
// be filtered out.

package xpipe

import (
)

//// thisdoc
//
//  Selects the document of nodes and sends them to the sink.  Any documents encountered
//  will pass through unmodified.  All other datums are filtered out.
//
type SelectDocumentProcess struct {
}

// Configures the process using the arguments from pipeline definition
func (p *SelectDocumentProcess) Config(args []ConfigArg) error {
    return nil
}

func (p *SelectDocumentProcess) Open(ctx *ProcessContext, sink ProcessSink) error {
    return SendOpen(sink, ctx)
}

func (p *SelectDocumentProcess) Close(ctx *ProcessContext, sink ProcessSink) error {
    return SendClose(sink, ctx)
}

// Applies the process with the specific datum.
func (p *SelectDocumentProcess) Apply(ctx *ProcessContext, in Datum, sink ProcessSink) error {
    switch x := in.(type) {
    case NodeDatum:
        return SendToSink(sink, ctx, DocDatum{x.Node.MyDocument()})
    case DocDatum:
        return SendToSink(sink, ctx, in)
    default:
        return nil
    }
}
