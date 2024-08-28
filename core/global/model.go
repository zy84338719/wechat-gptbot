package global

import (
	"github.com/eatmoreapple/openwechat"
)

var AllowGroupUsernameMap = map[string]GroupMember{}

func JugleGroup(username string) (bool, GroupMember) {
	for s, member := range AllowGroupUsernameMap {
		if s == username {
			return true, member
		}
	}
	return false, GroupMember{}
}

type GroupMember struct {
	Username   string
	GroupName  string
	Members    map[string]*openwechat.User
	Group      *openwechat.Group
	AllowGroup AllowGroup
	TextDefine *TextDefine
	SignedIn   *SignedIn
}

type AllowGroup struct {
	Username          string      `json:"username"`
	GroupName         string      `json:"groupName"`
	Redisplay         bool        `json:"redisplay"`
	ExitCheck         bool        `json:"exitCheck"`
	IsTextDefault     bool        `json:"IsTextDefault"`
	AllowSignIn       bool        `json:"allowSignIn"`
	SaveMessage       bool        `json:"saveMessage"`
	Reply             bool        `json:"reply"`
	AiReply           bool        `json:"aiReply"`
	EnableAbility     bool        `json:"enableAbility"`
	EnableListAbility []string    `json:"enableListAbility"`
	TextDefine        *TextDefine `json:"textDefine"`
	SignedIn          *SignedIn   `json:"signedIn"`
}

type SignedIn struct {
	AllowStartTime string
	AllowEndTime   string

	SignInSuccess      string
	SignInFailSignedIn string
	SignInFailTimeout  string
	SignInError        string
}

type TextDefine struct {
	Redisplay string
	ExitCheck string
	TrickMe   string
	JoinGroup string
	Error     string
}

var SendObejctChannel = make(chan SendObject, 10)

type SendObject struct {
	Content    string
	SendGroup  *openwechat.Group
	SendFriend *openwechat.Friend

	UserType UserType
}

type UserType int

const (
	IsFriend UserType = iota
	IsGroup
	IsFriendHelper
)

func BuildSendObjectByGroup(content string, sendGroup *openwechat.Group) SendObject {

	return SendObject{
		Content:   content,
		SendGroup: sendGroup,
		UserType:  IsGroup,
	}
}

func BuildSendObjectByFriend(content string, sendFriend *openwechat.Friend) SendObject {
	return SendObject{
		Content:    content,
		SendFriend: sendFriend,
		UserType:   IsFriend,
	}
}

func BuildSendObjectByFriendHelper(content string) SendObject {
	return SendObject{
		Content:  content,
		UserType: IsFriendHelper,
	}
}
