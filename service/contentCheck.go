package service

import (
	"encoding/json"
	green20220302 "github.com/alibabacloud-go/green-20220302/v2/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
	"net/http"
	"note_app_server/global"
	"note_app_server/model/commentModel"
	"note_app_server/model/noteModel"
)

// CheckNoteContent 检查笔记内容是否合法
// return 1通过审核 2转入人工审核 3审核未通过
func CheckNoteContent(note *noteModel.ESNote) int {
	client := global.ContentCheckClient

	runtime := &util.RuntimeOptions{}
	runtime.ReadTimeout = tea.Int(10000)
	runtime.ConnectTimeout = tea.Int(10000)

	serviceParameters, _ := json.Marshal(
		map[string]interface{}{
			"content": note.Title + note.Content,
		},
	)
	request := green20220302.TextModerationPlusRequest{
		Service:           tea.String("chat_detection_pro"),
		ServiceParameters: tea.String(string(serviceParameters)),
	}

	result, _err := client.TextModerationPlus(&request)
	if _err != nil {
		panic(_err)
	}

	if *result.StatusCode != http.StatusOK {
		CheckNoteContent(note)
	}
	body := result.Body
	if *body.Code != http.StatusOK {
		CheckNoteContent(note)
	}

	if *body.Data.RiskLevel == "none" {
		return 1
	}
	data := body.Data.Result[0]
	if *data.Confidence >= 90 {
		// 高风险
		return 3
	} else if *data.Confidence >= 70 {
		// 中风险
		return 2
	} else {
		// 低风险
		return 1
	}
}

// CheckCommentContent 检查评论内容
// 评论内容不执行人工审核
func CheckCommentContent(cmt *commentModel.Comment) int {
	client := global.ContentCheckClient

	runtime := &util.RuntimeOptions{}
	runtime.ReadTimeout = tea.Int(10000)
	runtime.ConnectTimeout = tea.Int(10000)

	serviceParameters, _ := json.Marshal(
		map[string]interface{}{
			"content": cmt.Content,
		},
	)
	request := green20220302.TextModerationPlusRequest{
		Service:           tea.String("chat_detection_pro"),
		ServiceParameters: tea.String(string(serviceParameters)),
	}

	result, _err := client.TextModerationPlus(&request)
	if _err != nil {
		panic(_err)
	}

	if *result.StatusCode != http.StatusOK {
		return 2
	}
	body := result.Body
	if *body.Code != http.StatusOK {
		return 2
	}

	if *body.Data.RiskLevel == "none" {
		return 1
	}
	data := body.Data.Result[0]
	if *data.Confidence >= 90 {
		// 高风险
		return 2
	} else {
		// 低风险
		return 1
	}
}

// CheckUserAvatar 检查用户头像是否违规
func CheckUserAvatar(content string, repeat int, maxRepeat int) {
	if repeat >= maxRepeat {
		return
	}
	go func() {
		client := global.ContentCheckClient

		runtime := &util.RuntimeOptions{}
		runtime.ReadTimeout = tea.Int(10000)
		runtime.ConnectTimeout = tea.Int(10000)

		serviceParameters, _ := json.Marshal(
			map[string]interface{}{
				"content": content,
			},
		)
		request := green20220302.TextModerationPlusRequest{
			Service:           tea.String("profilePhotoCheck"),
			ServiceParameters: tea.String(string(serviceParameters)),
		}

		result, _err := client.TextModerationPlus(&request)
		if _err != nil {
			panic(_err)
		}

		if *result.StatusCode != http.StatusOK {
			CheckUserAvatar(content, repeat+1, maxRepeat)
			return
		}
		body := result.Body
		if *body.Code != http.StatusOK {
			CheckUserAvatar(content, repeat+1, maxRepeat)
			return
		}

		data := body.Data.Result[0]
		if *data.Confidence >= 60 {
			//打回
		} else {

		}
	}()
}

// CheckUsername 检查用户名是否违规
func CheckUsername(content string, repeat int, maxRepeat int) {
	if repeat >= maxRepeat {
		return
	}
	go func() {
		client := global.ContentCheckClient

		runtime := &util.RuntimeOptions{}
		runtime.ReadTimeout = tea.Int(10000)
		runtime.ConnectTimeout = tea.Int(10000)

		serviceParameters, _ := json.Marshal(
			map[string]interface{}{
				"content": content,
			},
		)
		request := green20220302.TextModerationPlusRequest{
			Service:           tea.String("nickname_detection"),
			ServiceParameters: tea.String(string(serviceParameters)),
		}

		result, _err := client.TextModerationPlus(&request)
		if _err != nil {
			panic(_err)
		}

		if *result.StatusCode != http.StatusOK {
			CheckUsername(content, repeat+1, maxRepeat)
			return
		}
		body := result.Body
		if *body.Code != http.StatusOK {
			CheckUsername(content, repeat+1, maxRepeat)
			return
		}

		data := body.Data.Result[0]
		if *data.Confidence >= 60 {
			//打回
		} else {

		}
	}()
}
