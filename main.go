//
// A go program to convert a PPM image to grayscale.
// By: Sergio Ambelis Diaz
// U. of Illinios, Chicago
// Spring 2022
//

package main

import (
  "fmt"
  "bufio"
  "os"
  "strconv"
  "strings"
  "time"
)

//
// ReadImageFile
//
// Reads a P3 PPM image file and returns the width,
// height, depth, and pixels.
//
func ReadImageFile(filename string) (int, int, int, [][]int) {
  file, err := os.Open(filename);
  if err != nil {
    panic(err);
  }
  defer file.Close();
  
  scanner := bufio.NewScanner(file);
  scanner.Split(bufio.ScanWords);
  //var result [] int;
  
  position := 1;
  
  width  := 0;
  height := 0;
  depth  := 0;
  
  var pixels [][]int;
  var row    []int;
  rowcount := 0;
  
  for scanner.Scan() {
  
    if position == 1 { position++; continue; }  // header, skip
  
    x, err := strconv.Atoi(scanner.Text())
    if err != nil {
      panic(err);
    }

    if position == 2 { width  = x; position++; continue; }
    if position == 3 { height = x; position++; continue; }
    if position == 4 { depth  = x; position++; continue; }
    
    // else RGB value:
    row = append(row, x);
    rowcount++;
    if rowcount == (width*3) {
      pixels = append(pixels, row);
      row = make([]int, 0);
      rowcount = 0;
    }
  }
  
  return width, height, depth, pixels;
}


//
// WriteImageFile
//
// Writes the image to the file as a human-readable
// PPM image ("P3" format).
//
func WriteImageFile(filename string, width int, height int, depth int, pixels [][]int) {
  file, err := os.Create(filename);
  if err != nil {
    panic(err);
  }
  defer file.Close();
  
  file.WriteString("P3\n");
  file.WriteString(fmt.Sprintf("%d %d\n", width, height));
  file.WriteString(fmt.Sprintf("%d\n", depth));
  
  for r := range pixels {
    for c := range pixels[r] {
      file.WriteString(fmt.Sprintf("%d ", pixels[r][c]));
    }
    file.WriteString("\n");
  }
}


//
// debug
//
// Writes the PPM image to the console.
//
func debug(width int, height int, depth int, pixels [][]int) {
  fmt.Println(width);
  fmt.Println(height);
  fmt.Println(depth);
  
  for r := range pixels {
    //fmt.Printf("Row: %v\n", r);
    fmt.Println(pixels[r]);
  }
}


//
// grawOneRow
//
// Converts one row of the image to grayscale.
// Each pixel's RGB value is set to the average of
// those values.
//
func grayOneRow(width int, row []int) {

    for c := 0; c < (width*3); c += 3 {
    
      avg := (row[c] + row[c+1] + row[c+2]) / 3;
      
      row[c]   = avg;
      row[c+1] = avg;
      row[c+2] = avg;
    } 

}


//
// grayscale
//
// Converts the given image to grayscale.
// Each pixel's RGB value is set to the average of
// those values.
//
func grayscale(width int, height int, depth int, pixels [][]int) {

  for r := range pixels {
    grayOneRow(width, pixels[r]);
  }
}


////////////////////////////////////////////////////////
//
// main
//
func main() { 
  fmt.Printf("PPM image file> ");
  
  var filename string;
  fmt.Scanln(&filename);
  
  fmt.Println(fmt.Sprintf("Reading '%s'...", filename));
  
  w, h, d, pixels := ReadImageFile(filename);
  
  fmt.Println("Converting to grayscale...");
  start := time.Now();
  grayscale(w, h, d, pixels);
  duration := time.Since(start);
  //fmt.Println(fmt.Sprintf("time: %v ms", duration.Milliseconds()));
  fmt.Println(fmt.Sprintf("time: %v us", duration.Microseconds()));
  //fmt.Println(fmt.Sprintf("time: %v ns", duration.Nanoseconds()));
 
  filename = strings.Replace(filename, ".ppm", "-grayscale.ppm", 1);
  
  fmt.Println(fmt.Sprintf("Writing '%s'...", filename));
  
  WriteImageFile(filename, w, h, d, pixels);
  
  fmt.Println("done");
}
