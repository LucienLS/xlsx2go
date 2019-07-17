# Xlsx2Go
[![Build Status](https://travis-ci.org/plandem/xlsx2go.svg?branch=master)](https://travis-ci.org/plandem/xlsx2go)
[![Code Coverage](https://codecov.io/gh/plandem/xlsx2go/branch/master/graph/badge.svg)](https://codecov.io/gh/plandem/xlsx2go) 
[![Go Report Card](https://goreportcard.com/badge/github.com/plandem/xlsx2go)](https://goreportcard.com/report/github.com/plandem/xlsx2go)
[![GoDoc](https://godoc.org/github.com/plandem/xlsx2go?status.svg)](https://godoc.org/github.com/plandem/xlsx2go)
[![License](http://img.shields.io/badge/license-MIT-red.svg?style=flat)](https://raw.githubusercontent.com/plandem/xlsx2go/master/LICENSE) 
[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Fplandem%2Fxlsx2go.svg?type=shield)](https://app.fossa.io/projects/git%2Bgithub.com%2Fplandem%2Fxlsx2go?ref=badge_shield)
[![Donate](https://img.shields.io/badge/Donate-PayPal-green.svg)](https://www.paypal.me/gayvoronsky)

**Note:** Github repository was renamed from `xlsx` to `xlsx2go` to make it more easier to distinct existing xlsx libraries. Previous address will be auto redirected, package will be named as before - xlsx.   

 
```go
package main

import (
	"fmt"
	"github.com/plandem/xlsx"
	"github.com/plandem/xlsx/format/conditional"
	"github.com/plandem/xlsx/format/conditional/rule"
	"github.com/plandem/xlsx/format/styles"
	"github.com/plandem/xlsx/types"
	"github.com/plandem/xlsx/types/comment"
	"github.com/plandem/xlsx/types/hyperlink"
)

func main() {
	xl, err := xlsx.Open("./test_files/example_simple.xlsx")
	if err != nil {
		panic(err)
	}

	defer xl.Close()

	redBoldYellow := xl.AddStyles(
		styles.New(
			styles.Font.Bold,
			styles.Font.Color("#ff0000"),
			styles.Fill.Type(styles.PatternTypeSolid),
			styles.Fill.Color("#FFFF00"),
		),
	)

	//iterating via indexes
	sheet := xl.Sheet(0)
	iMaxCol, iMaxRow := sheet.Dimension()
	for iRow := 0; iRow < iMaxRow; iRow++ {
		for iCol := 0; iCol < iMaxCol; iCol++ {
			if iRow % 2 == 0 && iCol % 2 == 0 {
				cell := sheet.Cell(iCol, iRow)
				cell.SetStyles(redBoldYellow)
			}
		}
	}

	//iterating via iterators
	for rows := sheet.Rows(); rows.HasNext(); {
		_, row := rows.Next()
		
		for cells := row.Cells(); cells.HasNext(); {
			iCol, iRow, cell := cells.Next()
			if iRow % 2 == 0 && iCol % 2 == 0 {
				cell.SetStyles(redBoldYellow)
			}
		}
	}
    	
	//walk through the range's cells
	for rows := sheet.Rows(); rows.HasNext(); {
		_, row := rows.Next()
		row.Walk(func(idx, iCol, iRow int, cell *xlsx.Cell) {
			if iRow % 2 == 0 && iCol % 2 == 0 {
				cell.SetStyles(redBoldYellow)
			}
		})
 	}

	//Add hyperlink and set value same time
	sheet.CellByRef("A1").SetValueWithHyperlink("Link To Google", "http://google.com")
	
	//Add hyperlink as string
	sheet.RangeByRef("B1:C3").SetHyperlink("spam@spam.it")	

	//Add hyperlink via helper type for advanced settings
	sheet.CellByRef("A7").SetHyperlink(hyperlink.New(
		hyperlink.ToFile("./example_simple.xlsx"),
		hyperlink.ToRef("C3", "Sheet1"),
		hyperlink.Tooltip("That's a tooltip"),
		hyperlink.Display("Something to display"), //Cell still holds own value
		hyperlink.Styles(redBoldYellow),
	))

	sheet.CellByRef("A1").RemoveHyperlink()
	
	//Merged Cells
	sheet.RangeByRef("A1:C3").Merge()
	sheet.RangeByRef("A1:C3").Split()
	
	//Rich Text
	sheet.CellByRef("F10").SetText(
		"plain text", 
		styles.New(
			styles.Font.Bold,
			styles.Font.Color("#ff0000"),
		),
		" red bold text ",
		"another plain text",
	)
	
	//Add comment as string
	sheet.CellByRef("F11").SetComment("A cell's comment")

	//Add comment via helper type for advanced settings
	sheet.CellByRef("F12").SetComment(comment.New(
		comment.Author("John Doe"),
		comment.Text(
			"plain text", 
			styles.New(
				styles.Font.Bold,
				styles.Font.Color("#ff0000"),
			),
			" red bold text ",
			"another plain text",
		),
	))
	
	sheet.CellByRef("F11").RemoveComment()
	
	//Conditional formatting
	sheet.AddConditional(conditional.New(
		conditional.AddRule(
			rule.Value.Between(35.5, 78, styles.New(
				styles.Font.Bold,
				styles.Font.Color("#ff0000"),
			)),
		),
		conditional.AddRule(
			rule.TimePeriod.Last7Days(styles.New(
				styles.Font.Bold,
				styles.Font.Color("#ff0000"),
			)),
		),
		conditional.AddRule(
			rule.IconSet.Type(rule.IconSetType3Arrows),
			rule.IconSet.Value(0, "90", ">="),
			rule.IconSet.Value(1, "50", ">"),
		),
	), "A1:A10", "B2", "C1:C10")
    	
	xl.SaveAs("test1.xlsx")
}
```

# Introduction
Why another library to work with Excel XLSX in GO? 

Truth be told, developing of any library starts with some personal goals of author. Someone wants simple library to read Excel files, someone wants to create a new file, other wants to add charts. 

So what were the goals that time? It's a great pity, but I could not get a library that:

1) respects existing data/formatting - no corrupted files or lost formatting

> What if I need to open a well formatted file created with my favorite desktop application and update only one value?! I must get almost same file with just one updated value. None of existing library was able to do it. Corrupted file or lost formatting is common issue.

2) works with big files - reasonable speed and memory footprint

> Same here, someone could not open, others took forever to open with anomaly memory usage.

3) consistent and as small API as possible with enough features set to do most common tasks - learning curve means something 

> Why?! Because it's not rocket science - open/create file, create/read/update/delete sheets/rows/cols and use styles. XLSX is quite simple format to read/write and GO has quite powerful xml encoder/decoder, so the hardest part - that API. 

4) easy to read/understand source code, easy to maintain, easy to contribute - no shadow places/hacks/magic, just read and understand

> I was trying to contribute to existing libraries, but...actually it's much faster to create it from ground zero than to refactor existing and get satisfied results or fix some issues.

# Benchmarks
It was not a goal to make best of the best, but the same time it's interesting to know pros/cons. 
For some cases this library is second, for other - best, but in case of reading huge files - **the only**. 

|                | tealeg | excelize | xlsx |
|----------------|:------:|:--------:|:----:|
| RandomGet      |   1!   |     3    |   2  |
| RandomSet      |   1!   |     3    |   2  |
| RandomSetStyle |   1!   |     3    |   2  |
| ReadBigFile    |    2   |     3    |   1  |
| UpdateBigFile  |    2!! |     3    |   1  |
| ReadHugeFile   |    -   |     -    |   1  |
| UpdateHugeFile |    -   |     -    |   1  |

* ! - does not mutate information directly, so faster get/set, but slower read/write files - sometimes it can take forever to open file.
* !! - corrupted file after saving, lost styles/formatting
 
[Benchmarks report](BENCHMARKS.md) 

# Documentation and Examples
For more detailed documentation and examples you can check [godoc.org](https://godoc.org/github.com/plandem/xlsx)

# Roadmap
- [ ] sheet: copy
- [x] sheet: read as stream
- [ ] sheet: custom filters
- [ ] sheet: write as stream
- [x] merged cells: merge/split for ranges, cols, rows
- [x] hyperlinks: for cells, ranges, cols, rows
- [x] range: copy
- [x] row: copy
- [x] col: copy
- [x] cell: comments
- [ ] cell: formulas
- [x] cell: typed getter/setter for values
- [x] other: conditional formatting
- [x] other: rich texts
- [ ] other: drawing
- [ ] other: unpack package to temp folder to reduce memory usage
- [x] other: more tests

# Contribution 
- To prevent mess, sources have strict separation of markup and functionality. Document that describes OOXML is quite huge (about 6K pages), but the same time - functionality is not.
- All markup resides inside of 'ml' folders, only marshal/unmarshal is allowed here, no any functionality.
- Not every 'ml object' has related 'functional object' and vice versa.
- If you want some functionality, then wrap 'ml object' and do what you want.

### OOXML edition
XML is compliant with part 1 of the [5th edition](http://www.ecma-international.org/publications/standards/Ecma-376.htm) of the ECMA-376 Standard for Office Open XML

## License
[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Fplandem%2Fxlsx.svg?type=large)](https://app.fossa.io/projects/git%2Bgithub.com%2Fplandem%2Fxlsx?ref=badge_large)
