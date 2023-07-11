package Zelig_IniFile

/*
   缩略的  在go语言下实现的 ini 操作,纯属为了自己的习惯上使用
   icy
   zelig.cn
*/

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
)

type TIniValue struct {
	Value   string
	Comment string
}

type TIniSection struct {
	items map[string]TIniValue
	Name  string
}

func NewIniSection() TIniSection {
	return TIniSection{
		items: make(map[string]TIniValue),
		Name:  "",
	}
}

func (this *TIniSection) Clear() {
	this.Name = ""
	this.items = make(map[string]TIniValue, 0)
}

func (this *TIniSection) HasIdent(Ident string) bool {
	if _, ok := this.items[Ident]; ok {
		return ok
	}

	return false
}

func (this *TIniSection) ToJson() string {
	if len(this.items) > 0 {
		data := make(map[string]string)
		for k, v := range this.items {
			data[k] = v.Value
		}
		if buff, err := json.Marshal(data); err == nil {
			return string(buff)
		}
	}
	return "{}"
}

func (this *TIniSection) Int(Ident string, defval int64) int64 {
	if this.HasIdent(Ident) {
		if v, err := strconv.ParseInt(this.items[Ident].Value, 10, 64); err == nil {
			return v
		}
	}

	return defval
}

func (this *TIniSection) Float(Ident string, defflo float64) float64 {
	if this.HasIdent(Ident) {
		if v, err := strconv.ParseFloat(this.items[Ident].Value, 64); err == nil {
			return v
		}
	}

	return defflo
}

func (this *TIniSection) String(Ident string, defstr string) string {
	if this.HasIdent(Ident) {
		return this.items[Ident].Value
	}
	return defstr
}

func (this *TIniSection) Bool(Ident string, defbool bool) bool {
	if this.HasIdent(Ident) {
		if v, err := strconv.ParseBool(this.items[Ident].Value); err == nil {
			return v
		}
	}
	return defbool
}

func (this *TIniSection) Add(Ident string, data TIniValue) {
	this.items[Ident] = data
}

func (this *TIniSection) setValue(Ident string, value interface{}) {
	data := TIniValue{
		Value:   "",
		Comment: "",
	}
	switch reflect.ValueOf(value).Kind() {
	case reflect.Bool:
		data.Value = func() string {
			if value.(bool) {
				return `1`
			} else {
				return `0`
			}
		}()
	case reflect.Int64:
		data.Value = strconv.FormatInt(value.(int64), 10)
	case reflect.Uint64:
		data.Value = strconv.FormatUint(value.(uint64), 10)
	case reflect.Float64:
		data.Value = strconv.FormatFloat(value.(float64), 'f', 2, 64)
	case reflect.String:
		data.Value = value.(string)
	}

	this.items[Ident] = data
}

func (this *TIniSection) SetInt(Ident string, Value int64) {
	this.setValue(Ident, Value)
}

func (this *TIniSection) SetString(Ident, Value string) {
	this.setValue(Ident, Value)
}

func (this *TIniSection) SetFloat(Ident string, Value float64) {
	this.setValue(Ident, Value)
}

func (this *TIniSection) SetBool(Ident string, Value bool) {
	if Value {
		this.setValue(Ident, int64(1))
	} else {
		this.setValue(Ident, int64(0))
	}
}

func (this *TIniSection) DeleteKey(Ident string) {
	if this.HasIdent(Ident) {
		delete(this.items, Ident)
	}
}

func (this *TIniSection) Comment(Ident string) string {
	if this.HasIdent(Ident) {
		return this.items[Ident].Comment
	}
	return ""
}

func (this *TIniSection) SetComment(Ident, Comment string) {
	if this.HasIdent(Ident) {
		P := this.items[Ident]
		P.Comment = Comment
		this.items[Ident] = P
	}
}

type iniKeyValue struct {
	key   string
	value string
}

func getIniKeyValue(v string) iniKeyValue {
	item := strings.Split(v, `,`)
	if len(item) > 1 {
		return iniKeyValue{
			key:   item[0],
			value: item[1],
		}
	}

	return iniKeyValue{
		key:   v,
		value: "",
	}
}

