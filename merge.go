package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
)

type vidioCollation struct {
	name     string
	year     string
	episodes map[string] /*season*/ string /*episode */
}

/*
自动获取当前目录下的
*/
func getTargetPath() []string {
	r := make([]string, 0)
	//获取文件或目录相关信息
	fileInfoList, err := ioutil.ReadDir("./")
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range fileInfoList {
		if strings.HasSuffix(file.Name(), ".json") {
			r = append(r, file.Name())
		}
	}
	return r
}

type reduceSuggestion struct {
	name             string
	storage          string
	redundantName    string
	redundantStorage string
}

/* create data from files*/
func getReduceData(r []string) map[string]map[string]vidioCollation {
	totalvcm := make(map[string]map[string]vidioCollation)

	for _, f := range r {
		fmt.Println(f)
		data, e := ioutil.ReadFile(f)
		if e != nil {
			fmt.Printf("ReadFile %s Error: %v", f, e)
			continue
		}
		vcm := make(map[string]vidioCollation)

		//fmt.Printf("ReadFile %s success: %v", f, string(data))

		s := json.Unmarshal(data, &vcm)
		if s != nil {
			//fmt.Printf("Unmarshal  Error: %v", s)

			continue
		}
		//fmt.Printf("%v\n", vcm)
		totalvcm[f] = vcm
	}
	return totalvcm

}

func searchVideoName(videoName string,
	totalvcm map[string]map[string]vidioCollation,
	skip_storage string) map[string]map[string]vidioCollation {

	nm := make(map[string]map[string]vidioCollation)

	for storage, data := range totalvcm {
		if storage == skip_storage {
			continue
		}

		for other_videoName, videoMetadata := range data {
			if strings.Contains(other_videoName, videoName) {
				othermap, ok := nm[storage]
				if ok {
					othermap[other_videoName] = videoMetadata
					nm[storage] = othermap
				} else {
					new_othermap := make(map[string]vidioCollation)
					new_othermap[other_videoName] = videoMetadata
					nm[storage] = new_othermap
				}
			}
		}
	}
	return nm
}
func main() {

	r := getTargetPath()
	totalvcm := getReduceData(r)

	/* create reduce */
	totalreduce := make(map[string][]reduceSuggestion)
	for storage, data := range totalvcm {
		for videoName := range data {
			nm := searchVideoName(videoName, totalvcm, storage)
			//fmt.Printf("================================\n")
			//fmt.Printf("searchVideoName for %s result is %v\n", videoName, nm)

			for redStrorage, redFileInfo := range nm {
				for onefile := range redFileInfo {
					totalreduce[videoName] = make([]reduceSuggestion, 0)
					totalreduce[videoName] = append(totalreduce[videoName], reduceSuggestion{
						name:             videoName,
						storage:          storage,
						redundantName:    onefile,
						redundantStorage: redStrorage,
					})

				}
			}
		}
	}
	fmt.Printf("================================\n")
	for onevideo, reduceResult := range totalreduce {
		fmt.Printf("## %s Should Reduce in:\n", onevideo)
		for _, d := range reduceResult {
			fmt.Printf("- %s\n", d.name)
			fmt.Printf("  - %s\n", d.storage)
			fmt.Printf("- %s\n", d.redundantName)
			fmt.Printf("  - %s\n", d.redundantStorage)
			fmt.Println()
		}
	}

}
