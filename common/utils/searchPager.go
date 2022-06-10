package utils

import (
	"fmt"
	"github.com/zeromicro/go-zero/core/jsonx"
	"github.com/zeromicro/go-zero/core/stringx"
	"gorm.io/gorm"
	"strings"
)

/**
 * @Author: zze
 * @Date: 2022/5/23 11:22
 * @Desc: 分页插件
 */

type Search struct {
	searchMap map[string]map[string]interface{}
	fields    []string
	value     string
}

func NewSearch(searchFields []string, value, searchMapJsonStr string) (*Search, error) {
	var (
		err       error
		searchMap map[string]map[string]interface{}
	)
	if stringx.NotEmpty(searchMapJsonStr) {
		searchMap = make(map[string]map[string]interface{})
		if err = jsonx.UnmarshalFromString(searchMapJsonStr, &searchMap); err != nil {
			return nil, err
		}
	}

	return &Search{
		searchMap: searchMap,
		fields:    searchFields,
		value:     value,
	}, nil
}

type SearchPager struct {
	Search   *Search
	PageNum  int
	PageSize int
}

func NewSearchPager(search *Search, pageNum, pageSize int) *SearchPager {
	return &SearchPager{
		Search:   search,
		PageNum:  pageNum,
		PageSize: pageSize,
	}
}

func (p *SearchPager) Paginate() func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if p.PageNum == 0 {
			p.PageNum = 1
		}
		if p.PageSize == 0 {
			p.PageSize = 10
		}

		offset := (p.PageNum - 1) * p.PageSize
		return db.Offset(offset).Limit(p.PageSize)
	}
}

func (p *SearchPager) FieldsSearcher(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if p.Search == nil {
			return db
		}

		for _, field := range p.Search.fields {
			db = db.Or(field+" like ?", "%"+p.Search.value+"%")
		}
		return db
	}(db)
}

func (p *SearchPager) MapSearcher(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		for field, vObj := range p.Search.searchMap {
			if strings.TrimSpace(fmt.Sprintf("%v", vObj["value"])) == "" || vObj["value"] == nil {
				continue
			}

			if vObj["rule"] == "like" {
				db = db.Where(fmt.Sprintf("%s like ?", field), "%"+vObj["value"].(string)+"%")
			} else if vObj["rule"] == "in" {
				db = db.Where(fmt.Sprintf("%s in (?)", field), vObj["value"])
			} else if vObj["rule"] == "not_in" {
				db = db.Where(fmt.Sprintf("%s not in (?)", field), vObj["value"])
			} else if vObj["rule"] == "between" {
				db = db.Where(fmt.Sprintf("%s between ? and ?", field), vObj["value"].(map[string]int)["start"], vObj["value"].(map[string]int)["end"])
			} else if vObj["rule"] == "not_between" {
				db = db.Where(fmt.Sprintf("%s not between ? and ?", field), vObj["value"].(map[string]int)["start"], vObj["value"].(map[string]int)["end"])
			} else if vObj["rule"] == "is_null" {
				db = db.Where(fmt.Sprintf("%s is null", field))
			} else if vObj["rule"] == "is_not_null" {
				db = db.Where(fmt.Sprintf("%s is not null", field))
			} else if vObj["rule"] == "is_empty" {
				db = db.Where(fmt.Sprintf("%s = ''", field))
			} else if vObj["rule"] == "is_not_empty" {
				db = db.Where(fmt.Sprintf("%s != ''", field))
			} else if vObj["rule"] == "not_equal" {
				db = db.Where(fmt.Sprintf("%s != ?", field), vObj["value"])
			} else {
				db = db.Where(fmt.Sprintf("%s = ?", field), vObj["value"])
			}
		}
		return db
	}(db)
}
