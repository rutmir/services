package v2

type Action string

const (
	Authorize Action = "authorize"
	Result Action = "result"
	RefreshToken Action = "refreshtoken"
	Synchronize Action = "synchronize"

	AddedToChat Action = "addedtochat"
	CreateChat Action = "createchat"
	ChatList Action = "getchatlist" //removed
	AddMessage Action = "addmessage"
	UpdateMessage Action = "updatemessage"
	RemovedFromChat Action = "removedfromchat"
	GetChat Action = "getchat"
	RemoveChat Action = "removechat"
	GetContactList Action = "getcontactlist" //removed
	AddContact Action = "addcontact"
	RemoveContact Action = "removecontact"
	ProfileList Action = "getprofilelist"
	GetProfile Action = "getprofile"
	UpdateProfile Action = "updateprofile"
	RequestAnswer Action = "requestanswer"
	RequestResult Action = "requestresult"
	SetAvailabilityStatus Action = "setavailabilitystatus"
	ProfileChanged Action = "profilechanged"
	MessageStatus Action = "messagestatus"
	UpdateChat Action = "updatechat"
	ChatChanged Action = "chatchanged"
	MessageList Action = "getmessagelist"
	MessageCount Action = "getmessagecount"
	DeleteMessage Action = "deletemessage"
	KeyBoardTyping Action = "keyboardtyping"
	ClientInit Action = "clientinit"

	FileUpload Action = "fileupload"
	GetFile Action = "getfile"
	GetPreview Action = "getpreview"
	GetFileList Action = "getfilelist"
	FileActionNotify Action = "fileactionnotify"
	SaveFile Action = "savefile"
	FileDelete Action = "filedelete"
)