func getFinalValue[T any](Default string, fv T) string {
	var pType any = fv
	typ := reflect.TypeOf(pType)
	res := ``

	switch typ.Kind() {
	case reflect.Int64:
		res = strconv.FormatInt(pType.(int64), 10)
	case reflect.Float64:
		res = strconv.FormatFloat(pType.(float64), 'f', 2, 64)
	case reflect.Uint64:
		res = strconv.FormatUint(pType.(uint64), 10)
	case reflect.String:
		res = pType.(string)
	case reflect.Bool:
		res = func() string {
			if pType.(bool) {
				return `1`
			}
			return `0`
		}()

	}

	if res == "" && Default != "" {
		return Default
	}

	return res
}

func (this *TIniSection) SetStruct(value interface{}) error {
	val := reflect.ValueOf(value)
	if val.Kind() == reflect.Struct {
		size := val.NumField()
		if size == 0 {
			return fmt.Errorf(`这是一个空的结构体`)
		}

		t := reflect.TypeOf(value)
		for i := 0; i < size; i++ {
			v := val.Field(i)

			kv := getIniKeyValue(t.Field(i).Tag.Get(`ini`))
			if kv.key == "" {
				kv.key = t.Field(i).Name
			}

			switch v.Kind() {
			case reflect.Bool:
				this.setValue(kv.key, getFinalValue[bool](kv.value, v.Bool()))
			case reflect.String:
				this.setValue(kv.key, getFinalValue[string](kv.value, v.String()))
			case reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64:
				this.setValue(kv.key, getFinalValue[int64](kv.value, v.Int()))
			case reflect.Uint, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				this.setValue(kv.key, getFinalValue[uint64](kv.value, v.Uint()))
			case reflect.Float32, reflect.Float64:
				this.setValue(kv.key, getFinalValue[float64](kv.value, v.Float()))
			default:
				continue
			}
		}
	} else {
		return fmt.Errorf(`不是一个有效的结构体`)
	}

	return nil
}

func (this *TIniSection) Struct(value interface{}) error {
	vtype := reflect.TypeOf(value).Elem() //获取type的真实类型
	val := reflect.ValueOf(value).Elem()

	if vtype.Kind() == reflect.Struct {
		size := vtype.NumField()
		if size == 0 {
			return fmt.Errorf(`这是一个空的结构体`)
		}
		if len(this.items) > 0 {
			doSet := func(field string, p TIniValue) {
				for i := 0; i < size; i++ {
					t := vtype.Field(i)
					s := t.Tag.Get(`ini`)
					next := false
					if len(s) > 0 {
						next = strings.EqualFold(s, field)
					} else {
						next = strings.EqualFold(t.Name, field)
					}
					if next {
						v := val.Field(i)
						switch v.Kind() {
						case reflect.Bool:
							v.SetBool(this.Bool(field, false))
						case reflect.String:
							v.SetString(this.String(field, ""))
						case reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64:
							v.SetInt(this.Int(field, 0))
						case reflect.Uint, reflect.Uint16, reflect.Uint32, reflect.Uint64:
							v.SetUint(uint64(this.Int(field, 0)))
						case reflect.Float32, reflect.Float64:
							v.SetFloat(this.Float(field, 0))
						}
					}
				}
			}

			for k, v := range this.items {
				doSet(k, v)
			}
		} else {
			return fmt.Errorf(`段内是控的.`)
		}
	}

	return nil
}

func (this *TIniSection) IdentNames() []string {
	names := make([]string, 0)
	if len(this.items) > 0 {
		for k, _ := range this.items {
			names = append(names, k)
		}
	}
	return names
}

func (this *TIniSection) Values() []string {
	values := make([]string, 0)
	if len(this.items) > 0 {
		for _, v := range this.items {
			values = append(values, v.Value)
		}
	}

	return values
}

type TIniConfig struct {
	Sections []TIniSection
}

func NewIniConfig() *TIniConfig {
	return &TIniConfig{Sections: make([]TIniSection, 0)}
}

