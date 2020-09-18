# Schedule2ICS

## Description

transfer your schedule(must be a .xlsx file) into .ics

## Usage

```
A tool to transfer your schedule into .ics
You can import the .ics file into your Google Calendar or Apple Calendar.

Usage:
  schedule2ics [flags]

Flags:
  -c, --config string     config file (default is ./schedule2ics.yaml)
  -h, --help              help for schedule2ics
  -o, --output string     output .ics file (default "schedule.ics")
  -s, --schedule string   schedule file(must be .xls or .xlsx)
```

## Config

The format of config file can be found in [`schedule2ics.yaml`](./schedule2ics.yaml)