package table

import (
	"strings"
	"testing"

	"github.com/fatih/color"
	"github.com/jedib0t/go-pretty/text"
	"github.com/stretchr/testify/assert"
	"unicode/utf8"
)

var (
	testAlign           = []text.Align{text.AlignDefault, text.AlignLeft, text.AlignLeft, text.AlignRight}
	testCaption         = "test-caption"
	testColor           = text.Colors{color.FgGreen}
	testColorBoW        = text.Colors{color.FgBlack, color.BgWhite}
	testColorHiRedBold  = text.Colors{color.FgHiRed, color.Bold}
	testColorHiBlueBold = text.Colors{color.FgHiBlue, color.Bold}
	testColorWoB        = text.Colors{color.FgWhite, color.BgBlack}
	testColors          = []text.Colors{testColor, testColor, testColor, testColor, {color.FgCyan}}
	testColorsFooter    = []text.Colors{{}, {}, testColorHiBlueBold, testColorHiBlueBold}
	testColorsHeader    = []text.Colors{testColorHiRedBold, testColorHiRedBold, testColorHiRedBold, testColorHiRedBold}
	testCSSClass        = "test-css-class"
	testFooter          = Row{"", "", "Total", 10000}
	testHeader          = Row{"#", "First Name", "Last Name", "Salary"}
	testRows            = []Row{
		{1, "Arya", "Stark", 3000},
		{20, "Jon", "Snow", 2000, "You know nothing, Jon Snow!"},
		{300, "Tyrion", "Lannister", 5000},
	}
	testRowMultiLine = Row{0, "Winter", "Is", 0, "Coming.\nThe North Remembers!"}
	testRowNewLines  = Row{0, "Valar", "Morghulis", 0, "Faceless\nMen"}
	testRowPipes     = Row{0, "Valar", "Morghulis", 0, "Faceless|Men"}
	testRowTabs      = Row{0, "Valar", "Morghulis", 0, "Faceless\tMen"}
)

type myMockOutputMirror struct {
	mirroredOutput string
}

func (t *myMockOutputMirror) Write(p []byte) (n int, err error) {
	t.mirroredOutput += string(p)
	return len(p), nil
}

func TestNewWriter(t *testing.T) {
	tw := NewWriter()
	assert.Nil(t, tw.Style())

	tw.SetStyle(StyleBold)
	assert.NotNil(t, tw.Style())
	assert.Equal(t, StyleBold, *tw.Style())
}

func TestTable_AppendFooter(t *testing.T) {
	table := Table{}
	assert.Equal(t, 0, len(table.rowsFooter))

	table.AppendFooter([]interface{}{})
	assert.Equal(t, 0, table.Length())
	assert.Equal(t, 1, len(table.rowsFooter))
	assert.Equal(t, 0, len(table.rowsHeader))

	table.AppendFooter([]interface{}{})
	assert.Equal(t, 0, table.Length())
	assert.Equal(t, 2, len(table.rowsFooter))
	assert.Equal(t, 0, len(table.rowsHeader))
}

func TestTable_AppendHeader(t *testing.T) {
	table := Table{}
	assert.Equal(t, 0, len(table.rowsHeader))

	table.AppendHeader([]interface{}{})
	assert.Equal(t, 0, table.Length())
	assert.Equal(t, 0, len(table.rowsFooter))
	assert.Equal(t, 1, len(table.rowsHeader))

	table.AppendHeader([]interface{}{})
	assert.Equal(t, 0, table.Length())
	assert.Equal(t, 0, len(table.rowsFooter))
	assert.Equal(t, 2, len(table.rowsHeader))
}

func TestTable_AppendRow(t *testing.T) {
	table := Table{}
	assert.Equal(t, 0, table.Length())

	table.AppendRow([]interface{}{})
	assert.Equal(t, 1, table.Length())
	assert.Equal(t, 0, len(table.rowsFooter))
	assert.Equal(t, 0, len(table.rowsHeader))

	table.AppendRow([]interface{}{})
	assert.Equal(t, 2, table.Length())
	assert.Equal(t, 0, len(table.rowsFooter))
	assert.Equal(t, 0, len(table.rowsHeader))
}

