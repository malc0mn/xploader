package xploader

// init populates the UnicodeToCP437 map from CP437ToUnicode and manually adds any special cases not covered by default
// reverse mapping (e.g. NBSP).
func init() {
	for code, r := range CP437ToUnicode {
		UnicodeToCP437[r] = code
	}

	// Explicitly ensure encoding fidelity
	UnicodeToCP437['\u00A0'] = 255 // NBSP -> slot 255
}

// CP437ToUnicode maps Code Page 437 (CP437) character codes (0–255) to their corresponding Unicode runes, based on
// REXPaint's default font. If a custom font is used, the visual output may not match these mappings.
//
// This is necessary because REXPaint uses CP437 code points internally, but most modern tools and text encodings expect
// Unicode.
//
// The map covers all 256 characters from CP437, using explicit \uXXXX escapes to preserve control characters and
// ambiguous glyphs. Comments alongside each entry show the visual rune for clarity.
//
// REXPaint font customization note:
// Glyphs 254 and 255 are replaced by REXPaint to what it calls "radio boxes" in its default font, even though CP437
// defines them as '■' and NBSP. These values are preserved in this map as visual approximations.
var CP437ToUnicode = map[int32]rune{
	0:   '\u0000', // Null, rendered as a visible space (same as code 32) by this library
	1:   '\u263A', // ☺
	2:   '\u263B', // ☻
	3:   '\u2665', // ♥
	4:   '\u2666', // ♦
	5:   '\u2663', // ♣
	6:   '\u2660', // ♠
	7:   '\u2022', // •
	8:   '\u25D8', // ◘
	9:   '\u25CB', // ○
	10:  '\u25D9', // ◙
	11:  '\u2642', // ♂
	12:  '\u2640', // ♀
	13:  '\u266A', // ♪
	14:  '\u266B', // ♫
	15:  '\u263C', // ☼
	16:  '\u25BA', // ►
	17:  '\u25C4', // ◄
	18:  '\u2195', // ↕
	19:  '\u203C', // ‼
	20:  '\u00B6', // ¶
	21:  '\u00A7', // §
	22:  '\u25AC', // ▬
	23:  '\u21A8', // ↨
	24:  '\u2191', // ↑
	25:  '\u2193', // ↓
	26:  '\u2192', // →
	27:  '\u2190', // ←
	28:  '\u221F', // ∟
	29:  '\u2194', // ↔
	30:  '\u25B2', // ▲
	31:  '\u25BC', // ▼
	32:  '\u0020', //
	33:  '\u0021', // !
	34:  '\u0022', // "
	35:  '\u0023', // #
	36:  '\u0024', // $
	37:  '\u0025', // %
	38:  '\u0026', // &
	39:  '\u0027', // '
	40:  '\u0028', // (
	41:  '\u0029', // )
	42:  '\u002A', // *
	43:  '\u002B', // +
	44:  '\u002C', // ,
	45:  '\u002D', // -
	46:  '\u002E', // .
	47:  '\u002F', // /
	48:  '\u0030', // 0
	49:  '\u0031', // 1
	50:  '\u0032', // 2
	51:  '\u0033', // 3
	52:  '\u0034', // 4
	53:  '\u0035', // 5
	54:  '\u0036', // 6
	55:  '\u0037', // 7
	56:  '\u0038', // 8
	57:  '\u0039', // 9
	58:  '\u003A', // :
	59:  '\u003B', // ;
	60:  '\u003C', // <
	61:  '\u003D', // =
	62:  '\u003E', // >
	63:  '\u003F', // ?
	64:  '\u0040', // @
	65:  '\u0041', // A
	66:  '\u0042', // B
	67:  '\u0043', // C
	68:  '\u0044', // D
	69:  '\u0045', // E
	70:  '\u0046', // F
	71:  '\u0047', // G
	72:  '\u0048', // H
	73:  '\u0049', // I
	74:  '\u004A', // J
	75:  '\u004B', // K
	76:  '\u004C', // L
	77:  '\u004D', // M
	78:  '\u004E', // N
	79:  '\u004F', // O
	80:  '\u0050', // P
	81:  '\u0051', // Q
	82:  '\u0052', // R
	83:  '\u0053', // S
	84:  '\u0054', // T
	85:  '\u0055', // U
	86:  '\u0056', // V
	87:  '\u0057', // W
	88:  '\u0058', // X
	89:  '\u0059', // Y
	90:  '\u005A', // Z
	91:  '\u005B', // [
	92:  '\u005C', // \
	93:  '\u005D', // ]
	94:  '\u005E', // ^
	95:  '\u005F', // _
	96:  '\u0060', // `
	97:  '\u0061', // a
	98:  '\u0062', // b
	99:  '\u0063', // c
	100: '\u0064', // d
	101: '\u0065', // e
	102: '\u0066', // f
	103: '\u0067', // g
	104: '\u0068', // h
	105: '\u0069', // i
	106: '\u006A', // j
	107: '\u006B', // k
	108: '\u006C', // l
	109: '\u006D', // m
	110: '\u006E', // n
	111: '\u006F', // o
	112: '\u0070', // p
	113: '\u0071', // q
	114: '\u0072', // r
	115: '\u0073', // s
	116: '\u0074', // t
	117: '\u0075', // u
	118: '\u0076', // v
	119: '\u0077', // w
	120: '\u0078', // x
	121: '\u0079', // y
	122: '\u007A', // z
	123: '\u007B', // {
	124: '\u007C', // |
	125: '\u007D', // }
	126: '\u007E', // ~
	127: '\u2302', // ⌂
	128: '\u00C7', // Ç
	129: '\u00FC', // ü
	130: '\u00E9', // é
	131: '\u00E2', // â
	132: '\u00E4', // ä
	133: '\u00E0', // à
	134: '\u00E5', // å
	135: '\u00E7', // ç
	136: '\u00EA', // ê
	137: '\u00EB', // ë
	138: '\u00E8', // è
	139: '\u00EF', // ï
	140: '\u00EE', // î
	141: '\u00EC', // ì
	142: '\u00C4', // Ä
	143: '\u00C5', // Å
	144: '\u00C9', // É
	145: '\u00E6', // æ
	146: '\u00C6', // Æ
	147: '\u00F4', // ô
	148: '\u00F6', // ö
	149: '\u00F2', // ò
	150: '\u00FB', // û
	151: '\u00F9', // ù
	152: '\u00FF', // ÿ
	153: '\u00D6', // Ö
	154: '\u00DC', // Ü
	155: '\u00A2', // ¢
	156: '\u00A3', // £
	157: '\u00A5', // ¥
	158: '\u20A7', // ₧
	159: '\u0192', // ƒ
	160: '\u00E1', // á
	161: '\u00ED', // í
	162: '\u00F3', // ó
	163: '\u00FA', // ú
	164: '\u00F1', // ñ
	165: '\u00D1', // Ñ
	166: '\u00AA', // ª
	167: '\u00BA', // º
	168: '\u00BF', // ¿
	169: '\u2310', // ⌐
	170: '\u00AC', // ¬
	171: '\u00BD', // ½
	172: '\u00BC', // ¼
	173: '\u00A1', // ¡
	174: '\u00AB', // «
	175: '\u00BB', // »
	176: '\u2591', // ░
	177: '\u2592', // ▒
	178: '\u2593', // ▓
	179: '\u2502', // │
	180: '\u2524', // ┤
	181: '\u2561', // ╡
	182: '\u2562', // ╢
	183: '\u2556', // ╖
	184: '\u2555', // ╕
	185: '\u2563', // ╣
	186: '\u2551', // ║
	187: '\u2557', // ╗
	188: '\u255D', // ╝
	189: '\u255C', // ╜
	190: '\u255B', // ╛
	191: '\u2510', // ┐
	192: '\u2514', // └
	193: '\u2534', // ┴
	194: '\u252C', // ┬
	195: '\u251C', // ├
	196: '\u2500', // ─
	197: '\u253C', // ┼
	198: '\u255E', // ╞
	199: '\u255F', // ╟
	200: '\u255A', // ╚
	201: '\u2554', // ╔
	202: '\u2569', // ╩
	203: '\u2566', // ╦
	204: '\u2560', // ╠
	205: '\u2550', // ═
	206: '\u256C', // ╬
	207: '\u2567', // ╧
	208: '\u2568', // ╨
	209: '\u2564', // ╤
	210: '\u2565', // ╥
	211: '\u2559', // ╙
	212: '\u2558', // ╘
	213: '\u2552', // ╒
	214: '\u2553', // ╓
	215: '\u256B', // ╫
	216: '\u256A', // ╪
	217: '\u2518', // ┘
	218: '\u250C', // ┌
	219: '\u2588', // █
	220: '\u2584', // ▄
	221: '\u258C', // ▌
	222: '\u2590', // ▐
	223: '\u2580', // ▀
	224: '\u03B1', // α
	225: '\u03B2', // β
	226: '\u0393', // Γ
	227: '\u03C0', // π
	228: '\u03A3', // Σ
	229: '\u03C3', // σ
	230: '\u00B5', // µ
	231: '\u03C4', // τ
	232: '\u03A6', // Φ
	233: '\u0398', // Θ
	234: '\u03A9', // Ω
	235: '\u03B4', // δ
	236: '\u221E', // ∞
	237: '\u00F8', // ø
	238: '\u03B5', // ε
	239: '\u2229', // ∩
	240: '\u2261', // ≡
	241: '\u00B1', // ±
	242: '\u2265', // ≥
	243: '\u2264', // ≤
	244: '\u2320', // ⌠
	245: '\u2321', // ⌡
	246: '\u00F7', // ÷
	247: '\u2248', // ≈
	248: '\u00B0', // °
	249: '\u2219', // ∙
	250: '\u00B7', // ·
	251: '\u221A', // √
	252: '\u207F', // ⁿ
	253: '\u00B2', // ²
	// REXPaint changes 254 and 255 to what it calls 'radio boxes', so we pick something that 'looks right'.
	254: '\u25A0', // ■ -> luckily, this is still default CP437
	255: '\u25A1', // □ -> not default CP437, should be NBSP \u00A0 (see init() where we rectify this)
}

// UnicodeToCP437 is the inverse of CP437ToUnicode. It maps Unicode runes to their CP437 code points. Its intended use
// is writing fully REXPaint compatible .xp files.
//
// Note: If a rune does not have a CP437 equivalent, it is not included in the map. The encoder fallback behavior is to
// return the rune unchanged as an int32.
var UnicodeToCP437 = map[rune]int32{}

// CP437Decoder translates a CP437 code point to its Unicode equivalent. If the input code is not found in the
// CP437ToUnicode map, the original code is returned as a rune.
//
// This is intended for use when reading .xp files and converting them to Unicode for display or further manipulation.
func CP437Decoder(code int32) rune {
	if r, ok := CP437ToUnicode[code]; ok {
		return r
	}
	return code
}

// CP437Encoder translates a Unicode rune to its CP437 code point. If the rune has no corresponding CP437 code, its rune
// value is returned as-is.
//
// This is typically used when preparing Unicode characters for use in REXPaint-compatible output. Glyphs not found in
// CP437 may render incorrectly or not at all in REXPaint.
func CP437Encoder(r rune) int32 {
	if code, ok := UnicodeToCP437[r]; ok {
		return code
	}
	return r
}
