package main

import (
	"fmt"
	"log"
	"os"

	"github.com/rwcarlsen/goexif/exif"
	"github.com/rwcarlsen/goexif/mknote"
)

func ExampleDecode() {
	fname := "sample1.jpg"

	f, err := os.Open(fname)
	if err != nil {
		log.Fatal(err)
	}

	// Optionally register camera makenote data parsing - currently Nikon and
	// Canon are supported.
	exif.RegisterParsers(mknote.All...)

	x, err := exif.Decode(f)
	if err != nil {
		log.Fatal(err)
	}

	camModel, _ := x.Get(exif.Model) // normally, don't ignore errors!
	fmt.Println(camModel.StringVal())

	focal, _ := x.Get(exif.FocalLength)
	numer, denom, _ := focal.Rat2(0) // retrieve first (only) rat. value
	fmt.Printf("%v/%v", numer, denom)

	// Two convenience functions exist for date/time taken and GPS coords:
	tm, _ := x.DateTime()
	fmt.Println("Taken: ", tm)

	lat, long, _ := x.LatLong()
	fmt.Println("lat, long: ", lat, ", ", long)

	lat2, _ := x.Get(exif.GPSLatitude)
	fmt.Println("Raw latitude", lat2.String())

	l0, err := lat2.Rat(0)
	if err != nil {
		fmt.Println("cannot get coord", err.Error())
	}
	l1, _ := lat2.Rat(1)
	l2, _ := lat2.Rat(2)
	ref, _ := x.Get(exif.GPSLatitudeRef)
	refStr, _ := ref.StringVal()

	fmt.Printf("GPS_Latitude: \"%d deg %d' %d %s--%f\"\n", l0.Num(), l1.Num(), l2.Num(), refStr, lat)
}

func main() {
	fmt.Println("Functions 2: ")
	ExampleDecode()

}