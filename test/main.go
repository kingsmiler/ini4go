package main

import (
    "reflect"
    "fmt"
)

func main() {
    myMap := map[string]string{
        "a" : "one",
        "b": "tow",
        "c":"three",
        "d":"four",
    }

    item := "c"
    if _, exists := myMap[item]; exists {
        delete(myMap, item)
    }


    DeleteInMap(myMap, "a")

    fmt.Println(myMap)
}

/**
 * 删除 slice 中的指定元素。
 *
 */
func DeleteInMap(targetMap interface{}, item interface{}) {

    value := reflect.Indirect(reflect.ValueOf(targetMap))

    if value.MapIndex(reflect.ValueOf(item)).IsValid() {
        value.SetMapIndex(reflect.ValueOf(item), reflect.Value{})
    }
}

/**
 * 删除 slice 中的指定元素。
 *
 */
func DeleteInSlice(slice interface{}, item interface{}) {
    index := -1
    ve := reflect.Indirect(reflect.ValueOf(slice))
    size := ve.Len()

    for i := 0; i < size; i++ {
        if reflect.DeepEqual(ve.Index(i).Interface(), item) {
            index = i
            break
        }
    }

    if index >= 0 {
        ve.Set(reflect.AppendSlice(ve.Slice(0, index), ve.Slice(index + 1, size)))
    }
}

/**
 * 检查 obj 是否存在于 collection 中， collection 类型可以为 slice、array和map。
 * 如果存在，返回其下标（slice和array），或者 1 （map）；
 * 如果不存在，返回-1。
 */
func Contains(collection interface{}, obj interface{}) int {
    targetValue := reflect.ValueOf(collection)
    existed := -1
    switch reflect.TypeOf(collection).Kind() {
    case reflect.Slice, reflect.Array:
        for i := 0; i < targetValue.Len(); i++ {
            if reflect.DeepEqual(targetValue.Index(i).Interface(), obj) {
                existed = i
                break
            }
        }
    case reflect.Map:
        if targetValue.MapIndex(reflect.ValueOf(obj)).IsValid() {
            existed = 1
        }
    }

    return existed
}