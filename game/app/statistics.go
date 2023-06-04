package app

import (
	"encoding/csv"
	"fmt"
	"github.com/fatih/color"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/joohnes/wp-sea/game/logger"
	"os"
	"strconv"
)

func (a *App) LoadStatistics() error {
	f, err := os.Open("statistics.csv")
	defer f.Close()
	if err != nil {
		_, err = os.Create("statistics.csv")
		if err != nil {
			return err
		}
		f, err = os.Open("statistics.csv")
		defer f.Close()
	}
	r := csv.NewReader(f)
	records, err := r.ReadAll()
	if err != nil {
		logger.GetLoggerInstance().Println("couldn't read statistics")
		return err
	}
	basemap := make(map[string]int)
	for _, x := range records {
		number, err := strconv.Atoi(x[1])
		if err != nil {
			logger.GetLoggerInstance().Printf("couldn't load %v\n", x)
			return err
		}
		basemap[x[0]] = number
	}
	a.statistics = basemap
	return nil
}

func (a *App) SaveStatistics() {
	f, err := os.OpenFile("statistics.csv", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 777)
	defer f.Close()
	if err != nil {
		logger.GetLoggerInstance().Println("couldn't open'")
		return
	}
	w := csv.NewWriter(f)
	defer w.Flush()
	for coord, occurrences := range a.statistics {
		err := w.Write([]string{coord, strconv.Itoa(occurrences)})
		if err != nil {
			logger.GetLoggerInstance().Printf("couldn't save %s, %d\n", coord, occurrences)
		}
	}
}

func (a *App) ShowStatistics() {
	var min, max int
	for i := range a.statistics {
		if a.statistics[i] < min {
			min = a.statistics[i]
		} else if a.statistics[i] > max {
			max = a.statistics[i]
		}

	}
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
		color.FgYellow,
		color.FgHiYellow,
		color.FgRed,
		color.FgHiRed,
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
