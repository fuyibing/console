// author: wsfuyibing <websearch@163.com>
// date: 2021-02-25

package base

import (
	"fmt"
)

const (
	SyntaxTopLeftCorner     = "┌"
	SyntaxTopRightCorner    = "┐"
	SyntaxTopSeparator      = "─"
	SyntaxMiddleTop         = "┬"
	SyntaxMiddleLeft        = "├"
	SyntaxMiddleRight       = "┤"
	SyntaxMiddleBottom      = "┴"
	SyntaxMiddleCrossing    = "┼"
	SyntaxBottomLeftCorner  = "└"
	SyntaxBottomRightCorner = "┘"
	SyntaxBottomSeparator   = "─"
	SyntaxSide              = "│"
)

type Align int

const (
	AlignLeft Align = iota
	AlignRight
)

type Color int

const (
	ColorDefault Color = iota
	ColorRed
	ColorBlue
	ColorGreen
	ColorLine = 30
	ColorTitle = 30
)

// ///////////////////////////////////////////////////////////////
// Table Body Manager											//
// ///////////////////////////////////////////////////////////////

// Table Body.
type Body struct {
	rows []RowInterface
}

type BodyInterface interface {
	Add(rows ...RowInterface) BodyInterface
	Rows() []RowInterface
}

func NewBody() BodyInterface {
	return &Body{rows: make([]RowInterface, 0)}
}

func (o *Body) Add(rows ...RowInterface) BodyInterface {
	o.rows = append(o.rows, rows...)
	return o
}

func (o *Body) Rows() []RowInterface {
	return o.rows
}

// ///////////////////////////////////////////////////////////////
// Table Cell Manager											//
// ///////////////////////////////////////////////////////////////

// Cell struct.
type Cell struct {
	align Align
	color Color
	value string
	width int
}

// Cell interface.
type CellInterface interface {
	Content(width int) string
	SetAlign(align Align) CellInterface
	SetColor(color Color) CellInterface
	Width() int
}

// New cell instance.
func NewCell(value string) CellInterface {
	o := &Cell{align: AlignLeft, color: ColorDefault, value: value, width: 0}
	for _, s := range value {
		if s > 127 {
			o.width += 2
		} else {
			o.width += 1
		}
	}
	return o
}

// Convert to print content.
func (o *Cell) Content(width int) (str string) {
	if o.width == width {
		str = o.value
	} else {
		// origin string.
		// Equal to specified width.
		n := 0
		num := 0
		for _, s := range o.value {
			// char length.
			if s > 127 {
				n = 2
			} else {
				n = 1
			}
			// depth.
			if (num + n) > width {
				break
			}
			// append.
			num += n
			str += string(s)
		}
		// append space.
		if num < width {
			for m := num; m < width; m++ {
				if o.align == AlignRight {
					str = " " + str
				} else {
					str += " "
				}
			}
		}
	}
	// append color.
	switch o.color {
	case ColorRed:
		str = fmt.Sprintf("%c[%d;%dm%s%c[0m", 0x1B, 0, 31, str, 0x1B)
	case ColorGreen:
		str = fmt.Sprintf("%c[%d;%dm%s%c[0m", 0x1B, 0, 32, str, 0x1B)
	case ColorBlue:
		str = fmt.Sprintf("%c[%d;%dm%s%c[0m", 0x1B, 0, 34, str, 0x1B)
	}
	// ended.
	return
}

// Set cell align.
func (o *Cell) SetAlign(align Align) CellInterface {
	o.align = align
	return o
}

// Set cell color.
func (o *Cell) SetColor(color Color) CellInterface {
	o.color = color
	return o
}

// Return cell width.
func (o *Cell) Width() int {
	return o.width
}

// ///////////////////////////////////////////////////////////////
// Table Header Manager											//
// ///////////////////////////////////////////////////////////////

type Head struct {
	row RowInterface
}

type HeadInterface interface {
	Add(cells ...CellInterface) HeadInterface
	Row() RowInterface
}

func NewHead() HeadInterface {
	return &Head{row: NewRow()}
}

func (o *Head) Add(cells ...CellInterface) HeadInterface {
	o.row.Add(cells...)
	return o
}

func (o *Head) Row() RowInterface {
	return o.row
}

// ///////////////////////////////////////////////////////////////
// Table Row Manager											//
// ///////////////////////////////////////////////////////////////

// Row struct.
type Row struct {
	cells []CellInterface
}

// Row interface.
type RowInterface interface {
	Add(cell ...CellInterface) *Row
	Cells() []CellInterface
}

// New row interface.
func NewRow() RowInterface {
	return &Row{cells: make([]CellInterface, 0)}
}

// Add cell.
func (o *Row) Add(cells ...CellInterface) *Row {
	o.cells = append(o.cells, cells...)
	return o
}

// Return all cells.
func (o *Row) Cells() []CellInterface {
	return o.cells
}

// ///////////////////////////////////////////////////////////////
// Table Manager												//
// ///////////////////////////////////////////////////////////////

// Table struct.
type Table struct {
	body       BodyInterface
	head       HeadInterface
	prefix     string
	title      string
	width      []int
	outerWidth int
}

