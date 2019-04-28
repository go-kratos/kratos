package service

import (
	"regexp"

	"go-common/library/log"
)

func (ins *tidbInstance) check() (err error) {
	for _, db := range ins.config.Databases {
		for _, ctable := range db.CTables {
			if _, err = regexp.Compile(ctable.Name); err != nil {
				log.Error("regexp.Compile(%s) error(%v)", ctable.Name, err)
				return
			}
		}
	}
	return
}

func (ins *tidbInstance) getTable(dbName, table string) *Table {
	if ins.ignoreTables[dbName] != nil && ins.ignoreTables[dbName][table] {
		return nil
	}
	if ins.tables[dbName] != nil && ins.tables[dbName][table] != nil {
		return ins.tables[dbName][table]
	}
	var regex *regexp.Regexp
	for _, db := range ins.config.Databases {
		if db.Schema != dbName {
			continue
		}
		for _, ctable := range db.CTables {
			regex, _ = regexp.Compile(ctable.Name)
			if !regex.MatchString(table) {
				continue
			}
			if ins.tables[dbName] == nil {
				ins.tables[dbName] = make(map[string]*Table)
			}
			t := &Table{
				PrimaryKey: ctable.PrimaryKey,
				OmitField:  make(map[string]bool),
				OmitAction: make(map[string]bool),
				name:       ctable.Name,
				ch:         make(chan *msg, 1024),
			}
			for _, action := range ctable.OmitAction {
				t.OmitAction[action] = true
			}
			for _, field := range ctable.OmitField {
				t.OmitField[field] = true
			}
			ins.waitTable.Add(1)
			go ins.proc(t.ch)
			ins.tables[dbName][table] = t
			return t
		}
	}
	if ins.ignoreTables[dbName] == nil {
		ins.ignoreTables[dbName] = make(map[string]bool)
	}
	ins.ignoreTables[dbName][table] = true
	return nil
}
