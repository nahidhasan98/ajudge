package toph

import (
	"io"
	"net/url"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/nahidhasan98/ajudge/errorhandling"
	"github.com/nahidhasan98/ajudge/model"
)

//ProbDes function for grabbing problem description
func ProbDes(pNum string) string {
	defer errorhandling.Recovery() //for panic() error Recovery

	//resetting previous value
	model.PTitle, model.PTimeLimit, model.PMemoryLimit, model.PSourceLimit, model.PDesSrcVJ, model.POrigin = "", "", "", "", "", ""

	//defining a variable for returning data
	var TophDes string

	apiURL := "https://toph.co/p/" + pNum
	response := GETRequest(apiURL)
	defer response.Body.Close()
	document, err := goquery.NewDocumentFromReader(response.Body)
	errorhandling.Check(err)

	model.PTitle = document.Find("span[class='caption']").Find("h1").Text()

	if model.PTitle != "" { //if desired problem exist
		timeMemoryMixed := document.Find("span[class='limits']").Text()

		need := ","
		index := strings.Index(timeMemoryMixed, need)

		if index == -1 {
			model.PTimeLimit = "-"
			model.PMemoryLimit = "-"
		} else {
			model.PTimeLimit = timeMemoryMixed[0:index]
			model.PMemoryLimit = timeMemoryMixed[index+1:]
			model.PSourceLimit = "32 kB" //default, maximum problem has this limit
		}

		TophDes, _ = document.Find("div[class='artifact']").Html()

		//removing extra text from problem caption
		TophDes = removeCaption(TophDes, `<span class="caption">`, `</span>`, 7) //7 character in </span>
		TophDes = removeCaption(TophDes, `<div><span class="text-muted">`, `</div>`, 6)
		TophDes = removeCaption(TophDes, `<div><span title=`, `</div>`, 6)
		//got Title,TimeLimit,MemoryLimit,Description

		//TophDes = beautifyToph(TophDes)
		TophDes = convertKatex(TophDes)
	}
	//getting image manually (if present in prob statement)
	TophDes = getImage(TophDes)

	return TophDes
}

