// (c) 2020 Gon Y Yi. <https://gonyyi.com/copyright.txt>

package alog

import "io"

// =====================================================================================================================
// INT TO []BYTE
// =====================================================================================================================

// itoa converts int to []byte
// if minLength == 0, it will print without padding 0
// due to limit on int type, 19 digit max; 18 digit is safe.
func itoa(dst *[]byte, i int, minLength int) {
	var b [20]byte
	var positiveNum bool = true
	if i < 0 {
		positiveNum = false
		i = -i // change the sign to positive
	}
	bIdx := len(b) - 1

	for i >= 10 || minLength > 1 {
		minLength--
		q := i / 10
		b[bIdx] = byte('0' + i - q*10)
		bIdx--
		i = q
	}

	b[bIdx] = byte('0' + i)
	if positiveNum == false {
		bIdx--
		b[bIdx] = '-'
	}
	*dst = append(*dst, b[bIdx:]...)
}

// =====================================================================================================================
// INT TO []BYTE
// =====================================================================================================================

// ftoa converts float64 to []byte
func ftoa(dst *[]byte, f float64, decPlace int) {
	if int(f) == 0 && f < 0 {
		*dst = append(*dst, '-')
	}
	itoa(dst, int(f), 0) // add full number first

	if decPlace > 0 {
		// if decPlace == 3, multiplier will be 1000
		// get nth power
		var multiplier int = 1
		for i := decPlace; i > 0; i-- {
			multiplier = multiplier * 10
		}
		*dst = append(*dst, '.')
		tmp := int((f - float64(int(f))) * float64(multiplier))
		if f > 0 { // 2nd num shouldn't include decimala
			itoa(dst, tmp, decPlace)
		} else {
			itoa(dst, -tmp, decPlace)
		}
	}
}

// =====================================================================================================================
// TINY BUFFER
// =====================================================================================================================
type tinyBuffer []byte

func (cb *tinyBuffer) Write(p []byte) (n int, err error) {
	*cb = append(*cb, p...)
	return len(p), nil
}

// =====================================================================================================================
// UTILS
// =====================================================================================================================
type devNull int

var Discard io.Writer = devNull(0)

func (devNull) Write(p []byte) (int, error) {
	return 0, nil
}
