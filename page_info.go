package eel

import (
	"strconv"
	"github.com/kfchen81/eel/handler"
)

//PageInfo 指示当前查询的数据的page信息
type PageInfo struct {
	Page         int
	FromId       int
	CountPerPage int
	Mode         string
	Direction    string
}

//type MobilePageInfo struct {
//	FromId int
//	CountPerPage int
//}

func (self *PageInfo) IsApiServerMode() bool {
	return self.Mode == "apiserver"
}

func (self *PageInfo) Desc() *PageInfo {
	self.Direction = "desc"
	if self.FromId == 0 {
		self.FromId = 99999999999
	}
	return self
}

func (self *PageInfo) Asc() *PageInfo {
	self.Direction = "asc"
	return self
}

func getInt(req *handler.Request, key string, def ...int) (int, error) {
	strv := req.Query(key)
	if len(strv) == 0 && len(def) > 0 {
		return def[0], nil
	}
	return strconv.Atoi(strv)
}

//ExtractPageInfoFromRequest 从Request中抽取page信息
func ExtractPageInfoFromRequest(ctx *handler.Context) *PageInfo {
	req := ctx.Request
	fromParam := req.Query("_p_from")
	if fromParam != "" {
		fromId, _ := getInt(req, "_p_from")
		countPerPage, _ := getInt(req, "_p_count", 20)
		return &PageInfo{
			Page:         -1,
			FromId:       fromId,
			CountPerPage: countPerPage,
			Mode:         "apiserver",
			Direction:    "asc",
		}
	} else {
		page, _ := getInt(req, "page", 1)
		countPerPage, _ := getInt(req, "count_per_page", 20)
		return &PageInfo{
			Page:         page,
			FromId:       0,
			CountPerPage: countPerPage,
			Mode:         "backend",
			Direction:    "asc",
		}
	}
}
