package helper

import "fmt"

func MoveCurrentIndexToFront(arr []string, index int) []string {
	if index < 0 || index >= len(arr) {
		fmt.Println("Error, index is out of range")
		return arr
	}

	currentElement := arr[index]
	newArr := []string{currentElement}
	newArr = append(newArr, arr[:index]...)
	newArr = append(newArr, arr[index+1:]...)

	return newArr;
}