func TestTable_AppendRows(t *testing.T) {
	table := Table{}
	assert.Equal(t, 0, table.Length())

	table.AppendRows([]Row{{}})
	assert.Equal(t, 1, table.Length())
	assert.Equal(t, 0, len(table.rowsFooter))
	assert.Equal(t, 0, len(table.rowsHeader))

	table.AppendRows([]Row{{}})
	assert.Equal(t, 2, table.Length())
	assert.Equal(t, 0, len(table.rowsFooter))
	assert.Equal(t, 0, len(table.rowsHeader))
}

func TestTable_Length(t *testing.T) {
	table := Table{}
	assert.Zero(t, table.Length())

	table.AppendRow(testRows[0])
	assert.Equal(t, 1, table.Length())
	table.AppendRow(testRows[1])
	assert.Equal(t, 2, table.Length())

	table.AppendHeader(testHeader)
	assert.Equal(t, 2, table.Length())
}

func TestTable_SetAlign(t *testing.T) {
	table := Table{}
	assert.Nil(t, table.align)

	table.SetAlign([]text.Align{})
	assert.NotNil(t, table.align)

	table.AppendRows(testRows)
	table.AppendRow(testRowMultiLine)
	table.SetAlign([]text.Align{text.AlignDefault, text.AlignLeft, text.AlignLeft, text.AlignRight, text.AlignRight})

	expectedOut := `+-----+--------+-----------+------+-----------------------------+
|   1 | Arya   | Stark     | 3000 |                             |
|  20 | Jon    | Snow      | 2000 | You know nothing, Jon Snow! |
| 300 | Tyrion | Lannister | 5000 |                             |
|   0 | Winter | Is        |    0 |                     Coming. |
|     |        |           |      |        The North Remembers! |
+-----+--------+-----------+------+-----------------------------+`
	assert.Equal(t, expectedOut, table.Render())
}

func TestTable_SetAllowedColumnLengths(t *testing.T) {
	table := Table{}
	table.AppendRows(testRows)
	table.SetStyle(styleTest)

	expectedOut := `(-----^--------^-----------^------^-----------------------------)
[<  1>|<Arya  >|<Stark    >|<3000>|<                           >]
[< 20>|<Jon   >|<Snow     >|<2000>|<You know nothing, Jon Snow!>]
[<300>|<Tyrion>|<Lannister>|<5000>|<                           >]
\-----v--------v-----------v------v-----------------------------/`
	assert.Empty(t, table.allowedColumnLengths)
	assert.Equal(t, expectedOut, table.Render())

	table.SetAllowedColumnLengths([]int{0, 1, 2, 3, 7})
	expectedOut = `(-----^---^----^-----^---------)
[<  1>|<A>|<St>|<300>|<       >]
[<   >|<r>|<ar>|<  0>|<       >]
[<   >|<y>|<k >|<   >|<       >]
[<   >|<a>|<  >|<   >|<       >]
[< 20>|<J>|<Sn>|<200>|<You kno>]
[<   >|<o>|<ow>|<  0>|<w nothi>]
[<   >|<n>|<  >|<   >|<ng, Jon>]
[<   >|< >|<  >|<   >|< Snow! >]
[<300>|<T>|<La>|<500>|<       >]
[<   >|<y>|<nn>|<  0>|<       >]
[<   >|<r>|<is>|<   >|<       >]
[<   >|<i>|<te>|<   >|<       >]
[<   >|<o>|<r >|<   >|<       >]
[<   >|<n>|<  >|<   >|<       >]
\-----v---v----v-----v---------/`
	assert.Equal(t, []int{0, 1, 2, 3, 7}, table.allowedColumnLengths)
	assert.Equal(t, expectedOut, table.Render())
}

