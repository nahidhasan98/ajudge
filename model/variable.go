package model

import (
	"html/template"
	"net/http"
	"net/http/cookiejar"
	"time"

	"github.com/gorilla/sessions"
)

//Tpl variable
var Tpl *template.Template

func init() {
	Tpl = template.Must(template.ParseGlob("frontend/*/*"))
}

var (
	//CookieJar variable
	CookieJar, _ = cookiejar.New(nil)

	//Client variable
	Client = &http.Client{
		Jar:     CookieJar,
		Timeout: 6 * time.Second, //limiting time to 3 seconds for http request
	}

	//Store variable
	Store = sessions.NewCookieStore([]byte("mysession"))

	//Info variable for sending Data from backend to frontend
	Info = map[string]interface{}{}

	//LastPage variable
	LastPage = "/"

	//PopUpCause variable
	PopUpCause = ""

	//ErrorType variable
	ErrorType = ""

	//PTitle variable
	PTitle = ""

	//PTimeLimit variable
	PTimeLimit = ""

	//PMemoryLimit variable
	PMemoryLimit = ""

	//PSourceLimit variable
	PSourceLimit = ""

	//PDesSrcVJ variable
	PDesSrcVJ = ""

	//POrigin variable
	POrigin = ""

	//LanguagePack variable
	LanguagePack = make(map[string]string)
)

//ProblemList variable for holding problem Name,number etc.(collected from search result)
type ProblemList struct {
	OJ    string `bson:"OJ"`
	PNum  string `bson:"pNum"`
	PName string `bson:"pName"`
}

//OJSet variable for clarifying total number of OJ
var OJSet = map[string]bool{
	"51Nod":         true,
	"ACdream":       true,
	"Aizu":          true,
	"AtCoder":       true,
	"CodeChef":      true,
	"CodeForces":    true,
	"CSU":           true,
	"DimikOJ":       true,
	"EIJudge":       true,
	"EOlymp":        true,
	"FZU":           true,
	"Gym":           true,
	"HackerRank":    true,
	"HDU":           true,
	"HihoCoder":     true,
	"HIT":           true,
	"HRBUST":        true,
	"HUST":          true,
	"HYSBZ":         true,
	"Kattis":        true,
	"LibreOJ":       true,
	"LightOJ":       true,
	"Minieye":       true,
	"NBUT":          true,
	"OpenJ_Bailian": true,
	"OpenJ_POJ":     true,
	"POJ":           true,
	"SCU":           true,
	"SGU":           true,
	"SPOJ":          true,
	"TopCoder":      true,
	"Toph":          true,
	"UESTC":         true,
	"UESTC_old":     true,
	"UniversalOJ":   true,
	"URAL":          true,
	"URI":           true,
	"UVA":           true,
	"UVALive":       true,
	"Z_trening":     true,
	"ZOJ":           true,
	"计蒜客":           true,
	"黑暗爆炸":          true,
}

