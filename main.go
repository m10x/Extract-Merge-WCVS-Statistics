package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strconv"
	"strings"
)

var (
	reset  = "\033[0m"
	red    = "\033[31m"
	yellow = "\033[33m"
	green  = "\033[32m"
)

var (
	sortedList []string

	general           General
	indicator         General
	cachebuster       General
	once_indicator    General
	once_cachebuster  General
	total_indicator   General
	total_cachebuster General
	timeFalseNeg      Time
	timeFalsePos      Time

	temp_general           General
	temp_indicator         General
	temp_cachebuster       General
	temp_once_indicator    General
	temp_once_cachebuster  General
	temp_total_indicator   General
	temp_total_cachebuster General
	temp_timeFalseNeg      Time
	temp_timeFalsePos      Time
)

type (
	General struct {
		stats map[string]int
	}
	Time struct {
		count      int
		once_count int
		values     []int
	}
)

func main() {
	if runtime.GOOS == "windows" {
		reset = ""
		red = ""
		yellow = ""
		green = ""
	}

	pathList := parseFlags()

	if pathList == "" {
		fmt.Printf("%sError: -path wasn't specified%s\n", red, reset)
		os.Exit(1)
	}

	//initialize structs
	general = General{stats: make(map[string]int)}
	indicator = General{stats: make(map[string]int)}
	cachebuster = General{stats: make(map[string]int)}
	once_indicator = General{stats: make(map[string]int)}
	once_cachebuster = General{stats: make(map[string]int)}
	total_indicator = General{stats: make(map[string]int)}
	total_cachebuster = General{stats: make(map[string]int)}
	timeFalseNeg = Time{}
	timeFalsePos = Time{}

	//sliceList := readLocalFile(pathList)

	filepath.Walk(pathList, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("error walk: " + err.Error())
		}
		if !info.IsDir() && !strings.HasSuffix(path, ".csv") {
			var compName string

			content := readLocalFile(path)

			for i, x := range content {
				if i == 0 {
					compName = strings.Split(x, "D:\\aWCVS\\tests\\")[1]
					compName = strings.Split(compName, "\\")[0]
				}
				if strings.HasPrefix(x, "General: ") {
					//initialize/reset structs
					temp_general = General{stats: make(map[string]int)}
					temp_indicator = General{stats: make(map[string]int)}
					temp_cachebuster = General{stats: make(map[string]int)}
					temp_once_indicator = General{stats: make(map[string]int)}
					temp_once_cachebuster = General{stats: make(map[string]int)}
					temp_total_indicator = General{stats: make(map[string]int)}
					temp_total_cachebuster = General{stats: make(map[string]int)}
					temp_timeFalseNeg = Time{}
					temp_timeFalsePos = Time{}
					extractStats(x, "General")
				} else if strings.HasPrefix(x, "Indicator: ") {
					extractStats(x, "Indicator")
				} else if strings.HasPrefix(x, "Cachebuster: ") {
					extractStats(x, "Cachebuster")
				} else if strings.HasPrefix(x, "Once_Indicator: ") {
					extractStats(x, "Once_Indicator")
				} else if strings.HasPrefix(x, "Once_Cachebuster: ") {
					extractStats(x, "Once_Cachebuster")
				} else if strings.HasPrefix(x, "Total_Indicator: ") {
					extractStats(x, "Total_Indicator")
				} else if strings.HasPrefix(x, "Total_Cachebuster: ") {
					extractStats(x, "Total_Cachebuster")
				} else if strings.HasPrefix(x, "TimeFalseNeg: ") {
					extractStats(x, "TimeFalseNeg")
				} else if strings.HasPrefix(x, "TimeFalsePos: ") {
					extractStats(x, "TimeFalsePos")
				}
			}

			for key, val := range temp_general.stats {
				general.stats[key] += val
			}
			for key, val := range temp_indicator.stats {
				indicator.stats[key] += val
			}
			for key, val := range temp_cachebuster.stats {
				cachebuster.stats[key] += val
			}
			for key, val := range temp_once_indicator.stats {
				once_indicator.stats[key] += val
			}
			for key, val := range temp_once_cachebuster.stats {
				once_cachebuster.stats[key] += val
			}
			for key, val := range temp_total_indicator.stats {
				total_indicator.stats[key] += val
			}
			for key, val := range temp_total_cachebuster.stats {
				total_cachebuster.stats[key] += val
			}
			timeFalseNeg.count += temp_timeFalseNeg.count
			timeFalseNeg.once_count += temp_timeFalseNeg.once_count
			timeFalseNeg.values = append(timeFalseNeg.values, temp_timeFalseNeg.values...)

			timeFalsePos.count += temp_timeFalsePos.count
			timeFalsePos.once_count += temp_timeFalsePos.once_count
			timeFalsePos.values = append(timeFalsePos.values, temp_timeFalsePos.values...)
		}
		return nil
	})

	writeStructsToFile(pathList)
}

