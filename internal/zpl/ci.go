package zpl

import "fmt"

// CI represents the ZPL Character Set command (^CI).
// The CI command selects the international character set to be used for the label.
type ci struct {
	// charset specifies the international character set number:
	//   0 = Single Byte Encoding - U.S.A. 1 (default)
	//   1 = Single Byte Encoding - U.S.A. 2
	//   2 = Single Byte Encoding - U.K.
	//   3 = Single Byte Encoding - Holland
	//   4 = Single Byte Encoding - Denmark/Norway
	//   5 = Single Byte Encoding - Sweden/Finland
	//   6 = Single Byte Encoding - Germany
	//   7 = Single Byte Encoding - France 1
	//   8 = Single Byte Encoding - France 2
	//   9 = Single Byte Encoding - Italy
	//  10 = Single Byte Encoding - Spain
	//  11 = Single Byte Encoding - Miscellaneous
	//  12 = Single Byte Encoding - Japan (ASCII with Yen symbol)
	//  13 = Zebra Code Page 850 (Latin)
	//  14 = Double Byte Asian Encodings
	//  15 = Shift-JIS
	//  16 = EUC-JP and EUC-CN
	//  17 = Deprecated (UCS-2 Big Endian)
	//  24 = Single Byte Asian Encoding
	//  26 = Multibyte Asian Encodings with ASCII Transparency
	//  27 = Zebra Code Page 1252 (Latin)
	//  28 = Unicode (UTF-8)
	//  29 = Unicode (UTF-16 Big-Endian)
	//  30 = Unicode (UTF-16 Little-Endian)
	//  31 = Zebra Code Page 1250 (Central/Eastern European) - G-Series only
	//  33 = Code Page 1251
	//  34 = Code Page 1253
	//  35 = Code Page 1254
	//  36 = Code Page 1255
	// for more detailed information see [zpl2 documentation](https://www.zebra.com/content/dam/support-dam/en/documentation/unrestricted/guide/software/zpl-zbi2-pg-en.pdf)
	charset uint
}

// CI creates a new Character Set command with the specified charset number.
// Valid values are 0-36 (with some gaps - see documentation).
// Note: Some character sets require specific firmware versions:
//   - 28-30 require firmware V60.14.x/V50.14.x or later
//   - 31-36 require firmware x.16.x or later
//
// Default: 0 (U.S.A. 1 Character Set)
//
// Example:
//
//	label.AddCommand(zpl2.CI(28)) // Sets UTF-8 character set
func CI(charset uint) *ci {
	return &ci{charset}
}

// ToZPL converts the CI command to its ZPL string representation.
// Returns the formatted ^CI command followed by the charset number.
//
// Example output: "^CI28" for UTF-8 character set
func (ci *ci) ToZPL() string {
	return fmt.Sprintf("^CI%d", ci.charset)
}
