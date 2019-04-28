package model

// const const value.
const (
	MngActionDelete = string("delete")
	MngActionIgnore = string("ignore")
	MngActionPunish = string("punish")
	MngActionHandle = string("handle")
	MngActionAdd    = string("add")

	MngModuleReport   = string("report")
	MngModuleSynonym  = string("synonym")
	MngModuleLimit    = string("whitelist")
	MngModuleResource = string("resource")
	MngModuleTag      = string("tag")
	MngModuleRelation = string("relation")
	MngModuleHot      = string("hot_tag")

	MngTypeNone = int(-1)

	MngOidNone = int64(-1)
)
