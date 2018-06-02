// Parser for pipe expressions
package xpipe

import (
    "io"
    "fmt"
    "strconv"
    "text/scanner"
)


// Scanner
type Scanner struct {
    S           scanner.Scanner
    Token       rune
    TokenText   string
}

// Creates a new scanner
func NewScanner(r io.Reader, filename string) *Scanner {
    s := scanner.Scanner{}
    s.Init(r)
    s.Position.Filename = filename

    tok := s.Scan()
    toktext := s.TokenText()

    return &Scanner{s, tok, toktext}
}

// Fetch the next token
func (s *Scanner) Scan() {
    s.Token = s.S.Scan()
    s.TokenText = s.S.TokenText()
}

// Returns true if the current token is an ident with a specific text.
func (s *Scanner) IsKeyword(value string) bool {
    return (s.Token == scanner.Ident) && (s.TokenText == value)
}

// Returns the position
func (s *Scanner) Position() string {
    return s.S.Pos().String()
}

// --------------------------------------------------------------------------
//

// A Parser
type Parser struct {
    Scanner     *Scanner
}

// Creates a scan error
func (p *Parser) Error(msg string) error {
    return fmt.Errorf("%s: error - %s", p.Scanner.Position(), msg)
}

// Consumes a token of a specific type.  If the token is not of that type, returns an error
func (p *Parser) Consume(tok rune) error {
    if (p.Scanner.Token != tok) {
        return p.Error(fmt.Sprintf("Expected %s but found %s",
                scanner.TokenString(tok), scanner.TokenString(p.Scanner.Token)))
    } else {
        p.Scanner.Scan()
        return nil
    }
}

// Creates a new parser 
func NewParser(r io.Reader, filename string) *Parser {
    return &Parser{ NewScanner(r, filename) }
}

// Parses a script.
//  <script> = <item> (";" <item>)*
func (p *Parser) ParseScript() (*AstScript, error) {
    items := make([]AstItem, 0)
    for (p.Scanner.Token != scanner.EOF) {
        item, err := p.parseItem()
        if err != nil {
            return nil, err
        }

        if (p.Scanner.Token == ';') {
            p.Scanner.Scan()
        } else if (p.Scanner.Token != scanner.EOF) {
            return nil, p.Error("Unexpected token")
        }

        items = append(items, item)
    }

    return &AstScript{items}, nil
}

// Parse a single item.
//  <item> = <namespaceMapping> | <pipeline>
func (p *Parser) parseItem() (AstItem, error) {
    if (p.Scanner.IsKeyword("ns")) {
        return p.parseNamespaceMapping()
    } else {
        return p.parsePipeline()
    }
}

// Parse a namespace mapping.
//  <namespaceMapping> = "ns" <IDENT> "=" <STRING>
func (p *Parser) parseNamespaceMapping() (AstItem, error) {
    var err error
    if !p.Scanner.IsKeyword("ns") {
        return nil, p.Error("Expected keyword 'ns'")
    }
    p.Scanner.Scan()

    prefix := p.Scanner.TokenText
    if err = p.Consume(scanner.Ident) ; err != nil {
        return nil, err
    }

    if err := p.Consume('=') ; err != nil {
        return nil, err
    }

    value := p.Scanner.TokenText
    if err = p.Consume(scanner.String) ; err != nil {
        return nil, err
    }
    if value, err = strconv.Unquote(value); err != nil {
        return nil, err
    }

    return &AstNamespaceMapping{prefix, value}, nil
}

// Parses a pipeline
//  <pipeline> = <process> ("|" <process>)
func (p *Parser) parsePipeline() (*AstPipeline, error) {
    var firstPr *AstProcess = nil
    var currPr *AstProcess = nil

    for p.Scanner.Token != scanner.EOF {
        pr, err := p.parseProcess()
        if err != nil {
            return nil, err
        }

        if firstPr == nil {
            firstPr = pr
        }
        if currPr != nil {
            currPr.Next = pr
        }
        currPr = pr

        // If the next token is not '|' then break
        if (p.Scanner.Token != '|') {
            break
        } else {
            p.Scanner.Scan()
        }
    }

    return &AstPipeline{firstPr}, nil
}

// Parses a single process
//  <process> = IDENT (<args>*)
func (p *Parser) parseProcess() (*AstProcess, error) {
    if p.Scanner.Token != scanner.Ident {
        return nil, p.Error("Expected identifier")
    }

    prName := p.Scanner.TokenText
    args := make([]AstProcessArg, 0)

    //pr := &AstProcess{p.Scanner.TokenText, nil}
    p.Scanner.Scan()

    for (p.Scanner.Token != '|') && (p.Scanner.Token != ';') && (p.Scanner.Token != scanner.EOF) {
        arg, err := p.parseProcessArg()
        if err != nil {
            return nil, err
        }

        args = append(args, arg)
    }

    return &AstProcess{prName, args, nil}, nil
}

// Parses a single process argument
//  <arg> = <string>
func (p *Parser) parseProcessArg() (AstProcessArg, error) {
    if (p.Scanner.Token == scanner.String) {
        sval, err := strconv.Unquote(p.Scanner.TokenText)
        if err != nil {
            return nil, err
        }

        p.Scanner.Scan()
        return &AstLiteralProcessArg{StringDatum(sval)}, nil
    } else {
        return nil, p.Error("Unreognised process argument type")
    }
}
