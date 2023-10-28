package model

type Holiday struct {
	Name string
	Href string
}

func NewHoliday(name string, href string) *Holiday {
	return &Holiday{
		Name: name,
		Href: href,
	}
}

type Postcard struct {
	Holiday string
	Href    string
}

func NewPostcard(holiday string, href string) *Postcard {
	return &Postcard{
		Holiday: holiday,
		Href:    href,
	}
}
