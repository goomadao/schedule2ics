package util

import (
	"bytes"
	"fmt"
	"io/ioutil"

	"github.com/jordic/goics"
)

// Classes2ICS transfer class events to ics format
func Classes2ICS(classes []ClassEvent, filePath string) {
	c := goics.NewComponent()
	c.SetType("VCALENDAR")

	for _, class := range classes {
		s := goics.NewComponent()
		s.SetType("VEVENT")
		k, v := goics.FormatDateTimeField("DTSTART", class.StartTime)
		s.AddProperty(k, v)
		k, v = goics.FormatDateTimeField("DTEND", class.EndTime)
		s.AddProperty(k, v)
		s.AddProperty("DESCRIPTION", class.Description)
		s.AddProperty("SUMMARY", fmt.Sprintf("%s-%s", class.Name, class.Teacher))
		s.AddProperty("LOCATION", class.Location)

		c.AddComponent(s)
	}

	ics := &ICSEvent{
		component: c,
	}
	buf := &bytes.Buffer{}
	enc := goics.NewICalEncode(buf)
	enc.Encode(ics)

	ioutil.WriteFile(filePath, buf.Bytes(), 0644)
}
