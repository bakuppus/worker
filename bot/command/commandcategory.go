package command

type Category string

const (
	General    Category = "âšī¸ General"
	Tickets    Category = "đŠ Tickets"
	Settings   Category = "đ§ Settings"
	Tags       Category = "âī¸ Tags"
	Statistics Category = "đ Statistics"
)

var Categories = []Category{
	General,
	Tickets,
	Settings,
	Tags,
	Statistics,
}
