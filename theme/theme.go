package theme

import (
	"github.com/charmbracelet/glamour/ansi"
)

func MarkdownTheme() ansi.StyleConfig {
	return ansi.StyleConfig{
		Document: ansi.StyleBlock{
			StylePrimitive: ansi.StylePrimitive{
				BlockPrefix: "\n",
				BlockSuffix: "\n",
				Color:       stringPtr("#FFFFFF"),
			},
			Margin: uintPtr(1),
		},
		Heading: ansi.StyleBlock{
			StylePrimitive: ansi.StylePrimitive{
				BlockSuffix: "\n",
				Color:       stringPtr("#FFCC66"),
				Bold:        boolPtr(true),
			},
		},
		H1: ansi.StyleBlock{
			StylePrimitive: ansi.StylePrimitive{
				Prefix:          "# ",
				Suffix:          "\n\n",
				Color:           stringPtr("#FFCC66"),
				BackgroundColor: stringPtr("#2B2B2B"),
				Bold:            boolPtr(true),
			},
			Margin: uintPtr(0),
		},
		H2: ansi.StyleBlock{
			StylePrimitive: ansi.StylePrimitive{
				Prefix: "## ",
				Suffix: "\n\n",
				Color:  stringPtr("#FFCC66"),
				Bold:   boolPtr(true),
			},
		},
		Paragraph: ansi.StyleBlock{
			StylePrimitive: ansi.StylePrimitive{
				BlockSuffix: "\n",
				Color:       stringPtr("#FFFFFF"),
			},
		},
		BlockQuote: ansi.StyleBlock{
			StylePrimitive: ansi.StylePrimitive{
				Prefix: string("┃ "[0]),
			},
			Indent:      uintPtr(1),
			IndentToken: stringPtr("┃ "),
		},
		List: ansi.StyleList{
			LevelIndent: 2,
			StyleBlock: ansi.StyleBlock{
				StylePrimitive: ansi.StylePrimitive{
					BlockPrefix: "• ",
					BlockSuffix: "\n",
					Color:       stringPtr("#FFFFFF"),
				},
			},
		},
		Text: ansi.StylePrimitive{
			BlockPrefix: "",
			BlockSuffix: "\n",
			Color:       stringPtr("#FFFFFF"),
		},
		Emph: ansi.StylePrimitive{
			Color:  stringPtr("#FFDD88"),
			Italic: boolPtr(true),
		},
		Strong: ansi.StylePrimitive{
			Color: stringPtr("#FFCC66"),
			Bold:  boolPtr(true),
		},
		Link: ansi.StylePrimitive{
			Underline: boolPtr(true),
			Color:     stringPtr("#66CCFF"),
		},
		LinkText: ansi.StylePrimitive{
			Underline: boolPtr(true),
			Color:     stringPtr("#66CCFF"),
		},
		Code: ansi.StyleBlock{
			StylePrimitive: ansi.StylePrimitive{
				BlockPrefix: "`",
				BlockSuffix: "`",
				Color:       stringPtr("#B0D0FF"),
			},
		},
		CodeBlock: ansi.StyleCodeBlock{
			StyleBlock: ansi.StyleBlock{
				StylePrimitive: ansi.StylePrimitive{
					BlockPrefix: "\n",
					BlockSuffix: "\n",
					Color:       stringPtr("#B0D0FF"),
				},
				Margin: uintPtr(1),
			},
			Theme: "monokai",
		},
		Table: ansi.StyleTable{
			CenterSeparator: stringPtr("┼"),
			ColumnSeparator: stringPtr("│"),
			RowSeparator:    stringPtr("─"),
		},
		HorizontalRule: ansi.StylePrimitive{
			BlockPrefix: "\n",
			BlockSuffix: "\n",
			Color:       stringPtr("#666666"),
		},
		Strikethrough: ansi.StylePrimitive{
			CrossedOut: boolPtr(true),
			Color:      stringPtr("#888888"),
		},
	}
}

func stringPtr(s string) *string { return &s }
func boolPtr(b bool) *bool       { return &b }
func uintPtr(u uint) *uint       { return &u }

func stringPtrOrEmpty(s string) *string { return &s }
