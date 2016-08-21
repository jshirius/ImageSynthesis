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
  "image/color"
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
  addimg := getIMG("photobyhikage.png");


  for i := 0; i < srcimg.Bounds().Max.Y; i++ {
        for j := 0; j < srcimg.Bounds().Max.X; j++ {
            //r, g, b, _ := addimg.At(j, i).RGBA()
            //srcimg.Set(j, i, addimg.At(j, i))
            //i := rgb2int(int(r), int(g), int(b))
            //hist[i]++
        }
    }

    //白黒に変換する
  //outputImage := convertToMonochromeImage(srcimg); // 変換
  
  //画像をマージする
  outputImage := convertSynthesisImage(srcimg, addimg);

  //fmt.Println("Width:", img.Width, "Height:", img.Height)



  //画像ファイルとして保存する
  saveImage(outputImage)

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
// 引数で与えられたイメージをモノクロに変換して返す関数
// --------------------------------------------------
func convertToMonochromeImage(inputImage image.Image) image.Image {
    rect            := inputImage.Bounds()
    width           := rect.Size().X
    height          := rect.Size().Y
    rgba            := image.NewRGBA(rect)
 
    for x := 0; x < width; x++ {
        for y := 0; y < height; y++ {
            var col color.RGBA
            // 座標(x,y)のR, G, B, α の値を取得
            r,g,b,a := inputImage.At(x ,y).RGBA()

            // それぞれを重み付けして足し合わせる(NTSC 系加重平均法)
            outR := float32(uint8(r)) * 0.298912
            outG := float32(uint8(g)) * 0.58611
            outB := float32(uint8(b)) * 0.114478
            mono := uint8(outR + outG + outB)
            col.R = mono
            col.G = mono
            col.B = mono
            col.A = uint8(a)
            rgba.Set(x, y, col)
        }
    }

    return rgba.SubImage(rect)
}


// --------------------------------------------------
// ２つの画像を一つの画像にマージする
//  synthesisは、srcImageより小さい画像であること
// --------------------------------------------------
func convertSynthesisImage(srcImage image.Image, synthesis image.Image) image.Image {
    rect            := srcImage.Bounds()
    width           := rect.Size().X
    height          := rect.Size().Y
    rgba            := image.NewRGBA(rect)
 
    //元画像のサイズ出力
    fmt.Printf("元画像 Width=%d, Height=%d\n", width, height)

    //合成用画像の出力サイズ
    fmt.Printf("合成用画像 Width=%d, Height=%d\n", synthesis.Bounds().Size().X, synthesis.Bounds().Size().Y)

    //座標位置は、左上基準
    for x := 0; x < width; x++ {
        for y := 0; y < height; y++ {
            var col color.RGBA
            // 座標(x,y)のR, G, B, α の値を取得
            r,g,b,a := srcImage.At(x ,y).RGBA()

            //if(synthesis.Bounds().Size().X > x && synthesis.Bounds().Size().Y > y ){
            //    r,g,b,a = synthesis.At(x ,y).RGBA()
            //}
            // それぞれを重み付けして足し合わせる(NTSC 系加重平均法)
            //outR := float32(uint8(r)) * 0.298912
            //outG := float32(uint8(g)) * 0.58611
            //outB := float32(uint8(b)) * 0.114478
            //mono := uint8(outR + outG + outB)
            col.R = uint8(r)
            col.G = uint8(g)
            col.B = uint8(b)
            col.A = uint8(a)
            rgba.Set(x, y, col)
        }
    }


    for x := 0; x < synthesis.Bounds().Size().X  ;x++ {
        for y := 0; y < synthesis.Bounds().Size().Y ; y++ {
            var col color.RGBA

            if(synthesis.Bounds().Size().X > x && synthesis.Bounds().Size().Y > y ){
                r,g,b,a := synthesis.At(x ,y).RGBA()
                if(a > 1){
                                 col.R = uint8(r)
                col.G = uint8(g)
                col.B = uint8(b)
                col.A = uint8(a)
                rgba.Set(x, y, col)   
                }

            }          
        }
    }


    return rgba.SubImage(rect)
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
