package ini4go

import (
    "strings"
    "errors"
    "sort"
    "github.com/cznic/sortutil"
)


const (
    HEADER_START rune = '['
    HEADER_END rune = ']'
)

var (
    max_name_length int = 1024
    invalid_name_chars = []rune{HEADER_START, HEADER_END}

    default_option_whitespace = []rune{' ', '\t' }
    default_option_delimiters = []rune{'=',':'}
    default_comment_delimiters = []rune{';', '#'}
    default_case_sensitivity = false
    default_option_format = "%s %s %s"
)

type Section struct {
    name                string
    isCaseSensitive     bool
    lines               []string
    options             map[string]Option
    optionDelims        []rune
    optionDelimsSorted  []rune
    commentDelims       []rune
    commentDelimsSorted []rune
}

func NewSection(section *Section) (*Section, error) {
    if ! section.validName() {
        return nil, errors.New("Invalid name: " + section.name)
    }

    section.optionDelims = default_option_delimiters
    section.optionDelimsSorted = make([]rune,len(section.optionDelims))
    copy(section.optionDelimsSorted, section.optionDelims)

    if len(section.commentDelims) == 0 {
        section.commentDelims = default_comment_delimiters
    }
    section.commentDelimsSorted = make([]rune,len(section.commentDelims))
    copy(section.commentDelimsSorted, section.commentDelims)

    sort.Sort(sortutil.RuneSlice(section.optionDelimsSorted))
    sort.Sort(sortutil.RuneSlice(section.commentDelimsSorted))

    return section, nil
}

/**
 * Checks a string for validity as a section name. It can't contain the
 * characters '[' and ']'. An empty string or one consisting only of
 * white space isn't allowed either.
 */
func (section Section) validName() bool {
    name := section.name
    valid := false
    name = strings.TrimSpace(name)

    if len(name) > 0 {
        for _, c := range invalid_name_chars {
            if strings.ContainsAny(name, string(c)) {
                valid = false
                break
            }
            valid = true
        }
    }

    return valid
}

/**
* Normalizes an arbitrary string for use as an option name, ie makes
* it lower-case (provided this section isn't case-sensitive) and trims
* leading and trailing white space.
*/
func (section Section) normOption(name string) string {
    if section.isCaseSensitive {
        name = strings.ToLower(name)
    }

    return strings.TrimSpace(name)
}

/**
 * Returns the bracketed header of this section as appearing in an
 * actual INI file.
 */
func (section Section)  header() string {

    return string(HEADER_START) + section.name + string(HEADER_END)
}

/**
* Returns an actual Option instance.
*/
func (section Section) getOption(optionName string) Option {

    return section.options[optionName]
}

/**
* Returns an actual Option instance.
*/
func (section Section) setOptionFormatString(formatString string)  {

    //section.options[optionName]
}