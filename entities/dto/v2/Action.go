package v2

type Action string

const (
	Action_Authorize Action = "authorize"
	Action_Result Action = "result"
	Action_RefreshToken Action = "refreshtoken"
	Action_Synchronize Action = "synchronize"

	Action_AddedToChat Action = "addedtochat"
	Action_CreateChat Action = "createchat"
	Action_ChatList Action = "getchatlist" //removed
	Action_AddMessage Action = "addmessage"
	Action_UpdateMessage Action = "updatemessage"
	Action_RemovedFromChat Action = "removedfromchat"
	Action_GetChat Action = "getchat"
	Action_RemoveChat Action = "removechat"
	Action_GetContactList Action = "getcontactlist" //removed
	Action_AddContact Action = "addcontact"
	Action_RemoveContact Action = "removecontact"
	Action_ProfileList Action = "getprofilelist"
	Action_GetProfile Action = "getprofile"
	Action_UpdateProfile Action = "updateprofile"
	Action_RequestAnswer Action = "requestanswer"
	Action_RequestResult Action = "requestresult"
	Action_SetAvailabilityStatus Action = "setavailabilitystatus"
	Action_ProfileChanged Action = "profilechanged"
	Action_MessageStatus Action = "messagestatus"
	Action_UpdateChat Action = "updatechat"
	Action_ChatChanged Action = "chatchanged"
	Action_MessageList Action = "getmessagelist"
	Action_MessageCount Action = "getmessagecount"
	Action_DeleteMessage Action = "deletemessage"
	Action_KeyBoardTyping Action = "keyboardtyping"
	Action_ClientInit Action = "clientinit"

	Action_FileUpload Action = "fileupload"
	Action_GetFile Action = "getfile"
	Action_GetPreview Action = "getpreview"
	Action_GetFileList Action = "getfilelist"
	Action_FileActionNotify Action = "fileactionnotify"
	Action_SaveFile Action = "savefile"
	Action_FileDelete Action = "filedelete"
)