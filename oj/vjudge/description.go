package vjudge

import (
	"encoding/json"
	"io"
	"net/url"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gomarkdown/markdown"
	"github.com/nahidhasan98/ajudge/errorhandling"
	"github.com/nahidhasan98/ajudge/model"
)

//ProbDes function for grabbing problem description
func ProbDes(OJ, pNum string) (string, bool, int) {
	defer errorhandling.Recovery() //for panic() error Recovery

	//resetting previous value
	model.PTitle, model.PTimeLimit, model.PMemoryLimit, model.PSourceLimit, model.PDesSrcVJ, model.POrigin = "", "", "", "", "", ""

	//defining a variable for returning data
	var VJDes = ""
	var allowSubmit bool
	var status int

	//finding problem title
	if OJ == "计蒜客" || OJ == "黑暗爆炸" {
		OJ = url.QueryEscape(OJ)
	}

	apiURL := "https://vjudge.net/problem/" + OJ + "-" + pNum
	response := GETRequest(apiURL)
	defer response.Body.Close()

	document, err := goquery.NewDocumentFromReader(response.Body)
	errorhandling.Check(err)

	model.PTitle = document.Find("div[id='prob-title']").Find("h2").Text()
	//got prob title

	if model.PTitle != "" { //if desired problem exist //getting time & memory Limit
		resProperties := document.Find("textarea[name='dataJson']").Text()

		type properties struct {
			Title   string `json:"title"`
			Content string `json:"content"`
		}
		type limitation struct {
			Properties  []properties `json:"properties"`
			AllowSubmit bool         `json:"allowSubmit"`
			Status      int          `json:"status"`
		}
		var limit limitation
		json.Unmarshal([]byte(resProperties), &limit)

		allowSubmit = limit.AllowSubmit
		status = limit.Status

		//some problem has 1(time/memory/source) limit, some has 2(time,memory) limits, some has 3(time,memory,source) limits
		for k := 0; k <= model.Min(2, len(limit.Properties)-1); k++ { //so we'll look for max 3 times
			text := limit.Properties[k].Title
			need1 := "limit"
			need2 := "Limit"
			index1 := strings.Index(text, need1)
			index2 := strings.Index(text, need2)

			if index1 != -1 || index2 != -1 {
				if text == "Time limit" || text == "Case time limit" {
					model.PTimeLimit = limit.Properties[k].Content
				} else if text == "Memory limit" {
					model.PMemoryLimit = limit.Properties[k].Content
				} else if text == "Code length Limit" {
					model.PSourceLimit = limit.Properties[k].Content
				}
			}
		}
		//got time, memory & source limit

		//getting problem description
		model.PDesSrcVJ, _ = document.Find("iframe").Attr("src")

		apiURL := "https://vjudge.net" + model.PDesSrcVJ //prob description link
		response := GETRequest(apiURL)
		defer response.Body.Close()

		document, err := goquery.NewDocumentFromReader(response.Body)
		errorhandling.Check(err)

		textArea := document.Find("textarea").Text()

		//VJ returns prob description through a (complex!) structure
		type Inner2 struct {
			Format  string `json:"format"`
			Content string `json:"content"`
		}
		type Inner struct {
			Title string `json:"title"`
			Value Inner2 `json:"value"`
		}
		type Res struct {
			Sections []Inner `json:"sections"`
		}
		var res Res
		json.Unmarshal([]byte(textArea), &res)
		//got problem description in a structured format

		//now taking it to a string
		for i := 0; i < len(res.Sections); i++ {
			//Eliminating Default CSS on Example-Input-Output
			//styleBody := res.Sections[i].Value.Content
			//content := removeStyle(styleBody)
			//content := styleBody

			// mp := map[string]interface{}{
			// 	"Title":   template.HTML(res.Sections[i].Title),
			// 	"Content": template.HTML(res.Sections[i].Value.Content),
			// }
			// VJProblem = append(VJProblem, mp)

			VJDes += `<h3>` + res.Sections[i].Title + `</h3>`
			VJDes += res.Sections[i].Value.Content
		}
	}

	//getting image manually (if present in prob statement)
	VJDes = GetImage(VJDes)

	if OJ == "CodeChef" { //need to format
		VJDes = VJCodeChef(VJDes)
	} else if OJ == "LibreOJ" { //need to format
		VJDes = VJLibreOJ(VJDes)
	} else if OJ == "LightOJ" { //need to format
		VJDes = VJCodeChef(VJDes)
	} else if OJ == "CodeForces" || OJ == "Gym" { //need to format
		VJDes = addLatexFormat(VJDes)
	}

	//fmt.Println(200, VJDes)
	return VJDes, allowSubmit, status
}

