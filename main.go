package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func main() {

	fmt.Print("請輸入新舊，0舊 1新:")
	var isNew string
	fmt.Scanln(&isNew)
	runType, err := strconv.Atoi(isNew)
	//尺寸表路徑
	data, err := ioutil.ReadFile("config.json")
	if err != nil {
		fmt.Print(err)
		os.Exit(3)
	}

	type config struct {
		WorkPath      string `json:"WorkPath"`
		NewWorkPath   string `json:"NewWorkPath"`
		SizeTablePath string `json:"SizeTablePath"`
		Leve3Dir      string `json:"Leve3Dir"`
		Logo          string `json:"Logo"`
		Story         string `json:"Story"`
		BeginCount    string `json:"BeginCount"`
		MaxCount      string `json:"MaxCount"`
		Copy1         string `json:"copy1"`
		Copy2         string `json:"Copy2"`
		Copy3         string `json:"Copy3"`
		Copy1Max      int    `json:"Copy1Max"`
		GroupDir      string `json:"GroupDir"`
		XXLDir        string `json:"XXLDir"`
		XXLFile       string `json:"XXLFile"`
	}
	var obj config
	err = json.Unmarshal(data, &obj)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(3)
	}
	DirPath := ""
	mkOutOrg := ""
	XXLFilePath := ""
	if runType == 0 {
		DirPath = strings.Replace(obj.WorkPath, "\\", "\\\\", -1)
		//複製到的資料夾
		mkOutOrg = obj.WorkPath + string(os.PathSeparator) + obj.Leve3Dir
		mkDir(mkOutOrg)
		XXLFilePath = XXLFilePath + obj.WorkPath + string(os.PathSeparator) + obj.XXLFile
	} else {
		DirPath = strings.Replace(obj.NewWorkPath+string(os.PathSeparator)+obj.Leve3Dir, "\\", "\\\\", -1)
		XXLFilePath = strings.Replace(obj.NewWorkPath+string(os.PathSeparator)+obj.Leve3Dir+string(os.PathSeparator)+obj.XXLFile, "\\", "\\\\", -1)
	}
	//XXL txt內文字
	XXLArray := readXXlFile(XXLFilePath)
	SpecPath := strings.Replace(obj.SizeTablePath, "\\", "\\\\", -1)
	//掃描DIR
	dirArr := scandir(DirPath)
	//fmt.Println(dirArr)
	//LOGO路徑
	//LogoPath := obj.SizeTablePath + string(os.PathSeparator) + obj.Logo
	BeginCount, _ := strconv.Atoi(obj.BeginCount)
	MaxCount, _ := strconv.Atoi(obj.MaxCount)
	imageIndex := 1
	group := 1
	var failSizeTable []string
	var needFillArr []string
	var keepGoodsSn [][]string
	var test1Darray []string
	var xxlTest1DArray []string
	var needToSize1000T1040 []string
	nowString := ""

	for _, fileDir := range dirArr {
		if runType == 0 {

			level1Dir := obj.WorkPath + string(os.PathSeparator) + fileDir
			if strings.Index(fileDir, ".txt") != -1 {
				continue
			}
			mkOut := mkOutOrg + string(os.PathSeparator) + obj.GroupDir + strconv.Itoa(group)
			mkDir(mkOut)
			inSideImageArr := scandir(level1Dir)
			copyTo1 := ""
			copyTo2 := ""
			copyTo3 := ""
			copyToStory := ""
			imageIndex2 := 1
			level2ImgDir := ""
			indexNext := 0
			hasA := false
			for _, fileImage := range inSideImageArr {
				if fileImage == "A.jpg" || fileImage == "B.jpg" || fileImage == "C.jpg" || fileImage == "a.jpg" || fileImage == "b.jpg" || fileImage == "c.jpg" {
					hasA = true
				}
			}
			for _, fileImage := range inSideImageArr {
				if strings.Index(fileImage, ".jpg") > -1 {
					level2ImgDir = level1Dir + string(os.PathSeparator) + fileImage

					if hasA {
						indexNext = imageIndex2 + 1
					} else {
						indexNext = imageIndex2
					}

					if inSideImageArr[0] == "A.jpg" || inSideImageArr[0] == "B.jpg" || inSideImageArr[0] == "C.jpg" || inSideImageArr[0] == "a.jpg" || inSideImageArr[0] == "b.jpg" || inSideImageArr[0] == "c.jpg" {
						indexNext = imageIndex2
					}

					if fileImage == "A.jpg" || fileImage == "B.jpg" || fileImage == "C.jpg" || fileImage == "a.jpg" || fileImage == "b.jpg" || fileImage == "c.jpg" {
						copyTo1 = mkOut + string(os.PathSeparator) + strconv.Itoa(BeginCount) + "_" + obj.Copy1 + "1.jpg"
						copyTo2 = mkOut + string(os.PathSeparator) + strconv.Itoa(BeginCount) + "_" + obj.Copy2 + "_" + strconv.Itoa(imageIndex) + "_1.jpg"
						copyTo3 = mkOut + string(os.PathSeparator) + strconv.Itoa(BeginCount) + "_" + obj.Copy3 + ".jpg"
						needFillArr = append(needFillArr, copyTo1)
						//needFillArr = append(needFillArr, copyTo2)
					} else {

						copyTo1 = mkOut + string(os.PathSeparator) + strconv.Itoa(BeginCount) + "_" + obj.Copy1 + strconv.Itoa(indexNext) + ".jpg"
						copyTo2 = mkOut + string(os.PathSeparator) + strconv.Itoa(BeginCount) + "_" + obj.Copy2 + "_" + strconv.Itoa(imageIndex) + "_" + strconv.Itoa(indexNext) + ".jpg"
						copyTo3 = ""
					}
					needToSize1000T1040 = append(needToSize1000T1040, copyTo2)
					if (indexNext) <= obj.Copy1Max || fileImage == "A.jpg" || fileImage == "B.jpg" || fileImage == "C.jpg" || fileImage == "a.jpg" || fileImage == "b.jpg" || fileImage == "c.jpg" {
						CopyFile(level2ImgDir, copyTo1)
						fmt.Printf("複製圖片 %s ， 到 %s 成功!!\n", level2ImgDir, copyTo1)
					}
					CopyFile(level2ImgDir, copyTo2)
					fmt.Printf("複製圖片 %s ， 到 %s 成功!!\n", level2ImgDir, copyTo2)
					if copyTo3 != "" {
						CopyFile(level2ImgDir, copyTo3)
						fmt.Printf("複製圖片 %s ， 到 %s 成功!!\n", level2ImgDir, copyTo3)
					}
					imageIndex2++
				}
			}

			if strings.Index(fileDir, "OUT") == -1 {
				styleNoPath := SpecPath + string(os.PathSeparator) + fileDir[0:2] + fileDir[4:8] + ".jpg"
				storyPath := SpecPath + string(os.PathSeparator) + obj.Story
				copyTo1 = mkOut + string(os.PathSeparator) + strconv.Itoa(BeginCount) + "_" + obj.Copy1 + strconv.Itoa(imageIndex2) + ".jpg"
				copyTo2 = mkOut + string(os.PathSeparator) + strconv.Itoa(BeginCount) + "_" + obj.Copy2 + "_" + strconv.Itoa(imageIndex) + "_" + strconv.Itoa(imageIndex2) + ".jpg"
				copyToStory = mkOut + string(os.PathSeparator) + strconv.Itoa(BeginCount) + "_" + obj.Copy2 + "_" + strconv.Itoa(imageIndex) + "_" + strconv.Itoa(imageIndex2+1) + ".jpg"
				if (imageIndex2) <= obj.Copy1Max {
					err3 := CopyFile(styleNoPath, copyTo1)
					if err3 != nil {
						fmt.Printf("複製圖片 %s ， 到 %s 失敗!!\n", styleNoPath, copyTo1)
					} else {
						fmt.Printf("複製圖片 %s ， 到 %s 成功!!\n", styleNoPath, copyTo1)
					}
				}
				err4 := CopyFile(styleNoPath, copyTo2)
				err5 := CopyFile(storyPath, copyToStory)
				needToSize1000T1040 = append(needToSize1000T1040, copyToStory)
				needToSize1000T1040 = append(needToSize1000T1040, copyTo2)
				if err4 != nil || err5 != nil {
					fmt.Printf("複製圖片 %s ， 到 %s 失敗!!\n", styleNoPath, copyTo2)
					fmt.Printf("複製圖片 %s ， 到 %s 失敗!!\n", storyPath, copyToStory)
					if err5 != nil {
						failSizeTable = append(failSizeTable, "品牌故事:"+storyPath+"失敗\n")
					} else {
						failSizeTable = append(failSizeTable, "複製尺寸表:"+styleNoPath+"失敗\n")
					}

				} else {
					fmt.Printf("複製尺寸表 %s\n", styleNoPath+".jpg")
					fmt.Printf("複製圖片 %s ， 到 %s 成功!!\n", styleNoPath, copyTo2)
					fmt.Printf("品牌故事 %s\n", storyPath+".jpg")
					fmt.Printf("複製圖片 %s ， 到 %s 成功!!\n", storyPath, copyToStory)
				}
				BeginCount++
				//imageIndex++
				if BeginCount > MaxCount {
					newbegin, _ := strconv.Atoi(obj.BeginCount)
					BeginCount = newbegin
					//imageIndex = 1
					group++
				}
			}
			//os.Exit(3)
		} else {
			if strings.Index(fileDir, ".jpg") > -1 {
				array := strings.Split(fileDir, "_")
				//keepGoodsSn[array[0]] = append(keepGoodsSn[array[0]], fileDir)
				if nowString != array[0] {
					nowString = array[0]
					test1Darray = append(test1Darray, nowString)
					if InStringSlice(XXLArray, array[0]) {
						xxlTest1DArray = append(xxlTest1DArray, array[0])
					}
				}
			}
		}
	}

	if runType == 1 {
		for _, goodsSn := range test1Darray {
			var test2Array []string
			for _, fileDir := range dirArr {
				if strings.Index(fileDir, goodsSn) > -1 {
					test2Array = append(test2Array, fileDir)
				}
			}
			keepGoodsSn = append(keepGoodsSn, test2Array)
		}

		if len(keepGoodsSn) > 0 {
			for _, snArr := range keepGoodsSn {
				mkOut := obj.NewWorkPath + string(os.PathSeparator) + obj.Leve3Dir + string(os.PathSeparator) + obj.GroupDir + strconv.Itoa(group)
				mkDir(mkOut)
				imageIndex2 := 1
				for _, imagePath := range snArr {
					copyTo1 := ""
					copyTo2 := ""
					copyTo3 := ""

					level2ImgDir := obj.NewWorkPath + string(os.PathSeparator) + obj.Leve3Dir + string(os.PathSeparator) + imagePath
					if strings.Index(imagePath, "_01.") > -1 || strings.Index(imagePath, "A.") > -1 || strings.Index(imagePath, "B.") > -1 || strings.Index(imagePath, "C.") > -1 {
						copyTo1 = mkOut + string(os.PathSeparator) + strconv.Itoa(BeginCount) + "_" + obj.Copy1 + "1.jpg"
						copyTo2 = mkOut + string(os.PathSeparator) + strconv.Itoa(BeginCount) + "_" + obj.Copy2 + "_" + strconv.Itoa(imageIndex) + "_1.jpg"
						copyTo3 = mkOut + string(os.PathSeparator) + strconv.Itoa(BeginCount) + "_" + obj.Copy3 + ".jpg"
						needFillArr = append(needFillArr, copyTo1)
						imageIndex2--
						//needFillArr = append(needFillArr, copyTo2)
					} else {
						copyTo1 = mkOut + string(os.PathSeparator) + strconv.Itoa(BeginCount) + "_" + obj.Copy1 + strconv.Itoa(imageIndex2+1) + ".jpg"
						copyTo2 = mkOut + string(os.PathSeparator) + strconv.Itoa(BeginCount) + "_" + obj.Copy2 + "_" + strconv.Itoa(imageIndex) + "_" + strconv.Itoa(imageIndex2+1) + ".jpg"
						copyTo3 = ""
					}
					needToSize1000T1040 = append(needToSize1000T1040, copyTo2)
					if (imageIndex2+1) <= obj.Copy1Max || (strings.Index(imagePath, "01.") > -1 || strings.Index(imagePath, "A.") > -1 || strings.Index(imagePath, "B.") > -1 || strings.Index(imagePath, "C.") > -1) {
						CopyFile(level2ImgDir, copyTo1)
						fmt.Printf("複製圖片 %s ， 到 %s 成功!!\n", level2ImgDir, copyTo1)
					}
					CopyFile(level2ImgDir, copyTo2)
					fmt.Printf("複製圖片 %s ， 到 %s 成功!!\n", level2ImgDir, copyTo2)

					if copyTo3 != "" {
						CopyFile(level2ImgDir, copyTo3)
						fmt.Printf("複製圖片 %s ， 到 %s 成功!!\n", level2ImgDir, copyTo3)
					}
					imageIndex2++
				}

				storyPath := SpecPath + string(os.PathSeparator) + obj.Story
				copyToStory := mkOut + string(os.PathSeparator) + strconv.Itoa(BeginCount) + "_" + obj.Copy2 + "_" + strconv.Itoa(imageIndex) + "_" + strconv.Itoa(imageIndex2+1) + ".jpg"
				needToSize1000T1040 = append(needToSize1000T1040, copyToStory)
				err3 := CopyFile(storyPath, copyToStory)
				if err3 != nil {
					fmt.Printf("複製圖片 %s ， 到 %s 失敗!!\n", storyPath, copyToStory)
					failSizeTable = append(failSizeTable, "品牌故事:"+storyPath+"失敗\n")
				} else {
					fmt.Printf("品牌故事 %s\n", storyPath+".jpg")
					fmt.Printf("複製圖片 %s ， 到 %s 成功!!\n", storyPath, copyToStory)
				}
				BeginCount++
				//imageIndex++
				if BeginCount > MaxCount {
					newbegin, _ := strconv.Atoi(obj.BeginCount)
					BeginCount = newbegin
					imageIndex = 1
					group++
				}
			}

		}
	}

	if len(XXLArray) > 0 {

		BeginCount, _ = strconv.Atoi(obj.BeginCount)
		MaxCount, _ = strconv.Atoi(obj.MaxCount)
		imageIndex = 1
		xxl := 1

		if runType == 0 {
			for _, pathName := range XXLArray {
				level1Dir := obj.WorkPath + string(os.PathSeparator) + pathName
				if _, err := os.Stat(level1Dir); os.IsNotExist(err) {
					continue
				}
				mkOut := mkOutOrg + string(os.PathSeparator) + obj.XXLDir + strconv.Itoa(xxl)
				mkDir(mkOut)
				inSideImageArr := scandir(level1Dir)
				copyTo1 := ""
				copyTo2 := ""
				copyTo3 := ""
				imageIndex2 := 1
				level2ImgDir := ""
				copyToStory := ""
				indexNext := 0
				hasA := false
				for _, fileImage := range inSideImageArr {
					if fileImage == "A.jpg" || fileImage == "B.jpg" || fileImage == "C.jpg" || fileImage == "a.jpg" || fileImage == "b.jpg" || fileImage == "c.jpg" {
						hasA = true
					}
				}

				for _, fileImage := range inSideImageArr {
					if strings.Index(fileImage, ".jpg") > -1 {
						level2ImgDir = level1Dir + string(os.PathSeparator) + fileImage

						if hasA {
							indexNext = imageIndex2 + 1
						} else {
							indexNext = imageIndex2
						}
						if inSideImageArr[0] == "A.jpg" || inSideImageArr[0] == "B.jpg" || inSideImageArr[0] == "C.jpg" || inSideImageArr[0] == "a.jpg" || inSideImageArr[0] == "b.jpg" || inSideImageArr[0] == "c.jpg" {
							indexNext = imageIndex2
						}

						if fileImage == "A.jpg" || fileImage == "B.jpg" || fileImage == "C.jpg" || fileImage == "a.jpg" || fileImage == "b.jpg" || fileImage == "c.jpg" {
							copyTo1 = mkOut + string(os.PathSeparator) + strconv.Itoa(BeginCount) + "_" + obj.Copy1 + "1.jpg"
							copyTo2 = mkOut + string(os.PathSeparator) + strconv.Itoa(BeginCount) + "_" + obj.Copy2 + "_" + strconv.Itoa(imageIndex) + "_1.jpg"
							copyTo3 = mkOut + string(os.PathSeparator) + strconv.Itoa(BeginCount) + "_" + obj.Copy3 + ".jpg"
							needFillArr = append(needFillArr, copyTo1)
							//needFillArr = append(needFillArr, copyTo2)
						} else {
							copyTo1 = mkOut + string(os.PathSeparator) + strconv.Itoa(BeginCount) + "_" + obj.Copy1 + strconv.Itoa(indexNext) + ".jpg"
							copyTo2 = mkOut + string(os.PathSeparator) + strconv.Itoa(BeginCount) + "_" + obj.Copy2 + "_" + strconv.Itoa(imageIndex) + "_" + strconv.Itoa(indexNext) + ".jpg"
							copyTo3 = ""
						}
						needToSize1000T1040 = append(needToSize1000T1040, copyTo2)
						if (indexNext) <= obj.Copy1Max || fileImage == "A.jpg" || fileImage == "B.jpg" || fileImage == "C.jpg" || fileImage == "a.jpg" || fileImage == "b.jpg" || fileImage == "c.jpg" {
							CopyFile(level2ImgDir, copyTo1)
							fmt.Printf("複製圖片 %s ， 到 %s 成功!!\n", level2ImgDir, copyTo1)
						}
						CopyFile(level2ImgDir, copyTo2)
						fmt.Printf("複製圖片 %s ， 到 %s 成功!!\n", level2ImgDir, copyTo2)
						if copyTo3 != "" {
							CopyFile(level2ImgDir, copyTo3)
							fmt.Printf("複製圖片 %s ， 到 %s 成功!!\n", level2ImgDir, copyTo3)
						}
						imageIndex2++
					}
				}

				if strings.Index(pathName, "OUT") == -1 {
					styleNoPath := SpecPath + string(os.PathSeparator) + pathName[0:2] + pathName[4:8] + ".jpg"
					storyPath := SpecPath + string(os.PathSeparator) + obj.Story

					copyTo2 = mkOut + string(os.PathSeparator) + strconv.Itoa(BeginCount) + "_" + obj.Copy2 + "_" + strconv.Itoa(imageIndex) + "_" + strconv.Itoa(imageIndex2) + ".jpg"
					copyToStory = mkOut + string(os.PathSeparator) + strconv.Itoa(BeginCount) + "_" + obj.Copy2 + "_" + strconv.Itoa(imageIndex) + "_" + strconv.Itoa(imageIndex2+1) + ".jpg"

					if (imageIndex2) <= obj.Copy1Max {
						copyTo1 = mkOut + string(os.PathSeparator) + strconv.Itoa(BeginCount) + "_" + obj.Copy1 + strconv.Itoa(imageIndex2) + ".jpg"
						err3 := CopyFile(styleNoPath, copyTo1)
						if err3 != nil {
							fmt.Printf("複製圖片 %s ， 到 %s 失敗!!\n", styleNoPath, copyTo1)
						} else {
							fmt.Printf("複製圖片 %s ， 到 %s 成功!!\n", styleNoPath, copyTo1)
						}
					}

					err4 := CopyFile(styleNoPath, copyTo2)
					err5 := CopyFile(storyPath, copyToStory)
					needToSize1000T1040 = append(needToSize1000T1040, copyTo2)
					needToSize1000T1040 = append(needToSize1000T1040, copyToStory)
					if err4 != nil || err5 != nil {
						fmt.Printf("複製圖片 %s ， 到 %s 失敗!!\n", styleNoPath, copyTo1)
						fmt.Printf("複製圖片 %s ， 到 %s 失敗!!\n", styleNoPath, copyTo2)
						fmt.Printf("複製圖片 %s ， 到 %s 失敗!!\n", storyPath, copyToStory)
						if err5 != nil {
							failSizeTable = append(failSizeTable, "品牌故事:"+storyPath+"失敗\n")
						} else {
							failSizeTable = append(failSizeTable, "複製尺寸表:"+styleNoPath+"失敗\n")
						}

					} else {
						fmt.Printf("複製尺寸表 %s\n", styleNoPath+".jpg")
						fmt.Printf("複製圖片 %s ， 到 %s 成功!!\n", styleNoPath, copyTo1)
						fmt.Printf("複製圖片 %s ， 到 %s 成功!!\n", styleNoPath, copyTo2)
						fmt.Printf("品牌故事 %s\n", storyPath+".jpg")
						fmt.Printf("複製圖片 %s ， 到 %s 成功!!\n", storyPath, copyToStory)
					}
					BeginCount++
					//imageIndex++
					if BeginCount > MaxCount {
						newbegin, _ := strconv.Atoi(obj.BeginCount)
						BeginCount = newbegin
						imageIndex = 1
						xxl++
					}
				}
			}
		} else {
			var keepXXlGoodsSn [][]string
			for _, goodsSn := range xxlTest1DArray {
				var test2Array []string
				for _, fileDir := range dirArr {
					if strings.Index(fileDir, goodsSn) > -1 {
						test2Array = append(test2Array, fileDir)
					}
				}
				keepXXlGoodsSn = append(keepXXlGoodsSn, test2Array)
			}
			if len(keepXXlGoodsSn) > 0 {
				for _, snArr := range keepXXlGoodsSn {
					mkOut := obj.NewWorkPath + string(os.PathSeparator) + obj.Leve3Dir + string(os.PathSeparator) + obj.XXLDir + strconv.Itoa(xxl)
					mkDir(mkOut)
					imageIndex2 := 1
					for _, imagePath := range snArr {
						copyTo1 := ""
						copyTo2 := ""
						copyTo3 := ""
						level2ImgDir := obj.NewWorkPath + string(os.PathSeparator) + obj.Leve3Dir + string(os.PathSeparator) + imagePath
						if strings.Index(imagePath, "_01.") > -1 || strings.Index(imagePath, "A.") > -1 || strings.Index(imagePath, "B.") > -1 || strings.Index(imagePath, "C.") > -1 {
							copyTo1 = mkOut + string(os.PathSeparator) + strconv.Itoa(BeginCount) + "_" + obj.Copy1 + "1.jpg"
							copyTo2 = mkOut + string(os.PathSeparator) + strconv.Itoa(BeginCount) + "_" + obj.Copy2 + "_" + strconv.Itoa(imageIndex) + "_1.jpg"
							copyTo3 = mkOut + string(os.PathSeparator) + strconv.Itoa(BeginCount) + "_" + obj.Copy3 + ".jpg"
							needFillArr = append(needFillArr, copyTo1)
							imageIndex2--
							//needFillArr = append(needFillArr, copyTo2)
						} else {
							copyTo1 = mkOut + string(os.PathSeparator) + strconv.Itoa(BeginCount) + "_" + obj.Copy1 + strconv.Itoa(imageIndex2+1) + ".jpg"
							copyTo2 = mkOut + string(os.PathSeparator) + strconv.Itoa(BeginCount) + "_" + obj.Copy2 + "_" + strconv.Itoa(imageIndex) + "_" + strconv.Itoa(imageIndex2+1) + ".jpg"
							copyTo3 = ""
						}
						needToSize1000T1040 = append(needToSize1000T1040, copyTo2)
						if (imageIndex2) <= obj.Copy1Max {
							CopyFile(level2ImgDir, copyTo1)
							fmt.Printf("複製圖片 %s ， 到 %s 成功!!\n", level2ImgDir, copyTo1)
						}
						CopyFile(level2ImgDir, copyTo2)
						fmt.Printf("複製圖片 %s ， 到 %s 成功!!\n", level2ImgDir, copyTo2)
						if copyTo3 != "" {
							CopyFile(level2ImgDir, copyTo3)
							fmt.Printf("複製圖片 %s ， 到 %s 成功!!\n", level2ImgDir, copyTo3)
						}
						imageIndex2++
					}

					storyPath := SpecPath + string(os.PathSeparator) + obj.Story
					copyToStory := mkOut + string(os.PathSeparator) + strconv.Itoa(BeginCount) + "_" + obj.Copy2 + "_" + strconv.Itoa(imageIndex) + "_" + strconv.Itoa(imageIndex2+1) + ".jpg"

					err3 := CopyFile(storyPath, copyToStory)
					needToSize1000T1040 = append(needToSize1000T1040, copyToStory)
					if err3 != nil {
						fmt.Printf("複製圖片 %s ， 到 %s 失敗!!\n", storyPath, copyToStory)
						failSizeTable = append(failSizeTable, "品牌故事:"+storyPath+"失敗\n")
					} else {
						fmt.Printf("品牌故事 %s\n", storyPath+".jpg")
						fmt.Printf("複製圖片 %s ， 到 %s 成功!!\n", storyPath, copyToStory)
					}
					BeginCount++
					//imageIndex++
					if BeginCount > MaxCount {
						newbegin, _ := strconv.Atoi(obj.BeginCount)
						BeginCount = newbegin
						imageIndex = 1
						xxl++
					}
				}

			}
		}
	}

	var txtString string
	if len(failSizeTable) > 0 {
		for _, errorMsg := range failSizeTable {
			txtString = txtString + errorMsg
		}
		content := []byte(txtString)
		err := ioutil.WriteFile("log.txt", content, 0666)
		if err != nil {
			fmt.Println("ioutil WriteFile error: ", err)
		}
	}

	txtString = ""

	if len(needFillArr) > 0 {
		for _, fillPath := range needFillArr {
			txtString = txtString + fillPath + ";\n"
		}
		fmt.Println(txtString)
		content := []byte(txtString)
		err := ioutil.WriteFile(".\\needFill.txt", content, 0777)
		if err != nil {
			fmt.Println("ioutil WriteFile error: ", err)
		}
	}
	txtString = ""

	if len(needToSize1000T1040) > 0 {
		for _, fillPath := range needToSize1000T1040 {
			txtString = txtString + fillPath + ";\n"
		}
		fmt.Println(txtString)
		content := []byte(txtString)
		err := ioutil.WriteFile(".\\needToSize1000T1040.txt", content, 0777)
		if err != nil {
			fmt.Println("ioutil WriteFile error: ", err)
		}
	}

	fmt.Println("執行完成")
	fmt.Scanln()
}

