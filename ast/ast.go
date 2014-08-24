package ast

type Stylesheet struct {
	Children []Rule
}

// TODO(ttacon): don't think we need RuleList type?
type RuleList struct {
	Rule Rule
}

type Rule interface {
}

// TODO(ttacon): CDO/CDC?

type AtRule struct {
	// TODO(ttacon): atkeyword and any should be nodes...
	AtKeyword string
	Any       string
	Block     *Block
	JustSemi  bool
}

type QualifiedRule struct {
	Components []*ComponentValue
	DeclList   *DeclarationList
}

type ComponentValue struct {
	Name string
}

type DeclarationList struct {
	Declarations []*Declaration
}

type Declaration struct {
	Ident      string
	Components []string
}

type Important struct {
}

type CurlyBlock struct {
}

type ParenBlock struct {
}

type SquareBlock struct {
}

type FunctionBlock struct {
}

type Block struct {
	// TODO(ttacon): this needs to be updated
	DeclList *DeclarationList
}
