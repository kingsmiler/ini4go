package ini4go

const (
    HEADER_START rune = '['
    HEADER_END rune = ']'
)

var (
    max_name_length int = 1024
    invalid_name_chars = []rune{HEADER_START, HEADER_END}

    default_option_whitespace = []rune{' ', '\t' }
    default_option_delimiters = []rune{'=', ':'}
    default_comment_delimiters = []rune{'#', ';'}
    default_case_sensitivity = false
    default_option_format = "%s %s %s"
)

type section struct {
    name            string
    delims          []rune
    isCaseSensitive bool
}

func NewSection(name string, delims []rune, isCaseSensitive bool) *section {
    if len(delims) == 0 {
        delims = default_comment_delimiters
    }

    return &section{
        name:name,
        delims: delims,
        isCaseSensitive : isCaseSensitive,
    }
}