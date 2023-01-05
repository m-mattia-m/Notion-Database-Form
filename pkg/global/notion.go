package global

var (
	User *NotionUser
)

type NotionUser struct {
	Id string `json:"id"`
	//Firstname string `json:"firstname"`
	//Lastname  string `json:"lastname"`
}
