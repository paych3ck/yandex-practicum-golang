package morse

import (
	"fmt"
	"strings"
	"unicode"
)

const (
	А  = ".-"
	Б  = "-..."
	В  = ".--"
	Г  = "--."
	Д  = "-.."
	Е  = "."
	Ж  = "...-"
	З  = "--.."
	И  = ".."
	Й  = ".---"
	К  = "-.-"
	Л  = ".-.."
	М  = "--"
	Н  = "-."
	О  = "---"
	П  = ".--."
	Р  = ".-."
	С  = "..."
	Т  = "-"
	У  = "..-"
	Ф  = "..-."
	Х  = "...."
	Ц  = "-.-."
	Ч  = "---."
	Ш  = "----"
	Щ  = "--.-"
	ЪЬ = "-..-"
	Ы  = "-.--"
	Э  = "..-.."
	Ю  = "..--"
	Я  = ".-.-"

	One   = ".----"
	Two   = "..---"
	Three = "...--"
	Four  = "....-"
	Five  = "....."
	Six   = "-...."
	Seven = "--..."
	Eight = "---.."
	Nine  = "----."
	Zero  = "-----"

	Period       = "......"
	Comma        = ".-.-.-"
	Colon        = "---..."
	QuestionMark = "..--.."
	Apostrophe   = ".----."
	Hyphen       = "-....-"
	Division     = "-..-."
	LeftBracket  = "-.--."
	RightBracket = "-.--.-"
	IvertedComma = ".-..-."
	DoubleHyphen = "-...-"
	Cross        = ".-.-."
	CommercialAt = ".--.-."

	Space = " "
)

type EncodingMap map[rune]string

const averageSize = 4.53

var DefaultMorse = EncodingMap{
	'А': А,
	'Б': Б,
	'В': В,
	'Г': Г,
	'Д': Д,
	'Е': Е,
	'Ж': Ж,
	'З': З,
	'И': И,
	'Й': Й,
	'К': К,
	'Л': Л,
	'М': М,
	'Н': Н,
	'О': О,
	'П': П,
	'Р': Р,
	'С': С,
	'Т': Т,
	'У': У,
	'Ф': Ф,
	'Х': Х,
	'Ц': Ц,
	'Ч': Ч,
	'Ш': Ш,
	'Щ': Щ,
	'Ь': ЪЬ,
	'Ы': Ы,
	'Ъ': ЪЬ,
	'Э': Э,
	'Ю': Ю,
	'Я': Я,

	'1': One,
	'2': Two,
	'3': Three,
	'4': Four,
	'5': Five,
	'6': Six,
	'7': Seven,
	'8': Eight,
	'9': Nine,
	'0': Zero,

	'.':  Period,
	',':  Comma,
	':':  Colon,
	'?':  QuestionMark,
	'\'': Apostrophe,
	'-':  Hyphen,
	'/':  Division,
	'(':  LeftBracket,
	')':  RightBracket,
	'"':  IvertedComma,
}

var reverseDefaultMorse = reverseEncodingMap(DefaultMorse)

type ErrNoEncoding struct{ Text string }

func (e ErrNoEncoding) Error() string {
	return fmt.Sprintf("No encoding for: %q", e.Text)
}

func RuneToMorse(r rune) string {
	r = unicode.ToUpper(r)
	return DefaultMorse[r]
}

func MorseToRune(morse string) rune {
	return reverseDefaultMorse[morse]
}

func reverseEncodingMap(encoding EncodingMap) map[string]rune {
	ret := make(map[string]rune, len(encoding))

	for k, v := range encoding {
		ret[v] = k
	}

	return ret
}

func (c Converter) ToText(morse string) string {
	out := make([]rune, 0, int(float64(len(morse))/averageSize))

	words := strings.Split(morse, c.charSeparator+Space+c.charSeparator)
	for _, word := range words {
		chars := strings.Split(word, c.charSeparator)

		for _, ch := range chars {
			text, ok := c.morseToRune[ch]
			if !ok {
				hand := []rune(c.Handling(ErrNoEncoding{string(ch)}))
				out = append(out, hand...)

				if len(hand) != 0 {
					out = append(out, []rune(c.charSeparator)...)
				}
				continue
			}
			out = append(out, text)
		}

		out = append(out, ' ')
	}

	if !c.trailingSeparator && len(out) >= len(c.charSeparator) {
		out = out[:len(out)-len(c.charSeparator)]
	}

	return string(out)
}

type ConverterOption func(Converter) Converter

type ErrorHandler func(error) string

type Converter struct {
	runeToMorse       map[rune]string
	morseToRune       map[string]rune
	charSeparator     string
	wordSeparator     string
	convertToUpper    bool
	trailingSeparator bool

	Handling ErrorHandler
}

func IgnoreHandler(error) string {
	return ""
}

func NewConverter(convertingMap EncodingMap, options ...ConverterOption) Converter {
	if convertingMap == nil {
		panic("Using a nil EncodingMap")
	}

	morseToRune := reverseEncodingMap(convertingMap)

	c := Converter{
		runeToMorse:       convertingMap,
		morseToRune:       morseToRune,
		charSeparator:     " ",
		wordSeparator:     "",
		convertToUpper:    false,
		trailingSeparator: false,

		Handling: IgnoreHandler,
	}

	for _, opt := range options {
		c = opt(c)
	}

	if c.wordSeparator == "" {
		sp, ok := c.runeToMorse[' ']
		if !ok {
			sp = Space
		}

		c.wordSeparator = c.charSeparator + sp + c.charSeparator
	}

	return c
}

func (c Converter) ToMorse(text string) string {
	out := make([]rune, 0, int(float64(len(text))*averageSize))

	for _, ch := range text {
		if c.convertToUpper {
			ch = unicode.ToUpper(ch)
		}

		if _, ok := c.runeToMorse[ch]; !ok {
			hand := []rune(c.Handling(ErrNoEncoding{string(ch)}))
			out = append(out, hand...)

			if len(hand) != 0 {
				out = append(out, []rune(c.charSeparator)...)
			}
			continue
		}

		out = append(out, []rune(c.runeToMorse[ch])...)
		out = append(out, []rune(c.charSeparator)...)
	}

	if !c.trailingSeparator && len(out) >= len(c.charSeparator) {
		out = out[:len(out)-len(c.charSeparator)]
	}

	return string(out)
}

var DefaultConverter = NewConverter(
	DefaultMorse,

	WithCharSeparator(" "),
	WithWordSeparator("   "),
	WithLowercaseHandling(true),
	WithHandler(IgnoreHandler),
	WithTrailingSeparator(false),
)

func ToText(morse string) string {
	return DefaultConverter.ToText(morse)
}

func ToMorse(text string) string {
	return DefaultConverter.ToMorse(text)
}

func WithHandler(handler ErrorHandler) ConverterOption {
	return func(c Converter) Converter {
		c.Handling = handler
		return c
	}
}

func WithLowercaseHandling(lowercaseHandling bool) ConverterOption {
	return func(c Converter) Converter {
		c.convertToUpper = lowercaseHandling
		return c
	}
}

func WithTrailingSeparator(trailingSpace bool) ConverterOption {
	return func(c Converter) Converter {
		c.trailingSeparator = trailingSpace
		return c
	}
}

func WithCharSeparator(charSeparator string) ConverterOption {
	return func(c Converter) Converter {
		c.charSeparator = charSeparator
		return c
	}
}

func WithWordSeparator(wordSeparator string) ConverterOption {
	return func(c Converter) Converter {
		c.wordSeparator = wordSeparator
		return c
	}
}
