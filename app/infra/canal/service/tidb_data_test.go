package service

import (
	"encoding/json"
	"reflect"
	"testing"

	"go-common/app/infra/canal/model"

	pb "github.com/pingcap/tidb-tools/tidb_binlog/slave_binlog_proto/go-binlog"
)

func Test_tidbMakeData(t *testing.T) {
	insertMsg, insertData := prepareInsertData()
	delMsg, delData := prepareDeleteData()
	updateMsg, updateData := prepareUpdateData()
	updateMsg2, updateData2 := prepareUpdateData2()
	type args struct {
		m *msg
	}
	tests := []struct {
		name     string
		args     args
		wantData *model.Data
		wantErr  bool
	}{
		{name: "insert", args: args{m: insertMsg}, wantData: insertData, wantErr: false},
		{name: "delete", args: args{m: delMsg}, wantData: delData, wantErr: false},
		{name: "update", args: args{m: updateMsg}, wantData: updateData, wantErr: false},
		{name: "update2", args: args{m: updateMsg2}, wantData: updateData2, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotData, err := tidbMakeData(tt.args.m)
			if (err != nil) != tt.wantErr {
				t.Errorf("tidbMakeData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			gotjson, _ := json.Marshal(gotData)
			wantjson, _ := json.Marshal(tt.wantData)
			if !reflect.DeepEqual(gotjson, wantjson) {
				t.Errorf("tidbMakeData() = %v, want %v", gotData, tt.wantData)
			}
		})
	}
}

func prepareInsertData() (*msg, *model.Data) {
	insertPb := &pb.Binlog{}
	json.Unmarshal([]byte(`{"type":0,"commit_ts":403846216359608325,"dml_data":{"tables":[{"schema_name":"bilibili_likes","table_name":"likes","column_info":[{"name":"id","mysql_type":"bigint","is_primary_key":false},{"name":"mtime","mysql_type":"timestamp","is_primary_key":false},{"name":"ctime","mysql_type":"timestamp","is_primary_key":false},{"name":"business_id","mysql_type":"int","is_primary_key":false},{"name":"origin_id","mysql_type":"bigint","is_primary_key":false},{"name":"message_id","mysql_type":"bigint","is_primary_key":false},{"name":"mid","mysql_type":"int","is_primary_key":false},{"name":"type","mysql_type":"tinyint","is_primary_key":false}],"mutations":[{"type":0,"row":{"columns":[{"uint64_value":1},{"string_value":"2018-10-26 18:50:57"},{"string_value":"2018-10-26 18:50:57"},{"uint64_value":5},{"uint64_value":0},{"uint64_value":1},{"uint64_value":8167601},{"uint64_value":1}]}}]}]}}`), insertPb)
	insertMsg := &msg{
		db:          "bilibili_likes",
		table:       "counts",
		tableRegexp: "counts",
		mu:          insertPb.DmlData.Tables[0].Mutations[0],
		ignore:      map[string]bool{"ctime": true},
		keys:        []string{"id", "mid"},
		columns:     insertPb.DmlData.Tables[0].ColumnInfo,
	}
	insertData := &model.Data{
		Action: "insert",
		Table:  "counts",
		Key:    "1,8167601",
		New: map[string]interface{}{
			"id":          1,
			"business_id": 5,
			"origin_id":   0,
			"message_id":  1,
			"mid":         8167601,
			"type":        1,
			"mtime":       "2018-10-26 18:50:57",
		},
	}
	return insertMsg, insertData
}

func prepareDeleteData() (*msg, *model.Data) {
	pbData := &pb.Binlog{}
	json.Unmarshal([]byte(`{"type":0,"commit_ts":403846189135953921,"dml_data":{"tables":[{"schema_name":"bilibili_likes","table_name":"likes","column_info":[{"name":"id","mysql_type":"bigint","is_primary_key":false},{"name":"mtime","mysql_type":"timestamp","is_primary_key":false},{"name":"ctime","mysql_type":"timestamp","is_primary_key":false},{"name":"business_id","mysql_type":"int","is_primary_key":false},{"name":"origin_id","mysql_type":"bigint","is_primary_key":false},{"name":"message_id","mysql_type":"bigint","is_primary_key":false},{"name":"mid","mysql_type":"int","is_primary_key":false},{"name":"type","mysql_type":"tinyint","is_primary_key":false}],"mutations":[{"type":2,"row":{"columns":[{"uint64_value":7},{"string_value":"2018-01-11 12:19:10"},{"string_value":"2018-01-11 12:19:10"},{"uint64_value":2},{"uint64_value":0},{"uint64_value":897},{"uint64_value":27515233},{"uint64_value":1}]}}]}]}}`), pbData)
	msg := &msg{
		db:          "bilibili_likes",
		table:       "counts",
		tableRegexp: "counts",
		mu:          pbData.DmlData.Tables[0].Mutations[0],
		ignore:      map[string]bool{"ctime": true},
		keys:        []string{"message_id"},
		columns:     pbData.DmlData.Tables[0].ColumnInfo,
	}
	data := &model.Data{
		Action: "delete",
		Table:  "counts",
		Key:    "897",
		New: map[string]interface{}{
			"id":          7,
			"business_id": 2,
			"origin_id":   0,
			"message_id":  897,
			"mid":         27515233,
			"type":        1,
			"mtime":       "2018-01-11 12:19:10",
		},
	}
	return msg, data
}

func prepareUpdateData() (*msg, *model.Data) {
	pbData := &pb.Binlog{}
	// update likes type from 1 to 0
	json.Unmarshal([]byte(`{"type":0,"commit_ts":403846165844459523,"dml_data":{"tables":[{"schema_name":"bilibili_likes","table_name":"likes","column_info":[{"name":"id","mysql_type":"bigint","is_primary_key":false},{"name":"mtime","mysql_type":"timestamp","is_primary_key":false},{"name":"ctime","mysql_type":"timestamp","is_primary_key":false},{"name":"business_id","mysql_type":"int","is_primary_key":false},{"name":"origin_id","mysql_type":"bigint","is_primary_key":false},{"name":"message_id","mysql_type":"bigint","is_primary_key":false},{"name":"mid","mysql_type":"int","is_primary_key":false},{"name":"type","mysql_type":"tinyint","is_primary_key":false}],"mutations":[{"type":1,"row":{"columns":[{"uint64_value":4},{"string_value":"2018-10-26 18:47:44"},{"string_value":"2017-12-22 15:05:29"},{"uint64_value":5},{"uint64_value":0},{"uint64_value":46997},{"uint64_value":88895031},{"uint64_value":0}]},"change_row":{"columns":[{"uint64_value":4},{"string_value":"2017-12-22 15:55:52"},{"string_value":"2017-12-22 15:05:29"},{"uint64_value":5},{"uint64_value":0},{"uint64_value":46997},{"uint64_value":88895031},{"uint64_value":1}]}}]}]}}`), pbData)
	msg := &msg{
		db:          "bilibili_likes",
		table:       "counts",
		tableRegexp: "counts",
		mu:          pbData.DmlData.Tables[0].Mutations[0],
		ignore:      map[string]bool{"ctime": true},
		keys:        []string{"mid"},
		columns:     pbData.DmlData.Tables[0].ColumnInfo,
	}
	data := &model.Data{
		Action: "update",
		Table:  "counts",
		Key:    "88895031",
		Old: map[string]interface{}{
			"id":          4,
			"business_id": 5,
			"origin_id":   0,
			"message_id":  46997,
			"mid":         88895031,
			"type":        1,
			"mtime":       "2017-12-22 15:55:52",
		},
		New: map[string]interface{}{
			"id":          4,
			"business_id": 5,
			"origin_id":   0,
			"message_id":  46997,
			"mid":         88895031,
			"type":        0,
			"mtime":       "2018-10-26 18:47:44",
		},
	}
	return msg, data
}

func prepareUpdateData2() (*msg, *model.Data) {
	muJson := `{"type":1,"row":{"columns":[{"uint64_value":0},{"string_value":"2018-11-03 17:07:44"},{"string_value":"2018-11-03 14:55:38"},{"uint64_value":3},{"uint64_value":0},{"uint64_value":88889},{"uint64_value":3},{"uint64_value":0},{"int64_value":0},{"int64_value":0},{"uint64_value":8167601}]},"change_row":{"columns":[{"uint64_value":0},{"string_value":"2018-11-03 16:36:39"},{"string_value":"2018-11-03 14:55:38"},{"uint64_value":3},{"uint64_value":0},{"uint64_value":88889},{"uint64_value":2},{"uint64_value":0},{"int64_value":0},{"int64_value":0},{"uint64_value":8167601}]}}`
	columnJson := `[{"name":"id","mysql_type":"bigint","is_primary_key":false},{"name":"mtime","mysql_type":"timestamp","is_primary_key":false},{"name":"ctime","mysql_type":"timestamp","is_primary_key":false},{"name":"business_id","mysql_type":"int","is_primary_key":false},{"name":"origin_id","mysql_type":"bigint","is_primary_key":false},{"name":"message_id","mysql_type":"bigint","is_primary_key":false},{"name":"likes_count","mysql_type":"int","is_primary_key":false},{"name":"dislikes_count","mysql_type":"int","is_primary_key":false},{"name":"likes_change","mysql_type":"bigint","is_primary_key":false},{"name":"dislikes_change","mysql_type":"bigint","is_primary_key":false},{"name":"up_mid","mysql_type":"int","is_primary_key":false}]`
	msg := &msg{
		db:          "bilibili_likes",
		table:       "counts",
		tableRegexp: "counts",
		keys:        []string{"message_id"},
	}
	json.Unmarshal([]byte(columnJson), &msg.columns)
	json.Unmarshal([]byte(muJson), &msg.mu)
	data := &model.Data{
		Action: "update",
		Table:  "counts",
		Key:    "88889",
		Old: map[string]interface{}{
			"ctime":           "2018-11-03 14:55:38",
			"origin_id":       0,
			"dislikes_count":  0,
			"up_mid":          8167601,
			"id":              0,
			"mtime":           "2018-11-03 16:36:39",
			"likes_count":     2,
			"likes_change":    0,
			"dislikes_change": 0,
			"business_id":     3,
			"message_id":      88889,
		},
		New: map[string]interface{}{
			"likes_count":     3,
			"dislikes_count":  0,
			"likes_change":    0,
			"id":              0,
			"mtime":           "2018-11-03 17:07:44",
			"ctime":           "2018-11-03 14:55:38",
			"origin_id":       0,
			"message_id":      88889,
			"business_id":     3,
			"dislikes_change": 0,
			"up_mid":          8167601,
		},
	}
	return msg, data
}
