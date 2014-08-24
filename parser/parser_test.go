package parser

import (
	"reflect"
	"testing"

	"github.com/ttacon/css/ast"
	"github.com/ttacon/css/scanner"
	"github.com/ttacon/pretty"
)

type cssTest struct {
	text string
	node *ast.Stylesheet
	err  error
}

func TestParse(t *testing.T) {
	var tests = []cssTest{
		cssTest{
			text: `.cool-name { display: none;}`,
			node: &ast.Stylesheet{
				Children: []ast.Rule{
					&ast.QualifiedRule{
						Components: []*ast.ComponentValue{
							&ast.ComponentValue{Name: ".cool-name"},
						},
						Block: &ast.Block{
							DeclList: &ast.DeclarationList{
								Declarations: []*ast.Declaration{
									&ast.Declaration{
										Ident:      "display",
										Components: []string{"none"},
									},
								},
							},
						},
					},
				},
			},
		},
		cssTest{
			text: `.cool-name { display: none;}`,
			node: &ast.Stylesheet{
				Children: []ast.Rule{
					&ast.QualifiedRule{
						Components: []*ast.ComponentValue{
							&ast.ComponentValue{Name: ".cool-name"},
						},
						Block: &ast.Block{
							DeclList: &ast.DeclarationList{
								Declarations: []*ast.Declaration{
									&ast.Declaration{
										Ident:      "display",
										Components: []string{"none"},
									},
								},
							},
						},
					},
				},
			},
		},
		cssTest{
			text: `#cool-name { display: none; color: #fff;}`,
			node: &ast.Stylesheet{
				Children: []ast.Rule{
					&ast.QualifiedRule{
						Components: []*ast.ComponentValue{
							&ast.ComponentValue{Name: "#cool-name"},
						},
						Block: &ast.Block{
							DeclList: &ast.DeclarationList{
								Declarations: []*ast.Declaration{
									&ast.Declaration{
										Ident:      "display",
										Components: []string{"none"},
									},
									&ast.Declaration{
										Ident:      "color",
										Components: []string{"#fff"},
									},
								},
							},
						},
					},
				},
			},
		},
		cssTest{
			text: `#cool-name, .cool-name { display: none;}`,
			node: &ast.Stylesheet{
				Children: []ast.Rule{
					&ast.QualifiedRule{
						Components: []*ast.ComponentValue{
							&ast.ComponentValue{Name: "#cool-name"},
							&ast.ComponentValue{Name: ".cool-name"},
						},
						Block: &ast.Block{
							DeclList: &ast.DeclarationList{
								Declarations: []*ast.Declaration{
									&ast.Declaration{
										Ident:      "display",
										Components: []string{"none"},
									},
								},
							},
						},
					},
				},
			},
		},
		cssTest{
			text: `th, #cool-name, .cool-name { display: none;}`,
			node: &ast.Stylesheet{
				Children: []ast.Rule{
					&ast.QualifiedRule{
						Components: []*ast.ComponentValue{
							&ast.ComponentValue{Name: "th"},
							&ast.ComponentValue{Name: "#cool-name"},
							&ast.ComponentValue{Name: ".cool-name"},
						},
						Block: &ast.Block{
							DeclList: &ast.DeclarationList{
								Declarations: []*ast.Declaration{
									&ast.Declaration{
										Ident:      "display",
										Components: []string{"none"},
									},
								},
							},
						},
					},
				},
			},
		},
		cssTest{
			text: `table tbody, #cool-name, .cool-name { display: none;}`,
			node: &ast.Stylesheet{
				Children: []ast.Rule{
					&ast.QualifiedRule{
						Components: []*ast.ComponentValue{
							&ast.ComponentValue{Name: "table tbody"},
							&ast.ComponentValue{Name: "#cool-name"},
							&ast.ComponentValue{Name: ".cool-name"},
						},
						Block: &ast.Block{
							DeclList: &ast.DeclarationList{
								Declarations: []*ast.Declaration{
									&ast.Declaration{
										Ident:      "display",
										Components: []string{"none"},
									},
								},
							},
						},
					},
				},
			},
		},
		cssTest{
			text: `#cool-name[name="hello"] { display: none; color: #fff;}`,
			node: &ast.Stylesheet{
				Children: []ast.Rule{
					&ast.QualifiedRule{
						Components: []*ast.ComponentValue{
							&ast.ComponentValue{Name: `#cool-name[name="hello"]`},
						},
						Block: &ast.Block{
							DeclList: &ast.DeclarationList{
								Declarations: []*ast.Declaration{
									&ast.Declaration{
										Ident:      "display",
										Components: []string{"none"},
									},
									&ast.Declaration{
										Ident:      "color",
										Components: []string{"#fff"},
									},
								},
							},
						},
					},
				},
			},
		},
		cssTest{
			text: `@charset "UTF-8";`,
			node: &ast.Stylesheet{
				Children: []ast.Rule{
					&ast.AtRule{
						AtKeyword: "@charset",
						Any:       "\"UTF-8\"",
						Block:     (*ast.Block)(nil),
						JustSemi:  true,
					},
				},
			},
		},
		cssTest{
			text: `
@media print {
  body {
    font-size: 12pt;
  }
}
`,
			node: &ast.Stylesheet{
				Children: []ast.Rule{
					&ast.AtRule{
						AtKeyword: "@media",
						Any:       "print",
						QualifiedRule: &ast.QualifiedRule{
							Components: []*ast.ComponentValue{
								&ast.ComponentValue{Name: "body"},
							},
							Block: &ast.Block{
								DeclList: &ast.DeclarationList{
									Declarations: []*ast.Declaration{
										&ast.Declaration{
											Ident:      "font-size",
											Components: []string{"12pt"},
										},
									},
								},
							},
						},
						Block:    (*ast.Block)(nil),
						JustSemi: false,
					},
				},
			},
		},
		cssTest{
			text: `
@media print {
  body {
    font-size: 12pt;
  }
}

.super-cool, #it-is-awesome {
  border: 1px red solid;
}
`,
			node: &ast.Stylesheet{
				Children: []ast.Rule{
					&ast.AtRule{
						AtKeyword: "@media",
						Any:       "print",
						QualifiedRule: &ast.QualifiedRule{
							Components: []*ast.ComponentValue{
								&ast.ComponentValue{Name: "body"},
							},
							Block: &ast.Block{
								DeclList: &ast.DeclarationList{
									Declarations: []*ast.Declaration{
										&ast.Declaration{
											Ident:      "font-size",
											Components: []string{"12pt"},
										},
									},
								},
							},
						},
						Block:    (*ast.Block)(nil),
						JustSemi: false,
					},
					&ast.QualifiedRule{
						Components: []*ast.ComponentValue{
							&ast.ComponentValue{Name: ".super-cool"},
							&ast.ComponentValue{Name: "#it-is-awesome"},
						},
						Block: &ast.Block{
							DeclList: &ast.DeclarationList{
								Declarations: []*ast.Declaration{
									&ast.Declaration{
										Ident:      "border",
										Components: []string{"1px", "red", "solid"},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	for _, test := range tests {
		s := scanner.New(test.text)
		p := New(s)
		nodes, err := p.Parse()
		if err != test.err {
			t.Errorf("expected err: %v, got %q", errVal(test.err), errVal(err))
		} else if !reflect.DeepEqual(nodes, test.node) {
			t.Errorf("expected did not equal output, expected: %s\ngot: %s\n",
				pretty.Sprintf("%s", test.node),
				pretty.Sprintf("%s", nodes),
			)
		}
	}
}

func errVal(err error) string {
	if err != nil {
		return err.Error()
	}
	return "<nil>"
}
