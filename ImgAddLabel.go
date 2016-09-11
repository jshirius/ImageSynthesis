/*
* ２つの画像を合成する
* 主に写真にラベルをつけるなどの機能のために使用する
*
*/
package main

import (
  "os"
  "fmt"
  "time"
  "image"
  "image/jpeg"
  _"image/png"
  "image/draw"
    "io/ioutil"
    "path/filepath"
)

func main() {
    path := ""
    start := time.Now()
    if len(os.Args) >= 2 {
        path = os.Args[1]
    } else {
        fmt.Println("コマンドライン引数に画像ファイルを指定してください")
        return
    }

  srcimg := getIMG(path);
  addimg := getIMG("photoLabel.png");

//合成処理
outputImage := convertDrawImage(srcimg, addimg);



  //画像ファイルとして保存する
  saveImage(outputImage)
//saveBmpImage(outputImage) 

  // ビルドエラー回避
  var _ = path
  var _ = start
  var _ = srcimg
  var _ = addimg
 

//ファイル検索
listFiles("/Users/takizawa/Pictures/20160807_koiwa_hanabi/", "/Users/takizawa/Pictures/20160807_koiwa_hanabi/") 
}



//画像を読み込む
func getIMG(path string) image.Image {
    file, err := os.Open(path)

    //fmt.Println(file);
    defer file.Close()
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }

    //ファイルの情報取得
    conf, _, err := image.DecodeConfig(file)
    if err != nil {
        //log.Fatal(err)
    }

    
    fmt.Printf("Width=%d, Height=%d\n", conf.Width, conf.Height)

    //ファイルポインタを一番最初に戻す
    file.Seek(0, 0)

    

    img, _, err := image.Decode(file)
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }



    return img
}

//画像を保存する
func saveImage(img image.Image) {
    out, err := os.Create("output.jpg")
    defer out.Close()
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
    option := &jpeg.Options{Quality: 100}
    err = jpeg.Encode(out, img,option)
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
}


// --------------------------------------------------
// ２つの画像を一つの画像にマージする
//  synthesisは、srcImageより小さい画像であること
// --------------------------------------------------
func convertDrawImage(srcImage image.Image, synthesis image.Image) image.Image {


    rect            := srcImage.Bounds()
    width           := rect.Size().X
    height          := rect.Size().Y
    rgba            := image.NewRGBA(rect)
 
    //imageをコピーする
    for x := 0; x < width; x++ {
        for y := 0; y < height; y++ {
             rgba.Set(x, y, srcImage.At(x, y))
        }
    }

    fmt.Printf("元画像の大きさ=%s", srcImage.Bounds())
    fmt.Printf("合成画像の大きさ=%s", synthesis.Bounds())

    //貼り付け位置調整
    var mergin int
    mergin = 20
    a := image.Pt((-1)*(srcImage.Bounds().Max.X -  synthesis.Bounds().Max.X - mergin),
        (-1)*(srcImage.Bounds().Max.Y -  synthesis.Bounds().Max.Y - mergin))
    //synRect.Max.X = 1200;

    //合成処理
    draw.Draw(rgba, rect, synthesis, a, draw.Over)

    return rgba
}






func listFiles(rootPath, searchPath string) {
    fis, err := ioutil.ReadDir(searchPath)

    if err != nil {
        panic(err)
    }

    for _, fi := range fis {
        fullPath := filepath.Join(searchPath, fi.Name())

        if fi.IsDir() {
            listFiles(rootPath, fullPath)
        } else {
            rel, err := filepath.Rel(rootPath, fullPath)

            if err != nil {
                panic(err)
            }

            fmt.Println(rel)
        }
    }
}