func createFile(pathList string, name string) *os.File {
	fileName := pathList + name + ".csv"
	_, err := os.Stat(fileName)

	var file *os.File
	defer file.Close()

	if !os.IsNotExist(err) {
		fmt.Printf("The file %s will be overwritten, as it already exists\n", fileName)
		file, err = os.OpenFile(fileName, os.O_WRONLY, 0666)
	} else {
		file, err = os.Create(fileName)
	}
	if err != nil {
		fmt.Print("CompletedURLs: " + err.Error() + "\n")
	}
	return file
}

func checkFix(key string) string {
	fix := ""
	for strings.Contains(key, fix) {
		fix += "\""
	}
	return fix
}

func writeStructsToFile(pathList string) {
	file := createFile(pathList, "general")
	first := true
	for key := range general.stats {
		if first {
			first = false
		} else {
			file.WriteString(";")
		}
		fixes := checkFix(key)
		file.WriteString(fixes + key + fixes)
	}
	file.WriteString("\n")
	first = true
	for _, val := range general.stats {
		if first {
			first = false
		} else {
			file.WriteString(";")
		}
		fixes := checkFix(fmt.Sprint(val))
		file.WriteString(fixes + fmt.Sprint(val) + fixes)
	}
	file.Close()

	/* Not needed
	file = createFile(pathList, "indicator")
	first = true
	for key := range indicator.stats {
		if first {
			first = false
		} else {
			file.WriteString(";")
		}
		fixes := checkFix(key)
		file.WriteString(fixes + key + fixes)
	}
	file.WriteString("\n")
	first = true
	for _, val := range indicator.stats {
		if first {
			first = false
		} else {
			file.WriteString(";")
		}
		fixes := checkFix(fmt.Sprint(val))
		file.WriteString(fixes + fmt.Sprint(val) + fixes)
	}
	file.Close()
	*/

	/* Not needed
	file = createFile(pathList, "cachebuster")
	first = true
	for key := range cachebuster.stats {
		if first {
			first = false
		} else {
			file.WriteString(";")
		}
		fixes := checkFix(key)
		file.WriteString(fixes + key + fixes)
	}
	file.WriteString("\n")
	first = true
	for _, val := range cachebuster.stats {
		if first {
			first = false
		} else {
			file.WriteString(";")
		}
		fixes := checkFix(fmt.Sprint(val))
		file.WriteString(fixes + fmt.Sprint(val) + fixes)
	}
	file.Close()
	*/

	file = createFile(pathList, "once_indicator")
	first = true
	for key := range once_indicator.stats {
		if first {
			first = false
		} else {
			file.WriteString(";")
		}
		fixes := checkFix(key)
		file.WriteString(fixes + key + fixes)
	}
	file.WriteString("\n")
	first = true
	for _, val := range once_indicator.stats {
		if first {
			first = false
		} else {
			file.WriteString(";")
		}
		fixes := checkFix(fmt.Sprint(val))
		file.WriteString(fixes + fmt.Sprint(val) + fixes)
	}
	file.Close()

	file = createFile(pathList, "once_cachebuster")
	first = true
	for key := range once_cachebuster.stats {
		if first {
			first = false
		} else {
			file.WriteString(";")
		}
		fixes := checkFix(key)
		file.WriteString(fixes + key + fixes)
	}
	file.WriteString("\n")
	first = true
	for _, val := range once_cachebuster.stats {
		if first {
			first = false
		} else {
			file.WriteString(";")
		}
		fixes := checkFix(fmt.Sprint(val))
		file.WriteString(fixes + fmt.Sprint(val) + fixes)
	}
	file.Close()

	file = createFile(pathList, "total_indicator")
	first = true
	for key := range total_indicator.stats {
		if first {
			first = false
		} else {
			file.WriteString(";")
		}
		fixes := checkFix(key)
		file.WriteString(fixes + key + fixes)
	}
	file.WriteString("\n")
	first = true
	for _, val := range total_indicator.stats {
		if first {
			first = false
		} else {
			file.WriteString(";")
		}
		fixes := checkFix(fmt.Sprint(val))
		file.WriteString(fixes + fmt.Sprint(val) + fixes)
	}
	file.Close()

	file = createFile(pathList, "total_cachebuster")
	first = true
	for key := range total_cachebuster.stats {
		if first {
			first = false
		} else {
			file.WriteString(";")
		}
		fixes := checkFix(key)
		file.WriteString(fixes + key + fixes)
	}
	file.WriteString("\n")
	first = true
	for _, val := range total_cachebuster.stats {
		if first {
			first = false
		} else {
			file.WriteString(";")
		}
		fixes := checkFix(fmt.Sprint(val))
		file.WriteString(fixes + fmt.Sprint(val) + fixes)
	}
	file.Close()

	file = createFile(pathList, "timeFalseNeg")
	file.WriteString("once_count,count,values\n")
	sort.Ints(timeFalseNeg.values)
	file.WriteString(fmt.Sprintf("%d,%d,%v", timeFalseNeg.once_count, timeFalseNeg.count, timeFalseNeg.values))
	file.Close()

	file = createFile(pathList, "timeFalsePos")
	file.WriteString("once_count,count,values\n")
	sort.Ints(timeFalsePos.values)
	file.WriteString(fmt.Sprintf("%d,%d,%v", timeFalsePos.once_count, timeFalsePos.count, timeFalsePos.values))
	file.Close()
}

