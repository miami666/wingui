package main

import (
	"github.com/hesahesa/pwdbro"
	"github.com/tadvi/winc"
	"math"
	"strconv"
	"sync"
	"time"
	"unicode"
)



const (
	colorWhite      = 0xFFFFFF
	colorBlack      = 0x000000
	colorDodgerBlue = 0x1E90FF
)

func Brute(charset []rune, minLen, maxLen, buffer int) (combos <-chan string, closer func()) {

	results := make(chan string, buffer)
	done := make(chan struct{})
	charlen := len(charset)
	once := new(sync.Once)

	closer = func() {
		once.Do(func() {
			close(done)
		})
	}

	if minLen == 0 {
		minLen = 1
	}

	go func() {
		defer close(results)
		defer closer()
		for k := minLen; k <= maxLen; k++ {
			carry := 0
			indices := make([]int, k)
			for {
				select {
				case <-done:
					return
				default:
					out := ""
					for j := 0; j < k; j++ {
						out += string(charset[indices[j]])
					}
					results <- out
				}
				carry = 1
				for i := k - 1; i >= 0; i-- {
					if carry == 0 {
						break
					}

					indices[i] += carry
					carry = 0

					if indices[i] == charlen {
						carry = 1
						indices[i] = 0
					}
				}
				if carry == 1 {
					break
				}
			}
		}

	}()
	return results, closer
}

func possibleNumberCharacters(password string) (int, int) {
	var uppercasePresent bool
	var lowercasePresent bool
	var numberPresent bool
	var pnc = 0
	var specialCharPresent bool
	const minPassLength = 8
	const maxPassLength = 64
	var passLen int

	for _, ch := range password {
		switch {
		case unicode.IsNumber(ch):
			numberPresent = true
			passLen++

		case unicode.IsUpper(ch):
			uppercasePresent = true
			passLen++
		case unicode.IsLower(ch):
			lowercasePresent = true
			passLen++
		case unicode.IsPunct(ch) || unicode.IsSymbol(ch):
			specialCharPresent = true
			passLen++
		case ch == ' ':
			passLen++
		}
	}
	if numberPresent {
		pnc += 10
	}
	if uppercasePresent {
		pnc += 26
	}
	if lowercasePresent {
		pnc += 26
	}
	if specialCharPresent {
		pnc += 32
	}

	return pnc, passLen
}

