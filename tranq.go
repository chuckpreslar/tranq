package tranq

import (
	"fmt"
	"reflect"
)

const (
	id = "ID"
)

// InvalidKindError is a error which is returned when an unexpected
// reflect.Kind is encounted.
type InvalidKindError struct {
	Kind reflect.Kind
}

// Error implements the `error` interface.
func (i InvalidKindError) Error() string {
	return fmt.Sprintf("an unsupported `reflect.Kind` was encoutered, was `%s`", i.Kind)
}

// UninterfaceabledValueError when a reflect.Value cannot have
// it's `Interface` method called without panicking.
type UninterfaceabledValueError struct {
	Value reflect.Value
}

// Error implements the `error` interface.
func (i UninterfaceabledValueError) Error() string {
	return fmt.Sprintf("failed to call `Interface` method on `reflect.Value` of `%v`", i.Value)
}

// Dereference attempts to dereference the provided paramter `i`
// from reflect.Kind's of reflect.Ptr and reflect.Interface.
func Dereference(i interface{}) (reflect.Value, reflect.Kind, reflect.Type, error) {
	var (
		v reflect.Value
		k reflect.Kind
		t reflect.Type
	)

	v = reflect.ValueOf(i)
	k = v.Kind()

	if k == reflect.Ptr || k == reflect.Interface {
		v = v.Elem()

		if v.CanInterface() {
			return Dereference(v.Interface())
		}

		return reflect.Value{}, reflect.Invalid, nil, UninterfaceabledValueError{v}
	}

	t = v.Type()

	return v, k, t, nil
}

// TypeName attempts to resolve parameter `i`'s type name,
// returning an InvalidKindError if `i` cannot be dereferenced
// into a reflect.Kind of reflect.Struct, reflect.Slice or reflect.Array.
func TypeName(i interface{}) (string, error) {
	var (
		o bool
		t reflect.Type
		k reflect.Kind
		e error
	)

	if t, o = i.(reflect.Type); !o {
		_, k, t, e = Dereference(i)
		if nil != e {
			return "", e
		}
	} else {
		k = t.Kind()
	}

	if k == reflect.Struct {
		return t.Name(), nil
	} else if k == reflect.Array || k == reflect.Slice || k == reflect.Ptr {
		return TypeName(t.Elem())
	}

	return "", InvalidKindError{k}
}

func isBaseKind(kind reflect.Kind) bool {
	switch kind {
	case reflect.Bool:
	case reflect.Int:
	case reflect.Int8:
	case reflect.Int16:
	case reflect.Int32:
	case reflect.Int64:
	case reflect.Uint:
	case reflect.Uint8:
	case reflect.Uint16:
	case reflect.Uint32:
	case reflect.Uint64:
	case reflect.Float32:
	case reflect.Float64:
	case reflect.Complex64:
	case reflect.Complex128:
	case reflect.String:
	default:
		return false
	}

	return true
}

func isStructKind(kind reflect.Kind) bool {
	return kind == reflect.Struct
}

func isCollectionKind(kind reflect.Kind) bool {
	return kind == reflect.Slice || kind == reflect.Array
}

func shouldDescend(kind reflect.Kind, currentDepth, maxDepth int) bool {
	if isBaseKind(kind) {
		return true
	}

	return ((isStructKind(kind) || isCollectionKind(kind)) && (currentDepth+1 <= maxDepth))
}

type compiler struct {
	*Tranq
	typeName string
}

func (cr compiler) compile(i interface{}, c, m int) (interface{}, error) {
	var v, k, t, e = Dereference(i)

	if nil != e {
		return nil, e
	}

	if isBaseKind(k) {
		if v.CanInterface() {
			return v.Interface(), nil
		}

		return nil, UninterfaceabledValueError{v}
	} else if isStructKind(k) {
		return cr.compileStruct(v, t, c, m)
	} else if isCollectionKind(k) {
		return cr.compileCollection(v, c, m)
	}

	panic(InvalidKindError{k})
}

func (cr compiler) compileStruct(v reflect.Value, t reflect.Type, c, m int) (interface{}, error) {
	var payload = make(Payload)

	for i := 0; i < v.NumField(); i++ {
		var (
			fv = v.Field(i)
			fs = t.Field(i)
			ft = fv.Type()
			fk = fv.Kind()
		)

		if cr.strategy.ShouldSkipStructField(fs) {
			continue
		} else if !fv.CanInterface() {
			return nil, UninterfaceabledValueError{fv}
		} else if cr.strategy.ShouldLinkStructField(fs) {
			if err := cr.linkStructField(payload, fv, ft, fk, fs); nil != err {
				return nil, err
			}
		} else if shouldDescend(fk, c, m) {
			var (
				n           = cr.strategy.FormatAttributeName(fs.Name)
				result, err = cr.compile(fv.Interface(), c+1, m)
			)

			if nil != err {
				return nil, err
			}

			payload[n] = result
		}
	}

	return payload, nil
}

func (cr compiler) linkStructField(p Payload, v reflect.Value, t reflect.Type, k reflect.Kind, f reflect.StructField) error {
	var result, err = cr.compile(v.Interface(), 0, cr.strategy.GetMaxLinkDepth())

	if nil != err {
		return err
	}

	var l = Link{
		Interface:   result,
		Value:       v,
		Type:        t,
		Kind:        k,
		StructField: f,
		IDFormat:    cr.id,
	}

	return cr.strategy.LinkStructField(p, l)
}

func (cr compiler) compileCollection(v reflect.Value, c, m int) (interface{}, error) {
	var collection = make([]interface{}, 0, 0)

	for i := 0; i < v.Len(); i++ {
		var e = v.Index(i)

		if e.CanInterface() {
			var result, err = cr.compile(e.Interface(), c+1, m)

			if nil != err {
				return nil, err
			}

			collection = append(collection, result)
		} else {
			panic(UninterfaceabledValueError{e})
		}
	}

	return collection, nil
}

// Tranq allows for JSON serialization based on a Strategy intended
// to follow the JSON API standard.
type Tranq struct {
	strategy Strategy
	id       string
}

// CompilePayload uses the Tranq instances Strategy
// to serialize a Payload.
func (tq *Tranq) CompilePayload(i interface{}) (Payload, error) {
	var t, e = TypeName(i)

	if nil != e {
		return nil, e
	}

	var (
		p = make(Payload)
		m = tq.strategy.GetMaxMapDepth()
		n = tq.strategy.GetTopLevelNamespace(t)
		c = compiler{tq, t}
	)

	tq.strategy.SetPayloadRoot(p)
	if p[n], e = c.compile(i, 0, m); nil != e {
		return nil, e
	}

	return p, nil
}

// New returns a pointer to a Tranq struct.
func New(s Strategy) (t *Tranq) {
	t = new(Tranq)
	t.strategy = s
	t.id = s.FormatAttributeName(id)

	return
}
