package ini4go

import (
    "errors"
    "strings"
    "regexp"
)

type Line interface {
    ToString() string
}

type Option struct {
    name      string
    value     string
    separator rune
    format    OptionFormat
}

func NewOption(name string, value string, separator rune, format OptionFormat) (*Option, error) {

    option := &Option{name : name, value : value, separator: separator, format:format}

    if (!option.validName()) {
        return nil, errors.New("Illegal option name:" + name)
    }

    return option, nil
}

func (option Option) Set(value string) {

    if (len(value) > 0) {
        re := regexp.MustCompile("["+ ILLEGAL_VALUE_CHARS + "]")
        value = re.ReplaceAllLiteralString(value, "")
    }

    option.value = value
}

func (option Option) ToString() string  {

    return option.format.Format(option.name, option.value, option.separator)
}

func (option Option) validName() bool {
    name := strings.TrimSpace(option.name)
    if len(name) == 0 || strings.ContainsAny(name, string(option.separator)) {
        return false
    }

    return true;
}

const (
    EXPECTED_TOKENS int = 4
    ILLEGAL_VALUE_CHARS = "\n\r"
)

type OptionFormat struct {
    formatTokens []string
}

func NewOptionFormat(formatString string) *OptionFormat {
    of := &OptionFormat{}
    of.compileFormat(formatString)

    return of
}

func (of OptionFormat) Format(name string, value string, separator rune) string {
    t := of.formatTokens;

    return string(t[0]) + name + string(t[1]) + string(separator) + string(t[2]) + value + string(t[3]);
}

func (of OptionFormat) compileFormat(formatString string) ([]string, error) {
    tokens := []string{"", "", "", "" }
    tokenCount := 0
    seenPercent := false

    token := []string{}

    for _, c := range formatString {
        switch c {
        case '%':
            if (seenPercent) {
                token = append(token, "%")
                seenPercent = false;
            } else {
                seenPercent = true;
            }
        case 's':
            if (seenPercent) {
                if tokenCount >= EXPECTED_TOKENS {
                    return nil, errors.New("Illegal option format. Too many %s placeholders.")
                }

                tokens[tokenCount] = strings.Join(token, "")
                tokenCount++
                token = []string{}
                seenPercent = false
            } else {
                token = append(token, "s")
            }
        default:
            if (seenPercent) {
                return nil, errors.New("Illegal option format. Unknown format specifier.");
            }
            token = append(token, string(c))
        }
    }

    if (tokenCount != EXPECTED_TOKENS - 1) {
        return nil, errors.New("Illegal option format. Not enough %s placeholders.");
    }

    tokens[tokenCount] = strings.Join(token, "")
    return tokens, nil;
}

