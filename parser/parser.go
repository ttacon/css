package parser

import (
	"fmt"
	"strings"

	"github.com/ttacon/css/ast"
	"github.com/ttacon/css/scanner"
)

type Parser struct {
	s     *scanner.Scanner
	cache []*scanner.Token
}

func New(s *scanner.Scanner) *Parser {
	return &Parser{
		s: s,
	}
}

func isEnd(t *scanner.Token) bool {
	return t.Type == scanner.TokenError || t.Type == scanner.TokenEOF
}

func isBlockOpen(t *scanner.Token) bool {
	if t.Type != scanner.TokenChar {
		return false
	}
	return t.Value == "{" || t.Value == "[" || t.Value == "("
}

func isSemiColon(t *scanner.Token) bool {
	return t.Type == scanner.TokenChar && t.Value == ";"
}

func (p *Parser) Parse() (*ast.Stylesheet, error) {
	var (
		t     = p.nextNonWhitespaceToken()
		rules []ast.Rule
	)
	// TODO(ttacon): change to use channels/consumption
	for ; !isEnd(t); t = p.nextNonWhitespaceToken() {
		switch t.Type {
		case scanner.TokenS: // just consume it
		case scanner.TokenAtKeyword:
			// TODO(ttacon): pull out to own method
			// at-rule     : ATKEYWORD S* any* [ block | ';' S* ];
			// TODO(ttacon): this needs to actually comsume any
			any := p.nextNonWhitespaceToken()

			// sniff it
			var block *ast.Block
			var err error
			next := p.nextNonWhitespaceToken()
			semiOnly := false
			if isBlockOpen(next) {
				block, err = p.parseBlock(next)
				if err != nil {
					return nil, err
				}
			} else if isSemiColon(next) {
				semiOnly = true
			} else {
				// this is a parse error
				return nil, fmt.Errorf("expected opening block or ';', found %v", t)
			}

			rules = append(rules, &ast.AtRule{
				AtKeyword: t.Value,
				Any:       any.Value,
				Block:     block,
				JustSemi:  semiOnly,
			})
		case scanner.TokenChar:
			// this should just be '.'
			if t.Value != "." {
				return nil, fmt.Errorf("expected '.', got %q", t.Value)
			}

			newRule, err := p.parseQualifiedRule(t)
			if err != nil {
				return nil, err
			}

			rules = append(rules, newRule)
		case scanner.TokenHash:
			// this should just be '#'
			if !strings.HasPrefix(t.Value, "#") {
				return nil, fmt.Errorf("expected '#', got %q", t.Value)
			}

			newRule, err := p.parseQualifiedRule(t)
			if err != nil {
				return nil, err
			}

			rules = append(rules, newRule)
		case scanner.TokenIdent:
			newRule, err := p.parseQualifiedRule(t)
			if err != nil {
				return nil, err
			}
			rules = append(rules, newRule)
		default:
			fmt.Println("default: ", t)
		}
	}
	return &ast.Stylesheet{rules}, nil
}

func (p *Parser) parseQualifiedRule(entry *scanner.Token) (ast.Rule, error) {
	var (
		t    = entry
		name string
	)

	if entry.Type == scanner.TokenChar {
		if entry.Value != "." {
			// invalid starting token for rule/component
			return nil, fmt.Errorf("expecting '.', got %s", entry.Value)
		}
		name = entry.Value
		t = p.nextNonWhitespaceToken()
	}

	if entry.Type == scanner.TokenHash {
		if !strings.HasPrefix(entry.Value, "#") {
			// invalid starting token for rule/component
			return nil, fmt.Errorf("expecting '#', got %s", entry.Value)
		}
	}

	name = name + t.Value
	var names []string
	t = p.nextNonWhitespaceToken()
	var err error
	for (t.Type == scanner.TokenChar && t.Value == ",") ||
		t.Type == scanner.TokenIdent {

		if t.Type == scanner.TokenIdent {
			name = name + " " + t.Value
			t = p.nextNonWhitespaceToken()
			continue
		}

		names = append(names, name)

		t, err = p.componentValue()
		if err != nil {
			return nil, err
		}

		name = t.Value
		t = p.nextNonWhitespaceToken()
	}

	names = append(names, name)

	if t.Value != "{" {
		return nil, fmt.Errorf("expected '{', got %q", t.Value)
	}

	decls, err := p.parseDeclarations()
	if err != nil {
		return nil, err
	}

	var components = make([]*ast.ComponentValue, len(names))
	for i, name := range names {
		components[i] = &ast.ComponentValue{name}
	}

	return &ast.QualifiedRule{
		Components: components,
		DeclList:   decls,
	}, nil
}

