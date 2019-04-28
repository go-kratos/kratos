package model

const (
	// type
	FieldFav      = "folder"
	FieldArc      = "video"
	FieldResource = "resource"
	// action
	ActionAdd                 = "add"
	ActionDel                 = "del"
	ActionMove                = "move"
	ActionCopy                = "copy"
	ActionMdel                = "mdel"
	ActionIndef               = "indef"
	ActionIncol               = "incol"
	ActionClean               = "clean"
	ActionInitRelationFids    = "initRelationFids"
	ActionInitFolderRelations = "initFolderRelations"
	ActionMultiAdd            = "multiAdd"
	ActionMultiDel            = "multiDel"
	ActionFolderAdd           = "folderAdd"
	ActionFolderDel           = "folderDel"
)

type Message struct {
	Field         string  `json:"field"`
	Action        string  `json:"action"`
	Oid           int64   `json:"oid"`
	Type          int8    `json:"type"`
	Mid           int64   `json:"mid"`
	OldMid        int64   `json:"old_mid"`
	Fid           int64   `json:"fid"`
	FidState      int8    `json:"fid_state"`
	FolderAttr    int8    `json:"folder_attr"`
	OldFolderAttr int8    `json:"old_folder_attr"`
	NewFolderAttr int8    `json:"new_folder_attr"`
	Aid           int64   `json:"aid"`
	OldFid        int64   `json:"old_fid"`
	OldFidState   int8    `json:"old_fid_state"`
	NewFid        int64   `json:"new_fid"`
	NewFidState   int8    `json:"new_fid_state"`
	Affected      int64   `json:"affected"`
	Aids          []int64 `json:"aids"`
	Oids          []int64 `json:"oids"`
	Mids          []int64 `json:"mids"`
	FTime         int64   `json:"ftime"`
}
