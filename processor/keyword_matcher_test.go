package processor

import (
	"testing"
	"time"

	"gotest.tools/v3/assert"
)

func TestNewKeywordHandler(t *testing.T) {
	timeLoc := &time.Location{}
	got := NewKeywordHandler(timeLoc)

	assert.Equal(t, len(got.replacers), 2)

	for _, r := range got.replacers {
		switch val := r.(type) {
		case *NowReplacer:
			tw := val.timeWrapper.(Time)
			assert.Equal(t, tw.location, timeLoc)
		case *CurrentMonthDayReplacer:
			tw := val.timeWrapper.(Time)
			assert.Equal(t, tw.location, timeLoc)
		default:
			t.Error("unhandled replacer")
		}
	}
}

func TestTimeWrapper_Now(t *testing.T) {
	t.Run("returns local time", func(t *testing.T) {
		tm := Time{}
		assert.Equal(t, tm.Now().Truncate(time.Second), time.Now().Local().Truncate(time.Second))
	})

	t.Run("returns time in a specific location", func(t *testing.T) {
		loc, err := time.LoadLocation("America/New_York")
		assert.NilError(t, err)

		tm := Time{location: loc}
		assert.Equal(t,
			tm.Now().UTC().In(loc).Truncate(time.Second),
			time.Now().UTC().In(loc).Truncate(time.Second),
		)
	})
}

func TestNowReplacer_HasMatched(t *testing.T) {
	goodCases := []string{"NOW", "NOW+1", "NOW-1", "NOW-20", "NOW+20"}
	for _, in := range goodCases {
		t.Run("returns true when input is "+in, func(t *testing.T) {
			nr := NowReplacer{value: in}
			assert.Assert(t, nr.HasMatched())
		})
	}

	badCases := []string{"now", "NOW+", "NOW-", "NOW20"}
	for _, in := range badCases {
		t.Run("returns false when input is "+in, func(t *testing.T) {
			nr := NowReplacer{value: in}
			assert.Equal(t, false, nr.HasMatched())
		})
	}
}

func TestNowReplacer_SetValue(t *testing.T) {
	t.Run("sets the value", func(t *testing.T) {
		nr := NowReplacer{}
		nr.SetValue("NOW-1")

		assert.Equal(t, nr.value, "NOW-1")
	})
}

func TestNowReplacer_ComputedValue(t *testing.T) {
	t.Run("returns the original value when it doesn't have the expected matches", func(t *testing.T) {
		nr := NowReplacer{value: "NOW-", timeWrapper: mockTimeWrapper{}}
		assert.Equal(t, nr.ComputedValue(), "NOW-")
	})

	t.Run("returns the current date", func(t *testing.T) {
		nr := NowReplacer{value: "NOW", timeWrapper: mockTimeWrapper{}}
		assert.Equal(t, nr.ComputedValue(), "2005-02-04")
	})

	t.Run("returns yesterday date", func(t *testing.T) {
		nr := NowReplacer{value: "NOW-1", timeWrapper: mockTimeWrapper{}}
		assert.Equal(t, nr.ComputedValue(), "2005-02-03")
	})

	t.Run("returns tomorrows date", func(t *testing.T) {
		nr := NowReplacer{value: "NOW+1", timeWrapper: mockTimeWrapper{}}
		assert.Equal(t, nr.ComputedValue(), "2005-02-05")
	})

	t.Run("returns the 10 days ago", func(t *testing.T) {
		nr := NowReplacer{value: "NOW-10", timeWrapper: mockTimeWrapper{}}
		assert.Equal(t, nr.ComputedValue(), "2005-01-25")
	})
}

func TestCurrentMonthDayReplacer_HasMatched(t *testing.T) {
	t.Run("returns true", func(t *testing.T) {
		cm := CurrentMonthDayReplacer{value: "CURRENT_MONTH_DAY1"}
		assert.Assert(t, cm.HasMatched())
	})

	t.Run("returns false", func(t *testing.T) {
		cm := CurrentMonthDayReplacer{value: "CURRENT_MONTH_DAY2"}
		assert.Equal(t, false, cm.HasMatched())
	})

	t.Run("returns false", func(t *testing.T) {
		cm := CurrentMonthDayReplacer{value: "CURRENT_MONTH_DAY"}
		assert.Equal(t, false, cm.HasMatched())
	})
}

func TestCurrentMonthDayReplacer_SetValue(t *testing.T) {
	t.Run("sets the value", func(t *testing.T) {
		cm := CurrentMonthDayReplacer{}
		cm.SetValue("NOW-1")

		assert.Equal(t, cm.value, "NOW-1")
	})
}

func TestCurrentMonthDayReplacer_ComputedValue(t *testing.T) {
	t.Run("returns the original value when it doesn't have the expected matches", func(t *testing.T) {
		nr := CurrentMonthDayReplacer{value: "CURRENT_MONTH_DAY", timeWrapper: mockTimeWrapper{}}
		assert.Equal(t, nr.ComputedValue(), "CURRENT_MONTH_DAY")
	})

	t.Run("returns 1st of the current date", func(t *testing.T) {
		nr := CurrentMonthDayReplacer{value: "CURRENT_MONTH_DAY1", timeWrapper: mockTimeWrapper{}}
		assert.Equal(t, nr.ComputedValue(), "2005-02-01")
	})
}

type mockTimeWrapper struct {
	now time.Time
}

func (m mockTimeWrapper) Now() time.Time {
	if !m.now.IsZero() {
		return m.now
	}

	return time.Date(2005, 2, 4, 6, 8, 0, 0, time.UTC)
}