func TestTable_SetAllowedRowLength(t *testing.T) {
	table := Table{}
	table.AppendRows(testRows)
	table.SetStyle(styleTest)

	expectedOutWithNoRowLimit := `(-----^--------^-----------^------^-----------------------------)
[<  1>|<Arya  >|<Stark    >|<3000>|<                           >]
[< 20>|<Jon   >|<Snow     >|<2000>|<You know nothing, Jon Snow!>]
[<300>|<Tyrion>|<Lannister>|<5000>|<                           >]
\-----v--------v-----------v------v-----------------------------/`
	assert.Zero(t, table.allowedRowLength)
	assert.Equal(t, expectedOutWithNoRowLimit, table.Render())

	table.SetAllowedRowLength(utf8.RuneCountInString(table.style.BoxUnfinishedRow))
	assert.Equal(t, utf8.RuneCountInString(table.style.BoxUnfinishedRow), table.allowedRowLength)
	assert.Equal(t, "", table.Render())

	table.SetAllowedRowLength(5)
	expectedOutWithRowLimit := `( ~~~
[ ~~~
[ ~~~
[ ~~~
\ ~~~`
	assert.Equal(t, 5, table.allowedRowLength)
	assert.Equal(t, expectedOutWithRowLimit, table.Render())

	table.SetAllowedRowLength(30)
	expectedOutWithRowLimit = `(-----^--------^---------- ~~~
[<  1>|<Arya  >|<Stark     ~~~
[< 20>|<Jon   >|<Snow      ~~~
[<300>|<Tyrion>|<Lannister ~~~
\-----v--------v---------- ~~~`
	assert.Equal(t, 30, table.allowedRowLength)
	assert.Equal(t, expectedOutWithRowLimit, table.Render())

	table.SetAllowedRowLength(300)
	assert.Equal(t, 300, table.allowedRowLength)
	assert.Equal(t, expectedOutWithNoRowLimit, table.Render())
}

func TestTable_SetAutoIndex(t *testing.T) {
	table := Table{}
	table.AppendRows(testRows)
	table.SetStyle(styleTest)

	expectedOut := `(-----^--------^-----------^------^-----------------------------)
[<  1>|<Arya  >|<Stark    >|<3000>|<                           >]
[< 20>|<Jon   >|<Snow     >|<2000>|<You know nothing, Jon Snow!>]
[<300>|<Tyrion>|<Lannister>|<5000>|<                           >]
\-----v--------v-----------v------v-----------------------------/`
	assert.False(t, table.autoIndex)
	assert.Equal(t, expectedOut, table.Render())

	table.SetAutoIndex(true)
	expectedOut = `(---^-----^--------^-----------^------^-----------------------------)
[< >|< A >|<   B  >|<    C    >|<  D >|<             E             >]
{---+-----+--------+-----------+------+-----------------------------}
[<1>|<  1>|<Arya  >|<Stark    >|<3000>|<                           >]
[<2>|< 20>|<Jon   >|<Snow     >|<2000>|<You know nothing, Jon Snow!>]
[<3>|<300>|<Tyrion>|<Lannister>|<5000>|<                           >]
\---v-----v--------v-----------v------v-----------------------------/`
	assert.True(t, table.autoIndex)
	assert.Equal(t, expectedOut, table.Render())

	table.AppendHeader(testHeader)
	expectedOut = `(-----^------------^-----------^--------^-----------------------------)
[<  #>|<FIRST NAME>|<LAST NAME>|<SALARY>|<                           >]
{-----+------------+-----------+--------+-----------------------------}
[<  1>|<Arya      >|<Stark    >|<  3000>|<                           >]
[< 20>|<Jon       >|<Snow     >|<  2000>|<You know nothing, Jon Snow!>]
[<300>|<Tyrion    >|<Lannister>|<  5000>|<                           >]
\-----v------------v-----------v--------v-----------------------------/`
	assert.True(t, table.autoIndex)
	assert.Equal(t, expectedOut, table.Render())
}

func TestTable_SetCaption(t *testing.T) {
	table := Table{}
	assert.Empty(t, table.caption)

	table.SetCaption(testCaption)
	assert.NotEmpty(t, table.caption)
	assert.Equal(t, testCaption, table.caption)
}

func TestTable_SetColors(t *testing.T) {
	table := Table{}
	assert.Empty(t, table.colors)
	assert.Empty(t, table.colorsFooter)
	assert.Empty(t, table.colorsHeader)

	table.SetColors([]text.Colors{testColorWoB, testColorBoW})
	assert.NotEmpty(t, table.colors)
	assert.Empty(t, table.colorsFooter)
	assert.Empty(t, table.colorsHeader)
	assert.Equal(t, 2, len(table.colors))
}

func TestTable_SetColorsFooter(t *testing.T) {
	table := Table{}
	assert.Empty(t, table.colors)
	assert.Empty(t, table.colorsFooter)
	assert.Empty(t, table.colorsHeader)

	table.SetColorsFooter([]text.Colors{testColorWoB, testColorBoW})
	assert.Empty(t, table.colors)
	assert.NotEmpty(t, table.colorsFooter)
	assert.Empty(t, table.colorsHeader)
	assert.Equal(t, 2, len(table.colorsFooter))
}

