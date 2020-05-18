package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	configFile = "config.json"
)

type configuration struct {
	API          string
	Year         int
	CountryCode  string
	Location     string
	APILayout    string
	OutputLayout string
	APITimeout   uint
	UseFakeAPI   bool
	UseFakeDate  bool
	Date         string
}

type holiday struct {
	Date        string
	LocalName   string
	Name        string
	CountryCode string
	Fixed       bool
	Global      bool
	Counties    []string
	LaunchYear  int
	Type        string
}

func exitWithErrorMessage(err error) {
	fmt.Fprintf(os.Stderr, "error: %v\n", err)
	os.Exit(1)
}

func loadConfiguration(useConfigFile bool) configuration {
	var cfg configuration
	if useConfigFile {
		fp, err := os.Open(configFile)
		if err != nil {
			exitWithErrorMessage(err)
		}
		defer fp.Close()
		if err = json.NewDecoder(fp).Decode(&cfg); err != nil {
			exitWithErrorMessage(err)
		}
	} else {
		cfg.API = "https://date.nager.at/api/v2/publicholidays/"
		cfg.Year = 2020
		cfg.CountryCode = "UA"
		cfg.Location = "Europe/Kiev"
		cfg.APILayout = "2006-01-02"
		cfg.OutputLayout = "Jan 02"
		cfg.APITimeout = 10
		cfg.UseFakeAPI = false
		cfg.UseFakeDate = false
		cfg.Date = "2020-05-18"
	}
	return cfg
}

func loadHolidaysAPIData(cfg configuration) []holiday {
	var holidays []holiday
	if cfg.UseFakeAPI {
		holidays = []holiday{
			holiday{"2020-01-01", "Новий Рік", "New Year's Day", "UA", true, true, []string{}, 0, "Public"},
			holiday{"2020-01-07", "Різдво", "(Julian) Christmas", "UA", true, true, []string{}, 0, "Public"},
			holiday{"2020-03-08", "Міжнародний жіночий день", "International Women's Day", "UA", true, true, []string{}, 0, "Public"},
			holiday{"2020-04-19", "Великдень", "(Julian) Easter Sunday", "UA", false, true, []string{}, 0, "Public"},
			holiday{"2020-05-01", "День праці", "International Workers' Day", "UA", true, true, []string{}, 0, "Public"},
			holiday{"2020-05-09", "День перемоги над нацизмом у Другій світовій війні", "Victory day over Nazism in World War II", "UA", true, true, []string{}, 0, "Public"},
			holiday{"2020-06-07", "Трійця", "(Julian) Pentecost", "UA", false, true, []string{}, 0, "Public"},
			holiday{"2020-06-28", "День Конституції", "Constitution Day", "UA", true, true, []string{}, 0, "Public"},
			holiday{"2020-08-24", "День Незалежності", "Independence Day", "UA", true, true, []string{}, 0, "Public"},
			holiday{"2020-10-14", "День захисника України", "Defender of Ukraine Day", "UA", true, true, []string{}, 0, "Public"},
			holiday{"2020-12-25", "Різдво", "(Gregorian and Revised Julian) Christmas", "UA", true, true, []string{}, 0, "Public"},
		}
	} else {
		httpClient := &http.Client{Timeout: time.Duration(cfg.APITimeout) * time.Second}
		resp, err := httpClient.Get(fmt.Sprintf("%s/%d/%s", strings.Trim(cfg.API, "/"), cfg.Year, cfg.CountryCode))
		if err != nil {
			exitWithErrorMessage(err)
		}
		defer resp.Body.Close()
		if err = json.NewDecoder(resp.Body).Decode(&holidays); err != nil {
			exitWithErrorMessage(err)
		}
	}
	return holidays
}

