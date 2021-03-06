package main

import (
	"encoding/json"
	"fmt"
	"image/color"
	"image/jpeg"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/nfnt/resize"
	gim "github.com/ozankasikci/go-image-merge"
	"github.com/xuri/excelize"
)

type config struct {
	WorkPath         string                 `json:"WorkPath"`
	SizeTablePath    string                 `json:"SizeTablePath"`
	TryTablePath     string                 `json:"TryTablePath"`
	GetDir           string                 `json:"GetDir"`
	Leve3Dir         string                 `json:"Leve3Dir"`
	BlankImg         string                 `json:"BlankImg"`
	TryMapping       map[string]interface{} `json:"TryMapping"`
	TryPicName       string                 `json:"TryPicName"`
	ListPicName      string                 `json:"ListPicName"`
	SizePicName      string                 `json:"SizePicName"`
	OnePackageNumber int                    `json:"OnePackageNumber"`
	ExcelTitle       map[string]interface{} `json:"ExcelTitle"`
}

func main() {

	//設定擋路徑
	data, err := ioutil.ReadFile("config.json")
	if err != nil {
		fmt.Print(err)
		os.Exit(3)
	}
	//讀取設定檔
	var config config
	err = json.Unmarshal(data, &config)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(3)
	}

	var ExcelTitleArray = config.ExcelTitle
	//一個dir 內要幾套圖
	var OnePackageNumber = config.OnePackageNumber
	//套圖計算 index
	var beginCount = 1
	//外層資料夾 分套圖別名
	var outerDirCount = 1
	//創建資料夾路徑
	var mkOutOrg = config.WorkPath + string(os.PathSeparator) + config.Leve3Dir + strconv.Itoa(outerDirCount)
	//試穿表json名稱讀取
	var tryPicName = config.TryPicName
	//列表圖json名稱讀取
	var ListPicName = config.ListPicName
	//尺寸圖json名稱讀取
	var SizePicName = config.SizePicName
	//取代反斜線
	var DirPath = strings.Replace(config.WorkPath, "\\", "\\\\", -1)
	//尺寸表路徑 (用款號組成)
	var SizeTablePath = strings.Replace(config.SizeTablePath, "\\", "\\\\", -1)
	//試穿表路徑 (用料號前兩個字)
	var TryTablePath = strings.Replace(config.TryTablePath, "\\", "\\\\", -1)
	//試穿表對應圖片
	var TryMap = config.TryMapping
	//小圖路徑
	var SmallPath = config.GetDir

	//創建資料夾
	mkDir(mkOutOrg)

	//掃描DIR
	dirArr := scandir(DirPath)
	//需要縮放圖片位子陣列
	var needResize []string
	var sizePicInfoArray []string
	var tryPicInfoArray []string
	var mainPhoto = make(map[string]string)
	var otherProduct = make(map[string]string)
	var advertising = make(map[string]string)

	// fmt.Println("掃描路徑: ", dirArr)

	var needUnsetDir = 0

	for index, fileDir := range dirArr {
		if fileDir == config.Leve3Dir+strconv.Itoa(outerDirCount) {
			needUnsetDir = index
		}
	}
	var Slice2 = dirArr[:needUnsetDir]
	dirArr = Slice2
	var dirArrLastOneIndex = len(dirArr) - 1
	// fmt.Println(dirArr)
	// os.Exit(3)

	for fileIndex, fileDir := range dirArr {

		//切割資料夾變成陣列
		dirCutArr := strings.Split(fileDir, "_")

		fmt.Println("陣列大小: ", len(dirCutArr))

		if len(dirCutArr) >= 2 {

			if beginCount > OnePackageNumber {
				beginCount = 1
				outerDirCount++
				mkOutOrg = config.WorkPath + string(os.PathSeparator) + config.Leve3Dir + strconv.Itoa(outerDirCount)
			}

			fmt.Println("切割完: ", dirCutArr)

			//取得料號
			goodsSn := dirCutArr[0]
			//料號前兩碼
			twoCode := fmt.Sprintf("%v", TryMap[goodsSn[0:2]])

			//小圖路徑
			var smallPath = DirPath + string(os.PathSeparator) + fileDir + string(os.PathSeparator) + SmallPath

			fmt.Println("小圖路徑: ", smallPath)
			//創建料號資料夾
			// tvMallChildDir := mkOutOrg + string(os.PathSeparator) + goodsSn
			tvMallChildDir := mkOutOrg + string(os.PathSeparator)
			fmt.Println("創建料號資料夾: ", tvMallChildDir)
			mkDir(tvMallChildDir)

			//尺寸表圖路徑
			sizePicPath := SizeTablePath + string(os.PathSeparator) + goodsSn[0:2] + goodsSn[4:8] + ".jpg"
			fmt.Println("抓取尺寸表圖片路徑: ", sizePicPath)

			SizeStr := ""
			TryStr := ""

			if _, err := os.Stat(sizePicPath); os.IsNotExist(err) {
				sizePicInfoArray = append(sizePicInfoArray, sizePicPath)
			} else {
				fmt.Println(smallPath + string(os.PathSeparator) + goodsSn[0:2] + goodsSn[4:8] + ".jpg")
				fmt.Println(sizePicPath)
				SizeStr = goodsSn + "_" + SizePicName + ".jpg"
				CopyFile(sizePicPath, tvMallChildDir+string(os.PathSeparator)+goodsSn+"_"+SizePicName+".jpg")
			}

			//試穿表路徑
			tryPicPath := TryTablePath + string(os.PathSeparator) + twoCode
			if _, err := os.Stat(tryPicPath); os.IsNotExist(err) {
				tryPicInfoArray = append(tryPicInfoArray, tryPicPath)
			} else {
				TryStr = goodsSn + "_" + tryPicName + ".jpg"
				CopyFile(tryPicPath, tvMallChildDir+string(os.PathSeparator)+goodsSn+"_"+tryPicName+".jpg")
			}

			productPhoto := ""

			//掃描小圖 資料夾
			smallPicDirArray := scandir(smallPath)
			//fmt.Println(smallPicDirArray)
			for index, picDir := range smallPicDirArray {
				orgSmallPicPath := smallPath + string(os.PathSeparator) + picDir
				picIndex := strconv.Itoa(index + 1)
				if (index + 1) < 10 {
					picIndex = "0" + picIndex
				}
				CopySmallPicPath := tvMallChildDir + string(os.PathSeparator) + goodsSn + "_" + picIndex + ".jpg"
				CopyFile(orgSmallPicPath, CopySmallPicPath)
				productPhoto += goodsSn + "_" + picIndex + ".jpg,"
				if index == 0 {
					CopySmallPicPath = tvMallChildDir + string(os.PathSeparator) + goodsSn + "_" + ListPicName + ".jpg"
					CopyFile(orgSmallPicPath, CopySmallPicPath)

					mainPhoto[goodsSn] = goodsSn + "_" + ListPicName + ".jpg"

					CopySmallPicPath = tvMallChildDir + string(os.PathSeparator) + goodsSn + "_" + "250X250.jpg"
					CopyFile(orgSmallPicPath, CopySmallPicPath)
					needResize = append(needResize, CopySmallPicPath)
					advertising[goodsSn] = goodsSn + "_" + "250X250.jpg"
				}
			}

			productPhoto += TryStr + "," + SizeStr
			otherProduct[goodsSn] = productPhoto
			beginCount++

			if beginCount > OnePackageNumber || fileIndex == dirArrLastOneIndex {

				f := excelize.NewFile()
				// 創建一個工作表
				index := f.NewSheet("Sheet1")
				// 設定活頁簿的默認工作表
				f.SetActiveSheet(index)
				f.SetColWidth("Sheet1", "A", "D", 30)
				for column, value := range ExcelTitleArray {
					f.SetCellValue("Sheet1", column, value)
					f.SetCellValue("Sheet1", column, value)
					f.SetCellValue("Sheet1", column, value)
					f.SetCellValue("Sheet1", column, value)
				}

				keys := make([]string, 0, len(mainPhoto))
				for k := range mainPhoto {
					keys = append(keys, k)
				}
				sort.Strings(keys)

				beginRows := 2
				for _, value := range keys {

					f.SetCellValue("Sheet1", "A"+strconv.Itoa(beginRows), mainPhoto[value])
					f.SetCellValue("Sheet1", "B"+strconv.Itoa(beginRows), otherProduct[value])
					f.SetCellValue("Sheet1", "C"+strconv.Itoa(beginRows), otherProduct[value])
					f.SetCellValue("Sheet1", "D"+strconv.Itoa(beginRows), advertising[value])
					beginRows++
				}

				if err := f.SaveAs(mkOutOrg + string(os.PathSeparator) + config.Leve3Dir + strconv.Itoa(outerDirCount) + ".xlsx"); err != nil {
					fmt.Println(err)
				}

				mainPhoto = make(map[string]string)
				otherProduct = make(map[string]string)
				advertising = make(map[string]string)
			}

		}
	}

	//做resize
	if len(needResize) > 0 {
		bgImg := config.BlankImg
		for _, resizeImgPath := range needResize {
			//fmt.Println(resizeImgPath)
			other := resizeImgPath + "RESIZE"
			resizeImg(resizeImgPath, other, 250)
			moveFile(other, resizeImgPath)
			meragePic(bgImg, resizeImgPath, resizeImgPath)
		}
	}

	//沒找到的尺寸表
	if len(sizePicInfoArray) > 0 {
		var txtString string
		for _, errorMsg := range sizePicInfoArray {
			txtString = txtString + errorMsg
		}
		content := []byte(txtString)
		err := ioutil.WriteFile("找不到的尺寸表.txt", content, 0666)
		if err != nil {
			fmt.Println("ioutil WriteFile error: ", err)
		}
	}
	//沒找到的試穿表
	if len(tryPicInfoArray) > 0 {
		var txtString string
		for _, errorMsg := range tryPicInfoArray {
			txtString = txtString + errorMsg
		}
		content := []byte(txtString)
		err := ioutil.WriteFile("找不到的試穿表.txt", content, 0666)
		if err != nil {
			fmt.Println("ioutil WriteFile error: ", err)
		}
	}
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
		fmt.Println(src + "不存在!!")
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

