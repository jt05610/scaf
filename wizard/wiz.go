package wizard

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/jt0610/scaf/context"
	"go.uber.org/zap"
	"os"
	"reflect"
	"strconv"
)

type Wizard[T any] struct{}

func Prompt[T any](ctx context.Context) (*T, error) {
	var t T
	ptr := reflect.New(reflect.TypeOf(t))
	value := ptr.Elem()
	tType := value.Type()
	reader := bufio.NewReader(os.Stdin)
	for i := 0; i < value.NumField(); i++ {
		field := tType.Field(i)
		prompt := field.Tag.Get("prompt")
		if prompt != "" {
			ctx.Logger.Info("prompting", zap.String("field", tType.Field(i).Name))
			defaultValue := field.Tag.Get("default")
			options := field.Tag.Get("options")
			fmt.Printf("%s: ", prompt)
			if options != "" {
				fmt.Printf("[%s] ", options)
			}
			text, _ := reader.ReadString('\n')
			if text == "" {
				if defaultValue != "" {
					ctx.Logger.Info("using default value", zap.String("field", tType.Field(i).Name), zap.String("value", defaultValue))
					text = defaultValue
				} else {
					return nil, errors.New("no default value provided")
				}
			} else {
				ctx.Logger.Info("got value", zap.String("field", tType.Field(i).Name), zap.String("value", text))
			}
			err := Set(ctx, value.Field(i), field, field.Name, text)
			if err != nil {
				return nil, err
			}
		} else {
			ctx.Logger.Info("skipping", zap.String("field", tType.Field(i).Name))
		}
	}
	return ptr.Interface().(*T), nil
}

func Set(ctx context.Context, v reflect.Value, f reflect.StructField, name, value string) error {
	switch f.Type.Kind() {
	case reflect.String:
		ctx.Logger.Info("setting string", zap.String("field", name), zap.String("value", value))
		v.SetString(value)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		ctx.Logger.Info("setting int", zap.String("field", name), zap.String("value", value))
		val, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return err
		}
		v.SetInt(val)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		ctx.Logger.Info("setting uint", zap.String("field", name), zap.String("value", value))
		val, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return err
		}
		v.SetUint(val)
	case reflect.Bool:
		ctx.Logger.Info("setting bool", zap.String("field", name), zap.String("value", value))
		if value == "y" || value == "yes" || value == "true" {
			v.SetBool(true)
		} else {
			v.SetBool(false)
		}
	case reflect.Float32, reflect.Float64:
		ctx.Logger.Info("setting float", zap.String("field", name), zap.String("value", value))
		val, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return err
		}
		v.SetFloat(val)
	default:
		return errors.New(fmt.Sprintf("invalid type %s for field %s", f.Type.Kind(), name))
	}
	return nil
}

func (w *Wizard[T]) Run(ctx context.Context) (*T, error) {
	ctx.Logger.Info("starting wizard")
	return Prompt[T](ctx)
}