//TerminalVerdict variable
var TerminalVerdict = map[string][]string{
	"51Nod":         {"Accepted", "Wrong Answer", "Time Limit Exceed", "Memory Limit Exceed", "Runtime Error", "Compile Error"},
	"ACdream":       {"Accepted", "Presentation Error", "Wrong Answer", "Time Limit Exceeded", "Memory Limit Exceeded", "Output Limit Exceeded", "Runtime Error", "Compilation Error"},
	"Aizu":          {"Accepted", "Presentation Error", "Wrong Answer", "Time Limit Exceeded", "Memory Limit Exceeded", "Output Limit Exceeded", "Runtime Error", "Compile Error"},
	"AtCoder":       {"Accepted", "Wrong Answer", "Time Limit Exceeded", "Output Limit Exceeded", "Runtime Error", "Compilation Error"},
	"CodeChef":      {"Accepted", "Wrong Answer", "Time Limit Exceed", "Runtime Error", "Compile Error", "Code Length Exceeded"},
	"CodeForces":    {"Accepted", "Happy New Year!", "Wrong answer", "Time limit exceeded", "Memory limit exceeded", "Runtime error", "Compilation error", "Source Code Error", "Language Rejected"},
	"CSU":           {"Accepted", "Presentation Error", "Wrong Answer", "Memory Limit Exceed", "Output Limit Exceed"},
	"DimikOJ":       {"Accepted", "Wrong Answer", "Compilation Error", "Time Limit Exceeded", "Runtime Error"},
	"EIJudge":       {"Accepted", "Presentation error", "Wrong answer", "Time limit", "Memory limit", "Runtime error", "Compilation error"},
	"EOlymp":        {"Accepted", "Wrong Answer", "Time Limit Exceeded", "Memory Limit Exceeded", "Runtime Error", "Compilation Error"},
	"FZU":           {"Accepted", "Presentation Error", "Wrong Answer", "Time Limit Exceed", "Memory Limit Exceed", "Output Limit Exceed", "Runtime Error", "Compile Error", "Restrict Function Call"},
	"Gym":           {"Accepted", "Happy New Year!", "Wrong answer", "Time limit exceeded", "Memory limit exceeded", "Runtime error", "Compilation error", "Source Code Error", "Language Rejected"},
	"HackerRank":    {"Accepted", "Wrong Answer", "Terminated due to timeout", "Runtime Error", "Segmentation Fault", "Compilation error", "Source Code Error"},
	"HDU":           {"Accepted", "Presentation Error", "Wrong Answer", "Time Limit Exceeded", "Memory Limit Exceeded", "Output Limit Exceeded", "Runtime Error", "Compilation Error"},
	"HihoCoder":     {"Accepted", "Presentation Error", "Wrong Answer", "Time Limit Exceeded", "Memory Limit Exceeded", "Output Limit Exceeded", "Runtime Error", "Compile Error"},
	"HIT":           {"Accepted", "Presentation Error", "Wrong Answer", "Time Limit Exceed", "Memory Limit Exceed", "Output Limit Exceed", "Runtime Error", "Restricted Function", "Compilation Error"},
	"HRBUST":        {"Accepted", "Presentation Error", "Wrong Answer", "Time Limit Exceeded", "Memory Limit Exceeded", "Runtime Error", "Compile Error", "Restricted Function"},
	"HUST":          {"Accepted", "Presentation Error", "Wrong Answer", "Time Limit Exceed", "Memory Limit Exceed", "Output Limit Exceed", "Runtime Error", "Compile Error"},
	"HYSBZ":         {"Accepted", "Wrong Answer", "Time Limit Exceed", "Output Limit Exceed", "Runtime Error", "Compile Error"},
	"Kattis":        {"Accepted", "Wrong Answer", "Time Limit Exceeded", "Memory Limit Exceeded", "Output Limit Exceeded", "Run Time Error", "Compile Error", "Source Code Error"},
	"LibreOJ":       {"Accepted", "Wrong Answer", "File Error", "Time Limit Exceeded", "Memory Limit Exceeded", "Output Limit Exceeded", "Runtime Error", "Compile Error"},
	"LightOJ":       {"Accepted", "Presentation Error", "Wrong Answer", "Time Limit Exceeded", "Memory Limit Exceeded", "Output Limit Exceeded", "Runtime Error", "Compilation Error", "Restricted Function"},
	"Minieye":       {"Accepted", "Wrong Answer", "CPU Time Limit Exceeded", "Runtime Error", "Compile Error"},
	"NBUT":          {"Accepted", "Presentation error", "Wrong answer", "Time limit exceeded", "Memory limit exceeded", "Output limit exceeded", "Runtime error", "Compilation error", "Dangerous code"},
	"OpenJ_Bailian": {"Accepted", "Presentation Error", "Wrong Answer", "Time Limit Exceeded", "Memory Limit Exceeded", "Output Limit Exceeded", "Runtime Error", "Compile Error"},
	"OpenJ_POJ":     {"Accepted", "Presentation Error", "Wrong Answer", "Time Limit Exceeded", "Memory Limit Exceeded", "Output Limit Exceeded", "Runtime Error", "Compile Error"},
	"POJ":           {"Accepted", "Presentation Error", "Wrong Answer", "Time Limit Exceeded", "Memory Limit Exceeded", "Output Limit Exceeded", "Runtime Error", "Compile Error"},
	"SCU":           {"Accepted", "Presentation Error", "Wrong Answer", "Time Limit Exceeded", "Memory Limit Exceeded", "Output Limit Exceeded", "Runtime Error", "Compilation Error"},
	"SGU":           {"Accepted", "Wrong answer", "Time limit exceeded", "Memory limit exceeded", "Runtime error", "Compilation error", "Source Code Error", "Dangerous Code Error", "Language Rejected"},
	"SPOJ":          {"Accepted", "Wrong answer", "Time limit exceeded", "Runtime error", "Compilation error"},
	"TopCoder":      {"Accepted", "Wrong Answer", "Time Limit Exceeded", "Compile Error"},
	"Toph":          {"Accepted", "Passed", "Wrong answer", "CPU limit exceeded", "Memory limit exceeded", "Compilation error", "Runtime error"},
	"UESTC":         {"Accepted", "Presentation Error", "System Error", "Wrong Answer", "Time Limit Exceeded", "Memory Limit Exceeded", "Output Limit Exceeded", "Restricted Function"},
	"UESTC_old":     {"Accepted", "Presentation Error", "Wrong Answer", "Time Limit Exceeded", "Memory Limit Exceeded", "Output Limit Exceeded", "Runtime Error", "Compile Error"},
	"UniversalOJ":   {"Accepted", "Wrong Answer", "Extra Test Failed", "Dangerous Syscalls", "Time Limit Exceeded", "Memory Limit Exceeded", "Runtime Error", "Compile Error"},
	"URAL":          {"Accepted", "Wrong answer", "Time limit exceeded", "Memory limit exceeded", "Output limit exceeded", "Runtime error", "Compilation error", "Language Rejected"},
	"URI":           {"Accepted", "Compilation error", "Runtime error", "Time limit exceeded", "Presentation error", "Wrong answer", "Closed", "Possible runtime error", "Memory limit exceeded"},
	"UVA":           {"Accepted", "Presentation error", "Wrong answer", "Time limit exceeded", "Output limit exceeded", "Runtime error", "Compilation error", "Compile error", "Code Length Exceeded"},
	"UVALive":       {"Accepted", "Presentation error", "Wrong answer", "Time limit exceeded", "Output limit exceeded", "Runtime error", "Compilation error", "Compile error", "Code Length Exceeded"},
	"Z_trening":     {"Accepted", "WRONG ANSWER", "System error", "Time limit exceeded", "Memory limit exceeded", "Runtime error", "Compile Error", "Error"},
	"ZOJ":           {"Accepted", "Presentation Error", "Wrong Answer", "Time Limit Exceeded", "Memory Limit Exceeded", "Output Limit Exceeded", "Runtime Error", "Segmentation Fault", "Float Point Exception", "Compilation Error"},
	"计蒜客":           {"Accepted", "Presentation Error", "Wrong Answer", "Time Limit Exceeded", "Memory Limit Exceeded", "Output Limit Exceeded", "Runtime Error", "Segmentation Fault", "Arithmetical Error", "Compilation Error"},
	"黑暗爆炸":          {"Accepted", "Wrong Answer", "Dangerous Syscalls", "Time Limit Exceeded", "Checker Time Limit Exceeded", "Memory Limit Exceeded", "Runtime Error", "Compile Error"},
}

