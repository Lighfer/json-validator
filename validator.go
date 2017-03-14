package json

import (
	"errors"
	"fmt"
	"unicode"
)

const (
	OBJ_START = '{'
	OBJ_END   = '}'
	ARR_START = '['
	ARR_END   = ']'
	SEP_COLON = ':'
	SEP_COMMA = ','

	BOOL_T = 't'
	BOOL_F = 'f'

	NULL_START = 'n'

	CONTROL_CHARACTER = 0x20
)

const (
	QUOTATION_MARK          = '"'
	REVERSE_SOLIDUS         = '\\'
	SOLIDUS                 = '/'
	BACKSPACE               = 'b'
	FORMFEED                = 'f'
	NEWLINE                 = 'n'
	CARRIAGE_RETURN         = 'r'
	HORIZONTAL_TAB          = 't'
	FOUR_HEXADECIMAL_DIGITS = 'u'
)

const (
	NUMBER_DOT   = '.'
	NUMBER_e     = 'e'
	NUMBER_E     = 'E'
	NUMBER_PLUS  = '+'
	NUMBER_MINUS = '-'
	NUMBER_ZERO  = '0'
)

var (
	ErrInvalidJSON   = errors.New("invalid json format")
	ErrUnexpectedEOF = errors.New("unexpected end of JSON")
	ErrStringEscape  = errors.New("get an invalid escape character")
)

type JSON struct {
	jsonBytes   []byte
	position    uint
	maxPosition uint
}

func (j *JSON) len() int {
	return len(j.jsonBytes)
}

func (j *JSON) validateLen(x uint) {
	if j.maxPosition <= j.position {
		fmt.Println("asd")
		panic(ErrJSON{
			err:  ErrUnexpectedEOF,
			part: getPartOfJSON(j),
		})
	}
}

func (j *JSON) moveX(x uint) *JSON {
	if x == 0 {
		return j
	}

	j.validateLen(x)

	j.jsonBytes = j.jsonBytes[x:]
	j.position += x
	return j
}

func (j *JSON) moveOne() *JSON {
	return j.moveX(1)
}

func (j *JSON) byteX(x uint) byte {
	j.validateLen(x)

	return j.jsonBytes[x]
}

func (j *JSON) firstByte() byte {
	return j.byteX(0)
}

type ErrJSON struct {
	err        error
	additional string
	part       string
}

func (e ErrJSON) Error() string {
	return e.String()
}

func (e ErrJSON) String() string {
	return fmt.Sprintf("error:\n\t%s\nadditional:\n\t%s\n"+
		"occur at:\n\t %s\n", e.err, e.additional, e.part)
}

func Expect(b byte, data *JSON) {
	if data.firstByte() != b {
		panic(ErrJSON{
			err:        ErrInvalidJSON,
			additional: fmt.Sprintf("expect character: %c", b),
			part:       getPartOfJSON(data),
		})
	}
	TrimLeftSpace(data.moveOne())
	return
}

func Validate(jsonStr string) (err error) {
	defer func() {
		if e := recover(); e != nil {
			if e, ok := e.(error); ok {
				err = e.(error)
			} else {
				panic(e)
			}
		}
	}()

	data := &JSON{[]byte(jsonStr), 0, uint(len(jsonStr))}

	TrimLeftSpace(data)
	if data.firstByte() == OBJ_START {
		ValidateObj(data)

		if TrimLeftSpace(data).len() == 0 {
			return nil
		}
	} else if data.firstByte() == ARR_START {
		ValidateArr(data)

		if TrimLeftSpace(data).len() == 0 {
			return nil
		}
	}

	return ErrJSON{
		err:        ErrInvalidJSON,
		additional: "extra characters after parsing",
		part:       getPartOfJSON(data),
	}
}

func ValidateObj(data *JSON) {
	Expect(OBJ_START, data)

	if TrimLeftSpace(data).firstByte() == OBJ_END {
		data.moveOne()
		return
	}

	for {
		ValidateStr(TrimLeftSpace(data))

		Expect(SEP_COLON, TrimLeftSpace(data))

		ValidateValue(TrimLeftSpace(data))

		TrimLeftSpace(data)

		if data.firstByte() == SEP_COMMA {
			data.moveOne()
		} else if data.firstByte() == OBJ_END {
			data.moveOne()
			return
		} else {
			panic(ErrJSON{
				err:        ErrInvalidJSON,
				additional: `expect any one of the following characters: ','  '}'`,
				part:       getPartOfJSON(data),
			})
		}
	}
}

func ValidateArr(data *JSON) {
	Expect(ARR_START, data)

	if TrimLeftSpace(data).firstByte() == ARR_END {
		data.moveOne()
		return
	}

	for {
		ValidateValue(TrimLeftSpace(data))

		TrimLeftSpace(data)
		if data.firstByte() == SEP_COMMA {
			data.moveOne()
		} else if data.firstByte() == ARR_END {
			data.moveOne()
			return
		} else {
			panic(ErrJSON{
				err:        ErrInvalidJSON,
				additional: `expect any one of the following characters: ','  ']'`,
				part:       getPartOfJSON(data),
			})
		}
	}
}

func ValidateStr(data *JSON) {
	Expect(QUOTATION_MARK, data)

	var needEsc bool

RE_VALID:
	for idx, r := range data.jsonBytes {
		if needEsc {
			ValidateEsc(data.moveX(uint(idx)))
			needEsc = false
			goto RE_VALID
		}

		switch {
		case r == QUOTATION_MARK:
			data.moveX(uint(idx + 1))
			return
		case r == REVERSE_SOLIDUS:
			needEsc = true
		case r < CONTROL_CHARACTER:
			panic(ErrJSON{
				err:        ErrInvalidJSON,
				additional: "control characters are not allowed in string type(< 0x20)",
				part:       getPartOfJSON(data),
			})
		}
	}

	panic(ErrJSON{
		err:  ErrUnexpectedEOF,
		part: getPartOfJSON(data),
	})
}

