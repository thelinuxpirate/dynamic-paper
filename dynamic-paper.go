package main

import (
  "fmt"
  "os"
  "log"
  "errors"
  "os/exec"
  "time"
  "sort"
  
  "github.com/sevlyar/go-daemon"
  "github.com/urfave/cli/v2" 
)

var DesktopSession string = os.Getenv("XDG_SESSION_TYPE")

type LocalTimes struct {
    sunrise int
    day int
    sunset int
    night int
}

func finalizeTime(usrTimes []int) {
    localTime = LocalTimes{} 

    const (
        sunriseHour int = 6
        dayHour int     = 11
        sunsetHour int  = 19
        nightHour int   = 20
    )
}


func processWallpapers(wallPaths []string) {
    if len(wallPaths) != 4 {
        fmt.Println("Please provide exactly 4 wallpaper paths.")
        return
    }

    for _, wallPath := range Wallpaths {
        fmt.Println("Loaded wallpaper:", Wallpath)
    }
}

func setWallpaper() bool {
    currentHour := time.Now().Hour()

    var wallprogram, wallpaper string

    if DesktopSession == "x11" {
        wallprogram = "feh"
        switch {
        case currentHour >= nightHour:
            wallpaper = "~/Pictures/Wallpapers/etc/outset-island/night.jpg"
            exec.Command(wallprogram, "--bg-fill", wallpaper)
            fmt.Println("Wallpaper set to:", wallpaper)
        case currentHour >= sunsetHour:
            wallpaper = "~/Pictures/Wallpapers/etc/outset-island/beforeYafter.png"
            exec.Command(wallprogram, "--bg-fill", wallpaper)
            fmt.Println("Wallpaper set to:", wallpaper)
        case currentHour >= dayHour:
            wallpaper = "~/Pictures/Wallpapers/etc/outset-island/day.png"
            exec.Command(wallprogram, "--bg-fill", wallpaper)
            fmt.Println("Wallpaper set to:", wallpaper)
        case currentHour >= sunriseHour:
            wallpaper = "~/Pictures/Wallpapers/etc/outset-island/beforeYafter.png"
            exec.Command(wallprogram, "--bg-fill", wallpaper)
            fmt.Println("Wallpaper set to:", wallpaper)
        }
    } else if DesktopSession == "wayland" {
        wallprogram = "swaybg"
        switch {
        case currentHour >= nightHour:
            wallpaper = "~/Pictures/Wallpapers/etc/outset-island/night.jpg"
            exec.Command(wallprogram, "--bg-fill", wallpaper)
            fmt.Println("Wallpaper set to:", wallpaper)
        case currentHour >= sunsetHour:
            wallpaper = "~/Pictures/Wallpapers/etc/outset-island/beforeYafter.png"
            exec.Command(wallprogram, "--bg-fill", wallpaper)
            fmt.Println("Wallpaper set to:", wallpaper)
        case currentHour >= dayHour:
            wallpaper = "~/Pictures/Wallpapers/etc/outset-island/day.png"
            exec.Command(wallprogram, "--bg-fill", wallpaper)
            fmt.Println("Wallpaper set to:", wallpaper)
        case currentHour >= sunriseHour:
            wallpaper = "~/Pictures/Wallpapers/etc/outset-island/beforeYafter.png"
            exec.Command(wallprogram, "--bg-fill", wallpaper)
            fmt.Println("Wallpaper set to:", wallpaper)
        }
    } else {
        errors.New("Unable to determine $XDG_SESSION_TYPE...\n(Do you have a 'xdg-desktop-portal' installed?)")
        return false
    }

    return true
}

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

    ticker := time.NewTicker(1 * time.Hour)
    defer ticker.Stop()

    for {
        select {
        case <-ticker.C:
            setWallpaper()
            err := setWallpaper()
            if err == false {
                errors.New("Error executing command...")
            }
        }
    }
}

// make sure to echo PID & make a function that kills the program using the current PID
// develop an environment variable for Wallpaper path (DP_WALLPATH)
func main() {
    app := &cli.App{
        Name:  "dynamic-paper",
        Usage: "Define a wallpaper for the time of day",
        Action: func(*cli.Context) error {
            return errors.New("Please rerun the program with an argument...\n(dynamic-paper --help)")
        },
        Commands: []*cli.Command{
            {
                Name:    "daemon",
                Aliases: []string{"dm"},
                Usage:   "Activates the Dynamic-Paper Daemon",
                Action: func(*cli.Context) error {
                    activateDaemon()
                    return nil
                },
            },
            {
                Name:    "load",
                Aliases: []string{"l"},
                Usage:   "Provide a List of 4 Wallpaper Paths for Usage",
                Action: func(cCtx *cli.Context) error {
                    if cCtx.NArg() != 4 {
                        return errors.New("Please provide exactly 4 wallpaper paths")
                    }

                    paths := cCtx.Args().Slice()
                    processWallpapers(paths)
                    return nil
                },
            },
            {
                Name:    "set-time",
                Aliases: []string{"st"},
                Usage:   "Provide a List of 4 Hours (Order: Sunrise, Day, Sunset, Night)",
                Action: func(cCtx *cli.Context) error {
                    if cCtx.NArg() != 4 {
                        return errors.New("Please provide exactly 4 wallpaper paths")
                    }

                    paths := cCtx.Args().Slice()
                    processWallpapers(paths)
                    return nil
                },
            },
        },
    }

    sort.Sort(cli.FlagsByName(app.Flags))
    sort.Sort(cli.CommandsByName(app.Commands))

    if err := app.Run(os.Args); err != nil {
        log.Fatal(err)
    }
}
