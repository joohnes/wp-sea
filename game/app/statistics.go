package app

import (
	"encoding/csv"
	"fmt"
	"os"
	"sort"
	"strconv"

	"github.com/fatih/color"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/joohnes/wp-sea/game/logger"
)

func (a *App) LoadStatistics() error {
	f, err := os.Open("statistics.csv")
	if err != nil {
		_, err = os.Create("statistics.csv")
		if err != nil {
			return err
		}
		f, err = os.Open("statistics.csv")
		if err != nil {
			return err
		}
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			logger.GetLoggerInstance().Error.Println(err)
		}
	}(f)
	r := csv.NewReader(f)
	read, err := r.Read()
	if err != nil {
		return err
	}
	a.games, err = strconv.Atoi(read[0])
	if err != nil {
		return err
	}
	a.won, err = strconv.Atoi(read[1])
	if err != nil {
		return err
	}

	records, err := r.ReadAll()
	if err != nil {
		logger.GetLoggerInstance().Error.Println("couldn't read statistics")
		return err
	}
	basemap := make(map[string]int)
	for _, x := range records {
		number, err := strconv.Atoi(x[1])
		if err != nil {
			logger.GetLoggerInstance().Error.Printf("couldn't load %v\n", x)
			return err
		}
		basemap[x[0]] = number
	}
	a.statistics = basemap
	return nil
}

func (a *App) SaveStatistics() {
	f, err := os.OpenFile("statistics.csv", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777)
	if err != nil {
		logger.GetLoggerInstance().Error.Println("couldn't open")
		return
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			logger.GetLoggerInstance().Error.Println(err)
		}
	}(f)
	w := csv.NewWriter(f)
	defer w.Flush()
	err = w.Write([]string{strconv.Itoa(a.games), strconv.Itoa(a.won)})
	if err != nil {
		logger.GetLoggerInstance().Error.Println("couldn't write")
		return
	}
	for coord, occurrences := range a.statistics {
		err := w.Write([]string{coord, strconv.Itoa(occurrences)})
		if err != nil {
			logger.GetLoggerInstance().Error.Printf("couldn't save %s, %d\n", coord, occurrences)
		}
	}
}

func (a *App) ShowStatistics() {
	min := getMin(a.statistics)
	max := getMax(a.statistics)

	t := table.NewWriter()
	t.SetTitle("Heatmap")
	t.AppendHeader(table.Row{"#", "A", "B", "C", "D", "E", "F", "G", "H", "I", "J"})
	for i := 1; i < 11; i++ {
		h := table.Row{i}
		for _, x := range []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j"} {
			number := strconv.Itoa(i)
			c := color.New(GetColor(a.statistics[x+number], min, max)).SprintFunc()
			h = append(h, c(strconv.Itoa(a.statistics[x+number])))
		}
		t.AppendRow(h)
	}
	fmt.Println(t.Render())
	fmt.Println("Press enter to go back to the menu")
	_, _ = fmt.Scanln()
}

func getMax(arr map[string]int) (max int) {
	for i := range arr {
		if arr[i] > max {
			max = arr[i]
		}
	}
	return
}

func getMin(arr map[string]int) (min int) {
	min = 99999999 // hope its enough
	for i := range arr {
		if arr[i] < min {
			min = arr[i]
		}
	}
	return
}

func GetColor(x, min, max int) color.Attribute {
	diff := (max - min) / 6
	if diff < 1 {
		diff = 1
	}
	var basemap []int
	colors := []color.Attribute{
		color.FgHiGreen,
		color.FgGreen,
		color.FgHiYellow,
		color.FgHiYellow,
		color.FgHiRed,
		color.FgRed,
	}
	basemap = append(basemap, min)
	for i := 0; i < 5; i++ {
		basemap = append(basemap, basemap[i]+diff)
	}
	for j := 0; j < 5; j++ {
		if x >= basemap[j] && x < basemap[j+1] {
			return colors[j]
		}
	}
	return colors[len(colors)-1]
}

func (a *App) getSortedStatistics() PairList {
	p := make(PairList, 100)

	i := 0
	for k, v := range a.statistics {
		p[i] = Pair{Key: k, Value: v}
		i++
	}
	sort.Sort(sort.Reverse(p))
	return p
}
