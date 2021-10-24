package discord

import (
	"fmt"
	"strings"

	"github.com/nahidhasan98/ajudge/model"
)

type discordInterfacer interface {
	SendMessage(data interface{}, notifier string)
	EditMessage(data interface{}, notifier string)
}

type discordStruct struct {
	repoService repoInterfacer
}

func (ds discordStruct) SendMessage(data interface{}, notifier string) {
	jobs := make(chan int, 5)
	go sendWorker(jobs, data, notifier, ds)

	jobs <- 1
	close(jobs)
}

func (ds discordStruct) EditMessage(data interface{}, notifier string) {
	jobs := make(chan int, 5)
	go editWorker(jobs, data, notifier, ds)

	jobs <- 1
	close(jobs)
}

func prepareSubmissionMessage(data model.SubmissionData, formattedTime string) string {
	disMsg := "```md\n"
	disMsg += "# " + data.OJ + " - " + data.Username + "\n"
	disMsg += "Submission ID" + getSpace(3) + ": " + fmt.Sprintf("%v", data.SubID) + "\n"
	disMsg += "Remote ID" + getSpace(7) + ": " + data.VID + "\n"
	disMsg += "Username" + getSpace(8) + ": " + data.Username + "\n"
	disMsg += "OJ" + getSpace(14) + ": " + data.OJ + "\n"
	disMsg += "Problem Number" + getSpace(2) + ": " + data.PNum + "\n"
	disMsg += "Problem Name" + getSpace(4) + ": " + data.PName + "\n"
	disMsg += "Language" + getSpace(8) + ": " + data.Language + "\n"
	disMsg += "Time" + getSpace(12) + ": " + data.TimeExec + "\n"
	disMsg += "Memory" + getSpace(10) + ": " + data.MemoryExec + "\n"
	disMsg += "Submitted At" + getSpace(4) + ": " + formattedTime + "\n"
	disMsg += "Verdict" + getSpace(9) + ": " + data.Verdict + "\n"
	disMsg += "Contest ID" + getSpace(6) + ": " + fmt.Sprintf("%v", data.ContestID) + "\n"
	disMsg += "Serial Index" + getSpace(4) + ": " + data.SerialIndex + "\n"
	disMsg += "```"

	return disMsg
}

func prepareLoginRegMessage(data UserModel, formattedTime string) string {
	disMsg := "```md\n"
	disMsg += "# " + data.Username + "\n"
	disMsg += "Username" + getSpace(4) + ": " + data.Username + "\n"
	disMsg += "Email" + getSpace(7) + ": " + data.Email + "\n"
	disMsg += "Full Name" + getSpace(3) + ": " + data.FullName + "\n"
	disMsg += "Verified" + getSpace(4) + ": " + fmt.Sprintf("%v", data.IsVerified) + "\n"
	disMsg += "Member Since" + getSpace(0) + ": " + formattedTime + "\n"
	disMsg += "Total Solved" + getSpace(0) + ": " + fmt.Sprintf("%v", data.TotalSolved) + "\n"
	disMsg += "```"

	return disMsg
}

func prepareContestMessage(data model.ContestData, formattedTime string) string {
	disMsg := "```md\n"
	disMsg += "# " + fmt.Sprintf("%v", data.ContestID) + " - " + data.Title + "\n"
	disMsg += "Contest ID" + getSpace(2) + ": " + fmt.Sprintf("%v", data.ContestID) + "\n"
	disMsg += "Title" + getSpace(7) + ": " + data.Title + "\n"
	disMsg += "Start At" + getSpace(4) + ": " + formattedTime + "\n"
	disMsg += "Duration" + getSpace(4) + ": " + fmt.Sprintf("%v", data.Duration) + "\n"
	disMsg += "Author" + getSpace(6) + ": " + data.Author + "\n"
	disMsg += "```"

	return disMsg
}

func prepareResetMessage(data model.UserData) string {
	disMsg := "```md\n"
	disMsg += "Username" + getSpace(0) + ": " + data.Username + "\n"
	disMsg += "```"

	return disMsg
}

func prepareFeedbackMessage(data model.FeedbackData) string {
	disMsg := "```md\n"
	disMsg += "# " + data.Email + "\n"
	disMsg += "Name" + getSpace(4) + ": " + data.Name + "\n"
	disMsg += "Email" + getSpace(3) + ": " + data.Email + "\n"
	disMsg += "Message" + getSpace(1) + ": " + data.Message + "\n"
	disMsg += "```"

	return disMsg
}

func prepareSubmissionEditedMessage(old string, data model.SubmissionData) string {
	// tt := "Time1234567891234: ---\nMemory"
	idx1 := strings.Index(old, "Time")
	idx2 := strings.Index(old, "Memory")
	disMsgV1 := old[:idx1+18] + data.TimeExec + "\n" + old[idx2:]

	// tt := "Memory12345678912: ---\nSubmitted"
	idx1 = strings.Index(disMsgV1, "Memory")
	idx2 = strings.Index(disMsgV1, "Submitted")
	disMsgV2 := disMsgV1[:idx1+18] + data.MemoryExec + "\n" + disMsgV1[idx2:]

	// tt := "Verdict1234567891: ---\nContest"
	idx1 = strings.Index(disMsgV2, "Verdict")
	idx2 = strings.Index(disMsgV2, "Contest")
	disMsg := disMsgV2[:idx1+18] + data.Verdict + "\n" + disMsgV2[idx2:]

	return disMsg
}

func prepareContestEditedMessage(old string, data model.ContestData) string {
	// tt := "- some contest title\nContest ID"
	idx1 := strings.Index(old, "-")
	idx2 := strings.Index(old, "Contest ID")
	disMsgV1 := old[:idx1+2] + data.Title + "\n" + old[idx2:]

	// tt := "Title1234567: some contest title\nStart At"
	idx1 = strings.Index(disMsgV1, "Title")
	idx2 = strings.Index(disMsgV1, "Start At")
	disMsgV2 := disMsgV1[:idx1+14] + data.Title + "\n" + disMsgV1[idx2:]

	// tt := "Duration1234: --:--\nAuthor"
	idx1 = strings.Index(disMsgV2, "Duration")
	idx2 = strings.Index(disMsgV2, "Author")
	disMsg := disMsgV2[:idx1+14] + data.Duration + "\n" + disMsgV2[idx2:]

	return disMsg
}

// for indentation purpose
func getSpace(num int) string {
	str := ""

	for i := 0; i < num; i++ {
		str += " "
	}

	return str
}

func newDiscordService(repo repoInterfacer) discordInterfacer {
	return &discordStruct{
		repoService: repo,
	}
}
