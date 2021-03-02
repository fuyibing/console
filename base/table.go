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
	body   BodyInterface
	head   HeadInterface
	prefix string
	width  []int
}

// Table instance.
type TableInterface interface {
	Body() BodyInterface
	Head() HeadInterface
	Print()
	SetPrefix(string) TableInterface
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
	o.printHead()
	o.printBody()
	o.printBottom()
}

func (o *Table) SetPrefix(prefix string) TableInterface {
	o.prefix = prefix
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
				cs += SyntaxSide
			}
			// content.
			cs += " " + cell.Content(o.width[n]) + " "
			// end cell
			cs += SyntaxSide
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
	println(o.prefix + bs)
}

// Print table Head.
func (o *Table) printHead() {
	bs, ts, cs := "", "", ""
	max := len(o.head.Row().Cells())
	for n, cell := range o.head.Row().Cells() {
		// separator
		if n == 0 {
			bs += SyntaxMiddleLeft
			ts += SyntaxTopLeftCorner
			cs += SyntaxSide
		}
		// content.
		cs += " " + cell.Content(o.width[n]) + " "
		// append linear.
		for x := 0; x < o.width[n]+2; x++ {
			bs += SyntaxBottomSeparator
			ts += SyntaxTopSeparator
		}
		// end cell
		cs += SyntaxSide
		if n == (max - 1) {
			bs += SyntaxMiddleRight
			ts += SyntaxTopRightCorner
		} else {
			bs += SyntaxMiddleCrossing
			ts += SyntaxMiddleTop
		}
	}
	println(o.prefix + ts)
	println(o.prefix + cs)
	println(o.prefix + bs)
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
}
