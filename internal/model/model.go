package model

type Holiday struct {
	Name string
	Date string
	Href string
}

func NewHoliday(name string, date string, href string) *Holiday {
	return &Holiday{
		Name: name,
		Date: date,
		Href: href,
	}
}

type Postcard struct {
	Holiday    string
	Page       string
	Href       string
	YouTube    bool
	Downloaded bool
	Path       string
	Name       string
}

func NewPostcard(holiday string, page string) *Postcard {
	return &Postcard{
		Holiday: holiday,
		Page:    page,
	}
}
