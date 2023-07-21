package main

import (
	"fmt"
	"io/ioutil"
	"net/url"
	"strconv"
	"strings"
	"sync"

	"github.com/aiomonitors/godiscord"

	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/fatih/color"
	http "github.com/saucesteals/fhttp"
	"github.com/saucesteals/mimic"
)

var (
	latestVersion = mimic.MustGetLatestVersion(mimic.PlatformWindows)
)

var previouslySent sync.Map

type Course struct {
	Subject   string
	Course    string
	Section   string
	Professor string
	Seats     int
	Times     string
	Url       string
}

func main() {
	m, _ := mimic.Chromium(mimic.BrandChrome, latestVersion)
	client := &http.Client{Transport: m.ConfigureTransport(&http.Transport{})}
	go getSeats("https://app.testudo.umd.edu/soc/202301/CMSC/CMSC132", client, "https://discord.com/api/webhooks/1050460562628292679/2n47GlZTG3g4OQit5fm_VBzLRDFRqlAHNo0vSGF_QSP5WYhgaOmRZC33rWmmXEs6pWQD")
	go getSeats("https://app.testudo.umd.edu/soc/202301/MATH/MATH240", client, "https://discord.com/api/webhooks/1050460668320546856/Fo9P3u3tb6WsVYfLyWm9uvrz5RaBut4MsYgNZZCxd7ERvak89R0lzQbNC1AjNOCCFYyQ")
	go getSeats("https://app.testudo.umd.edu/soc/202301/AASP/AASP298L", client, "https://discord.com/api/webhooks/1050460762717552720/6LNp1OOTx9mVlI2yJhcFU18FQULAF15owKMmimVoZ7JWDehxxZ8e4vDES487CTZwwQZE")
	go getSeats("https://app.testudo.umd.edu/soc/202301/MATH/MATH141", client, "https://discord.com/api/webhooks/1050464761860608091/3yC7VVTcPrStEICPxE6-T6einVP3HRJMOTz1DIFP3OvLgBJ34Jz-QYRMlu0vrDUqon4z")
	go getSeats("https://app.testudo.umd.edu/soc/202301/COMM/COMM107", client, "https://discord.com/api/webhooks/1050465096574451832/2vRfaUuUJmwUjF4P9FtZ2eCXw56iRIaVbBCbZs-CC9unAsaEPeDt_YJ6NBvgeX1Qc7kH")

	select {}
}

func getSeats(url string, client *http.Client, webhook string) {
	for {

		m, _ := mimic.Chromium(mimic.BrandChrome, latestVersion)
		req, _ := http.NewRequest("GET", url, nil)
		req.Header = http.Header{
			"cache-control":      {"no-cache"},
			"sec-ch-ua":          {`"Not?A_Brand";v="8", "Chromium";v="108", "Google Chrome";v="108"`},
			"pragma":             {"no-cache"},
			"if-modified-since":  {"Sat, 01 Dec 2001 00:00:00 GMT"},
			"sec-ch-ua-mobile":   {"?0"},
			"user-agent":         {"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 Safari/537.36"},
			"sec-ch-ua-platform": {`"macOS"`},
			"accept":             {"*/*"},
			"sec-fetch-site":     {"same-origin"},
			"sec-fetch-mode":     {"cors"},
			"sec-fetch-dest":     {"empty"},
			"referer":            {"https://app.testudo.umd.edu/"},
			"accept-encoding":    {"gzip, deflate, br"},
			"accept-language":    {"en-US,en;q=0.9"},
			"dnt":                {"1"},
			http.HeaderOrderKey: {
				"cache-control",
				"sec-ch-ua",
				"pragma",
				"if-modified-since",
				"sec-ch-ua-mobile",
				"user-agent",
				"sec-ch-ua-platform",
				"accept",
				"sec-fetch-site",
				"sec-fetch-mode",
				"sec-fetch-dest",
				"referer",
				"accept-encoding",
				"accept-language",
				"cookie",
				"dnt",
			},
			http.PHeaderOrderKey: m.PseudoHeaderOrder(),
		}
		// fmt.Println("Getting data")
		resp, _ := client.Do(req)
		courseName := strings.Split(url, "/")[len(strings.Split(url, "/"))-1]
		logAction(courseName, "Got data with "+fmt.Sprint(resp.StatusCode)+" status code")

		// fmt.Println("Got response:", resp.StatusCode)

		bd, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println(err)
			continue
		}
		bodystr := string(bd)
		//parse with goquery
		document, err := goquery.NewDocumentFromReader(strings.NewReader(bodystr))
		if err != nil {
			fmt.Println(err)
			continue
		}
		//get table rows
		rows := document.Find(".section.delivery-f2f")

		//get course name, section number, and instructor, and open seats
		rows.Each(func(i int, s *goquery.Selection) {
			//get course name (last element in strings.Split(url, "/"))
			//get section number
			sectionNumber := strings.TrimSpace(s.Find("span.section-id").Text())
			//get instructor
			instructor := s.Find("span.section-instructor").Text()
			//get open seats
			openSeats := strings.TrimSpace(s.Find("span.open-seats").Text())
			//split by :[1] and parse to int
			split := strings.Split(strings.TrimSpace(strings.Split(openSeats, ":")[1]), ",")[0]
			seatsInt := 0
			// fmt.Println(split)
			seatsIntval, err := strconv.Atoi(split)
			if err == nil {
				seatsInt = seatsIntval
			} else {
				// fmt.Println(err)
			}
			times := ""

			//get days
			s.Find("div.section-day-time-group").Find("span.section-days").Each(func(z int, x *goquery.Selection) {
				days := x.Text()
				//get start time with same index
				startTime := s.Find("div.section-day-time-group").Find("span.class-start-time").Eq(z).Text()
				//get end time with same index
				endTime := s.Find("div.section-day-time-group").Find("span.class-end-time").Eq(z).Text()

				times += days + " " + startTime + "-" + endTime + " "

			})

			course := Course{
				Subject:   courseName,
				Course:    courseName,
				Section:   sectionNumber,
				Professor: instructor,
				Seats:     seatsInt,
				Times:     times,
				Url:       url,
			}
			// fmt.Println(course.Seats)
			go pingCourse(course, webhook)
			// fmt.Println(course)

		})
		resp.Body.Close()
		time.Sleep(3 * time.Second)
	}

}

