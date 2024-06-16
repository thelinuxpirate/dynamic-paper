package main

import (
  "errors"
  "fmt"
  "log"
  "os"
  "os/exec"
  "sort"
  "strconv"
  "strings"
  "time" 

  "github.com/sevlyar/go-daemon"
  "github.com/urfave/cli/v2" 
)

type LocalTimes struct {
    sunrise int
    day int
    sunset int
    night int
}

type localPapers struct {
    sunrise string
    day string
    sunset string
    night string
}

var DesktopSession string = os.Getenv("XDG_SESSION_TYPE")
var localTimes LocalTimes

var defaultTimes = LocalTimes{
    sunrise: 6,
    day:     11,
    sunset:  19,
    night:   20,
}

func detectDefaultWallpapers() {
    var wallPath string = os.Getenv("DP_WALLPATH")
    if wallPath == "" {
       setEnv(0) 
    } else {
        fmt.Println("W")
    }
}

func setEnv(x int) {
    if x == 0 {
        fmt.Println("DP_WALLPATH not detected, would you like to set it?")
    } else {
        fmt.Println("Error handeling")
    }
}

func finalizeTime(usrTimes string) error {
    times := strings.Split(usrTimes, ",")
    if len(times) != 4 {
        return errors.New("Please provide exactly 4 time values")
    }

    var err error
    localTimes.sunrise, err = strconv.Atoi(times[0])
    if err != nil {
        return err
    }
    localTimes.day, err = strconv.Atoi(times[1])
    if err != nil {
        return err
    }
    localTimes.sunset, err = strconv.Atoi(times[2])
    if err != nil {
        return err
    }
    localTimes.night, err = strconv.Atoi(times[3])
    if err != nil {
        return err
    }
    
    return nil
}

func processWallpapers(wallPaths string) {
    paths := strings.Split(wallPaths, ",")
    if len(paths) != 4 {
        fmt.Println("Please provide exactly 4 wallpaper paths")
        return
    }

    for _, wallPath := range paths {
      fmt.Println("Loaded wallpaper:", wallPath)
      log.Print("Loaded wallpaper:", wallPath)
    }
}

func setWallpaper() bool {
    currentHour := time.Now().Hour()

    var wallprogram, wallpaper string

    if DesktopSession == "x11" {
        wallprogram = "feh"
        switch {
        case currentHour >= localTimes.night:
            wallpaper = "~/Pictures/Wallpapers/etc/outset-island/night.jpg"
            exec.Command(wallprogram, "--bg-fill", wallpaper)
            fmt.Println("Wallpaper set to:", wallpaper)
        case currentHour >= localTimes.sunset:
            wallpaper = "~/Pictures/Wallpapers/etc/outset-island/beforeYafter.png"
            exec.Command(wallprogram, "--bg-fill", wallpaper)
            fmt.Println("Wallpaper set to:", wallpaper)
        case currentHour >= localTimes.day:
            wallpaper = "~/Pictures/Wallpapers/etc/outset-island/day.png"
            exec.Command(wallprogram, "--bg-fill", wallpaper)
            fmt.Println("Wallpaper set to:", wallpaper)
        case currentHour >= localTimes.sunrise:
            wallpaper = "~/Pictures/Wallpapers/etc/outset-island/beforeYafter.png"
            exec.Command(wallprogram, "--bg-fill", wallpaper)
            fmt.Println("Wallpaper set to:", wallpaper)
        }
    } else if DesktopSession == "wayland" {
        wallprogram = "swaybg"
        switch {
        case currentHour >= localTimes.night:
            wallpaper = "~/Pictures/Wallpapers/etc/outset-island/night.jpg"
            exec.Command(wallprogram, "--bg-fill", wallpaper)
            fmt.Println("Wallpaper set to:", wallpaper)
        case currentHour >= localTimes.sunset:
            wallpaper = "~/Pictures/Wallpapers/etc/outset-island/beforeYafter.png"
            exec.Command(wallprogram, "--bg-fill", wallpaper)
            fmt.Println("Wallpaper set to:", wallpaper)
        case currentHour >= localTimes.day:
            wallpaper = "~/Pictures/Wallpapers/etc/outset-island/day.png"
            exec.Command(wallprogram, "--bg-fill", wallpaper)
            fmt.Println("Wallpaper set to:", wallpaper)
        case currentHour >= localTimes.sunrise:
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
    localTimes = defaultTimes
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
                    if cCtx.NArg() != 1 {
                        return errors.New("Please provide a comma-separated list of 4 wallpaper paths")
                    }

                    paths := cCtx.Args().Get(0)
                    processWallpapers(paths)
                    return nil
                },
            },
            {
                Name:    "set-time",
                Aliases: []string{"st"},
                Usage:   "Provide a List of 4 Hours (Order: Sunrise, Day, Sunset, Night)",
                Action: func(cCtx *cli.Context) error {
                    if cCtx.NArg() != 1 {
                        return errors.New("Please provide a comma-separated list of 4 time values")
                    }

                    usrTimes := cCtx.Args().Get(0)
                    if err := finalizeTime(usrTimes); err != nil {
                        return err
                    }
                    return nil
                },
            },
            {
                Name:    "run-me",
                Aliases: []string{"r"},
                Usage:   "Im a demo",
                Action: func(cCtx *cli.Context) error {
                    detectDefaultWallpapers()
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