func adjustWeekend(t time.Time, d time.Time, cfg configuration) string {
	var w string
	extendedWeekend := ", and the weekend will last 3 days: %s - %s"
	remainingWeekend := ", and the weekend will last to %s"
	if t.Weekday() == time.Friday || t.Weekday() == time.Saturday {
		w = fmt.Sprintf(extendedWeekend, t.Format(cfg.OutputLayout), t.AddDate(0, 0, 2).Format(cfg.OutputLayout))
	} else if t.Weekday() == time.Sunday {
		if t.Equal(d) {
			w = fmt.Sprintf(remainingWeekend, t.AddDate(0, 0, 1).Format(cfg.OutputLayout))
		} else {
			w = fmt.Sprintf(extendedWeekend, t.AddDate(0, 0, -1).Format(cfg.OutputLayout), t.AddDate(0, 0, 1).Format(cfg.OutputLayout))
		}
	} else if t.Weekday() == time.Monday {
		if t.AddDate(0, 0, -1).Equal(d) {
			w = fmt.Sprintf(remainingWeekend, t.AddDate(0, 0, 1).Format(cfg.OutputLayout))
		} else if !t.Equal(d) {
			w = fmt.Sprintf(extendedWeekend, t.AddDate(0, 0, -2).Format(cfg.OutputLayout), t.Format(cfg.OutputLayout))
		}
	}
	return w
}

func timeToDate(t time.Time) time.Time {
	y, m, d := t.Date()
	return time.Date(y, m, d, 0, 0, 0, 0, t.Location())
}

func loadLocalDate(cfg configuration, holidays []holiday) (time.Time, []holiday) {
	loc, err := time.LoadLocation(cfg.Location)
	if err != nil {
		exitWithErrorMessage(err)
	}
	var dateInLoc time.Time
	if cfg.UseFakeDate {
		dateInLoc, err = time.ParseInLocation(cfg.APILayout, cfg.Date, loc)
		if err != nil {
			exitWithErrorMessage(err)
		}
	} else {
		dateInLoc = timeToDate(time.Now().In(loc))
	}
	if !cfg.UseFakeAPI && cfg.Year != dateInLoc.Year() {
		fmt.Fprintf(os.Stderr, "warning: API data has wrong date, retrieving actual data...\n")
		cfg.Year = dateInLoc.Year()
		holidays = loadHolidaysAPIData(cfg)
	}
	return dateInLoc, holidays
}

func showNextHoliday(cfg configuration, holidays []holiday) {
	dateInLoc, holidays := loadLocalDate(cfg, holidays)
	for i, h := range holidays {
		t, err := time.ParseInLocation(cfg.APILayout, h.Date, dateInLoc.Location())
		if err != nil {
			exitWithErrorMessage(err)
		}
		if t.Before(dateInLoc) {
			if i+1 == len(holidays) {
				if cfg.UseFakeDate || cfg.UseFakeAPI {
					fmt.Fprintf(os.Stderr, "warning: wrong configuration or no more holidays this year, trying to retrieve first holiday in next year...\n")
				}
				cfg.Year++
				holidays = loadHolidaysAPIData(cfg)
				h = holidays[0]
				t, err = time.ParseInLocation(cfg.APILayout, h.Date, dateInLoc.Location())
				if err != nil {
					exitWithErrorMessage(err)
				}
				fmt.Printf("The next holiday is %s, %s%s\n", h.Name, t.Format(cfg.OutputLayout), adjustWeekend(t, dateInLoc, cfg))
			}
			continue
		} else if t.Equal(dateInLoc) {
			fmt.Printf("Today is %s, %s%s\n", h.Name, t.Format(cfg.OutputLayout), adjustWeekend(t, dateInLoc, cfg))
		} else {
			fmt.Printf("The next holiday is %s, %s%s\n", h.Name, t.Format(cfg.OutputLayout), adjustWeekend(t, dateInLoc, cfg))
		}
		break
	}
}

func parseArgumetns() bool {
	useConfigFile := flag.Bool("conf", false, "load configuration from 'config.json' file")
	flag.Parse()
	return *useConfigFile
}

func main() {
	useConfigFile := parseArgumetns()
	cfg := loadConfiguration(useConfigFile)
	holidays := loadHolidaysAPIData(cfg)
	showNextHoliday(cfg, holidays)
}
