# LoL Monitor: Stay in control of your gaming session

![](assets/lolmonitor_wide_azure.jpg)

LoL Monitor is a tool to help League of Legends players manage their gaming sessions, enforcing cooldowns between games or maximum number of games played in a session.

This tool is currently only designed for use by Windows users.

## Quick start

### Download the program

The easiest way to get started is to download [lolmonitor.exe](/build/lolmonitor.exe). Clicking this link takes you to the Github page for the file, in which you can download it (by clicking the "..." elipses in the top right, or using Ctrl+Shift+S). Save it wherever you like and run it. This will create a config.json file in the same location.

### Default settings

When you run the program, it will:
1. Close the Lobby after each game. There will be a 30 second delay, so you have time to honor your team mates or take a quick look at the damage chart.
2. Allow you to play 3 games in a session (a ["3 block"](https://www.youtube.com/watch?v=6K1xBJCjxi0&ab_channel=BrokenByConcept)) with a short 5 minute break between each game.
3. Ignore games under 15 minutes in duration (remakes).
4. Require a longer break of 1 hour between sessions.
5. Load automatically whenever you log into windows.

## Customise your experience

You can change things to suit your preferences, by editing the values in the config.json file, using any sort of text editor, such as Notepad. The table below shows the values available in the file:

| Parameter                     | Description                                                                 | Default Value | Alternative Values | Value to Disable |
|------------------------------|-----------------------------------------------------------------------------|---------------|----------------------------|------------------|
| breakBetweenGamesMinutes     | Minimum enforced break (in minutes) between individual games.              | 5             | 1, 15                     | 0                |
| breakBetweenSessionsMinutes  | Minimum enforced break (in minutes) between gaming sessions.               | 60            | 30, 120                     | 0                |
| gamesPerSession              | Maximum number of games allowed in a single session.                       | 3             | 1, 2                       | 0                |
| minimumGameDurationMinutes   | Minimum duration (in minutes) a game must last to count towards a session. | 15            | 5, 20                     | 0                |
| lobbyCloseDelaySeconds       | Delay (in seconds) before closing the lobby window.             | 30            | 0, 60                     | 0                |
| dailyStartTime               | Earliest time of the day when you can open the lobby (24-hour format).          | 00:00         | 04:00, 06:00               | 00:00            |
| dailyEndTime                 | Latest time of the day when you can open the lobby (24-hour format).                    | 00:00         | 22:00, 23:00               | 00:00            |
| loadOnStartup                | Whether the program should auto-load when you log in.                    | true          | false                      | n/a              |

A few examples:

* Let's say you want to only play one game each day, so the game is more memorable and you have time to review it in detail and reflect about it. You would adjust gamesPerSession to 1.

* Let's say you are happy to play as many games as you like (particularly on the weekend), so long as you don't queue after 10pm, because you don't like being overtired - it makes you play worse the next day. You disable the session games limit by setting gamesPerSession to 0 (and may also disable breakBetweenGamesMinutes so you can play continuously) and set the dailyEndTime value to "22:00".

* I personally can't decide whether it's good to have a delay before the lobby is closed. I like the idea of honoring my team, but I feel there's some residual risk of questionable behaviour, such as checking profiles to see if someone was a smurf or reading the post game chat where someone might make a toxic comment about you. I also suspect there's quite a few people out there who would rather not see their rank, to help reduce *ranked anxiety*. If any of this resonates with you, you can set the lobbyCloseDelaySeconds to 0, so the lobby closes instantly.

## Frequently Asked Questions

* What if something goes wrong? I'm 30 minutes into a game and my internet cuts out. What will I do then?

If anything goes wrong, you can always just open Task Manager (right click the task bar at the bottom of the screen and it's towards the bottom of the context menu that appears), find lolmonitor.exe in the list and right click > End task.

* Now that I know I can close down lolmonitor.exe (or edit the config.json file), won't I just cheat and do this to keep playing?

You're welcome to do this! However, you may find that having the lobby close automatically after each game is surprisingly effective. The main goal is just to get the emotional autopilot monkey brain out of the way. If you find that the post game lobby environment makes you want to keep playing longer than you should, try setting the lobbyCloseDelaySeconds to 0 so you never have to deal with it.

* What happens if I am in the queue or in champion select when my dailyEndTime occurs?

The dailyEndTime stops you from opening the lobby, but it won't close the lobby if it's already open.

* What if I don't trust you?

This project is open source, so you're free to look through the code - or ask ChatGPT to do so for you - to gain some comfort. Then you can compile the program yourself using the instructions at the top of [main.go](/cmd/lolmonitor/main.go). If you haven't already, you'll need to have some sort of IDE such as [vscode](https://code.visualstudio.com/) and [install Go](https://go.dev/dl/).

* That's a cool logo?

Thanks! It was made using Generative AI and any resemblances to the actual League of Legends logo are purely coincidental 😛. The three golden squares symbolise the ["3 block"](https://www.youtube.com/watch?v=6K1xBJCjxi0&ab_channel=BrokenByConcept) process, and how this helps contain the League of Legends gaming experience.

## Repo structure

The main folders and files in this repo are as follows:

```bash
lolmonitor/
├── build/                      # Compiled binary and config file
│   ├── config.json             
│   └── lolmonitor.exe          
├── cmd/                        
│   └── lolmonitor/
│       └── main.go             # Entry point for the application
├── devlog.md                   # Development documentation
├── go.mod                      # Go module dependencies
├── internal/                   # Internal packages that are specific to this application
│   ├── application/
│   │   ├── orchestrate.go      # The main loop, monitors for windows and takes actions
│   │   └── orchestrate_test.go 
│   ├── config/                 
│   │   ├── config.go           # Handles interactions with the config.json configuration file
│   │   └── config_test.go      
│   ├── infrastructure/         
│   │   └── startup/
│   │       ├── startup.go      # Uses Windows Registry to auto-load at startup
│   │       └── startup_test.go 
│   └── interfaces/             
│       └── notifications/
│           ├── toast.go        # Displays windows toast notifications
│           └── toast_test.go   
└── pkg/                        # Public package that can be imported by other modules
    └── window/                 
        ├── close.go            # Forcibly closes windows
        ├── close_test.go       
        ├── monitor.go          # Monitors for the opening and closing of windows
        ├── monitor_test.go     
        └── wmi.go              # Windows Management Instrumentation (WMI) Functions and Types
```

## Contributing

If any Mac users want to try porting this application to the Mac environment, this could be useful for other Mac users. You can open a pull request if you get it working!

You're also welcome to provide feedback! You can send it to [alexgouldblog@gmail.com](mailto:alexgouldblog@gamil.com).
