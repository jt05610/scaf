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

func Prompt[T any](ctx context.Context, t T) error {
	value := reflect.ValueOf(t)
	tType := value.Type()
	reader := bufio.NewReader(os.Stdin)
	for i := 0; i < value.NumField(); i++ {
		ctx.Logger.Info("prompting", zap.String("field", tType.Field(i).Name))
		field := tType.Field(i)
		prompt := field.Tag.Get("prompt")
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
				return errors.New("no default value provided")
			}
		} else {
			ctx.Logger.Info("got value", zap.String("field", tType.Field(i).Name), zap.String("value", text))
			err := Set(ctx, value.Field(i), field, field.Name, text)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func Set(ctx context.Context, v reflect.Value, f reflect.StructField, name, value string) error {
	switch f.Type.Kind() {
	case reflect.String:
		ctx.Logger.Info("setting string", zap.String("field", name), zap.String("value", value))
		v.SetString(value)
	case reflect.Int:
		ctx.Logger.Info("setting int", zap.String("field", name), zap.String("value", value))
		val, err := strconv.Atoi(value)
		if err != nil {
			return err
		}
		v.SetInt(int64(val))
	case reflect.Bool:
		ctx.Logger.Info("setting bool", zap.String("field", name), zap.String("value", value))
		if value == "y" || value == "yes" || value == "true" {
			v.SetBool(true)
		} else {
			v.SetBool(false)
		}
	default:
		return errors.New("invalid type")
	}
	return nil
}

func (w *Wizard[T]) Run(ctx context.Context, t T) error {
	ctx.Logger.Info("starting wizard")
	return Prompt[T](ctx, t)
}
