package model

type Data struct {
	Users  []FriendsInfo `json:"users"`
	Groups []GroupsInfo  `json:"groups"`
}

type FriendsInfo struct {
	Username   string `json:"username"`
	Nickname   string `json:"nickname"`
	RemarkName string `json:"remark_name"`
	HeadImgUrl string `json:"head_img_url"`
}

type GroupsInfo struct {
	Username   string `json:"username"`
	Nickname   string `json:"nickname"`
	RemarkName string `json:"remark_name"`
	HeadImgUrl string `json:"head_img_url"`
}
