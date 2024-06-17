package main

import (
  "errors"
  "fmt"
  "log"
  "os"
  "os/exec"
  "os/user"
  "path/filepath" 
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

type LocalPapers struct {
    sunrise string
    day string
    sunset string
    night string
}

var DesktopSession string = os.Getenv("XDG_SESSION_TYPE")
var localTimes LocalTimes
var localPapers LocalPapers
var wallPath string

var defaultTimes = LocalTimes{
    sunrise: 6,
    day:     11,
    sunset:  19,
    night:   20,
}

func detectDefaultWallpapers() {
    wallPath = os.Getenv("DP_WALLPATH")
    if wallPath == "" {
        fmt.Println("$DP_WALLPATH not detected, please set this variable or use the \"load\" command")
        fmt.Println("$ export DP_WALLPATH=\"$HOME/YOURPATH/\"\nOr set via your shellrc file")
        errors.New("ABORTING PROGRAM")
    } else {
        log.Print("Using wallpaper path from $DP_WALLPATH:", wallPath)
        fmt.Println("Using wallpaper path from $DP_WALLPATH:", wallPath)
        loadDefaultWallpapers()
        setWallpaper()
    }
}

func loadDefaultWallpapers() {
    err := filepath.Walk(wallPath, func(path string, info os.FileInfo, err error) error {
        if err != nil {
        return err
    }

    if !info.IsDir() {
        lowerName := strings.ToLower(info.Name())
        switch {
        case strings.Contains(lowerName, "sunrise"):
            localPapers.sunrise = path
        case strings.Contains(lowerName, "day"):
            localPapers.day = path
        case strings.Contains(lowerName, "sunset"):
            localPapers.sunset = path
        case strings.Contains(lowerName, "night"):
            localPapers.night = path
        }
    }
    return nil
  })

    if err != nil {
        log.Fatalf("Error loading default wallpapers: %v", err)
    }
    log.Print("Loaded wallpapers: %+v\n", localPapers)
    fmt.Printf("Loaded wallpapers: %+v\n", localPapers)
}


func expandPath(path string) string {
    if path[:2] == "~/" {
        usr, err := user.Current()
        if err != nil {
            log.Fatal(err)
        }
        return filepath.Join(usr.HomeDir, path[2:])
    }
    return path
}

func processWallpapers(wallPaths string) {
    paths := strings.Split(wallPaths, ",")
    if len(paths) != 4 {
        fmt.Println("Please provide exactly 4 wallpaper paths")
        return
    }

    for i := range paths {
        paths[i] = strings.TrimSpace(paths[i])
    }

    if paths[0] != "" {
        localPapers.sunrise = paths[0]
    }
    if paths[1] != "" {
        localPapers.day = paths[1]
    }
    if paths[2] != "" {
        localPapers.sunset = paths[2]
    }
    if paths[3] != "" {
        localPapers.night = paths[3]
    }

    for i, wallPath := range paths {
        fmt.Println("Loaded wallpaper:", wallPath)
        log.Print("Loaded wallpaper:", i, "", wallPath)
    }
    setWallpaper()
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

func setWallpaper() bool {
    currentHour := time.Now().Hour()

    var wallprogram, args, wallpaper string

    if DesktopSession == "x11" {
        wallprogram = "feh"
        switch {
        case currentHour >= localTimes.night:
            wallpaper = expandPath(localPapers.night) 
        case currentHour >= localTimes.sunset:
            wallpaper = expandPath(localPapers.sunset)
        case currentHour >= localTimes.day:
            wallpaper = expandPath(localPapers.day)
        case currentHour >= localTimes.sunrise:
            wallpaper = expandPath(localPapers.sunrise)
        }
    } else if DesktopSession == "wayland" {
        wallprogram = "swaybg"
        switch {
        case currentHour >= localTimes.night:
            wallpaper = expandPath(localPapers.night) 
        case currentHour >= localTimes.sunset:
            wallpaper = expandPath(localPapers.sunset)
        case currentHour >= localTimes.day:
            wallpaper = expandPath(localPapers.day)
        case currentHour >= localTimes.sunrise:
            wallpaper = expandPath(localPapers.sunrise)
        }
    } else {
        errors.New("Unable to determine $XDG_SESSION_TYPE...\n(Do you have a 'xdg-desktop-portal' installed?)")
        return false
    }

    if wallprogram == "feh" {
        args = "--bg-fill"
    } else if wallprogram == "swaybg" {
        args = "-i"
    } else {
        errors.New("Unable to detect your wallpaper program")
        log.Println("Honesly dont know what happened... UNEXPECTED ERROR")
    }

    cmd := exec.Command(wallprogram, args, wallpaper)
    cmdOutput, err := cmd.CombinedOutput()
    if err != nil {
        log.Println("Error executing command:", err)
        log.Println("Command output:", string(cmdOutput))
        fmt.Println("Error executing command:", err)
        fmt.Println("Command output:", string(cmdOutput))
        return false
    }
    log.Print("Wallpaper set to:", wallpaper)
    fmt.Println("Wallpaper set to:", wallpaper)
    return true
}

func activateDaemon() {
    homeDir, err := os.UserHomeDir()
    if err != nil {
        log.Fatal("Unable to get user home directory: ", err)
    }

    daemonDir := filepath.Join(homeDir, ".local", "share", "dynamic-paper")

    if _, err := os.Stat(daemonDir); os.IsNotExist(err) {
        err := os.MkdirAll(daemonDir, 0755)
        if err != nil {
            log.Fatal("Unable to create daemon directory: ", err)
        }
    }

    logFile := filepath.Join(daemonDir, "dynamic-paper.log")
    pidFile := filepath.Join(daemonDir, "dynamic-paper.pid")

    cntxt := &daemon.Context{
        PidFileName: pidFile,
        PidFilePerm: 0644,
        LogFileName: logFile,
        LogFilePerm: 0640,
        WorkDir:     daemonDir,
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
            detectDefaultWallpapers() 
        }
    }
}

func killDaemon() {
    homeDir, err := os.UserHomeDir()
    daemonDir := filepath.Join(homeDir, ".local", "share", "dynamic-paper")
    pidFile := filepath.Join(daemonDir, "dynamic-paper.pid")
    data, err := os.ReadFile(pidFile)
    if err != nil {
        log.Fatalf("Unable to read PID file: %v", err)
    }

    pid, err := strconv.Atoi(string(data))
    if err != nil {
        log.Fatalf("Invalid PID found in file: %v", err)
    }

    proc, err := os.FindProcess(pid)
    if err != nil {
        log.Fatalf("Unable to find process: %v", err)
    }

    err = proc.Kill()
    if err != nil {
        log.Fatalf("Failed to kill process: %v", err)
    }

    fmt.Printf("Killed daemon with PID %d\n", pid)
    os.Remove(pidFile)
}

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
                Aliases: []string{"ad"},
                Usage:   "Activates the daemon",
                Action: func(*cli.Context) error {
                    activateDaemon()
                    return nil
                },
            },
            {
                Name:    "kill-daemon",
                Aliases: []string{"kd"},
                Usage:   "Kills the running daemon",
                Action: func(*cli.Context) error {
                    killDaemon()
                    return nil
                },
            },
            {
                Name:    "load",
                Aliases: []string{"l"},
                Usage:   "Loads 4 wallpapers for usage (Order: Sunrise, Day, Sunset, Night)",
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
                Usage:   "Sets a custom time for each time of day (Order: Sunrise, Day, Sunset, Night)",
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
                Name:    "run",
                Aliases: []string{"r"},
                Usage:   "Sets the wallpaper for the current time of day",
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
