package utils

import (
    "log"
    "image/jpeg"
    "mime/multipart"
    "github.com/disintegration/imaging"
)

func ConvertToGrayScale(file multipart.File) ([16]int, error) {
    // load original image
    file.Seek(0, 0) // seek back to the beginning of the file
    src, err := jpeg.Decode(file)
    if err != nil {
        log.Println(err)
        return [16]int{}, err
    }

    // grayscale the image
    grayimg := imaging.Grayscale(src)

    // Calculate a 16-bin histogram for gray component
    bounds := grayimg.Bounds()
    var histogram [16]int
    for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
        for x := bounds.Min.X; x < bounds.Max.X; x++ {
            r, _, _, _ := grayimg.At(x, y).RGBA()
            // A color's RGBA method returns values in the range [0, 65535].
            // Shifting by 12 reduces this to the range [0, 15].
            histogram[r>>12]++
        }
    }

    // Print the results.
    // for _, x := range histogram {
    //     log.Printf("%6d\n", x)
    // }

    return histogram, nil
}