func ValidateEsc(data *JSON) {
	switch data.firstByte() {
	case QUOTATION_MARK, REVERSE_SOLIDUS, SOLIDUS, BACKSPACE, FORMFEED,
		NEWLINE, CARRIAGE_RETURN, HORIZONTAL_TAB:
		TrimLeftSpace(data.moveOne())
		return
	case FOUR_HEXADECIMAL_DIGITS:
		for i := 1; i <= 4; i++ {
			switch {
			case data.byteX(uint(i)) >= '0' && data.byteX(uint(i)) <= '9':
			case data.byteX(uint(i)) >= 'A' && data.byteX(uint(i)) <= 'F':
			case data.byteX(uint(i)) >= 'a' && data.byteX(uint(i)) <= 'f':
			default:
				panic(ErrJSON{
					err:        ErrStringEscape,
					additional: `expect to get unicode characters consisting of \u and 4 hexadecimal digits`,
					part:       getPartOfJSON(data),
				})
			}
		}
		TrimLeftSpace(data.moveX(5))
	default:
		panic(ErrJSON{
			err:        ErrStringEscape,
			additional: `expect to get unicode characters consisting of \u and 4 hexadecimal digits, or any one of the following characters: '"'  '\'  '/'  'b'  'f'  'n'  'r'  't'`,
			part:       getPartOfJSON(data),
		})
	}
	return
}

func ValidateValue(data *JSON) {
	b := data.firstByte()
	switch {
	case b == QUOTATION_MARK:
		ValidateStr(data)
	case b == OBJ_START:
		ValidateObj(data)
	case b == ARR_START:
		ValidateArr(data)
	case b == BOOL_T:
		if data.byteX(1) != 'r' || data.byteX(2) != 'u' ||
			data.byteX(3) != 'e' {
			panic(ErrJSON{
				err:        ErrInvalidJSON,
				additional: "expect a bool value: true",
				part:       getPartOfJSON(data),
			})
		}
		data.moveX(4)
		return
	case b == BOOL_F:
		if data.byteX(1) != 'a' || data.byteX(2) != 'l' ||
			data.byteX(3) != 's' || data.byteX(4) != 'e' {
			panic(ErrJSON{
				err:        ErrInvalidJSON,
				additional: "expect a bool value: false",
				part:       getPartOfJSON(data),
			})
		}
		data.moveX(5)
		return
	case b == NULL_START:
		if data.byteX(1) != 'u' || data.byteX(2) != 'l' ||
			data.byteX(3) != 'l' {
			panic(ErrJSON{
				err:        ErrInvalidJSON,
				additional: "expect a null value: null",
				part:       getPartOfJSON(data),
			})
		}
		data.moveX(4)
		return
	case b == NUMBER_MINUS || b == NUMBER_ZERO || (b >= '1' && b <= '9'):
		ValidateNumber(data)
	default:
		panic(ErrJSON{
			err:        ErrInvalidJSON,
			additional: `expect any one of the following characters: '"'  '{'  '['  't'  'f'  'n'  '-'  '0'  '1'  '2'  '3'  '4'  '5'  '6'  '7'  '8'  '9'`,
			part:       getPartOfJSON(data),
		})
	}

	return
}

func ValidateNumber(data *JSON) {
	if data.firstByte() == NUMBER_MINUS {
		data.moveOne()
	}

	if data.firstByte() == NUMBER_ZERO {
		data.moveOne()
		// do nothing, maybe need read continuous '0' character
	} else if data.firstByte() >= '1' || data.firstByte() <= '9' {
		data.moveOne()

		if data.firstByte() >= '0' && data.firstByte() <= '9' {
			ValidateDigit(data)
		}
	} else {
		panic(ErrJSON{
			err:        ErrInvalidJSON,
			additional: `expect any one of the following characters: '-'  '0'  '1'  '2'  '3'  '4'  '5'  '6'  '7'  '8'  '9'`,
			part:       getPartOfJSON(data),
		})
	}

	if data.firstByte() == NUMBER_DOT {
		ValidateDigit(data.moveOne())
	}

	if data.firstByte() != NUMBER_e && data.firstByte() != NUMBER_E {
		return
	}

	data.moveOne()

	if data.firstByte() == NUMBER_PLUS || data.firstByte() == NUMBER_MINUS {
		data.moveOne()
	}

	ValidateDigit(data)

	return
}

func ValidateDigit(data *JSON) {
	if data.firstByte() < '0' || data.firstByte() > '9' {
		panic(ErrJSON{
			err:        ErrInvalidJSON,
			additional: "expect any one of the following characters: '0'  '1'  '2'  '3'  '4'  '5'  '6'  '7'  '8'  '9'",
			part:       getPartOfJSON(data),
		})
	}

	data.moveOne()

	for idx, b := range data.jsonBytes {
		if b < '0' || b > '9' {
			data.moveX(uint(idx))
			return
		}
	}

	panic(ErrJSON{
		err:  ErrUnexpectedEOF,
		part: getPartOfJSON(data),
	})
}

func TrimLeftSpace(data *JSON) *JSON {
	for idx, r := range data.jsonBytes {
		if !unicode.IsSpace(rune(r)) {
			return data.moveX(uint(idx))
		}
	}
	return data.moveX(uint(data.len()))
}

func getPartOfJSON(data *JSON) string {
	return string([]rune(string(data.jsonBytes[:160]))[:40])
}
