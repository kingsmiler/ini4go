package ini4go

import (
    "strings"
    "errors"
    "sort"
    "github.com/cznic/sortutil"
    "reflect"
)

const (
    HEADER_START rune = '['
    HEADER_END rune = ']'
)

var (
    max_name_length int = 1024
    invalid_name_chars = []rune{HEADER_START, HEADER_END}

    default_option_whitespace = []rune{' ', '\t' }
    default_option_delimiters = []rune{'=', ':'}
    default_comment_delimiters = []rune{ '#', ';' }
    default_case_sensitivity = false
    default_option_format = "%s %s %s"
)

type Section struct {
    name                string
    isCaseSensitive     bool
    lines               []Line
    options             map[string]Option
    optionDelims        []rune
    optionDelimsSorted  []rune
    commentDelims       []rune
    commentDelimsSorted []rune

    optionFormat        *OptionFormat
}

func NewSection(section *Section) (*Section, error) {
    if ! section.validName() {
        return nil, errors.New("Invalid name: " + section.name)
    }

    section.optionDelims = default_option_delimiters
    section.optionDelimsSorted = make([]rune, len(section.optionDelims))
    copy(section.optionDelimsSorted, section.optionDelims)

    if len(section.commentDelims) == 0 {
        section.commentDelims = default_comment_delimiters
    }
    section.commentDelimsSorted = make([]rune, len(section.commentDelims))
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
 * Sets the option format for this section to the given string. Options
 * in this section will be rendered according to the given format
 * string. The string must contain <code>%s</code> three times, these
 * will be replaced with the option name, the option separator and the
 * option value in this order. Literal percentage signs must be escaped
 * by preceding them with another percentage sign (i.e., <code>%%</code>
 * corresponds to one percentage sign). The default format string is
 * <code>"%s %s %s"</code>.
 *
 * Option formats may look like format strings as supported by Java 1.5,
 * but the string is in fact parsed in a custom fashion to guarantee
 * backwards compatibility. So don't try clever stuff like using format
 * conversion types other than <code>%s</code>.
 *
 */
func (section Section) setOptionFormatString(formatString string) {

    section.optionFormat = NewOptionFormat(formatString)
}

/**
* Sets the option format for this section. Options will be rendered
* according to the given format when printed.
*
*/
func (section Section) setOptionFormat(format *OptionFormat) {

    section.optionFormat = format
}

/**
* Returns the names of all options in this section.
*/

func (section Section) OptionNames() []string {
    optNames := []string{}

    for _, v := range section.lines {
        switch inst := v.(type){
        case Option:
            optNames = append(optNames, inst.name)
        }
    }

    return optNames
}

/**
 * Checks whether a given option exists in this section.
 *
 * @param name the name of the option to test for
 * @return true if the option exists in this section
 */
func (section Section) HasOption(name string) bool {

    _, exists := section.options[section.normOption(name)]

    return exists
}

// Returns an option's value.
func (section Section) GetOptionValue(name string) string {
    normed := section.normOption(name);
    if (section.HasOption(normed)) {
        return section.getOption(normed).value
    }

    return "";
}

// Sets an option's value and creates the option if it doesn't exist.
func (section Section) SetOptionValue(optionName string, value string) {
    section.SetOptionValueWithDelimiter(optionName, value, section.optionDelims[0])
}

// Sets an option's value and creates the option if it doesn't exist.
func (section Section) SetOptionValueWithDelimiter(optionName string, value string, delimiter rune) {
    normed := section.normOption(optionName);
    if section.HasOption(normed) {
        option := section.getOption(normed)
        option.value = value
    } else {
        // Option constructor might throw IllegalArgumentException
        option, _ := NewOption(normed, value, delimiter, *section.optionFormat)

        section.options[normed] = *option
        section.lines = append(section.lines, option)
    }
}

/**
 * Removes an option if it exists.
 */
func (section Section) RemoveOption(optionName string) bool {
    normed := section.normOption(optionName)

    if section.HasOption(normed) {
        DeleteInMap(section.options, normed)
        DeleteInSlice(section.lines, section.getOption(normed))

        return true
    } else {
        return false
    }
}


/**
 * 删除 map 中的指定元素。
 *
 */
func DeleteInMap(targetMap interface{}, item interface{}) {

    value := reflect.Indirect(reflect.ValueOf(targetMap))

    if value.MapIndex(reflect.ValueOf(item)).IsValid() {
        value.SetMapIndex(reflect.ValueOf(item), reflect.Value{})
    }
}

/**
 * 删除 slice 中的指定元素。
 *
 */
func DeleteInSlice(slice interface{}, item interface{}) {
    index := -1
    ve := reflect.Indirect(reflect.ValueOf(slice))
    size := ve.Len()

    for i := 0; i < size; i++ {
        if reflect.DeepEqual(ve.Index(i).Interface(), item) {
            index = i
            break
        }
    }

    if index >= 0 {
        ve.Set(reflect.AppendSlice(ve.Slice(0, index), ve.Slice(index + 1, size)))
    }
}

/**
 * 检查 obj 是否存在于 collection 中， collection 类型可以为 slice、array和map。
 * 如果存在，返回其下标（slice和array），或者 1 （map）；
 * 如果不存在，返回-1。
 */
func Contains(collection interface{}, obj interface{}) int {
    targetValue := reflect.ValueOf(collection)
    existed := -1
    switch reflect.TypeOf(collection).Kind() {
    case reflect.Slice, reflect.Array:
        for i := 0; i < targetValue.Len(); i++ {
            if reflect.DeepEqual(targetValue.Index(i).Interface(), obj) {
                existed = i
                break
            }
        }
    case reflect.Map:
        if targetValue.MapIndex(reflect.ValueOf(obj)).IsValid() {
            existed = 1
        }
    }

    return existed
}