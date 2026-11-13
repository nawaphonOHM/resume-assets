package handler

import (
	"regexp"
	"strconv"
	"strings"
)

const (
	headerControlBoundary = 0x20
	headerDeleteByte      = 0x7f
	headerIntegerBitSize  = 64
)

var headerNamePattern = regexp.MustCompile("^[!#$%&'*+-.^_`|~0-9A-Za-z]+$")
var decimalPattern = regexp.MustCompile("^[0-9]+$")

func validHeaderName(name string) bool {
	return headerNamePattern.MatchString(name)
}

func isHopByHopHeader(name string) bool {
	switch strings.ToLower(name) {
	case "connection", "keep-alive", "proxy-authenticate", "proxy-authorization", "proxy-connection",
		"te", "trailer", "transfer-encoding", "upgrade":
		return true
	default:
		return false
	}
}

func validHeaderValue(value string) bool {
	for index := 0; index < len(value); index++ {
		if !isHeaderValueByte(value[index]) {
			return false
		}
	}
	return true
}

func isHeaderValueByte(character byte) bool {
	return character == '\t' ||
		character >= headerControlBoundary && character != headerDeleteByte
}

func validContentLength(value string) bool {
	if value == "" {
		return false
	}
	_, err := strconv.ParseInt(value, decimalBase, headerIntegerBitSize)
	return err == nil && isDecimal(value)
}

func isDecimal(value string) bool {
	return decimalPattern.MatchString(value)
}
