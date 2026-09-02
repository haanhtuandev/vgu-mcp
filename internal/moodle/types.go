package moodle

type SiteInfo struct {
	UserID   int    `json:"userid"`
	Username string `json:"username"`
	Fullname string `json:"fullname"`
	SiteName string `json:"sitename"`
	SiteURL  string `json:"siteurl"`
}

type Course struct {
	ID        int    `json:"id"`
	FullName  string `json:"fullname"`
	ShortName string `json:"shortname"`
	Category  int    `json:"categoryid"`
}

type Module struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	ModName  string `json:"modname"`
	URL      string `json:"url"`
	Instance int    `json:"instance"`
}

type Section struct {
	ID      int      `json:"id"`
	Name    string   `json:"name"`
	Modules []Module `json:"modules"`
}
