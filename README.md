# Dynamic-Paper
A dynamic wallpaper setter for Unix systems.

This is the first version of Dynamic-Paper that fully works as intended!

## Installation
Check out the pre-built binary [here](https://github.com/thelinuxpirate/dynamic-paper/releases)

## Usage
There are multiple commands for you to use & non of them are dependenant on one-another.
```
$ dynamic-paper --help
NAME:
   dynamic-paper - Define a wallpaper for the time of day

USAGE:
   dynamic-paper [global options] command [command options]

COMMANDS:
   daemon, ad       Activates the daemon
   kill-daemon, kd  Kills the running daemon
   load, l          Loads 4 wallpapers for usage (Order: Sunrise, Day, Sunset, Night)
   run, r           Sets the wallpaper for the current time of day
   set-time, st     Sets a custom time for each time of day (Order: Sunrise, Day, Sunset, Night)
   help, h          Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h  show help
```

### The DP_WALLPATH Variable
Dynamic-Paper by default will look for a PATH that can be specified by setting the `DP_WALLPATH` variable globally.
This enables usage with the program to be set at the start of your system like in your `.xinitrc` for example.
#### How to set
There are multiple ways to set an environment variable but I'll document two here:
- In your shell rc (recommended)
You can set this path in your shell's rc being `.bashrc` if you use Bash or `.zshrc` if you use the Zsh shell.
```sh
# Code goes here
export DP_WALLPATH="$HOME/YOURFILEPATHHERE/"
# Maybe some here too
```
- Locally in your current terminal session
By doing this you set the variable only for use in your current terminal session:
```
$ export DP_WALLPATH="$HOME/YOURFILEPATHHERE/"
```
#### How to use the 'load' command
Enter a filepath for this order: Sunrise, Day, Sunset, Night.
Each file is seperated by a `,`
```
$ dynamic-paper load "~/Pictures/Wallpapers/dp_wallpapers/Sunrise.png,~/Pictures/Wallpapers/dp_wallpapers/Day.png,~/Pictures/Wallpapers/dp_wallpapers/Sunset.png,~/Pictures/Wallpapers/dp_wallpapers/Night.jpg"
```
(Same thing is done for the 'set-time' command)

## Release Notes
This release of Dynamic-Paper supports:
- Reading wallpapers from the 'DP_WALLPATH' environment variable
- Auto-detecting whether you are in a X11 or Wayland session (X11 uses feh & Wayland uses Swaybg) 
- Changing the default time integers to ones of your liking
- Loading a wallpaper for each time of day
- Activating the program as a daemon (Log & PID file stored in `~/.local/share/dynamic-paper/`)
- Killing the daemon via command
- And just run the program as it is
