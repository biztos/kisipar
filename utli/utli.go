// utli.go - utilities we were sort of hoping to not have to write.
// -------

// Package utli assembles various utility functions useful in and around the
// Kisipar web-server project.  Its name is a play on "gúgli" which is the
// Hungarian slang term for "Google."
package utli

import (
	"time"
)

// TIME_PARSING_FORMAT_STRINGS contains an ordered list of time format strings
// used for parsing times.
var TIME_PARSING_FORMAT_STRINGS = []string{

	// I Can't Believe They're Not Constants!
	// (These are not only useful, they are probably more common than the
	// standard set of fomats.)
	"2006-01-02",              //  obvious way to do a European timestamp
	"2006.01.02",              // it's also done like this quite often
	"20060102",                // did they seriously forget this?
	"2006-01-02 15:04:05 MST", // very useful format!
	"2006-01-02 15:04:05",     // ...also without the TZ!

	// WTF Golang?  No constant for your own default time stringification?
	"2006-01-02 15:04:05.999999999 -0700 MST", // time.String()

	// Golang standards:
	time.ANSIC,       // "Mon Jan _2 15:04:05 2006"
	time.UnixDate,    // "Mon Jan _2 15:04:05 MST 2006"
	time.RubyDate,    // "Mon Jan 02 15:04:05 -0700 2006"
	time.RFC822,      // "02 Jan 06 15:04 MST"
	time.RFC822Z,     // "02 Jan 06 15:04 -0700", - RFC822 with numeric zone
	time.RFC850,      // "Monday, 02-Jan-06 15:04:05 MST"
	time.RFC1123,     // "Mon, 02 Jan 2006 15:04:05 MST"
	time.RFC1123Z,    // "Mon, 02 Jan 2006 15:04:05 -0700" - RFC1123 w/num.z.
	time.RFC3339,     // "2006-01-02T15:04:05Z07:00"
	time.RFC3339Nano, // "2006-01-02T15:04:05.999999999Z07:00"
	time.Kitchen,     // "3:04PM"
	time.Stamp,       // "Jan _2 15:04:05" - "Handy time stamp" per Go. :-)
	time.StampMilli,  // "Jan _2 15:04:05.000" - ditto
	time.StampMicro,  // "Jan _2 15:04:05.000000" - ditto
	time.StampNano,   // "Jan _2 15:04:05.000000000" - ditto

}

// ParseTimeString attempts to parse s into a time, trying each format in the
// TIME_PARSING_FORMAT_STRINGS in order.  If no parse is successful, nil is
// returned.
func ParseTimeString(s string) *time.Time {

	for _, f := range TIME_PARSING_FORMAT_STRINGS {
		t, err := time.Parse(f, s)
		if err == nil {
			return &t
		}
	}
	return nil
}
