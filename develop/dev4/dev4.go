package main

import (
	"cmp"
	"fmt"
	"slices"
	"sort"
	"strings"
)

/*
Написать функцию поиска всех множеств анаграмм по словарю.


Например:
'пятак', 'пятка' и 'тяпка' - принадлежат одному множеству,
'листок', 'слиток' и 'столик' - другому.


Требования:
Входные данные для функции: ссылка на массив, каждый элемент которого - слово на русском языке в кодировке utf8
Выходные данные: ссылка на мапу множеств анаграмм
Ключ - первое встретившееся в словаре слово из множества. Значение - ссылка на массив, каждый элемент которого,
слово из множества.
Массив должен быть отсортирован по возрастанию.
Множества из одного элемента не должны попасть в результат.
Все слова должны быть приведены к нижнему регистру.
В результате каждое слово должно встречаться только один раз.

*/

var dictionaryArray = []string{"пятка", "ПЯтАК", "тяпка", "Слиток", "листок", "пятак", "т"}

var inputData = &dictionaryArray

func removeDupes(s []string) []string {
	uniques := make(map[string]struct{})

	for _, word := range s {
		uniques[word] = struct{}{}
	}
	result := make([]string, len(uniques))

	i := 0
	for key := range uniques {
		result[i] = key
		i++
	}

	return result
}

func anagrams(dict *[]string) map[string]*[]string {
	anagramsMap := make(map[string][]string)

	for _, word := range *dict {
		wordBytes := []rune(strings.ToLower(word))

		if len(wordBytes) < 2 {
			continue
		}

		sort.Slice(wordBytes, func(i, j int) bool {
			return wordBytes[i] < wordBytes[j]
		})
		sortedBytes := string(wordBytes)
		anagramsMap[sortedBytes] = append(anagramsMap[sortedBytes], strings.ToLower(word))
	}

	anagramsResult := make(map[string]*[]string)

	for _, value := range anagramsMap {

		value = removeDupes(value)

		slices.SortStableFunc(value, func(a, b string) int {
			return cmp.Compare(a, b)
		})
		anagramsResult[value[0]] = &value
	}

	return anagramsResult
}

func main() {
	result := anagrams(inputData)
	for key, val := range result {
		fmt.Println(key, *val)
	}
}
