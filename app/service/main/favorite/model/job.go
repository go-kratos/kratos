package model

const (
	// type
	FieldFav      = "folder"
	FieldArc      = "video"
	FieldResource = "resource"
	// action
	ActionAdd                    = "add"
	ActionDel                    = "del"
	ActionMove                   = "move"
	ActionCopy                   = "copy"
	ActionMdel                   = "mdel"
	ActionIndef                  = "indef"
	ActionIncol                  = "incol"
	ActionClean                  = "clean"
	ActionInitRelationFids       = "initRelationFids"
	ActionInitFolderRelations    = "initFolderRelations"
	ActionInitAllFolderRelations = "initAllFolderRelations"
	ActionMultiAdd               = "multiAdd"
	ActionMultiDel               = "multiDel"
	ActionFolderAdd              = "folderAdd"
	ActionFolderDel              = "folderDel"
	ActionSortFavs               = "sortFavs"
)

type Message struct {
	Field         string    `json:"field,omitempty"`
	Action        string    `json:"action,omitempty"`
	Oid           int64     `json:"oid,omitempty"`
	Otype         int8      `json:"otype,omitempty"`
	Type          int8      `json:"type,omitempty"`
	Mid           int64     `json:"mid,omitempty"`
	OldMid        int64     `json:"old_mid,omitempty"`
	Fid           int64     `json:"fid,omitempty"`
	FidState      int8      `json:"fid_state,omitempty"`
	FolderAttr    int32     `json:"folder_attr,omitempty"`
	OldFolderAttr int32     `json:"old_folder_attr,omitempty"`
	NewFolderAttr int32     `json:"new_folder_attr,omitempty"`
	Aid           int64     `json:"aid,omitempty"`
	OldFid        int64     `json:"old_fid,omitempty"`
	OldFidState   int8      `json:"old_fid_state,omitempty"`
	NewFid        int64     `json:"new_fid,omitempty"`
	NewFidState   int8      `json:"new_fid_state,omitempty"`
	Affected      int64     `json:"affected,omitempty"`
	Aids          []int64   `json:"aids,omitempty"`
	Oids          []int64   `json:"oids,omitempty"`
	Mids          []int64   `json:"mids,omitempty"`
	FTime         int64     `json:"ftime,omitempty"`
	SortFavs      []SortFav `json:"sort_favs,omitempty"`
}

type SortFav struct {
	Pre    *Resource `json:"preID,omitempty"`
	Insert *Resource `json:"id,omitempty"`
}