func TestTable_SetColorsHeader(t *testing.T) {
	table := Table{}
	assert.Empty(t, table.colors)
	assert.Empty(t, table.colorsFooter)
	assert.Empty(t, table.colorsHeader)

	table.SetColorsHeader([]text.Colors{testColorWoB, testColorBoW})
	assert.Empty(t, table.colors)
	assert.Empty(t, table.colorsFooter)
	assert.NotEmpty(t, table.colorsHeader)
	assert.Equal(t, 2, len(table.colorsHeader))
}

func TestTable_SetHTMLCSSClass(t *testing.T) {
	table := Table{}
	table.AppendRow(testRows[0])
	expectedHTML := `<table class="` + DefaultHTMLCSSClass + `">
  <tbody>
  <tr>
    <td align="right">1</td>
    <td>Arya</td>
    <td>Stark</td>
    <td align="right">3000</td>
  </tr>
  </tbody>
</table>`
	assert.Equal(t, "", table.htmlCSSClass)
	assert.Equal(t, expectedHTML, table.RenderHTML())

	table.SetHTMLCSSClass(testCSSClass)
	assert.Equal(t, testCSSClass, table.htmlCSSClass)
	assert.Equal(t, strings.Replace(expectedHTML, DefaultHTMLCSSClass, testCSSClass, -1), table.RenderHTML())
}

func TestTable_SetOutputMirror(t *testing.T) {
	table := Table{}
	table.AppendRow(testRows[0])
	expectedOut := `+---+------+-------+------+
| 1 | Arya | Stark | 3000 |
+---+------+-------+------+`
	assert.Equal(t, nil, table.outputMirror)
	assert.Equal(t, expectedOut, table.Render())

	mockOutputMirror := &myMockOutputMirror{}
	table.SetOutputMirror(mockOutputMirror)
	assert.Equal(t, mockOutputMirror, table.outputMirror)
	assert.Equal(t, expectedOut, table.Render())
	assert.Equal(t, expectedOut+"\n", mockOutputMirror.mirroredOutput)
}

func TestTable_SetVAlign(t *testing.T) {
	table := Table{}
	assert.Nil(t, table.vAlign)

	table.SetVAlign([]text.VAlign{})
	assert.NotNil(t, table.vAlign)

	table.AppendRow(testRowMultiLine)
	table.SetVAlign([]text.VAlign{text.VAlignTop, text.VAlignMiddle, text.VAlignBottom, text.VAlignDefault})

	expectedOut := `+---+--------+----+---+----------------------+
| 0 | Winter |    | 0 | Coming.              |
|   |        | Is |   | The North Remembers! |
+---+--------+----+---+----------------------+`

	assert.Equal(t, expectedOut, table.Render())
}

func TestTable_SetStyle(t *testing.T) {
	table := Table{}
	assert.Nil(t, table.Style())

	table.SetStyle(StyleDefault)
	assert.NotNil(t, table.Style())
	assert.Equal(t, &StyleDefault, table.Style())
}

func TestTable_ShowBorder(t *testing.T) {
	table := Table{}
	assert.False(t, table.disableBorder)

	table.ShowBorder(false)
	assert.True(t, table.disableBorder)

	table.AppendRow(testRows[0])
	out := table.Render()

	assert.NotEmpty(t, out)
	assert.Equal(t, 0, strings.Count(out, "\n"))
	assert.Equal(t, " 1 | Arya | Stark | 3000 ", out)
}

func TestTable_ShowSeparators(t *testing.T) {
	table := Table{}
	assert.False(t, table.enableSeparators)

	table.ShowSeparators(true)
	assert.True(t, table.enableSeparators)

	table.AppendRows(testRows)

	expectedOut := `+-----+--------+-----------+------+-----------------------------+
|   1 | Arya   | Stark     | 3000 |                             |
+-----+--------+-----------+------+-----------------------------+
|  20 | Jon    | Snow      | 2000 | You know nothing, Jon Snow! |
+-----+--------+-----------+------+-----------------------------+
| 300 | Tyrion | Lannister | 5000 |                             |
+-----+--------+-----------+------+-----------------------------+`

	assert.Equal(t, expectedOut, table.Render())
}