func (this *TIniConfig) ReadFromBytes(data []byte) int {
	Buildr := strings.Builder{}
	Comment := strings.Builder{}
	var HasSection, HasComment bool
	BeginIndex := 0
	Section := NewIniSection()

	doAdd := func() {
		item := strings.Split(Buildr.String(), "=")
		if len(item) == 2 {
			Section.Add(strings.Trim(item[0], " "), TIniValue{
				Value:   strings.Trim(item[1], " "),
				Comment: Comment.String(),
			})
		}

		Buildr = strings.Builder{}
		Comment = strings.Builder{}
	}

	doAddSection := func() {
		if HasSection {
			doAdd()
			this.Sections = append(this.Sections, Section)
			Section.Clear()
		}

	}

	for i, c := range data {
		switch c {
		case 91:
			{
				doAddSection()
				BeginIndex = i
			}
			break
		case 93:
			{
				HasSection = true
				HasComment = false
				Section.Name = string(data[BeginIndex+1 : i])
			}
			break
		case 59:
			HasComment = true
			break

		case 10:
			{
				doAdd()
				Buildr = strings.Builder{}
				HasComment = false
			}
			break
		case 13:
			break
		default:
			if HasSection && !HasComment {
				Buildr.WriteByte(c)
			}

			if HasComment {
				Comment.WriteByte(c)
			}
			break
		}
	}

	doAddSection()
	return len(this.Sections)
}

func (this *TIniConfig) ReadFromString(value string) int {
	return this.ReadFromBytes([]byte(value))
}

func (this *TIniConfig) ReadFromFile(fileName string) int {
	if buff, err := os.ReadFile(fileName); err == nil {
		return this.ReadFromBytes(buff)
	}

	return 0
}

func (this *TIniConfig) GetSection(name string) (*TIniSection, error) {
	if len(this.Sections) > 0 {
		for i, p := range this.Sections {
			if strings.EqualFold(p.Name, name) {
				return &this.Sections[i], nil
			}
		}
	}

	return nil, fmt.Errorf(`未找到段落`)
}

func (this *TIniConfig) Struct(name string, data interface{}) error {
	if sec, err := this.GetSection(name); err == nil {
		return sec.Struct(data)
	} else {
		return err
	}
}

func (this *TIniConfig) SetStruct(name string, data interface{}) error {
	Sec, err := this.GetSection(name)
	if err == nil {
		return Sec.SetStruct(data)
	}

	return err
}

func (this *TIniConfig) SaveToBytes(data *[]byte) {
	var text string
	this.SaveToString(&text)
	*data = []byte(text)
}

func (this *TIniConfig) SaveToString(text *string) {
	if len(this.Sections) > 0 {
		Builder := strings.Builder{}
		for _, sec := range this.Sections {
			Builder.WriteString("[" + sec.Name + "]\n")
			if len(sec.items) > 0 {
				for k, v := range sec.items {
					Builder.WriteString(k + " = " + v.Value)
					if v.Comment != "" {
						Builder.WriteString(" ; " + v.Comment)
					}
					Builder.WriteString("\n")
				}
			}
		}

		*text = Builder.String()
	}
}

func (this *TIniConfig) SaveToFile(fileName string) {
	if f, err := os.Create(fileName); err == nil {
		defer f.Close()

		var text string
		this.SaveToString(&text)
		f.WriteString(text)
	} else {
		fmt.Println(err)
	}
}

func (this *TIniConfig) DeleteSection(index int) {
	this.Sections = append(this.Sections[:index], this.Sections[index+1:]...)
}

func (this *TIniConfig) ClearSection() {
	this.Sections = make([]TIniSection, 0)
}

func (this *TIniConfig) SectionNames() []string {
	names := make([]string, 0)
	if len(this.Sections) > 0 {
		for _, p := range this.Sections {
			names = append(names, p.Name)
		}
	}
	return names
}

func (this *TIniConfig) AddSection(name string) *TIniSection {
	this.Sections = append(this.Sections, NewIniSection())
	sec := &this.Sections[len(this.Sections)-1]
	sec.Name = name
	return sec
}

func (this *TIniConfig) HasSection(name string) bool {
	names := this.SectionNames()
	if len(names) > 0 {
		for _, v := range names {
			if v == name {
				return true
			}
		}
	}

	return false
}
