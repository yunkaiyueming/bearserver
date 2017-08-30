package db

import (
	"fmt"
	"strings"
)

type ModelBase struct {
	BatchCnt           int // 默认100
	Keys               []string
	FieldsDefaultValue []map[string]interface{}
}

func (this *ModelBase) UpsertBatch(game, table string, data []map[string]interface{}, onlyUpdateKeys []string) bool {
	var filtedData []string
	for _, row := range data {
		keysValue, flag := keysPrepared(row, this.Keys)
		if !flag {
			return false
		}

		fieldsValue := getFieldsValue(row, this.FieldsDefaultValue)

		filtedData = append(filtedData, "("+keysValue+","+fieldsValue+")")
	}

	fields := mergeKeysFields(this.Keys, this.FieldsDefaultValue)

	// 检查批量处理数
	if this.BatchCnt <= 0 {
		this.BatchCnt = 100
	}

	cnt := len(filtedData)
	var filted []string
	// Batch this baby
	for i := 0; i < cnt; i += this.BatchCnt {
		if i+this.BatchCnt > cnt {
			filted = filtedData[i:cnt]
		} else {
			filted = filtedData[i : i+this.BatchCnt]
		}

		sql := getUpsertSql(table, filted, onlyUpdateKeys, fields)
		//		fmt.Println(sql)

		getOrm(game).Raw(sql).Exec()
	}
	getOrm(game).Commit()

	return true
}

func keysPrepared(data map[string]interface{}, keys []string) (string, bool) {
	var str string
	for _, v := range keys {
		if _, ok := data[v]; !ok {
			return "", false
		}

		str += fmt.Sprintf("'%v',", data[v])
	}

	return str[0 : len(str)-1], true
}

func getFieldsValue(data map[string]interface{}, fieldsDefaultValue []map[string]interface{}) string {
	var str string
	for _, field := range fieldsDefaultValue {
		for key, val := range field {
			if v, ok := data[key]; ok {
				str += fmt.Sprintf("'%v',", v)
			} else {
				str += fmt.Sprintf("'%v',", val)
			}
		}
	}

	return str[0 : len(str)-1]
}

func mergeKeysFields(keys []string, fields []map[string]interface{}) []string {
	var ret []string
	for _, field := range fields {
		for k, _ := range field {
			ret = append(ret, k)
		}
	}

	return append(keys, ret...)
}

func getUpsertSql(table string, data []string, onlyUpdateKeys, fields []string) string {
	var updateFields []string

	// 为空时更新全部
	if onlyUpdateKeys == nil {
		onlyUpdateKeys = fields
	}
	for _, onlykey := range onlyUpdateKeys {
		updateFields = append(updateFields, onlykey+"=VALUES("+onlykey+")")
	}

	return fmt.Sprintf("INSERT INTO `%s` (%s) VALUES %s ON DUPLICATE KEY UPDATE %s", table, strings.Join(fields, ","), strings.Join(data, ","), strings.Join(updateFields, ","))
}

func getSingleUpsertSql(table string, fields, values, onlyUpdateKeys []string) string {
	var updateFields []string

	// 为空时更新全部
	if onlyUpdateKeys == nil {
		onlyUpdateKeys = fields
	}
	for _, onlykey := range onlyUpdateKeys {
		updateFields = append(updateFields, onlykey+"=VALUES("+onlykey+")")
	}

	return fmt.Sprintf("INSERT INTO `%s` (%s) VALUES (%s) ON DUPLICATE KEY UPDATE %s", table, strings.Join(fields, ","), strings.Join(values, ","), strings.Join(updateFields, ","))
}