func pingCourse(course Course, webhook string) {
	_, ok := previouslySent.Load(course.Subject + "x|x" + fmt.Sprint(course.Seats) + "x|x" + course.Section)
	if ok {
		return
	}
	//log to console wiht COLOR and fancy formatting
	cyan := color.New(color.FgCyan).SprintFunc()
	blue := color.New(color.FgBlue).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()
	magenta := color.New(color.FgMagenta).SprintFunc()
	output := "---------------------------------------------"
	output += "\n" + red(fmt.Sprint(course.Seats)+" seats") + " " + cyan("open in "+course.Subject)
	output += "\n" + "=============================================	"
	output += "\n" + blue(course.Course) + " | " + green(course.Section) + " | " + yellow(course.Professor)
	output += "\n" + magenta(course.Times)
	output += "\n" + "---------------------------------------------"
	fmt.Println(output)
	previouslySent.Store(course.Subject+"x|x"+fmt.Sprint(course.Seats)+"x|x"+course.Section, true)

	//send discord webhook
	if webhook != "" {
		fmt.Println("Sending webhook to " + webhook)
		embed := godiscord.NewEmbed("UMD Course Monitor", (fmt.Sprint(course.Seats)+" seats")+" "+("open in "+course.Subject), course.Url)
		embed.AddField("Class", (course.Course), true)
		embed.AddField("Section", (course.Section), true)
		embed.AddField("Professor", (course.Professor), true)
		embed.AddField("Times", (course.Times), true)
		embed.SetColor("7BBF1C")
		if course.Seats == 0 {
			embed.SetColor("E41E26")
			embed.SetThumbnail("https://thumbs.dreamstime.com/b/cartoon-image-sad-turtle-has-to-leave-its-polluted-environment-194538165.jpg")
			embed.SetFooter("UMD Course Monitor", "https://cdn.discordapp.com/attachments/1050460744019365908/1050461586042658937/download.jpeg")

		} else {
			embed.SetColor("7BBF1C")
			embed.SetFooter("UMD Course Monitor", "https://cf.ltkcdn.net/reptiles/turtles-and-tortoises/images/std-xs/325001-340x227-turtle-eating-strawberry.jpg")

			embed.SetThumbnail("https://cdn.discordapp.com/attachments/1050460543091232869/1050462885110878338/download_1.jpeg")
		}
		embed.SendToWebhook(webhook)
	}
}

func logAction(coursename string, value string) {
	green := color.New(color.FgGreen).SprintFunc()
	cyan := color.New(color.FgCyan).SprintFunc()
	time := "[" + time.Now().Format("15:04:05") + "] "

	fmt.Println(green(time+"["+coursename+"]") + " " + cyan(value))

}
