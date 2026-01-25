package web

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/andreasphil/one/adapter"
	"github.com/andreasphil/one/lib"
	"github.com/andreasphil/one/util"
)

func urlForDate(date time.Time) string {
	return fmt.Sprintf("/calendar/%d/%02d", date.Year(), int(date.Month()))
}

func getCalendar() http.Handler {
	return http.RedirectHandler(urlForDate(time.Now()), http.StatusTemporaryRedirect)
}

func getCalendarMonth(finder adapter.NotesByDateFinder) http.HandlerFunc {
	type calendarPageData struct {
		Month            time.Time
		Calendar         map[int][]lib.Note
		PreviousMonthUrl string
		NextMonthUrl     string
		DaysBeforeFirst  int
		DaysAfterLast    int
	}

	render := newRenderFunc[calendarPageData]("get_calendar.html")

	return func(w http.ResponseWriter, r *http.Request) {
		year, err := strconv.ParseInt(r.PathValue("year"), 10, 0)
		if err != nil {
			util.Errorf("cannot convert value for year '%v' to int", r.PathValue("year"))
			http.NotFound(w, r)
			return
		}

		month, err := strconv.ParseInt(r.PathValue("month"), 10, 0)
		if err != nil {
			util.Errorf("cannot convert value for month '%v' to int", r.PathValue("month"))
			http.NotFound(w, r)
			return
		}

		currentMonth := time.Date(int(year), time.Month(int(month)), 1, 0, 0, 0, 0, time.UTC)
		lastDayOfCurrentMonth := currentMonth.AddDate(0, 1, 0).AddDate(0, 0, -1)
		daysInMonth := lastDayOfCurrentMonth.Day()

		calendar := make(map[int][]lib.Note)
		for i := range daysInMonth {
			calendar[i+1] = []lib.Note{}
		}

		notes := finder.FindNotesByDate(currentMonth, lastDayOfCurrentMonth)
		for _, note := range notes {
			if !note.IsDailyNote() {
				continue
			}

			day := note.Date.Day()
			calendar[day] = append(calendar[day], note)
			calendar[day] = append(calendar[day], note.Children...)
		}

		// Normalize weekdays to start on Monday (Weekday() starts on Sunday):
		// Sunday should be 6 (so that there's an empty week visible before Sunday),
		// shift everything else back one day.
		previousMonthDays := int(currentMonth.Weekday()) - 1
		if previousMonthDays == -1 {
			previousMonthDays = 6
		}

		err = render(w, data[calendarPageData]{
			Title:      currentMonth.Format("January 2006"),
			CurrentUrl: r.URL.Path,

			Data: calendarPageData{
				Month:            currentMonth,
				NextMonthUrl:     urlForDate(currentMonth.AddDate(0, 1, 0)),
				DaysAfterLast:    42 - previousMonthDays - daysInMonth, // 42 = 6 weeks
				PreviousMonthUrl: urlForDate(currentMonth.AddDate(0, -1, 0)),
				DaysBeforeFirst:  previousMonthDays,
				Calendar:         calendar,
			},
		})

		if err != nil {
			util.HttpErrorf(w, http.StatusInternalServerError, "failed to render page template: %v", err)
			return
		}
	}
}