func extractStats(content string, prefix string) {
	if !strings.HasPrefix(content, "Time") {
		content = strings.Split(content, "[")[1]
		content = strings.TrimSuffix(content, "]")

		for _, x := range strings.Split(content, " ") {
			values := strings.Split(x, ":")
			writeToStruct(values[0], values[1], prefix)
		}
	} else {
		contentSlice := strings.Split(content, "[")
		values := strings.TrimSuffix(contentSlice[2], "]")
		for _, x := range strings.Split(values, " ") {
			valueint, _ := strconv.Atoi(x)
			if prefix == "TimeFalseNeg" {
				temp_timeFalseNeg.values = append(temp_timeFalseNeg.values, valueint)
			} else if prefix == "TimeFalsePos" {
				temp_timeFalsePos.values = append(temp_timeFalseNeg.values, valueint)
			} else {
				fmt.Println("This shouldnt happen... times: " + x)
				os.Exit(4)
			}
		}
		counts := strings.Split(contentSlice[1], "]")
		countsSlice := strings.Split(counts[0], " ")
		once_count, _ := strconv.Atoi(countsSlice[0])
		count, _ := strconv.Atoi(countsSlice[1])
		if prefix == "TimeFalseNeg" {
			temp_timeFalseNeg.once_count += once_count
			temp_timeFalseNeg.count += count
		} else if prefix == "TimeFalsePos" {
			temp_timeFalsePos.once_count += once_count
			temp_timeFalsePos.count += count
		} else {
			fmt.Println("This shouldnt happen... times: count")
			os.Exit(5)
		}
	}
}

func writeToStruct(name string, value string, prefix string) {
	if prefix == "General" {
		valueint, _ := strconv.Atoi(value)
		temp_general.stats[name] = valueint
	} else if prefix == "Indicator" {
		valueint, _ := strconv.Atoi(value)
		temp_indicator.stats[name] = valueint
	} else if prefix == "Cachebuster" {
		valueint, _ := strconv.Atoi(value)
		temp_cachebuster.stats[name] = valueint
	} else if prefix == "Once_Indicator" {
		valueint, _ := strconv.Atoi(value)
		temp_once_indicator.stats[name] = valueint
	} else if prefix == "Once_Cachebuster" {
		valueint, _ := strconv.Atoi(value)
		temp_once_cachebuster.stats[name] = valueint
	} else if prefix == "Total_Indicator" {
		valueint, _ := strconv.Atoi(value)
		temp_total_indicator.stats[name] = valueint
	} else if prefix == "Total_Cachebuster" {
		valueint, _ := strconv.Atoi(value)
		temp_total_cachebuster.stats[name] = valueint
	} else {
		fmt.Println("This shouldnt happen... " + name + ":" + value)
		os.Exit(2)
	}
}

func readLocalFile(path string) []string {

	w, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Printf("%sError while Reading list %s: %s%s\n", red, path, err.Error(), reset)
		os.Exit(3)
	}

	return strings.Split(string(w), "\n")
}

func parseFlags() string {
	var pathList string

	flag.StringVar(&pathList, "path", "", "path to folder containing the log files")

	flag.Parse()

	return pathList
}
