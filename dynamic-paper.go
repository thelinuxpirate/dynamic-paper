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
    "syscall"
    "time"

    "github.com/sevlyar/go-daemon"
    "github.com/urfave/cli/v2"
)

type LocalTimes struct {
    sunrise int
    day     int
    sunset  int
    night   int
}

type LocalPapers struct {
    sunrise string
    day     string
    sunset  string
    night   string
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
        log.Println("$DP_WALLPATH not detected, please set this variable or use the \"load\" command")
        log.Println("$ export DP_WALLPATH=\"$HOME/YOURPATH/\"\nOr set via your shellrc file")
        errors.New("ABORTING PROGRAM")
    } else {
        log.Println("Using wallpaper path from $DP_WALLPATH:", wallPath)
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
    log.Printf("Loaded wallpapers: %+v\n", localPapers)
}

func expandPath(path string) string {
    if strings.HasPrefix(path, "~/") {
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
        log.Printf("Loaded wallpaper %d: %s", i, wallPath)
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
        log.Println("Unable to determine $XDG_SESSION_TYPE. Do you have 'xdg-desktop-portal' installed?")
        return false
    }

    if wallprogram == "feh" {
        args = "--bg-fill"
    } else if wallprogram == "swaybg" {
        args = "-i"
    } else {
        log.Println("Unable to detect your wallpaper program")
        return false
    }

    cmd := exec.Command(wallprogram, args, wallpaper)
    cmdOutput, err := cmd.CombinedOutput()
    log.Printf("%s %s %s", wallprogram, args, wallpaper)
    if err != nil {
        log.Printf("Error executing command: %v", err)
        log.Printf("Command output: %s", string(cmdOutput))
        return false
    }
    log.Printf("Wallpaper set to: %s", wallpaper)
    return true
}


func isDesktopSessionActive() bool {
    cmd := exec.Command("loginctl", "show-session", os.Getenv("XDG_SESSION_ID"), "-p", "Active")
    output, err := cmd.CombinedOutput()
    if err != nil {
        log.Println("Error checking session active status:", err)
        return false
    }
    return strings.Contains(string(output), "yes")
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

    detectDefaultWallpapers()
    setWallpaper() 

    initialTicker := time.NewTicker(10 * time.Second)
    stopChan := make(chan bool)

    go func() {
        for i := 0; i < 10; i++ {
            <-initialTicker.C
            setWallpaper()
        }
        initialTicker.Stop()
        stopChan <- true
    }()

    <-stopChan

    hourlyTicker := time.NewTicker(1 * time.Hour)
    defer hourlyTicker.Stop()

    for {
        select {
        case <-hourlyTicker.C:
            detectDefaultWallpapers()
            setWallpaper()
        }
    }
}

func killDaemon() {
    homeDir, err := os.UserHomeDir()
    if err != nil {
        log.Fatal("Unable to get user home directory: ", err)
    }

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

    err = proc.Signal(syscall.Signal(0))
    if err == nil {
        err = proc.Kill()
        if err != nil {
            log.Fatalf("Failed to kill process: %v", err)
        }
        fmt.Printf("Killed daemon with PID %d\n", pid)
        os.Remove(pidFile)
    } else if err.Error() == "no such process" {
        fmt.Println("Daemon process already terminated")
        os.Remove(pidFile)
    } else {
        log.Fatalf("Error checking process: %v", err)
    }
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
