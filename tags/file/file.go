// package file provides a ValueSetter for file imput
// note different behavior for File compared to others

// Package file provides a tags.ValueSetter for file input.
// The tag value is used as the filename.
// Values of type *os.File will be set normally.
// Other values will be depending on the file type.
// File type is derived from the file extension, or optionally overridden by a tag option.
// Supported types include: txt, json, xml, gob
package file
import (
	"github.com/go-modules/modules/tags"
	"reflect"
	"os"
	"path/filepath"
	"fmt"
	"io/ioutil"
	"encoding/json"
	"encoding/xml"
	"encoding/gob"
)

var ValueSetter = tags.ValueSetterFunc(valueSetterFunc)

func valueSetterFunc(value reflect.Value, tagValue string) (bool, error) {
	if typeOfFile.AssignableTo(value.Type()) || typeOfFile.ConvertibleTo(value.Type()) {
		file, err := os.Open(tagValue)
		if err != nil {
			return false, err
		}
		fmt.Printf("%s\n", value)
		value.Elem().Set(reflect.ValueOf(file).Elem())
		return true, nil
	} else {
		tag, optionalType := tags.ParseTag(tagValue)
		file, err := os.Open(tag)
		if err != nil {
			return false, err
		}
		fileType := string(optionalType)
		if fileType == "" {
			fileType = filepath.Ext(tag)[1:]
		}
		if fileType == "" {
			return false, fmt.Errorf("no extension or type option given for file: %s", tag)
		}
		switch fileType {
		case "txt":
			if bytes, err := ioutil.ReadAll(file); err != nil {
				return false, err
			} else {
				value.SetString(string(bytes))
			}
		case "json":
			if err := json.NewDecoder(file).Decode(value.Interface()); err != nil {
				return false, err
			}
		case "xml":
			if err := xml.NewDecoder(file).Decode(value.Interface()); err != nil {
				return false, err
			}
		case "gob":
			if err := gob.NewDecoder(file).Decode(value.Interface()); err != nil {
				return false, err
			}
		default:
			return false, fmt.Errorf("unable to read file %s, unrecognized file type %s", tag, fileType)
		}
		return true, nil
	}
}

var typeOfFile = reflect.TypeOf(new(os.File))