func main() {

	mainWindow := winc.NewForm(nil)
	//window:=winc.CreateWindow("test",mainWindow,0,0)
	//b:= winc.NewSolidColorBrush(0xff0000)
    //pen:=winc.NewPen(0, 1, b  )
    //win.WglGetCurrentDC()



	f := winc.NewFont("Verdana", 10, 0)



	mainWindow.SetSize(240, 300) // (width, height)



	label := winc.NewLabel(mainWindow)

	label.SetText("Tragen Sie hier das Passwort ein:")
	label.SetPos(10, 10)
	label.SetSize(200, 30)



	edt := winc.NewEdit(mainWindow)
	/*	right := winc.NewEdit(mainWindow)
		dock := winc.NewSimpleDock(mainWindow)
		dock.Dock(right, winc.Fill)
		dock.Dock(edt,winc.Top)*/

	edt.SetFont(f)

	edt.SetPos(10, 30)

	//cb:=winc.NewSolidColorBrush(25500)

	// Most Controls have default size unless SetSize is called.
	edt.SetText("")

	btn := winc.NewPushButton(mainWindow)
	btn.SetText("benötigte Dauer berechnen")
	btn.SetPos(10, 80)   // (x, y)
	btn.SetSize(200, 40) // (width, height)
	btn.OnClick().Bind(func(e *winc.Event) {

		pnc, passLen := possibleNumberCharacters(edt.Text())
		//fmt.Println(password, " ", pnc)
		pc := math.Pow(float64(pnc), float64(passLen))
		println(int64(pc))
		timeRequired := float64(pc) / float64(100000000000)
		str := strconv.FormatFloat(timeRequired, 'f', 2, 64)
		if timeRequired <= 60 {
			//winc.MsgBoxOkCancel(mainWindow,str,str)
			winc.Printf(mainWindow, "Das Knacken dieses Passwortes würde mit einem heutigen handeslsüblichen Computer ungefähr "+str+" Sekunden benötigen.")
		}
		if timeRequired > 60 {
			timeRequiredMinutes := timeRequired / 60
			strMin := strconv.FormatInt(int64(timeRequiredMinutes), 10)
			winc.Printf(mainWindow, "Das Knacken dieses Passwortes würde mit einem heutigen handeslsüblichen Computer ungefähr "+strMin+" Minuten benötigen.")
		}

		println(str)
		/*if edt.Visible() {
			edt.Hide()


		} else {
			edt.Show()
		}*/
	})
	/*bf, err := gobf.New(
		gobf.WithNumber(false),
		gobf.WithUpper(false),
		gobf.WithLower(true),
		gobf.WithSize(len(edt.Text())),
		gobf.WithConcrencyLimit(1000),
	)
	if err != nil {
		log.Fatal(err)
	}*/

	/*	btn2 := winc.NewPushButton(mainWindow)
		btn2.SetText("Brute Force Simulation")
		btn2.SetPos(10, 130)
		btn2.SetSize(200, 40)
		btn2.OnClick().Bind(func(e *winc.Event) {
	        log.Println(len(edt.Text()))
			log.Println("start to search pattern: "+edt.Text())
			start := time.Now()
			p:=edt.Text()

			err = bf.Do(context.Background(), func(pattern string) {


				if pattern == p {

					elapsed := time.Since(start)
					log.Println(pattern + "gefunden in "+elapsed.String()+" Sekunden")

				}
			})
			if err != nil {
				log.Fatal(err)
			}

		})*/
	btn3 := winc.NewPushButton(mainWindow)
	btn3.SetText("Brute Force")
	btn3.SetPos(10, 130)
	btn3.SetSize(200, 40)

	btn3.OnClick().Bind(func(e *winc.Event) {
		start := time.Now()
		pr := winc.NewProgressBar(mainWindow)
		pr.Show()
		pr.SetPos(10, 200)

		for i := 0; i < 100; i++ {

			time.Sleep(10 * time.Millisecond)
			pr.SetValue(i)
		}

		//measuring the durration of the for loop
		/*		for index := 0; index < 10; index++ {
				time.Sleep(500 * time.Millisecond)
			}*/

		characterSer := []rune("abcdefghijklmnopqrstuvwxyz")
		var comb string
		minLen := len(edt.Text())
		maxLen := 4
		b, _ := Brute(characterSer, minLen, maxLen, 4)
		for combination := range b {

			comb += "\t" + combination
			if combination == edt.Text() {
				elapsed := time.Since(start)
				pr.Hide()
				winc.Printf(mainWindow, combination+"\t"+elapsed.String()+"\n"+comb)
				break

			}
		}

	})
	btn4 := winc.NewPushButton(mainWindow)
	btn4.SetText("Pwned Passwort Check")
	btn4.SetPos(10, 180)
	btn4.SetSize(200, 40)
	btn4.OnClick().Bind(func(e *winc.Event) {
		pwdbro := pwdbro.NewDefaultPwdBro()
		status, _ := pwdbro.RunParallelChecks(edt.Text())

		for _, resp := range status {
			/*			println("=======")
						println(resp.Safe)
						println(resp.Method)
						println(resp.Message)
						println(resp.Error)*/
			method := resp.Method

			if method == "pwnedpasswords API" {
				if resp.Safe {
					winc.Printf(mainWindow, "noch nicht geknackt")
				} else {
					winc.Printf(mainWindow, "You have been pwned!")

				}

			}
		}

	})

	mainWindow.Center()
	mainWindow.Show()
	mainWindow.OnClose().Bind(wndOnClose)

	winc.RunMainLoop() // Must call to start event loop.
}

func wndOnClose(arg *winc.Event) {
	winc.Exit()
}
