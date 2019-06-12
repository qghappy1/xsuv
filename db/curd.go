
package db

import (
	"fmt"
	"errors"
	"database/sql"
)

type DataRow map[string]interface{};
 
type RowCollection []DataRow;
 
type DataColumn []string;
 
type DataTable struct {
    Rows RowCollection;
    Columns DataColumn;
}

//  
func GetTable(rows *sql.Rows) (*DataTable,error){
    dt := new(DataTable)
    columns ,err := rows.Columns();
    if(err !=nil){
        return dt,err;
    }
    dt.Columns = columns;
    count := len(columns);
    values := make([]interface{},count);
    valuePtrs := make([]interface{},count);
 
    for rows.Next(){
        for i:=0;i< count;i++{
            valuePtrs[i] = &values[i];
        }
        rows.Scan(valuePtrs...);
        entry := make(map[string]interface{});
        for i, col := range columns {
            var v interface{}
 
            val := values[i]
            b, ok := val.([]byte)
            if ok {
                v = string(b)
            } else {
                v = val
            }
            entry[col] = v
        }
        dt.Rows = append(dt.Rows,entry);
    }
    return dt,nil;
}

// 查询
func SqlSelect(db *sql.DB, seq string) (dt *DataTable, err error) {
	rows, err := db.Query(seq)
	if err != nil {
		e := fmt.Sprintf("seq1:%s err:%s", seq, err.Error())
		err = errors.New(e)
		return dt, err
	}
	defer rows.Close()
	
	dt, err = GetTable(rows)
	if err != nil {
		e := fmt.Sprintf("seq2:%s err:%s", seq, err.Error())
		err = errors.New(e)
		return dt, err
	}
	return dt, err
}

// 插入
func SqlInsert(db *sql.DB, seq string) (int64, error) {
	ret, err := db.Exec(seq)
	if err != nil {
		e := fmt.Sprintf("seq1:%s err:%s", seq, err.Error())
		err = errors.New(e)
		return -1, err
	}
	id, err2 := ret.LastInsertId()
	if err2 != nil {
		e := fmt.Sprintf("seq2:%s err:%s", seq, err2.Error())
		err = errors.New(e)
		return -1, err2	
	}
	return id, nil 
}

// 插入
func SqlUpdate(db *sql.DB, seq string) (int64, error) {
	ret, err := db.Exec(seq)
	if err != nil {
		e := fmt.Sprintf("seq1:%s err:%s", seq, err.Error())
		err = errors.New(e)
		return -1, err
	}
	id, err2 := ret.RowsAffected()
	if err2 != nil {
		e := fmt.Sprintf("seq2:%s err:%s", seq, err2.Error())
		err = errors.New(e)
		return -1, err2	
	}
	return id, nil 
}

// 删除
func SqlDelete(db *sql.DB, seq string) (int64, error) {
	ret, err := db.Exec(seq)
	if err != nil {
		e := fmt.Sprintf("seq1:%s err:%s", seq, err.Error())
		err = errors.New(e)
		return -1, err
	}
	id, err2 := ret.RowsAffected()
	if err2 != nil {
		e := fmt.Sprintf("seq2:%s err:%s", seq, err2.Error())
		err = errors.New(e)
		return -1, err2	
	}
	return id, nil 
}