func resizeImg(imgPath string, outPath string, change int) {
	file, err := os.Open(imgPath)
	if err != nil {
		log.Fatal(err)
	}

	// decode jpeg into image.Image
	img, err := jpeg.Decode(file)
	if err != nil {
		log.Fatal(err)
	}
	file.Close()

	width := img.Bounds().Dx()  //
	height := img.Bounds().Dy() //

	// fmt.Println(width)
	// fmt.Println(height)
	// os.Exit(3)
	if width > height {
		width = change
		height = 0
	} else {
		width = 0
		height = change
	}

	m := resize.Resize(uint(width), uint(height), img, resize.Lanczos3)
	out, err := os.Create(outPath)
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()

	jpeg.Encode(out, m, &jpeg.Options{Quality: 100})
}

func meragePic(bgPath string, topPath string, outPath string) {

	file, err := os.Open(topPath)
	if err != nil {
		log.Fatal(err)
	}

	// decode jpeg into image.Image
	img, err := jpeg.Decode(file)
	if err != nil {
		log.Fatal(err)
	}
	file.Close()

	width := (250 - img.Bounds().Dx()) / 2 //
	grids := []*gim.Grid{
		{
			ImageFilePath:   bgPath,
			BackgroundColor: color.White,
			// these grids will be drawn on top of the first grid
			Grids: []*gim.Grid{
				{
					ImageFilePath: topPath,
					OffsetX:       width, OffsetY: 0,
				},
			},
		},
	}
	rgba, err := gim.New(grids, 1, 1).Merge()
	if err != nil {
		log.Fatal(err)
	}
	// save the output to jpg or png
	file, err2 := os.Create(outPath)
	if err2 != nil {
		log.Fatal(err2)
	}
	err = jpeg.Encode(file, rgba, &jpeg.Options{Quality: 100})
}
