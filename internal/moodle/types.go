package moodle

// SiteInfo represents the current authenticated user and site.
type SiteInfo struct {
	UserID   int    `json:"userid"`
	Username string `json:"username"`
	Fullname string `json:"fullname"`
	SiteName string `json:"sitename"`
	SiteURL  string `json:"siteurl"`
}

// Course represents an enrolled Moodle course.
type Course struct {
	ID        int    `json:"id"`
	FullName  string `json:"fullname"`
	ShortName string `json:"shortname"`
	IDNumber  string `json:"idnumber"`
	StartDate int64  `json:"startdate"`
	EndDate   int64  `json:"enddate"`
	Category  int    `json:"categoryid"`
}

// Content represents a file or link inside a module.
type Content struct {
	Type     string `json:"type"`
	Filename string `json:"filename"`
	FileSize int64  `json:"filesize"`
	FileURL  string `json:"fileurl"`
	MimeType string `json:"mimetype"`
}

// Module represents a learning activity or resource inside a section.
type Module struct {
	ID       int       `json:"id"`
	Name     string    `json:"name"`
	ModName  string    `json:"modname"`
	URL      string    `json:"url"`
	Instance int       `json:"instance"`
	Contents []Content `json:"contents"`
}

// Section represents a weekly/topic section inside a course.
type Section struct {
	ID      int      `json:"id"`
	Name    string   `json:"name"`
	Section int      `json:"section"`
	Modules []Module `json:"modules"`
}

// CalendarCourse is the course field embedded in a calendar event.
type CalendarCourse struct {
	ID        int    `json:"id"`
	FullName  string `json:"fullname"`
	ShortName string `json:"shortname"`
}

// CalendarAction is the action field embedded in a calendar event.
type CalendarAction struct {
	Name       string `json:"name"`
	URL        string `json:"url"`
	Actionable bool   `json:"actionable"`
}

// CalendarEvent represents a deadline/action event from the calendar API.
type CalendarEvent struct {
	ID          int            `json:"id"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	TimeSort    int64          `json:"timesort"`
	TimeStart   int64          `json:"timestart"`
	Course      CalendarCourse `json:"course"`
	Action      CalendarAction `json:"action"`
}

// GradeItem represents a single grade entry for a student.
type GradeItem struct {
	ID                  int     `json:"id"`
	ItemName            string  `json:"itemname"`
	ItemType            string  `json:"itemtype"`
	GradeMin            float64 `json:"grademin"`
	GradeMax            float64 `json:"grademax"`
	GradeFormatted      string  `json:"gradeformatted"`
	PercentageFormatted string  `json:"percentageformatted"`
	Feedback            string  `json:"feedback"`
}

// UserGrade wraps grade items for a specific user in a course.
type UserGrade struct {
	CourseID   int         `json:"courseid"`
	UserID     int         `json:"userid"`
	GradeItems []GradeItem `json:"gradeitems"`
}

// Forum represents a Moodle forum.
type Forum struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}

// Discussion represents a post in a forum.
type Discussion struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	Message      string `json:"message"`
	Created      int64  `json:"created"`
	UserFullName string `json:"userfullname"`
	NumReplies   int    `json:"numreplies"`
}
