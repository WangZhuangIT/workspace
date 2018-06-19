package dbdao

import (
	"github.com/go-xorm/xorm"
)

type DbBaseDao struct {
	Engine  *xorm.Engine
	Session *xorm.Session
}

type Param interface{}
type ParamNil struct{}
type ParamDesc bool
type ParamIn []interface{}
type ParamRange struct {
	Min interface{}
	Max interface{}
}
type ParamInDesc ParamIn
type ParamRangeDesc ParamRange

func CastToParamIn(input interface{}) ParamIn {
	params := make(ParamIn, 0)
	switch v := input.(type) {
	case []interface{}:
		for _, param := range v {
			params = append(params, param)
		}
	case []int64:
		for _, param := range v {
			params = append(params, param)
		}
	case []int:
		for _, param := range v {
			params = append(params, param)
		}
	case []int32:
		for _, param := range v {
			params = append(params, param)
		}
	case []int8:
		for _, param := range v {
			params = append(params, param)
		}
		return params
	case []string:
		for _, param := range v {
			params = append(params, param)
		}
	default:
		params = append(params, 0)
	}
	return params
}

func CastToParamInDesc(input interface{}) ParamInDesc {
	return ParamInDesc(CastToParamIn(input))
}

func (this *DbBaseDao) buildQuery(input Param, name string) {
	name = this.Engine.Quote(name)
	this.Session = this.Engine.Where(nil)
	switch val := input.(type) {
	case ParamDesc:
		if val {
			this.Session = this.Session.Desc(name)
		}
	case ParamIn:
		if len(val) == 1 {
			this.Session = this.Session.And(name+"=?", val[0])
		} else {
			this.Session = this.Session.In(name, val)
		}
	case ParamInDesc:
		if len(val) == 1 {
			this.Session = this.Session.And(name+"=?", val[0])
		} else {
			this.Session = this.Session.In(name, val)
		}
		this.Session = this.Session.Desc(name)
	case ParamRange:
		if val.Min != nil {
			this.Session = this.Session.And(name+">=?", val.Min)
		}
		if val.Max != nil {
			this.Session = this.Session.And(name+"<?", val.Max)
		}
	case ParamRangeDesc:
		if val.Min != nil {
			this.Session = this.Session.And(name+">=?", val.Min)
		}
		if val.Max != nil {
			this.Session = this.Session.And(name+"<?", val.Max)
		}
		this.Session = this.Session.Desc(name)
	case ParamNil:
	case nil:
	default:
		this.Session = this.Session.And(name+"=?", val)
	}
}

func (this *DbBaseDao) UpdateEngine(v ...interface{}) {
	if len(v) == 0 {
		this.Engine = GetDefault("reader").Engine
		this.Session = nil
	} else if len(v) == 1 {
		param := v[0]
		if engine, ok := param.(*xorm.Engine); ok {
			this.Engine = engine
			this.Session = nil
		} else if session, ok := param.(*xorm.Session); ok {
			this.Session = session
			this.Engine = nil
		} else if tpe, ok := param.(bool); ok {
			cluster := "reader"
			if tpe == true {
				cluster = "writer"
			}
			this.Engine = GetDefault(cluster).Engine
			this.Session = nil
		}
	}
}

func (this *DbBaseDao) Create(bean interface{}) (int64, error) {
	return this.Engine.Insert(bean)
}

func (this *DbBaseDao) Update(bean interface{}) (int64, error) {
	return this.Engine.Id(this.Engine.IdOf(bean)).AllCols().Update(bean)
}

func (this *DbBaseDao) UpdateCols(bean interface{}, cols ...string) (int64, error) {
	return this.Engine.Id(this.Engine.IdOf(bean)).Cols(cols...).Update(bean)
}

func (this *DbBaseDao) Delete(bean interface{}) (int64, error) {
	return this.Engine.Id(this.Engine.IdOf(bean)).Delete(bean)
}