func (p *Parser) nextNonWhitespaceToken() *scanner.Token {
	var t = p.s.Next()
	for t.Type == scanner.TokenS {
		t = p.s.Next()
	}
	return t
}

func (p *Parser) parseDeclarations() (*ast.DeclarationList, error) {
	// sniff @-rule vs decl
	tok := p.nextNonWhitespaceToken()
	var decls []*ast.Declaration
	for ; tok.Value != "}"; tok = p.nextNonWhitespaceToken() {
		if tok.Type == scanner.TokenAtKeyword {
			// TODO(ttacon): do it
			continue
		}

		decl, err := p.parseDeclaration(tok)
		if err != nil {
			return nil, err
		}

		decls = append(decls, decl)
	}

	return &ast.DeclarationList{decls}, nil
}

func (p *Parser) parseDeclaration(ident *scanner.Token) (*ast.Declaration, error) {
	tok := p.nextNonWhitespaceToken()
	if tok.Type != scanner.TokenChar || tok.Value != ":" {
		return nil, fmt.Errorf("expected ':', got %s", tok.Value)
	}

	var components []string
	tok = p.nextNonWhitespaceToken()
	for ; tok.Type != scanner.TokenError && tok.Type != scanner.TokenEOF; tok = p.nextNonWhitespaceToken() {
		if tok.Value == "!important" {
			return &ast.Declaration{
				Ident:      ident.Value,
				Components: append(components, tok.Value),
			}, nil
		}
		if tok.Value == ";" {
			break
		}

		components = append(components, tok.Value)
	}

	if len(components) == 0 {
		return nil, fmt.Errorf("expected components, none found")
	}

	return &ast.Declaration{
		Ident:      ident.Value,
		Components: components,
	}, nil
}

func (p *Parser) componentValue() (*scanner.Token, error) {
	// TODO(ttacon): this can't return string ({}, (), [], func too)
	// TODO(ttacon): this whole function is a dirty, dirty hack... :(
	t := p.nextNonWhitespaceToken()
	if t.Type == scanner.TokenChar { // it's a '.'
		// TODO(ttacon): this should just be Next() and check the type is
		// an identifier
		var next *scanner.Token
		var err error
		switch t.Value {
		case ".":
			next = p.nextNonWhitespaceToken()
		case "{": // is this valid here?
			next, err = p.componentValue() // TODO(ttacon): is this right?
			if err != nil {
				return nil, err
			}
			next.Value = t.Value + next.Value
			t = next
			next = p.nextNonWhitespaceToken()
			if next.Type != scanner.TokenChar || next.Value != "}" {
				return nil, fmt.Errorf("expected '}', got %s", next.Value)
			}
			next.Value = t.Value + next.Value
			t = next
		case "[":
			next, err = p.squareBlock()
			if err != nil {
				return nil, err
			}
			next.Value = t.Value + next.Value
			t = next
			next = p.nextNonWhitespaceToken()
			if next.Type != scanner.TokenChar || next.Value != "]" {
				return nil, fmt.Errorf("expected ']', got %s", next.Value)
			}
			next.Value = t.Value + next.Value
			t = next
		case "(":
			next, err = p.parenBlock()
			if err != nil {
				return nil, err
			}
			next.Value = t.Value + next.Value
			t = next
			next = p.nextNonWhitespaceToken()
			if next.Type != scanner.TokenChar || next.Value != ")" {
				return nil, fmt.Errorf("expected ')', got %s", next.Value)
			}
			next.Value = t.Value + next.Value
			t = next
		}
		next.Value = t.Value + next.Value
		t = next

	}

	return t, nil
}

func (p *Parser) squareBlock() (*scanner.Token, error) {
	var t = p.nextNonWhitespaceToken()
	for t.Type != scanner.TokenError && t.Type != scanner.TokenEOF {
		fmt.Println(t)
		t = p.nextNonWhitespaceToken()
	}
	return nil, nil
}

func (p *Parser) parenBlock() (*scanner.Token, error) {
	return nil, nil
}

func (p *Parser) parseBlock(t *scanner.Token) (*ast.Block, error) {
	var vals []string
	t = p.nextNonWhitespaceToken()

	for ; !isEnd(t) && !isClosingBrace(t); t = p.nextNonWhitespaceToken() {
		if !isSpace(t) {
			vals = append(vals, t.Value)
		}
	}

	if isEnd(t) {
		return nil, fmt.Errorf("hit EOF/Error while parsing block")
	}

	return &ast.Block{
		Components: vals,
	}, nil
}

func isClosingBrace(t *scanner.Token) bool {
	return t.Type == scanner.TokenChar && t.Value == "}"
}

func isSpace(t *scanner.Token) bool {
	return t.Type == scanner.TokenS
}