//function used above by this particular file
func removeCaption(body, startPoint, endPoint string, point int) string {
	index1 := strings.Index(body, startPoint)
	index2 := strings.Index(body, endPoint)

	if index1 != -1 {
		part1 := body[0:index1]
		part2 := body[index2+point:]

		body = part1 + part2
	}
	//fmt.Println(body)
	return body
}
func beautifyToph(body string) string {
	for k := 1; k <= 100; k++ { //finding <code> part, assuming at most about 100 times present
		need1 := "<code>"
		need2 := "</code>"
		index1 := strings.Index(body, need1)
		index2 := strings.Index(body, need2)

		if index1 != -1 {
			part1 := body[0:(index1 + 6)]
			middle := body[index1+6 : index2]
			part2 := body[index2:]

			middle = convertToSpecial(middle) //code part will be customised
			body = part1 + middle + part2

			body = strings.Replace(body, "<code>", "<strong>", 1)   //code part will be bold and special font
			body = strings.Replace(body, "</code>", "</strong>", 1) //1 for 1(first element) time replace
			//fmt.Println(body)
		} else {
			break
		}
	}

	return body
}
func convertToSpecial(middle string) string {
	middle = strings.ReplaceAll(middle, "$", "")
	middle = strings.ReplaceAll(middle, `\texttt`, "")
	middle = strings.ReplaceAll(middle, `\times`, "*")
	middle = strings.ReplaceAll(middle, "{", "")
	middle = strings.ReplaceAll(middle, "}", "")
	middle = strings.ReplaceAll(middle, "\\lt", "<")
	middle = strings.ReplaceAll(middle, "\\leq", "<=")
	middle = strings.ReplaceAll(middle, "\\le", "<=")
	middle = strings.ReplaceAll(middle, "\\gt", ">")
	middle = strings.ReplaceAll(middle, "\\geq", ">=")
	middle = strings.ReplaceAll(middle, "\\ge", ">=")
	middle = strings.ReplaceAll(middle, "\\;", "")
	middle = strings.ReplaceAll(middle, "\\#", "#")

	return middle
}
func convertKatex(body string) string {
	body = strings.ReplaceAll(body, "<code>\\", `<code>$\`)        //some problem missing $ sign,so putting that sign
	body = strings.ReplaceAll(body, "<code>$", `\(`)               //converting Katex format($EQUATION$ or \(EQUATION\))
	body = strings.ReplaceAll(body, "$</code>", `\)`)              //<code>$EQUATION$</code> --> <code>\(EQUATION\)</code>
	body = strings.ReplaceAll(body, `\(</code>`, `<code>$</code>`) //this is real $ sign, not for Katex (we changed previous in line by assuming its for Katex)
	body = strings.ReplaceAll(body, `<code>\)`, `<code>$</code>`)  //<code>\)</code> --> <code>$</code>

	//adding mathjax script for math equations/Katex
	body = body + `<script src="https://polyfill.io/v3/polyfill.min.js?features=es6"></script>
	<script id="MathJax-script" async src="https://cdn.jsdelivr.net/npm/mathjax@3.0.1/es5/tex-mml-chtml.js"></script>`

	return body
}

//getImage function for grabbing image of problem statement from OJ
func getImage(body string) string {
	defer errorhandling.Recovery() //for panic() error Recovery

	need1 := `data-src="//uploads.drafts.toph.co/drafts-images/`
	need2 := `data-src="//uploads.toph.co/arena-images/`
	need3 := `data-src="https://uploads.toph.co/`

	for k := 0; k < len(body)-49; k++ {
		subStr1 := body[k : k+49] //`data-src="//uploads.drafts.toph.co/drafts-images/` it has 49 character
		subStr2 := body[k : k+41] //`data-src="//uploads.toph.co/arena-images/` it has 41 character
		subStr3 := body[k : k+34] //`data-src="https://uploads.toph.co/` it has 34 character

		if subStr1 == need1 {
			var part1, part2, middle string

			part1 = body[0:k]
			middle = body[k : k+151] //link is like: data-src="//uploads.drafts.toph.co/drafts-images/5829a5ef04469e48383b360f-1575827552649040444-8344245478832197509-e25c8b977b5c931df7b07dacc6c1f209.jpg"
			part2 = body[k+151:]

			imageLink := middle[10 : len(middle)-1] //removing first 10 characters like: data-src="
			imageLink = "https:" + imageLink

			newImageSrc := `src="../` + fileOperation(imageLink, "-", "") + `" class="TophImage"`

			//replacing image source to our local source
			body = part1 + newImageSrc + part2
		} else if subStr2 == need2 {
			var part1, part2, middle string

			part1 = body[0:k]
			middle = body[k : k+143] //link is like: data-src="//uploads.drafts.toph.co/arena-images/5829a5ef04469e48383b360f-1575827552649040444-8344245478832197509-e25c8b977b5c931df7b07dacc6c1f209.jpg"
			part2 = body[k+143:]

			imageLink := middle[10 : len(middle)-1]
			imageLink = "https:" + imageLink

			newImageSrc := `src="../` + fileOperation(imageLink, "-", "") + `" class="TophImage"`

			body = part1 + newImageSrc + part2
		} else if subStr3 == need3 {
			var part1, part2, middle string

			part1 = body[0:k]
			middle = body[k : k+251] //link is like: data-src="https://uploads.toph.co/kusaC1dxlkPI0laIlmN4R0v78Kq3A-wBY2F3EyaivAY/resize:fit:768:0:0/dpr:2/czM6Ly90b3BoLXBsYXRmb3JtLXVwbG9hZHMvaW1hZ2VzLzE1ODIwOTU0MjcwMjgzMDE3NzEtMjI3NjkyNDA1OTAxODI4MzI5NC1iN2M4NDllZDNhNDJhZThhN2M1ZGI1NjExYjllZWQ4YS5qcGc"
			part2 = body[k+251:]

			imageLink := middle[10 : len(middle)-1]

			newImageSrc := `src="../` + fileOperation(imageLink, "/", ".png") + `" class="TophImage"`

			body = part1 + newImageSrc + part2
		}
	}
	return body
}

var fileName string

func fileOperation(imageLink, segmentSeparator, fileExtension string) string {
	// Build fileName from fullPath
	fileURL, err := url.Parse(imageLink)
	errorhandling.Check(err)
	path := fileURL.Path
	nameSegments := strings.Split(path, segmentSeparator)
	fileName = nameSegments[len(nameSegments)-1]

	//creating directory if does not exist
	nameDir := "assets/temp"
	_, err = os.Stat(nameDir)
	if os.IsNotExist(err) {
		//fmt.Println("Dir does not exist")
		os.Mkdir("assets/temp", 0755) //(owner:7=rwx group:5=r-x other:5=r-x) This means that the directory has the default permissions -rwxr-xr-x (represented in octal notation as 0755).
	}

	fullFileName := "assets/temp/" + fileName + fileExtension

	// Create blank file
	var file *os.File
	_, err = os.Stat(fullFileName)
	if os.IsNotExist(err) {
		//fmt.Println("file does not exist")
		file, err = os.Create(fullFileName)
		errorhandling.Check(err)
	}

	// Put content on file
	// client2 := http.Client{
	// 	CheckRedirect: func(r *http.Request, via []*http.Request) error {
	// 		r.URL.Opaque = r.URL.Path
	// 		return nil
	// 	},
	// }
	resp, err := model.Client.Get(imageLink)
	errorhandling.Check(err)
	defer resp.Body.Close()
	io.Copy(file, resp.Body) //copying response image to my just created empty file
	defer file.Close()
	//fmt.Println("Just Downloaded a file %s with size %d", fileName, size)

	return fullFileName
}
