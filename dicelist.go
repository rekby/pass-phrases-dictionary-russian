// Command dicelist превращает список слов-кандидатов в нумерованный
// словарь для составления парольных фраз бросками кубика (diceware).
//
// На вход подаётся файл, где слова разделены запятыми и/или переносами
// строк (как candidates.txt). Скрипт:
//   - приводит все слова к нижнему регистру;
//   - ставит каждое слово на отдельную строку;
//   - берёт наибольшее число слов вида 6^k и нумерует их k-значными
//     кодами кубика (цифры 1..6), по одной цифре на бросок;
//   - оставшиеся слова, не влезшие в степень шести, выводит в блок
//     «Запасные».
//
// Зависимостей кроме стандартной библиотеки нет.
//
// Использование:
//
//	go run dicelist.go [-i candidates.txt] [-o dictionary.txt]
//
// По умолчанию читает candidates.txt и пишет в stdout.
package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
	"unicode"
)

func main() {
	in := flag.String("i", "candidates.txt", "входной файл со словами-кандидатами")
	out := flag.String("o", "", "выходной файл (по умолчанию stdout)")
	flag.Parse()

	data, err := os.ReadFile(*in)
	if err != nil {
		fmt.Fprintln(os.Stderr, "не удалось прочитать входной файл:", err)
		os.Exit(1)
	}

	words := parseWords(string(data))
	if len(words) == 0 {
		fmt.Fprintln(os.Stderr, "во входном файле не найдено ни одного слова")
		os.Exit(1)
	}

	k := dicePower(len(words))
	used := pow6(k)

	w := os.Stdout
	if *out != "" {
		f, err := os.Create(*out)
		if err != nil {
			fmt.Fprintln(os.Stderr, "не удалось создать выходной файл:", err)
			os.Exit(1)
		}
		defer f.Close()
		w = f
	}

	bw := bufio.NewWriter(w)

	// Основной словарь: 6^k слов с кодами кубика.
	for i := range used {
		fmt.Fprintf(bw, "%s %s\n", diceCode(i, k), words[i])
	}

	// Запасные слова, не влезшие в степень шести.
	if used < len(words) {
		fmt.Fprintf(bw, "\nЗапасные (%d):\n", len(words)-used)
		for _, word := range words[used:] {
			fmt.Fprintln(bw, word)
		}
	}

	if err := bw.Flush(); err != nil {
		fmt.Fprintln(os.Stderr, "ошибка записи вывода:", err)
		os.Exit(1)
	}

	fmt.Fprintf(os.Stderr,
		"всего слов: %d; в словаре: %d (6^%d, коды по %d броска(ов)); в запасе: %d\n",
		len(words), used, k, k, len(words)-used)
}

// parseWords разбивает текст на слова по запятым и пробельным символам,
// приводит их к нижнему регистру и убирает пустые токены.
func parseWords(s string) []string {
	fields := strings.FieldsFunc(s, func(r rune) bool {
		return r == ',' || unicode.IsSpace(r)
	})
	words := make([]string, 0, len(fields))
	for _, f := range fields {
		f = strings.TrimSpace(f)
		if f != "" {
			words = append(words, strings.ToLower(f))
		}
	}
	return words
}

// dicePower возвращает наибольшее k, при котором 6^k <= n.
func dicePower(n int) int {
	k := 0
	for pow6(k+1) <= n {
		k++
	}
	return k
}

// pow6 возвращает 6 в степени k.
func pow6(k int) int {
	p := 1
	for range k {
		p *= 6
	}
	return p
}

// diceCode преобразует индекс i в k-значный код кубика (цифры 1..6).
// i=0 -> "11..1", i=1 -> "11..2" и т.д.
func diceCode(i, k int) string {
	digits := make([]byte, k)
	for pos := k - 1; pos >= 0; pos-- {
		digits[pos] = byte('1' + i%6)
		i /= 6
	}
	return string(digits)
}
