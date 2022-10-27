package processor

import (
	"regexp"
	"strconv"
	"time"
)

const dateFormat = "2006-01-02"

var (
	nowSubRegexp          = regexp.MustCompile(`NOW(\-|\+)(\d+)`)
	currentMonthDayRegexp = regexp.MustCompile(`CURRENT_MONTH_DAY1`)
)

type TimeWrapper interface {
	Now() time.Time
}

type Time struct {
	location *time.Location
}

func (nt Time) Now() time.Time {
	if nt.location == nil {
		return time.Now().Local()
	}

	return time.Now().UTC().In(nt.location)
}

type KeywordReplacer interface {
	HasMatched() bool
	ComputedValue() string
	SetValue(string)
}

type NowReplacer struct {
	value       string
	timeWrapper TimeWrapper
}

type CurrentMonthDayReplacer struct {
	value       string
	timeWrapper TimeWrapper
}

type KeywordHandler struct {
	value     string
	replacer  KeywordReplacer
	replacers []KeywordReplacer
}

func NewKeywordHandler(timeLocation *time.Location) *KeywordHandler {
	timeNow := Time{location: timeLocation}

	return &KeywordHandler{
		replacers: []KeywordReplacer{
			&NowReplacer{timeWrapper: timeNow},
			&CurrentMonthDayReplacer{timeWrapper: timeNow},
		},
	}
}

func (n *NowReplacer) HasMatched() bool {
	return n.value == "NOW" || nowSubRegexp.MatchString(n.value)
}

func (n *NowReplacer) ComputedValue() string {
	matches := nowSubRegexp.FindStringSubmatch(n.value)

	if n.value == "NOW" {
		return n.timeWrapper.Now().Format(dateFormat)
	}

	if len(matches) != 3 {
		return n.value
	}

	parsedNum, _ := strconv.Atoi(matches[2])
	days := time.Duration(parsedNum*24) * time.Hour

	switch matches[1] {
	case "-":
		return n.timeWrapper.Now().Add(-days).Format(dateFormat)
	default:
		// Default to adding day
		return n.timeWrapper.Now().Add(days).Format(dateFormat)
	}
}

func (n *NowReplacer) SetValue(value string) {
	n.value = value
}

func (c *CurrentMonthDayReplacer) HasMatched() bool {
	return currentMonthDayRegexp.MatchString(c.value)
}

func (c *CurrentMonthDayReplacer) ComputedValue() string {
	if !c.HasMatched() {
		return c.value
	}

	now := c.timeWrapper.Now()
	pointOfTime := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.Local)

	return pointOfTime.Format(dateFormat)
}

func (n *CurrentMonthDayReplacer) SetValue(value string) {
	n.value = value
}

func (r *KeywordHandler) HasMatched() bool {
	for _, kr := range r.replacers {
		kr.SetValue(r.value)

		if kr.HasMatched() {
			r.replacer = kr
			return true
		}
	}

	return false
}

func (r *KeywordHandler) ComputedValue() string {
	return r.replacer.ComputedValue()
}

func (r *KeywordHandler) SetValue(value string) {
	r.value = value
}
