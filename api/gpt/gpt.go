package gpt

import (
	"chatgpt-web-go/global"
	enum "chatgpt-web-go/global/enum/gpt"
	model "chatgpt-web-go/model/api/gpt"
	"chatgpt-web-go/model/api/gpt/request"
	models "chatgpt-web-go/model/api/gpt/response"
	result "chatgpt-web-go/model/common/response"
	"chatgpt-web-go/service/gpt"
	"chatgpt-web-go/utils"
	"context"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/sashabaranov/go-openai"
	"go.uber.org/zap"
	"io"
	"net/http"
	"strconv"
)

func ChatConversationProcess(c *gin.Context) {
	c.Header("Content-type", "application/octet-stream")
	var req request.ChatProcessRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, result.Fail.WithMessage(err.Error()))
		return
	}

	chatMessageService := gpt.NewChatMessageService()
	chatConversationService := gpt.NewChatConversationService()

	// 通过Conversation来判断是否是新的会话,获取会话
	chatConversation := new(model.ChatConversation)
	if err := chatConversationService.InitChatConversation(chatConversation, req); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, result.Fail.WithMessage(err.Error()).WithData(err))
		return
	}
	aiRequest, err := chatMessageService.GetOpenAiRequest(req, chatConversation)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK, result.Fail.WithMessage(err.Error()))
		return
	}
	chatConversation.Question.IP = c.RemoteIP()
	chatConversation.Answer.IP = c.RemoteIP()
	if err := chatMessageService.SaveChatMessage(chatConversation.Question); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, result.Fail.WithMessage(err.Error()))
		return
	}
	if err := chatMessageService.SaveChatMessage(chatConversation.Answer); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, result.Fail.WithMessage(err.Error()))
		return
	}
	processChatCompletionStream(c, chatConversation, &chatConversationService, &chatMessageService, aiRequest)
}

func AddChatRoom(c *gin.Context) {
	chatRoomService := gpt.NewChatRoomService()
	if room, err := chatRoomService.CreateChatRoom(); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, result.Fail.WithMessage(err.Error()))
		return
	} else {
		c.JSON(http.StatusOK, result.OK.WithData(room))
		return
	}
}

func saveAnswer(chatConversation *model.ChatConversation, chatConversationService *gpt.ChatConversationService, chatMessageService *gpt.ChatMessageService, response openai.ChatCompletionStreamResponse, status enum.ChatMessageStatusEnum, resText string, c *gin.Context) {
	chatConversation.Answer.Content = resText
	chatConversation.Answer.Status = status
	chatConversation.Answer.TotalTokens = utils.NumTokensFromText(resText, openai.GPT3Dot5Turbo)

	if response.Choices != nil {
		response.Choices[0].Delta.Content = resText
		jsonV, _ := json.Marshal(response)
		chatConversation.Answer.OriginalData = string(jsonV)
	} else {
		response.Choices = []openai.ChatCompletionStreamChoice{
			{
				Delta: openai.ChatCompletionStreamChoiceDelta{
					Content: resText,
				},
			},
		}
		jsonV, _ := json.Marshal(response)
		chatConversation.Answer.OriginalData = string(jsonV)
	}

	if err := (*chatMessageService).UpdateChatMessage(chatConversation.Answer); err != nil {
		global.Gzap.Error("chatMessageRepo.CreateChatMessage", zap.Error(err))
		c.AbortWithStatusJSON(http.StatusOK, result.Fail.WithMessage(err.Error()))
	}
	if err := (*chatConversationService).CreateConversation(chatConversation); err != nil {
		global.Gzap.Error("chatMessageRepo.CreateChatMessage", zap.Error(err))
		c.AbortWithStatusJSON(http.StatusOK, result.Fail.WithMessage(err.Error()))
	}
}

func handleResponse(responseCount int, chatConversationService *gpt.ChatConversationService, chatConversation *model.ChatConversation, chatMessageService *gpt.ChatMessageService, response *openai.ChatCompletionStreamResponse, resText *string, c *gin.Context, w io.Writer) bool {
	*resText = *resText + response.Choices[0].Delta.Content
	chatReplyMessageVO := new(models.ChatReplyMessage)
	chatReplyMessageVO.Id = strconv.FormatUint(chatConversation.Id, 10)
	chatReplyMessageVO.Role = ""
	chatReplyMessageVO.ParentMessageId = strconv.FormatUint(chatConversation.ParentId, 10)
	chatReplyMessageVO.Text = *resText
	re, _ := json.Marshal(chatReplyMessageVO)
	if responseCount != 0 {
		re = append([]byte("\n"), re...)
	}
	_, writeErr := w.Write(re)
	if writeErr != nil {
		saveAnswer(chatConversation, chatConversationService, chatMessageService, *response, enum.PART_SUCCESS, *resText, c)
		c.AbortWithStatusJSON(http.StatusOK, result.Fail.WithMessage(writeErr.Error()))
		return false
	}
	if f, ok := w.(http.Flusher); ok {
		f.Flush()
	} else {
		saveAnswer(chatConversation, chatConversationService, chatMessageService, *response, enum.PART_SUCCESS, *resText, c)
		c.AbortWithStatusJSON(http.StatusOK, result.Fail.WithMessage("Unable to flush response"))
		return false
	}
	return true
}

func processChatCompletionStream(c *gin.Context, chatConversation *model.ChatConversation, chatConversationService *gpt.ChatConversationService, chatMessageService *gpt.ChatMessageService, aiRequest openai.ChatCompletionRequest) {
	stream, err := global.GPTClient.CreateChatCompletionStream(context.Background(), aiRequest)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK, result.Fail.WithMessage(err.Error()))
		return
	}
	defer stream.Close()

	resText := ""
	responseCount := 0

	c.Stream(func(w io.Writer) bool {
		for {
			response, err := stream.Recv()
			if errors.Is(err, io.EOF) || (response.Choices != nil && response.Choices[0].FinishReason == "stop") {
				// 存入 response
				saveAnswer(chatConversation, chatConversationService, chatMessageService, response, enum.COMPLETE_SUCCESS, resText, c)
				return false
			}
			if err != nil {
				saveAnswer(chatConversation, chatConversationService, chatMessageService, response, enum.PART_SUCCESS, resText, c)
				c.AbortWithStatusJSON(http.StatusOK, result.Fail.WithMessage(err.Error()))
				return false
			}

			if !handleResponse(responseCount, chatConversationService, chatConversation, chatMessageService, &response, &resText, c, w) {
				return false
			}

			responseCount++
		}
	})
}
