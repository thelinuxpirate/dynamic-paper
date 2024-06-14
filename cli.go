package main

import (
  "fmt"
  "os/exec"
  "time"
)

func determineTime() {
	const (
		sunriseHour int = 6
		dayHour int     = 11
		sunsetHour int  = 19
		nightHour int   = 20
	)

	now := time.Now()
	currentHour := now.Hour()

	var wallpaper string

	switch {
	case currentHour >= nightHour:
		wallpaper = "~/Pictures/Wallpapers/etc/outset-island/night.jpg"
    exec.Command("feh", "--bg-fill", wallpaper)
		fmt.Println("Wallpaper set to:", wallpaper)
	case currentHour >= sunsetHour:
		wallpaper = "~/Pictures/Wallpapers/etc/outset-island/beforeYafter.png"
    exec.Command("feh", "--bg-fill", wallpaper)
		fmt.Println("Wallpaper set to:", wallpaper)
	case currentHour >= dayHour:
		wallpaper = "~/Pictures/Wallpapers/etc/outset-island/day.png"
    exec.Command("feh", "--bg-fill", wallpaper)
		fmt.Println("Wallpaper set to:", wallpaper)
	case currentHour >= sunriseHour:
		wallpaper = "~/Pictures/Wallpapers/etc/outset-island/beforeYafter.png"
    exec.Command("feh", "--bg-fill", wallpaper)
		fmt.Println("Wallpaper set to:", wallpaper)
	}

}

func main() {
	determineTime()

	now := time.Now()
	fmt.Printf("Current hour: %d\n", now.Hour())
}