//掃描資料夾底下檔案
func scandir(dir string) []string {
	var files []string
	filelist, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}
	for _, f := range filelist {
		files = append(files, f.Name())
	}
	return files
}

func readXXlFile(XXLFilePath string) []string {
	var content []string
	f, err := os.Open(XXLFilePath)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		content = append(content, scanner.Text())
	}
	return content
}

//byte 轉 string
func BytesToString(data []byte) string {
	return string(data[:])
}

func moveFile(orgPath string, movePath string) {

	fmt.Println(movePath)
	path := filepath.Dir(movePath)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.Mkdir(path, 0777)
	}
	err := os.Rename(orgPath, movePath)
	if err != nil {
		fmt.Println("移動檔案失敗!!")
	}
}

// CopyFile copies a file from src to dst. If src and dst files exist, and are
// the same, then return success. Otherise, attempt to create a hard link
// between the two files. If that fail, copy the file contents from src to dst.
func CopyFile(src, dst string) (err error) {
	sfi, err := os.Stat(src)
	if err != nil {
		return
	}
	if !sfi.Mode().IsRegular() {
		// cannot copy non-regular files (e.g., directories,
		// symlinks, devices, etc.)
		return fmt.Errorf("CopyFile: non-regular source file %s (%q)", sfi.Name(), sfi.Mode().String())
	}
	dfi, err := os.Stat(dst)
	if err != nil {
		if !os.IsNotExist(err) {
			return
		}
	} else {
		if !(dfi.Mode().IsRegular()) {
			return fmt.Errorf("CopyFile: non-regular destination file %s (%q)", dfi.Name(), dfi.Mode().String())
		}
		if os.SameFile(sfi, dfi) {
			return
		}
	}
	if err = os.Link(src, dst); err == nil {
		return
	}
	err = copyFileContents(src, dst)
	return
}

// copyFileContents copies the contents of the file named src to the file named
// by dst. The file will be created if it does not already exist. If the
// destination file exists, all it's contents will be replaced by the contents
// of the source file.
func copyFileContents(src, dst string) (err error) {
	in, err := os.Open(src)
	if err != nil {
		return
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return
	}
	defer func() {
		cerr := out.Close()
		if err == nil {
			err = cerr
		}
	}()
	if _, err = io.Copy(out, in); err != nil {
		return
	}
	err = out.Sync()
	return
}

func mkDir(src string) (err error) {
	if _, err := os.Stat(src); os.IsNotExist(err) {
		os.Mkdir(src, 0755)
		fmt.Println("建立資料夾:" + src)
	}
	if err != nil {
		fmt.Print("建立錯誤:")
		fmt.Print(err)
	}
	return
}

func InStringSlice(haystack []string, needle string) bool {
	for _, e := range haystack {
		if e == needle {
			return true
		}
	}

	return false
}

func dd(data string) (err error) {
	fmt.Println(data)
	os.Exit(3)
	return
}