// Table instance.
type TableInterface interface {
	Body() BodyInterface
	Head() HeadInterface
	Print()
	SetPrefix(string) TableInterface
	SetTitle(string) TableInterface
}

// New table instance.
func NewTable() TableInterface {
	return &Table{
		body:  NewBody(),
		head:  NewHead(),
		width: make([]int, 0),
	}
}

// Get Body struct.
func (o *Table) Body() BodyInterface {
	return o.body
}

// Get Head struct.
func (o *Table) Head() HeadInterface {
	return o.head
}

// Print table.
func (o *Table) Print() {
	o.resetWidth()
	o.printTitle()
	o.printHead()
	o.printBody()
	o.printBottom()
}

// Set table prefix.
func (o *Table) SetPrefix(prefix string) TableInterface {
	o.prefix = prefix
	return o
}

// Set table title.
func (o *Table) SetTitle(title string) TableInterface {
	o.title = title
	return o
}

// Print table body.
func (o *Table) printBody() {
	// max := len(o.Head.row.Cells())
	for _, row := range o.body.Rows() {
		cs := ""
		for n, cell := range row.Cells() {
			// separator
			if n == 0 {
				cs += o.printLineColor(SyntaxSide)
			}
			// content.
			cs += " " + cell.Content(o.width[n]) + " "
			// end cell
			cs += o.printLineColor(SyntaxSide)
		}
		println(o.prefix + cs)
	}
}

// Print table bottom.
func (o *Table) printBottom() {
	bs := ""
	max := len(o.head.Row().Cells())
	for n, _ := range o.head.Row().Cells() {
		// separator
		if n == 0 {
			bs += SyntaxBottomLeftCorner
		}
		// append linear.
		for x := 0; x < o.width[n]+2; x++ {
			bs += SyntaxBottomSeparator
		}
		// end cell
		if n == (max - 1) {
			bs += SyntaxBottomRightCorner
		} else {
			bs += SyntaxMiddleBottom
		}
	}
	println(o.prefix + o.printLineColor(bs))
}

// Print table Head.
func (o *Table) printHead() {
	bs, ts, cs := "", "", ""
	max := len(o.head.Row().Cells())
	for n, cell := range o.head.Row().Cells() {
		// separator
		if n == 0 {
			cs += o.printLineColor(SyntaxSide)
			bs += SyntaxMiddleLeft
			// with head.
			if o.title == "" {
				ts += SyntaxTopLeftCorner
			} else {
				ts += SyntaxMiddleLeft
			}
		}
		// content.
		cs += " " + cell.Content(o.width[n]) + " "
		// append linear.
		for x := 0; x < o.width[n]+2; x++ {
			bs += SyntaxBottomSeparator
			ts += SyntaxTopSeparator
		}
		// end cell
		cs += o.printLineColor(SyntaxSide)
		if n == (max - 1) {
			bs += SyntaxMiddleRight
			// with head.
			if o.title == "" {
				ts += SyntaxTopRightCorner
			} else {
				ts += SyntaxMiddleRight
			}
		} else {
			ts += SyntaxMiddleTop
			bs += SyntaxMiddleCrossing
		}
	}
	println(o.prefix + o.printLineColor(ts))
	println(o.prefix + cs)
	println(o.prefix + o.printLineColor(bs))
}

// Print table title.
func (o *Table) printTitle() {
	if o.title == "" {
		return
	}
	// top separator.
	ts, _ := SyntaxTopLeftCorner, ""
	for n := 0; n < o.outerWidth; n++ {
		ts += SyntaxTopSeparator
	}
	ts += SyntaxTopRightCorner
	println(o.prefix + o.printLineColor(ts))
	// lines.
	render := func(str string, num int) {
		for n := num; n < o.outerWidth; n++ {
			str += " "
		}
		println(o.prefix + o.printLineColor(SyntaxSide) + fmt.Sprintf("%c[%d;%dm%s%c[0m", 0x1B, 0, ColorTitle, str, 0x1B) + o.printLineColor(SyntaxSide))
	}
	sn := 1
	ss := " "
	for _, c := range o.title {
		// char.
		cn := 0
		cs := string(c)
		if c > 127 {
			cn = 2
		} else {
			cn = 1
		}
		// out range.
		if (sn + cn) > (o.outerWidth - 1) {
			render(ss, sn)
			sn = cn + 1
			ss = " " + cs
			continue
		}
		// inner.
		sn += cn
		ss += cs
	}
	render(ss, sn)
}

func (o *Table) printLineColor(str string) string {
	return fmt.Sprintf("%c[%d;%dm%s%c[0m", 0x1B, 0, ColorLine, str, 0x1B)
}

// Reset cell width.
func (o *Table) resetWidth() {
	// Head width.
	max := 0
	for _, cell := range o.head.Row().Cells() {
		o.width = append(o.width, cell.Width())
		max++
	}
	// Body width
	for _, row := range o.body.Rows() {
		for n, cell := range row.Cells() {
			if n >= max {
				o.width = append(o.width, cell.Width())
				continue
			}
			if x := o.width[n]; x < cell.Width() {
				o.width[n] = cell.Width()
			}
		}
	}
	// outer width.
	for n, x := range o.width {
		if n > 0 {
			o.outerWidth += x + 3
		} else {
			o.outerWidth += x + 2
		}
	}
}
