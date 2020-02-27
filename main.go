package main

import (
	"flag"
	"fmt"
	"gocrud/engine"
	. "gocrud/global"
	"gocrud/helper"
	"gocrud/loaders"
	"gocrud/trace"
	"time"

	"github.com/eiannone/keyboard"
)

func main() {
	var err error

	configPath := flag.String("config", "config.json", "config file path")
	flag.Parse()

	engine := new(engine.Engine)
	engine.Config = loaders.LoadConfig(*configPath)

	basic := loaders.LoadBasic(engine.Config)

	engine.Dirs = trace.Directory(engine.Config)
	engine.Header = basic.Header
	engine.Env = helper.ParseEnv(basic.Env)
	engine.InitEnvArr(DIRECTORIES)

	err = keyboard.Open()
	if err != nil {
		panic(err)
	}
	defer keyboard.Close()

	//var ptr *[]string

	fmt.Println("Press ESC to quit, h for help")
	engine.Show()
	for {

		char, key, err := keyboard.GetKey()
		if err != nil {
			panic(err)
		} else if key == keyboard.KeyEsc {
			break
		} else if key == keyboard.KeyCtrlC {
			break
		}

		switch {
		case key == keyboard.KeySpace:
			if engine.SelectedDir != "" && engine.SelectedFile != "" {
				engine.Content = loaders.Open(engine.Config, engine.SelectedDir, engine.SelectedFile)
				engine.SendRequest()
			} else {
				fmt.Println("at first you should choose a file, Directories: d, Files: f")
				engine.ShowCurrentPath()
			}

		case key == keyboard.KeyEnter:
			fmt.Printf("...\n...\n...\n")

		case char == 112: //p
			if engine.SelectedDir != "" && engine.SelectedFile != "" {
				engine.Content = loaders.Open(engine.Config, engine.SelectedDir, engine.SelectedFile)
				engine.ShowCurrentPath()
			} else {
				fmt.Println("at first you should choose a file, Directories: d, Files: f")
				engine.ShowCurrentPath()
			}

		case char == 120: //x
			engine.Page--
			if engine.Page < 1 {
				engine.Page = 1
			}
		case char == 115:
			if engine.SelectedDir == "" {
				fmt.Println("NOTICE: at first choose directory, press d to see directoreis")
			}
			engine.Files = trace.Files(engine.Config, engine.SelectedDir)
			for i, v := range engine.Files {
				engine.Content = loaders.Open(engine.Config, engine.SelectedDir, v)
				engine.SendRequest()
				time.Sleep(100 * time.Millisecond)
				fmt.Printf("%v\n%v", i, v)
			}

			fmt.Println("you pressed s which means start process....", engine.SelectedDir, engine.Files)
		case char == 100: //d
			engine.Page = DIRECTORIES
			engine.Dirs = trace.Directory(engine.Config)
			engine.Show()
		case char == 102: //f
			if engine.SelectedDir == "" {
				fmt.Println("NOTICE: at first choose directory, press d to see directoreis")
			}
			engine.Files = trace.Files(engine.Config, engine.SelectedDir)
			engine.Page = FILES
			engine.Show()

		case (char > 47 && char < 58) || (char >= 65 && char <= 90):

			//fmt.Println("navigation 0-9 a-f has been pressed", ptr)
			selectedNum := int(SelectionKey[char])
			switch engine.Page {
			case DIRECTORIES:
				if selectedNum < len(engine.Dirs) {
					engine.PrePage = engine.Page
					engine.SelectedDir = engine.Dirs[selectedNum]
					engine.Files = trace.Files(engine.Config, engine.SelectedDir)
					engine.Page = FILES
					engine.Show()
				}
			case FILES:
				if selectedNum < len(engine.Files) {
					engine.PrePage = engine.Page
					engine.SelectedFile = engine.Files[selectedNum]
					engine.Content = loaders.Open(engine.Config, engine.SelectedDir, engine.SelectedFile)
					fmt.Printf("%v", engine.Content)

					engine.Page = CONTENT
					engine.Show()
				}
			case ENVIRONMENTS:
				if selectedNum < len(engine.EnvArr) {
					engine.Page = engine.PrePage
					engine.SelectedEnv = engine.EnvArr[selectedNum]
					fmt.Printf("Environment has been updated: %v", engine.SelectedEnv)
					engine.Show()
				}
			}

		case char == 101: //e
			engine.PrePage = engine.Page
			engine.Page = ENVIRONMENTS
			engine.Show()

		case char == 104: //h
			fmt.Printf("\nHELP\n1.Directories: d\n2.Files: f\n3.Environments: e\n" +
				"4.Request preview: p>\n5.Send request: <Space>\n6.Exit: <Esc>\n" +
				"7.Upper page: x\n8.Send All Reqeusts: s\n9.Print ...: <Enter> ")
		}

		// fmt.Printf("You pressed:%[1]T %[1]q   %[1]v\r\n", char)
		//fmt.Println(SelectionKey[char], engine.PrePage, engine.SelectedEnv)
	}

}