//Addition var for adding 2 int (used in go tmpl)
type Addition struct {
}

//Add func is a method of Addition struct/object
func (add Addition) Add(a, b int) int {
	return a + b
}

//ProblemSet variable for holding problemset of a contest
type ProblemSet struct {
	SerialIndex string `bson:"serialIndex"`
	OJ          string `bson:"OJ"`
	PNum        string `bson:"pNum"`
	PName       string `bson:"pName"`
	CustomName  string `bson:"customName"`
}

//Clarification variable for holding clarification of a contest
type Clarification struct {
	SerialIndex   string `bson:"serialIndex"`
	RequesterName string `bson:"requesterName"`
	RequestBody   string `bson:"requestBody"`
	RequestAt     string `bson:"requestAt"`
	AnsweredBy    string `bson:"answeredBy"`
	AnswerBody    string `bson:"answerBody"`
	AnswerAt      string `bson:"answerAt"`
	IsIgnored     bool   `bson:"isIgnored"`
}

//ContestData variable for holding data of a single contest
type ContestData struct {
	ContestID      int             `bson:"contestID"`
	Title          string          `bson:"title"`
	Date           string          `bson:"date"`
	Time           string          `bson:"time"`
	StartAt        int64           `bson:"startAt"`
	Duration       string          `bson:"duration"`
	FrozenTime     string          `bson:"frozenTime"`
	Author         string          `bson:"author"`
	ProblemSet     []ProblemSet    `bson:"problemSet"`
	Clarifications []Clarification `bson:"clarifications"`
}

//LastUsedID variable for holding the last ID used for user registration, problem submission & contest creation
type LastUsedID struct {
	LastUserID       int `bson:"lastUserID"`
	LastSubmissionID int `bson:"lastSubmissionID"`
	LastContestID    int `bson:"lastContestID"`
}

//SubmissionData variable for holding a single submission details
type SubmissionData struct {
	SubID           int    `bson:"subID"`
	Username        string `bson:"username"`
	OJ              string `bson:"OJ"`
	PNum            string `bson:"pNum"`
	PName           string `bson:"pName"`
	Language        string `bson:"language"`
	SubmittedAt     int64  `bson:"submittedAt"`
	VID             string `bson:"vID"`
	SourceCode      string `bson:"sourceCode"`
	Verdict         string `bson:"verdict"`
	MemoryExec      string `bson:"memoryExec"`
	TimeExec        string `bson:"timeExec"`
	TerminalVerdict bool   `bson:"terminalVerdict"`
	ContestID       int    `bson:"contestID"`
	SerialIndex     string `bson:"serialIndex"`
}

//UserData variable for holding a single user details
type UserData struct {
	UserID               int    `bson:"userID"`
	FullName             string `bson:"fullName"`
	Email                string `bson:"email"`
	Username             string `bson:"username"`
	Password             string `bson:"password"`
	CreatedAt            int64  `bson:"createdAt"`
	IsVerified           bool   `bson:"isVerified"`
	AccVerifyToken       string `bson:"accVerifyToken"`
	AccVerifyTokenSentAt int64  `bson:"accVerifyTokenSentAt"`
	PassResetToken       string `bson:"passResetToken"`
	PassResetTokenSentAt int64  `bson:"passResetTokenSentAt"`
	TotalSolved          int    `bson:"totalSolved"`
}

//FeedbackData variable for holding a feedback details
type FeedbackData struct {
	Name    string `bson:"name"`
	Email   string `bson:"email"`
	Message string `bson:"mesasge"`
	SentAt  int64  `bson:"sentAt"`
}
