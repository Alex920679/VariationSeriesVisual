package main

import (
	"bufio"
	"fmt"
	v1charts "github.com/go-echarts/go-echarts/charts"
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/xuri/excelize/v2"

	"log"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
)

func round(x, unit float64) float64 {
	return math.Round(x/unit) * unit
}

func readInput() string {
	textIn, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil {
		fmt.Print("Don't give up!")
		log.Fatal(err)
	}
	txt := strings.TrimSpace(textIn)
	return txt
}

func fillMap(intSlice []int) map[int]int {
	mp := map[int]int{}
	for _, elem := range intSlice {
		mp[elem] += 1
	}
	return mp
}

func runFillGetStr(mp map[int]int, uniqueSlice []int, serSize int) *SetFull {
	var structure1 Str
	var structure2 SetFull
	for _, elem := range uniqueSlice {
		ans := mp[elem]
		structure1.fillStruct(elem, ans, serSize)
		structure2.set = append(structure2.set, structure1)
	}
	return &structure2
}

func getSlice(txt string) []int {
	intSlice := make([]int, 0, 5)
	txtSlice := strings.Split(txt, "; ")
	for _, elem := range txtSlice {
		strInt, err := strconv.Atoi(elem)
		if err != nil {
			fmt.Println("Houston, we have a problem:", err)
			log.Fatal(err)
		}
		intSlice = append(intSlice, strInt)
	}
	return intSlice
}

func removeDuplicates(sl []int) []int {
	var res []int
	var el int
	for _, elem := range sl {
		if elem != el {
			el = elem
			res = append(res, elem)
		}
	}
	return res
}

type StrCriteria struct {
	max     int
	preMax  int
	postMin int
	min     int
}

func newStore() (s *StrCriteria) {
	s = &StrCriteria{}
	return s
}

func (s *StrCriteria) fillCriteria(sl []int) {
	lenSl := len(sl)
	s.max = sl[lenSl-1]
	s.preMax = sl[lenSl-2]
	s.min = sl[0]
	s.postMin = sl[1]
}

func (s *StrCriteria) getCriteria() (float64, float64) {
	res1 := float64(s.max-s.preMax) / float64(s.max-s.postMin)
	res2 := float64(s.postMin-s.min) / float64(s.preMax-s.min)
	return round(res1, 0.0005), round(res2, 0.0005)
}

var checkCriteria = map[int]float64{ // табличные данные для критериев. Максимальная выборка - 30
	4:  0.955,
	5:  0.807,
	6:  0.689,
	7:  0.610,
	8:  0.554,
	9:  0.512,
	10: 0.477,
	11: 0.450,
	12: 0.428,
	13: 0.410,
	14: 0.395,
	15: 0.381,
	16: 0.369,
	17: 0.359,
	18: 0.349,
	19: 0.341,
	20: 0.334,
	21: 0.327,
	22: 0.320,
	23: 0.314,
	24: 0.309,
	25: 0.304,
	26: 0.299,
	27: 0.295,
	28: 0.291,
	29: 0.287,
	30: 0.283,
}

type Str struct {
	series     int
	freq       int
	relVarFreq float64
}

type SetFull struct {
	set []Str
}

func (st *Str) fillStruct(ch1, ch2, lenSl int) {
	st.series = ch1
	st.freq = ch2
	st.relVarFreq = round((float64(ch2)/float64(lenSl))*100, 0.05)
}

func checker(structure2 SetFull) (int, float64) {
	sumFreq := 0
	sumRelValFreq := 0.0
	for _, elem := range structure2.set {
		sumFreq += elem.freq
		sumRelValFreq += elem.relVarFreq
	}
	return sumFreq, sumRelValFreq
}

func checkCrt1(crt1 float64, serSize int) bool {
	var b1 bool
	if crt1 > checkCriteria[serSize] {
		b1 = true
	}
	return b1
}

func checkCrt2(crt2 float64, serSize int) bool {
	var b2 bool
	if crt2 > checkCriteria[serSize] {
		b2 = true
	}
	return b2
}

func getRelValFreqSl(s *SetFull) []float64 {
	relValFreqSlice := make([]float64, 0, 5)
	for _, elem := range s.set {
		relValFreqSlice = append(relValFreqSlice, elem.relVarFreq)
	}
	return relValFreqSlice
}

func generateBarItems(lenUniqueSl int, relValFreqSlice []float64) []opts.BarData {
	items := make([]opts.BarData, 0)
	for i := 0; i < lenUniqueSl; i++ {
		items = append(items, opts.BarData{Value: relValFreqSlice[i]})
	}
	return items
}

func drawBarChart(uniqueSlice []int, relValFreqSlice []float64) { // рисовальщик гистограммы
	lenUniqueSlice := len(uniqueSlice)
	bar := charts.NewBar()
	bar.SetGlobalOptions(charts.WithTitleOpts(opts.Title{
		Title: "Гистограмма данных",
	}))
	bar.SetXAxis(uniqueSlice).
		AddSeries("Category A", generateBarItems(lenUniqueSlice, relValFreqSlice))
	f, err := os.Create("Гистограмма_данных.html")
	if err != nil {
		fmt.Println("Something went wrong", err)
		log.Fatal(err)
	}
	errN := bar.Render(f)
	if errN != nil {
		fmt.Println("Something went wrong", errN)
		log.Fatal(errN)
	}
}

func drawLine(uniqueSlice []int, relValFreqSlice []float64) *v1charts.Line { // рисовальщик полигона распределения
	line := v1charts.NewLine()
	line.SetGlobalOptions(v1charts.TitleOpts{Title: "Полигон распределения"})
	line.AddXAxis(uniqueSlice).AddYAxis("ОЧВ = f(В)", relValFreqSlice)
	return line
}