//GetImage function for grabbing image of problem statement from OJ
func GetImage(body string) string {
	defer errorhandling.Recovery() //for panic() error Recovery

	body = strings.ReplaceAll(body, "CDN_BASE_URL", "https://vj.ppsucxtt.cn")

	// //for pdf file
	// need1 := `<iframe src="https://vj.z180.cn/`
	// for k := 0; k < len(body)-32; k++ {
	// 	subStr := body[k : k+32] //`<iframe src="https://vj.z180.cn/` it has 32 character

	// 	if subStr == need1 {
	// 		var part1, part2, middle string

	// 		part1 = body[0 : k+8]
	// 		middle = body[k+8 : k+88] //link is like: src="https://vj.z180.cn/c00b3cf31e6f17859d4c4d0d5cd94757?v=1602540878#view=FitH"
	// 		part2 = body[k+88:]

	// 		pdfLink := middle[5:56] //removing first 5 characters like: src="
	// 		//fmt.Println(pdfLink)
	// 		newPdfSrc := `src="/` + fileOperation(pdfLink, "/", "") + `#view=FitH"` //#view=FitH for fitting 100% horizontally

	// 		//replacing pdf source to our local source
	// 		body = part1 + newPdfSrc + part2
	// 	}
	// }
	// need2 := `<iframe src="https://vj.ppsucxtt.cn/`
	// for k := 0; k < len(body)-36; k++ {
	// 	subStr := body[k : k+36] //`<iframe src="https://vj.ppsucxtt.cn/` it has 36 character

	// 	if subStr == need2 {
	// 		var part1, part2, middle string

	// 		part1 = body[0 : k+8]
	// 		middle = body[k+8 : k+92] //link is like: src="https://vj.ppsucxtt.cn/95ef86acbc0a8f031db5751693d6725f?v=1633807853#view=FitH"
	// 		part2 = body[k+92:]

	// 		pdfLink := middle[5:60] //removing first 5 characters like: src="
	// 		//fmt.Println(pdfLink)
	// 		newPdfSrc := `src="/` + fileOperation(pdfLink, "/", "") + `#view=FitH"` //#view=FitH for fitting 100% horizontally

	// 		//replacing pdf source to our local source
	// 		body = part1 + newPdfSrc + part2
	// 	}
	// }
	// //for image file
	// need3 := `src="https://vj.z180.cn/`
	// need4 := `SRC="https://vj.z180.cn/`
	// for k := 0; k < len(body)-24; k++ {
	// 	subStr := body[k : k+24] //`src="https://vj.z180.cn/` it has 24 character

	// 	if subStr == need3 || subStr == need4 {
	// 		var part1, part2, middle string

	// 		part1 = body[0:k]
	// 		middle = body[k : k+70] //link is like: src="https://vj.z180.cn/5fb7165c882f9f4835f0623e8c580bda?v=1600300763"
	// 		part2 = body[k+70:]

	// 		imageLink := middle[5:56] //removing first 5 characters like: src="
	// 		//fmt.Println(imageLink)
	// 		newImageSrc := `src="../` + fileOperation(imageLink, "/", "") + `"`

	// 		//replacing image source to our local source
	// 		body = part1 + newImageSrc + part2
	// 	}
	// }

	// //changing https://vj.z180.cn/5b29127c46fedc0e54ba9d20875c6899?v=1603194575 to /assets/temp/images/5b29127c46fedc0e54ba9d20875c6899)
	// need5 := `href="https://vj.z180.cn/`
	// need6 := `href='https://vj.z180.cn/`
	// for k := 0; k < len(body)-25; k++ {
	// 	subStr := body[k : k+25] //`href="https://vj.z180.cn/` it has 25 character

	// 	if subStr == need5 || subStr == need6 {
	// 		var part1, part2, middle string

	// 		part1 = body[0:k]
	// 		middle = body[k : k+71] //link is like: href="https://vj.z180.cn/5fb7165c882f9f4835f0623e8c580bda?v=1600300763"
	// 		part2 = body[k+71:]

	// 		imageLink := middle[6:57] //removing first 5 characters like: href="
	// 		//fmt.Println(imageLink)
	// 		newImageSrc := `href="../` + fileOperation(imageLink, "/", "") + `"`

	// 		body = part1 + newImageSrc + part2
	// 	}
	// }

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

	if resp != nil { //if response is nil (due to bad url/server side error/403), then return the blank file (fullFileName)
		errorhandling.Check(err)
		defer resp.Body.Close()
		io.Copy(file, resp.Body) //copying response image to my just created empty file
		defer file.Close()
		//fmt.Println("Just Downloaded a file %s with size %d", fileName, size)
	}

	return fullFileName
}

//VJCodeChef for CodeChef OJ
func VJCodeChef(body string) string {
	body = removeMathJaxConfig(body)
	body = convertMDtoHTML(body) //CodeChef descriptions are in .MD format

	body = strings.ReplaceAll(body, `&amp;`, `&`)
	body = strings.ReplaceAll(body, `&lt;`, `<`)
	body = strings.ReplaceAll(body, `&gt;`, `>`)
	body = strings.ReplaceAll(body, `&gt;`, `>`)

	return body
}
func removeMathJaxConfig(body string) string {
	need1 := `<script type="text/x-mathjax-config">`
	need2 := `</script>`

	index1 := strings.Index(body, need1)
	index2 := strings.Index(body, need2)
	if index1 != -1 {
		part1 := body[:index1]
		part2 := body[index2+9:]

		body = part1 + part2
	}
	return body
}
func convertMDtoHTML(body string) string {
	md := []byte(body)
	output := markdown.ToHTML(md, nil, nil)
	body = string(output)

	return body
}

//VJLibreOJ for LibreOJ
func VJLibreOJ(body string) string {
	body = strings.ReplaceAll(body, "####", `<h4>`)
	body = strings.ReplaceAll(body, "```plain", `</h4><pre><code>`)
	body = strings.ReplaceAll(body, "```", `</code></pre>`)
	body = strings.ReplaceAll(body, "* `", `<li><code>`)
	body = strings.ReplaceAll(body, "`，", `</code></li>, `)

	return body
}

func addLatexFormat(body string) string {
	// got from codeforces
	body += `<script type="text/x-mathjax-config">
    MathJax.Hub.Config({
      tex2jax: {inlineMath: [['$$$','$$$']], displayMath: [['$$$$$$','$$$$$$']]}
    });
   
    </script>
    <script type="text/javascript" async src="https://mathjax.codeforces.org/MathJax.js?config=TeX-AMS_HTML-full"></script>`

	return body
}
