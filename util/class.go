package util

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/viper"
	"github.com/tealeg/xlsx/v3"
)

// Xlsx2Classes transfer schedule.xlsx into class events
func Xlsx2Classes(filePath string) (classes []ClassEvent) {
	wb, err := xlsx.OpenFile(filePath)
	if err != nil {
		panic(err)
	}
	sh := wb.Sheets[0]
	startCol, startRow, err := xlsx.GetCoordsFromCellIDString(viper.GetString("BeginningOfSchedule"))
	if err != nil {
		panic(err)
	}
	for c := startCol + 1; c < sh.MaxCol; c++ {
		for r := startRow + 1; r < sh.MaxRow; r++ {
			cell, err := sh.Cell(r, c)
			if err != nil {
				panic(err)
			}
			if cell.Value != "" {
				classes = append(classes, getClassEvents(c-startCol, r-startRow, cell.Value)...)
			}
		}
	}
	return
}

func getClassEvents(day, index int, info string) (events []ClassEvent) {
	name, teacher, startWeek, endWeek, location := splitInfo(info)

	rangeOfClasses := viper.GetStringSlice("RangeOfClasses")
	if len(rangeOfClasses) < index {
		panic(errors.New("Class index is too large, check yout config file"))
	}
	match, err := regexp.Match("[0-9]{1,2}:[0-9]{1,2}-[0-9]{1,2}:[0-9]{1,2}", []byte(rangeOfClasses[index-1]))
	if err != nil {
		panic(err)
	}
	if !match {
		panic(errors.New("The format of RangeOfClasses is not correct, please check"))
	}
	timeRange := strings.Split(rangeOfClasses[index-1], "-")

	date := viper.GetString("MondayOfFirstWeek")
	match, err = regexp.Match("[0-9]{4}-[0-9]{1,2}-[0-9]{1,2}", []byte(date))
	if err != nil {
		panic(err)
	}
	if !match {
		panic(errors.New("The format of MondayOfFirstWeek is not correct, please check"))
	}
	timezone := viper.GetString("Timezone")
	tz, err := time.LoadLocation(timezone)
	if err != nil {
		panic(err)
	}
	baseStartTime, err := time.ParseInLocation("2006-01-02 15:04", date+" "+timeRange[0], tz)
	if err != nil {
		panic(err)
	}
	baseEndTime, err := time.ParseInLocation("2006-01-02 15:04", date+" "+timeRange[1], tz)

	for i := 1; i < day; i++ {
		baseStartTime = baseStartTime.Add(24 * time.Hour)
		baseEndTime = baseEndTime.Add(24 * time.Hour)
	}

	for i := 1; int64(i) < startWeek; i++ {
		baseStartTime = baseStartTime.Add(7 * 24 * time.Hour)
		baseEndTime = baseEndTime.Add(7 * 24 * time.Hour)
	}

	startTime := baseStartTime
	endTime := baseEndTime

	for i := startWeek; i <= endWeek; i++ {
		event := ClassEvent{
			StartTime:   startTime,
			EndTime:     endTime,
			Name:        name,
			Location:    location,
			Teacher:     teacher,
			Description: info,
		}
		events = append(events, event)
		startTime = startTime.Add(7 * 24 * time.Hour)
		endTime = endTime.Add(7 * 24 * time.Hour)
	}

	return
}

func splitInfo(info string) (name, teacher string, startWeek, endWeek int64, location string) {
	index, rightBracketCount, leftBracketCount := len(info)-2, 1, 0
	for rightBracketCount != leftBracketCount {
		index--
		if info[index] == '(' {
			leftBracketCount++
		} else if info[index] == ')' {
			rightBracketCount++
		}
	}
	name = info[:index]

	otherInfo := info[index+1 : len(info)-1]
	otherInfos := strings.Split(otherInfo, "；")
	teacher, weeks, location := otherInfos[0], otherInfos[2][:len(otherInfos[2])-3], otherInfos[3]
	if strings.Contains(weeks, "全") {
		startWeek = 1
		endWeek = viper.GetInt64("WeeksOfTerm")
		if endWeek < startWeek {
			panic(errors.New("please specify WeeksOfTerm in your config file"))
		}
	} else {
		weeksArray := strings.Split(weeks, "-")
		var err error
		startWeek, err = strconv.ParseInt(weeksArray[0], 10, 32)
		if err != nil {
			panic(err)
		}
		endWeek, err = strconv.ParseInt(weeksArray[1], 10, 32)
		if err != nil {
			panic(err)
		}
	}
	return
}
