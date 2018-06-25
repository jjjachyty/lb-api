package models

const avatarCN = "avatar"

type Avatar struct {
	Divider bool   `json:"divider"`
	Header  string `json:"header"`
	Id      string `json:"id"`
	Name    string `json:"name"`
	Group   string `json:"group"`
	Url     string `json:"url"`
}

//获取所有的头像数据
func FindAllAvatar() ([]Avatar, error) {
	var ruslt = make([]Avatar, 0)
	err := DB.C(avatarCN).Find(nil).All(&ruslt)
	return ruslt, err
}
