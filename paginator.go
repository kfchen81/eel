package eel

import (
	"reflect"
	

	"github.com/bitly/go-simplejson"
	"strconv"
	"encoding/json"
	"github.com/kfchen81/gorm"
	"github.com/kfchen81/eel/log"
)

//INextPageInfo
type INextPageInfo interface {
	ToMap() map[string]interface{}
}

//PaginateResult 分页的结果
type PaginateResult struct {
	HasPrev      bool
	HasNext      bool
	HasHead      bool
	HasTail      bool
	Prev         int
	Next         int
	CurPage      int
	MaxPage      int
	TotalCount   int64
	DisplayPages []int
	
	Offset      int
	CountInPage int
}

func (this PaginateResult) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"has_head":      this.HasHead,
		"has_tail":      this.HasTail,
		"has_prev":      this.HasPrev,
		"has_next":      this.HasNext,
		"next":          this.Next,
		"prev":          this.Prev,
		"total_count":   this.TotalCount,
		"cur_page":      this.CurPage,
		"display_pages": this.DisplayPages,
		"max_page":      this.MaxPage,
	}
}

func NewPaginateResultFromData(js *simplejson.Json) INextPageInfo {
	displayPages := make([]int, 0)
	for _, page := range js.Get("display_pages").MustArray() {
		pageValue, _ := strconv.Atoi(page.(json.Number).String())
		displayPages = append(displayPages, pageValue)
	}
	
	prev, _ := js.Get("prev").Int()
	next, _ := js.Get("next").Int()
	curPage, _ := js.Get("cur_page").Int()
	maxPage, _ := js.Get("max_page").Int()
	totalCount, _ := js.Get("total_count").Int64()
	return &PaginateResult{
		HasHead: js.Get("has_head").MustBool(),
		HasTail: js.Get("has_tail").MustBool(),
		HasPrev: js.Get("has_prev").MustBool(),
		HasNext: js.Get("has_next").MustBool(),
		Prev: prev,
		Next: next,
		TotalCount: totalCount,
		CurPage: curPage,
		DisplayPages: displayPages,
		MaxPage: maxPage,
	}
}

func getTotalPageCount(itemCount int64, itemCountPerPage int) int {
	var totalPage int64
	_itemCountPerPage := int64(itemCountPerPage)
	if itemCount%_itemCountPerPage == 0 {
		totalPage = itemCount / _itemCountPerPage
		if totalPage == 0 {
			totalPage = 1
		}
	} else {
		totalPage = (itemCount / _itemCountPerPage) + 1
	}
	
	return int(totalPage)
}

func getOffset(curPage int, itemCountPerPage int) int {
	return (curPage - 1) * itemCountPerPage
}

func _range(start int, end int) []int {
	result := make([]int, 0, 20)
	for index := start; index <= end; index++ {
		result = append(result, index)
	}
	
	return result
}

//Paginate 进行分页
func doPaginate(itemCount int64, curPage int, itemCountPerPage int) INextPageInfo {
	paginateResult := PaginateResult{}
	
	paginateResult.TotalCount = itemCount
	total := getTotalPageCount(itemCount, itemCountPerPage)
	
	//如果浏览页数超过最大页数，则显示最后一页数据
	if curPage > total {
		curPage = total
	}
	
	paginateResult.MaxPage = total
	paginateResult.CurPage = curPage
	
	if curPage == total {
		paginateResult.HasTail = false
	} else {
		paginateResult.HasTail = true
	}
	
	if curPage == 1 {
		paginateResult.HasPrev = false
		paginateResult.HasHead = false
	} else {
		paginateResult.Prev = curPage - 1
		paginateResult.HasHead = true
		paginateResult.HasPrev = true
	}
	
	if curPage >= total {
		paginateResult.HasNext = false
	} else {
		paginateResult.Next = curPage + 1
		paginateResult.HasNext = true
	}
	
	//计算需要显示的页数序列
	if paginateResult.MaxPage <= 5 {
		paginateResult.DisplayPages = _range(1, paginateResult.MaxPage+1)
	} else if curPage+2 <= paginateResult.MaxPage {
		if curPage >= 3 {
			paginateResult.DisplayPages = _range(curPage-2, curPage+3)
		} else {
			paginateResult.DisplayPages = _range(1, 6)
		}
	} else if paginateResult.MaxPage == 0 {
		paginateResult.DisplayPages = _range(1, paginateResult.MaxPage+1)
	} else {
		if curPage >= 5 {
			paginateResult.DisplayPages = _range(paginateResult.MaxPage-5, paginateResult.MaxPage+1)
		}
	}
	
	//获取当前page应包含的对象列表
	paginateResult.Offset = (curPage - 1) * itemCountPerPage
	paginateResult.CountInPage = itemCountPerPage
	
	return paginateResult
}

//Mobile模式的分页结果
type APIServiceNextPageInfo struct {
	HasNext    bool
	NextFromId int64
}

func (this APIServiceNextPageInfo) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"has_next":     this.HasNext,
		"next_from_id": this.NextFromId,
	}
}

//PaginateAndFill 进行分页，并获取填充数据
//func PaginateAndFill(objects orm.QuerySeter, curPage int, itemCountPerPage int, container interface{}) (INextPageInfo, error) {
//	nextPageInfo := doPaginate(objects, curPage, itemCountPerPage)
//	_, err := objects.Limit(nextPageInfo.(PaginateResult).CountInPage).Offset(nextPageInfo.(PaginateResult).Offset).All(container)
//
//	return nextPageInfo, err
//}

//PaginateAndFill 进行分页，并获取填充数据
func Paginate(db *gorm.DB, page *PageInfo, container interface{}) (INextPageInfo, error) {
	var err error
	var nextPageInfo INextPageInfo
	if page.IsApiServerMode() {
		//多取1个，用于进行分页判断
		_ = db.Limit(page.CountPerPage + 1).All(container)
		val := reflect.ValueOf(container)
		ind := reflect.Indirect(val)
		
		realItemCount := ind.Len()
		
		if realItemCount > page.CountPerPage {
			val = ind.Index(realItemCount - 2)
			
			lastItemId := val.Elem().FieldByName("Id").Int()
			nextPageInfo = &APIServiceNextPageInfo{
				HasNext:    true,
				NextFromId: lastItemId,
			}
			
			slice := reflect.MakeSlice(ind.Type(), 0, 0)
			for i := 0; i < realItemCount-1; i++ {
				slice = reflect.Append(slice, ind.Index(i))
			}
			ind.Set(slice)
		} else {
			nextPageInfo = &APIServiceNextPageInfo{
				HasNext:    false,
				NextFromId: -1,
			}
		}
	} else {
		itemCount, err := db.Count()
		if err != nil {
			log.Logger.Error(err)
			panic(err)
		}
		nextPageInfo = doPaginate(itemCount, page.Page, page.CountPerPage)
		err = db.Limit(nextPageInfo.(PaginateResult).CountInPage).Offset(nextPageInfo.(PaginateResult).Offset).All(container)
		return nextPageInfo, err
	}
	
	return nextPageInfo, err
}

//MockPaginate 模拟进行分页
func MockPaginate(itemCount int64, page *PageInfo) INextPageInfo {
	nextPageInfo := doPaginate(itemCount, page.Page, page.CountPerPage)
	return nextPageInfo
}