func renderLine(line *v1charts.Line) {
	f, err := os.Create("Полигон распределения.html")
	if err != nil {
		fmt.Println("Something went wrong", err)
		log.Fatal(err)
	}
	errN := line.Render(f)
	if errN != nil {
		fmt.Println("Something went wrong", errN)
		log.Fatal(errN)
	}
}

func createExcelTable(uniqueSlice []int, mp map[int]int, relValFreqSlice []float64) {
	f := excelize.NewFile()
	digitOfCell := 2
	// прописываю стиль заголовков
	style, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold:   true,
			Italic: false,
			Family: "Times New Roman",
			Size:   12,
			Color:  "#777777",
		},
	})
	if err != nil {
		fmt.Println(err)
	} // прописал, далее заполняю ячейки
	f.SetCellValue("Sheet1", "B2", "Вариант")
	f.SetCellValue("Sheet1", "C2", "Частота")
	f.SetCellValue("Sheet1", "D2", "ОЧВ")
	// наделяю заголовки стилем style
	err = f.SetCellStyle("Sheet1", "B2", "B2", style)
	err = f.SetCellStyle("Sheet1", "C2", "C2", style)
	err = f.SetCellStyle("Sheet1", "D2", "D2", style)
	// далее заполняю таблицу
	for i := 0; i < len(uniqueSlice); i++ {
		digitOfCell++
		d := strconv.Itoa(digitOfCell)
		s1 := "B" + d
		f.SetCellValue("Sheet1", s1, uniqueSlice[i])
		s2 := "C" + d
		f.SetCellValue("Sheet1", s2, mp[uniqueSlice[i]])
		s3 := "D" + d
		f.SetCellValue("Sheet1", s3, relValFreqSlice[i])
	}
	if err := f.SaveAs("Таблица данных.xlsx"); err != nil {
		fmt.Println(err)
	}
}

/*
	образцы вариационного ряда -

11; 8; 9; 10; 8; 6; 7; 7; 9; 11; 10; 6; 5; 11; 10; 7; 9; 11; 10
1; 2; 1; 2; 2; 4; 3; 3; 25; 2; 1; 25
5; 8; 8; 9; 9; 5; 8; 5; 2; 7; 9; 5; 2; 10; 8; 7; 5; 7; 3; 9; 5; 3; 9; 7
*/
func main() {
	fmt.Print("Введите вариационный ряд из положительных чисел.\n")
	fmt.Println("Обратите внимание, как должны отделяться числа друг от друга, например\n1; 2; 3")
	txt := readInput()        // считываем вариационный ряд как строку
	intSlice := getSlice(txt) // преобразуем в интовый слайс
	sort.Ints(intSlice)
	mp := fillMap(intSlice)                   // заполняем мапу, ключ: вариант, значение: частота его повторения
	serSize := len(intSlice)                  // serSize - это объем выборки
	uniqueSlice := removeDuplicates(intSlice) // получаем массив уникальных вариантов (которые не повторяются)
	sort.Ints(uniqueSlice)
	structure3 := newStore()             // конструктор структуры типа *StrCriteria
	structure3.fillCriteria(uniqueSlice) // высчитываем критерии отброса максимального и минимального значения
	fmt.Println("ВНИМАНИЕ, ОТВЕТ!")
	fmt.Printf("Объем выборки составил: %d.", serSize)
	fmt.Println("\nТаблица `вариант:частота:относительная частота варианта в процентах` сконвертирована программой в Excel.")
	structure2 := runFillGetStr(mp, uniqueSlice, serSize) // заполняем структру, содержащую данные о Варианте, его Частоте и относительной частоте варианта
	sumFreq, sumRelValFreq := checker(*structure2)        // проверка
	fmt.Printf("Объем выбоки %v должен быть равен %v, а сумма относительных частот %.0f равна 100.\n", serSize, sumFreq, sumRelValFreq)
	fmt.Println("Если что-то не так, сообщите разработчику!\nГистограмма данных и полигон распределения готовы!")
	// подача данных в рисовальщик гистограммы
	relValFreqSlice := getRelValFreqSl(structure2)
	drawBarChart(uniqueSlice, relValFreqSlice)
	// подача данных в эксель
	createExcelTable(uniqueSlice, mp, relValFreqSlice)
	// подача далее в рисовальщик полигона распределения
	line := drawLine(uniqueSlice, relValFreqSlice)
	renderLine(line)
	// подача закончена
	crt1, crt2 := structure3.getCriteria() // получаем высчитанные ранее критерии
	b1 := checkCrt1(crt1, serSize)         // проверяем первый критерий
	if b1 == true {
		fmt.Printf("Критерий К1 = %.4f > табличного критерия Кт = %.4f для объема выборки %v,\nпоэтому следует исключить вариант %v из вариационного ряда и прогнать обновленный ряд через программу.", crt1, checkCriteria[serSize], serSize, structure3.max)
	}
	b2 := checkCrt2(crt2, serSize) // проверяем второй критерий
	if b2 == true {
		fmt.Printf("Критерий К2 = %.4f > табличного критерия Кт = %.4f для объема выборки %v,\nпоэтому следует исключить вариант %v из вариационного ряда и прогнать обновленный ряд через программу.", crt2, checkCriteria[serSize], serSize, structure3.min)
	}
	if b1 == true || b2 == true {
		fmt.Println("\nЖду новый вариационный ряд!")
	} else {
		fmt.Print("Исключать варианты из вариационного ряда не нужно, ")
		fmt.Printf("т.к. твой критерий К1 = %.4f и К2 = %.4f меньше табличого %.4f.", crt1, crt2, checkCriteria[serSize])
	}
}
