package main

import (
  "fmt"
  "os"
  "log"
  "errors"
  "os/exec"
  "time"
  
  "github.com/sevlyar/go-daemon"
  "github.com/urfave/cli" 
)

type DayTimes struct {
    sunrise int
    day int
    sunset int
    night int
}

// func determineTime() {
//  const (
//    sunriseHour int = 6
//    dayHour int     = 11
//    sunsetHour int  = 19
//    nightHour int   = 20
//  )
//
//  now := time.Now()
//  currentHour := now.Hour()
//
//  var wallpaper string
//
//  switch {
//  case currentHour >= nightHour:
//    wallpaper = "~/Pictures/Wallpapers/etc/outset-island/night.jpg"
//     exec.Command("feh", "--bg-fill", wallpaper)
//    fmt.Println("Wallpaper set to:", wallpaper)
//  case currentHour >= sunsetHour:
//    wallpaper = "~/Pictures/Wallpapers/etc/outset-island/beforeYafter.png"
//     exec.Command("feh", "--bg-fill", wallpaper)
//    fmt.Println("Wallpaper set to:", wallpaper)
//  case currentHour >= dayHour:
//    wallpaper = "~/Pictures/Wallpapers/etc/outset-island/day.png"
//     exec.Command("feh", "--bg-fill", wallpaper)
//    fmt.Println("Wallpaper set to:", wallpaper)
//  case currentHour >= sunriseHour:
//    wallpaper = "~/Pictures/Wallpapers/etc/outset-island/beforeYafter.png"
//     exec.Command("feh", "--bg-fill", wallpaper)
//    fmt.Println("Wallpaper set to:", wallpaper)
//  }
//
// }
//
// func main() {
//  determineTime()
//
//  now := time.Now()
//  fmt.Printf("Current hour: %d\n", now.Hour())
// }

func activateDaemon() {
    cntxt := &daemon.Context{
        PidFileName: "dynamic-paper.pid",
        PidFilePerm: 0644,
        LogFileName: "dynamic-paper.log",
        LogFilePerm: 0640,
        WorkDir:     "./",
        Umask:       027,
    }

    d, err := cntxt.Reborn()
    if err != nil {
        log.Fatal("Unable to run: ", err)
    }
    if d != nil {
        return
    }
    defer cntxt.Release()

    log.Print("- - - - - - - - - - - - - - -")
    log.Print("Dynamic-Paper Daemon Activated")
    log.Print("- - - - - - - - - - - - - - -")
    log.Print("Reading the current local time")
    log.Print("- - - - - - - - - - - - - - -")

    ticker := time.NewTicker(1 * time.Minute)
    defer ticker.Stop()

    for {
        select {
        case <-ticker.C:
             cmd := exec.Command("echo", "hello world")
             cmd.Stdout = os.Stdout
             cmd.Stderr = os.Stderr
             err := cmd.Run()
             if err != nil {
                fmt.Println("Error executing command:", err)
             }
        }
    }
}

func main() {
    app := &cli.App{
        Name:  "dynamic-paper",
        Usage: "Define a wallpaper for the time of day",
        Action: func(*cli.Context) error {
            return errors.New("Please rerun the program with an argument...")
        },
    }

    if err := app.Run(os.Args); err != nil {
        log.Fatal(err)
    }
}
