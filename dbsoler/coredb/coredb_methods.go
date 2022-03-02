package coredb

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"time"
)

func DatToMap(value reflect.Value) (data map[string]string) {
	adItem := reflect.Indirect(value)
	adType := adItem.Type()

	if adType.Kind() == reflect.Ptr {
		adItem = adItem.Elem()
		adType = adItem.Type()
	}

	numFields := adType.NumField()

	data = make(map[string]string, 0)

	for index := 0; index < numFields; index++ {
		vf := adItem.Field(index)
		sf := adType.Field(index)

		switch vf.Kind() {
		case reflect.Int64, reflect.Int32, reflect.Int16:
			{
				data[sf.Name] = strconv.FormatInt(vf.Int(), 10)
			}

		case reflect.Int8, reflect.Uint8:
			{
				data[sf.Name] = string(vf.Bytes())
			}

		case reflect.String:
			{
				data[sf.Name] = vf.String()
			}

		case reflect.Struct:
			{
				if sf.Type == reflect.TypeOf(time.Time{}) {
					dtTime := vf.Interface().(time.Time)
					data[sf.Name] = dtTime.Format("02/01/2006")
				}
			}

		default:
			{
				data[sf.Name] = fmt.Sprintf("%v", vf)
			}

		}

	}

	return data
}

// GetRecordSize 456456456
func GetRecordSize(adType reflect.Type) (size int) {
	size = 0

	adRef := reflect.ValueOf(reflect.New(adType))
	kind := adRef.Kind()

	if kind != reflect.Struct {
		if kind == reflect.Ptr {
			adRef = adRef.Elem()
		} else {
			panic("source kind is unhandled")
		}
	}

	for index := 0; index < adType.NumField(); index++ {
		field := adType.Field(index)
		fieldLen, _ := strconv.Atoi(field.Tag.Get("length"))

		size += fieldLen
	}

	return size
}

// Unmarshal asdasd
func Unmarshal(ptr interface{}, data []byte, ptrType reflect.Type) (err error) {

	var offset int = 0

	vp := reflect.ValueOf(ptr)

	if vp.Kind() == reflect.Ptr {
		vp = reflect.Indirect(vp)
	}

	if vp.Kind() != reflect.Struct {
		panic("not struct")
	}

	numFields := ptrType.NumField()

	for index := 0; index < numFields && offset < len(data); index++ {
		vf := vp.Field(index)
		sf := ptrType.Field(index)

		fieldLen, err := strconv.Atoi(sf.Tag.Get("length"))
		if err != nil {
			panic(fmt.Sprintf("Could not read Tag from field '%s'", sf.Name))
		}

		fieldStart := offset
		fieldEnd := offset + fieldLen
		fieldData := data[fieldStart:fieldEnd]

		switch vf.Kind() {
		case reflect.String:
			{
				buffer := bytes.NewBuffer(fieldData)
				fieldValue, _ := buffer.ReadString(0x00)
				vf.SetString(fieldValue)
			}

		case reflect.Int:
			{
				fieldValue := binary.LittleEndian.Uint64(fieldData)
				vf.SetInt(int64(fieldValue))
			}

		case reflect.Int32:
			{
				fieldValue := binary.LittleEndian.Uint32(fieldData)
				vf.SetInt(int64(fieldValue))
			}

		case reflect.Int16:
			{
				fieldValue := binary.LittleEndian.Uint16(fieldData)
				vf.SetInt(int64(fieldValue))
			}

		case reflect.Int8:
			{
				fieldValue := int8(data[fieldStart])
				vf.SetInt(int64(fieldValue))
			}

		case reflect.Uint8:
			{
				fieldValue := uint8(data[fieldStart])
				vf.SetUint(uint64(fieldValue))
			}

		case reflect.Struct:
			{
				// Provavelmente time.Time e afins
				if sf.Type == reflect.TypeOf(time.Time{}) {
					fieldValue := binary.LittleEndian.Uint32(fieldData)
					var dtTime time.Time
					if fieldValue > 0 {
						dtTime = time.Unix(int64(fieldValue), 0)
					} else {
						dtTime, _ = time.Parse("02/01/2006", "31/12/3000")
					}

					vf.Set(reflect.ValueOf(dtTime))

				} else {
					panic(fmt.Sprintf("Don't know how to handle '%s'", sf.Name))
				}
			}

		default:
			{
				panic(fmt.Sprintf("Don't know how to handle '%s'", vf.Kind().String()))
			}
		}

		offset += fieldLen
	}

	return err
}

// ReadFile adsadasdasd
func ReadFile(fileName string, from int) (dados []reflect.Value, err error, limit int) {

	fileParts := strings.Split(filepath.Base(fileName), ".")
	adType := GetStructType(fileParts[0])

	if adType == nil {
		return nil, fmt.Errorf("Arquivo '%s' ainda nao implementado",
			filepath.Base(fileName)), 0
	}

	file, err := os.Open(fileName)
	if err != nil {
		return nil, err, 0
	}
	defer file.Close()

	recordSize := GetRecordSize(adType)
	data := make([]byte, recordSize)

	dados = make([]reflect.Value, 0)

	if from > 0 {
		file.Seek(int64((recordSize+32)*from), 1)
	}

	limit = 0
	for {
		count, err := file.Read(data)
		if err != nil {
			break
		}

		if count != int(recordSize) {
			continue
		}

		st := reflect.New(adType)
		err = Unmarshal(st.Interface(), data, adType)
		if err == nil {
			dados = append(dados, st)
		}

		// Skips 32 bytes nao sei porque
		file.Seek(32, 1)

		limit++

		if limit > 64 {
			break
		}
	}

	return dados, nil, limit
}
