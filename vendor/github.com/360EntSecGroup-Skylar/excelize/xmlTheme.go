package excelize

import "encoding/xml"

// xlsxTheme directly maps the theme element in the namespace
// http://schemas.openxmlformats.org/drawingml/2006/main
type xlsxTheme struct {
	ThemeElements     xlsxThemeElements     `xml:"themeElements"`
	ObjectDefaults    xlsxObjectDefaults    `xml:"objectDefaults"`
	ExtraClrSchemeLst xlsxExtraClrSchemeLst `xml:"extraClrSchemeLst"`
	ExtLst            *xlsxExtLst           `xml:"extLst"`
}

// objectDefaults element allows for the definition of default shape, line,
// and textbox formatting properties. An application can use this information
// to format a shape (or text) initially on insertion into a document.
type xlsxObjectDefaults struct {
	ObjectDefaults string `xml:",innerxml"`
}

// xlsxExtraClrSchemeLst element is a container for the list of extra color
// schemes present in a document.
type xlsxExtraClrSchemeLst struct {
	ExtraClrSchemeLst string `xml:",innerxml"`
}

// xlsxThemeElements directly maps the element defines the theme formatting
// options for the theme and is the workhorse of the theme. This is where the
// bulk of the shared theme information is contained and used by a document.
// This element contains the color scheme, font scheme, and format scheme
// elements which define the different formatting aspects of what a theme
// defines.
type xlsxThemeElements struct {
	ClrScheme  xlsxClrScheme  `xml:"clrScheme"`
	FontScheme xlsxFontScheme `xml:"fontScheme"`
	FmtScheme  xlsxFmtScheme  `xml:"fmtScheme"`
}

// xlsxClrScheme element specifies the theme color, stored in the document's
// Theme part to which the value of this theme color shall be mapped. This
// mapping enables multiple theme colors to be chained together.
type xlsxClrScheme struct {
	Name     string            `xml:"name,attr"`
	Children []xlsxClrSchemeEl `xml:",any"`
}

// xlsxFontScheme element defines the font scheme within the theme. The font
// scheme consists of a pair of major and minor fonts for which to use in a
// document. The major font corresponds well with the heading areas of a
// document, and the minor font corresponds well with the normal text or
// paragraph areas.
type xlsxFontScheme struct {
	Name      string        `xml:"name,attr"`
	MajorFont xlsxMajorFont `xml:"majorFont"`
	MinorFont xlsxMinorFont `xml:"minorFont"`
	ExtLst    *xlsxExtLst   `xml:"extLst"`
}

// xlsxMajorFont element defines the set of major fonts which are to be used
// under different languages or locals.
type xlsxMajorFont struct {
	Children []xlsxFontSchemeEl `xml:",any"`
}

// xlsxMinorFont element defines the set of minor fonts that are to be used
// under different languages or locals.
type xlsxMinorFont struct {
	Children []xlsxFontSchemeEl `xml:",any"`
}

// xlsxFmtScheme element contains the background fill styles, effect styles,
// fill styles, and line styles which define the style matrix for a theme. The
// style matrix consists of subtle, moderate, and intense fills, lines, and
// effects. The background fills are not generally thought of to directly be
// associated with the matrix, but do play a role in the style of the overall
// document. Usually, a given object chooses a single line style, a single
// fill style, and a single effect style in order to define the overall final
// look of the object.
type xlsxFmtScheme struct {
	Name           string             `xml:"name,attr"`
	FillStyleLst   xlsxFillStyleLst   `xml:"fillStyleLst"`
	LnStyleLst     xlsxLnStyleLst     `xml:"lnStyleLst"`
	EffectStyleLst xlsxEffectStyleLst `xml:"effectStyleLst"`
	BgFillStyleLst xlsxBgFillStyleLst `xml:"bgFillStyleLst"`
}

// xlsxFillStyleLst element defines a set of three fill styles that are used
// within a theme. The three fill styles are arranged in order from subtle to
// moderate to intense.
type xlsxFillStyleLst struct {
	FillStyleLst string `xml:",innerxml"`
}

// xlsxLnStyleLst element defines a list of three line styles for use within a
// theme. The three line styles are arranged in order from subtle to moderate
// to intense versions of lines. This list makes up part of the style matrix.
type xlsxLnStyleLst struct {
	LnStyleLst string `xml:",innerxml"`
}

// xlsxEffectStyleLst element defines a set of three effect styles that create
// the effect style list for a theme. The effect styles are arranged in order
// of subtle to moderate to intense.
type xlsxEffectStyleLst struct {
	EffectStyleLst string `xml:",innerxml"`
}

// xlsxBgFillStyleLst  element defines a list of background fills that are
// used within a theme. The background fills consist of three fills, arranged
// in order from subtle to moderate to intense.
type xlsxBgFillStyleLst struct {
	BgFillStyleLst string `xml:",innerxml"`
}

// xlsxClrScheme maps to children of the clrScheme element in the namespace
// http://schemas.openxmlformats.org/drawingml/2006/main - currently I have
// not checked it for completeness - it does as much as I need.
type xlsxClrSchemeEl struct {
	XMLName xml.Name
	SysClr  *xlsxSysClr    `xml:"sysClr"`
	SrgbClr *attrValString `xml:"srgbClr"`
}

// xlsxFontSchemeEl directly maps the major and minor font of the style's font
// scheme.
type xlsxFontSchemeEl struct {
	XMLName     xml.Name
	Script      string `xml:"script,attr,omitempty"`
	Typeface    string `xml:"typeface,attr"`
	Panose      string `xml:"panose,attr,omitempty"`
	PitchFamily string `xml:"pitchFamily,attr,omitempty"`
	Charset     string `xml:"charset,attr,omitempty"`
}

// xlsxSysClr element specifies a color bound to predefined operating system
// elements.
type xlsxSysClr struct {
	Val     string `xml:"val,attr"`
	LastClr string `xml:"lastClr,attr"`
}
