package gpt

import (
	"chatgpt-web-go/src/global"
	enum "chatgpt-web-go/src/global/enum/gpt"
	gptmodel "chatgpt-web-go/src/model/api/gpt"
	"chatgpt-web-go/src/model/api/gpt/request"
	models "chatgpt-web-go/src/model/api/gpt/response"
	result "chatgpt-web-go/src/model/common/response"
	"chatgpt-web-go/src/repository"
	"chatgpt-web-go/src/service/gpt"
	"chatgpt-web-go/src/utils"
	"context"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sashabaranov/go-openai"
	"go.uber.org/zap"
	"io"
	"net/http"
)

func GPT(c *gin.Context) {
	c.Header("Content-type", "application/octet-stream")
	var req request.ChatProcessRequest
	if err := c.BindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, result.Fail.WithMessage(err.Error()))
		return
	}

	// 组装chatCompletionMessages
	chatMessageService := gpt.NewChatMessageService()
	chatMessageDO, completionRequest, err := chatMessageService.GetOpenAiRequestReady(req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK, result.Fail.WithMessage(err.Error()))
		return
	}

	// 计算token是否超限
	numTokens := utils.NumTokensFromMessages(completionRequest.Messages, openai.GPT3Dot5Turbo)
	if numTokens+global.Cfg.GPT.MaxToken > 4000 {
		c.AbortWithStatusJSON(http.StatusOK, result.Fail.WithMessage("上文聊天记录已超限, 请新建聊天或不携带聊天记录"))
		return
	}

	// 流式输出
	stream, err := global.GPTClient.CreateChatCompletionStream(context.Background(), completionRequest)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK, result.Fail.WithMessage(err.Error()))
		return
	}
	defer stream.Close()
	resText := ""
	responseCount := 0

	// 第一次生成messageId和conversationId uuid
	if chatMessageDO.MessageID == "" {
		chatMessageDO.MessageID = uuid.New().String()
	}
	if chatMessageDO.ConversationID == "" {
		chatMessageDO.ConversationID = uuid.New().String()
	}

	// 存入request
	chatMessageRepo := repository.NewChatMessageRepository()
	var questionDO gptmodel.ChatMessageDO
	var answerDO gptmodel.ChatMessageDO
	if err := utils.DeepCopyByJson(&chatMessageDO, &questionDO); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, result.Fail.WithMessage(err.Error()))
		return
	}
	if err := utils.DeepCopyByJson(&chatMessageDO, &answerDO); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, result.Fail.WithMessage(err.Error()))
		return
	}
	questionDO.IP = c.ClientIP()
	questionDO.OriginalData = func() string {
		jsonV, _ := json.Marshal(completionRequest)
		return string(jsonV)
	}()
	questionDO.PromptTokens = numTokens
	questionDO.Status = enum.PART_SUCCESS
	questionDO.MessageType = enum.QUESTION
	questionDO.ParentAnswerMessageID = questionDO.ParentMessageID

	answerDO.ID = uint64(utils.GetSnowIdInt64())
	answerDO.MessageID = uuid.New().String()
	answerDO.ParentMessageID = questionDO.MessageID
	answerDO.ParentQuestionMessageID = questionDO.MessageID
	chatMessageRepo.CreateChatMessage(&questionDO)

	c.Stream(func(w io.Writer) bool {
		for {
			response, err := stream.Recv()
			if errors.Is(err, io.EOF) || (response.Choices != nil && response.Choices[0].FinishReason == "stop") {

				//存入response
				answerDO.Content = resText
				answerDO.MessageType = enum.ANSWER
				answerDO.Status = enum.COMPLETE_SUCCESS
				answerDO.PromptTokens = numTokens
				answerDO.CompletionTokens = utils.NumTokensFromText(resText, openai.GPT3Dot5Turbo)
				answerDO.TotalTokens = answerDO.PromptTokens + answerDO.CompletionTokens

				if response.Choices != nil {
					response.Choices[0].Delta.Content = resText
					jsonV, _ := json.Marshal(response)
					answerDO.OriginalData = string(jsonV)
				} else {
					response.Choices = []openai.ChatCompletionStreamChoice{
						{
							Delta: openai.ChatCompletionStreamChoiceDelta{
								Content: resText,
							},
						},
					}
					jsonV, _ := json.Marshal(response)
					answerDO.OriginalData = string(jsonV)
				}

				questionDO.Status = enum.COMPLETE_SUCCESS

				chatMessageRepo.UpdateChatMessage(&questionDO)
				err := chatMessageRepo.CreateChatMessage(&answerDO)
				if err != nil {
					global.Gzap.Error("chatMessageRepo.CreateChatMessage", zap.Error(err))
					c.AbortWithStatusJSON(http.StatusOK, result.Fail.WithMessage(err.Error()))
					return false
				}
				return false
			}
			if err != nil {
				c.JSON(http.StatusOK, result.OK.WithData(err))
				return false
			}
			resText = resText + response.Choices[0].Delta.Content
			chatReplyMessageVO := new(models.ChatReplyMessage)
			chatReplyMessageVO.ID = answerDO.MessageID
			chatReplyMessageVO.Role = ""
			chatReplyMessageVO.ParentMessageID = answerDO.ParentMessageID
			chatReplyMessageVO.ConversationID = answerDO.ConversationID
			chatReplyMessageVO.Text = resText
			re, _ := json.Marshal(chatReplyMessageVO)
			if responseCount != 0 {
				re = append([]byte("\n"), re...)
			}
			_, writeErr := w.Write(re)
			if writeErr != nil {
				c.AbortWithStatusJSON(http.StatusOK, result.Fail.WithMessage(writeErr.Error()))
				return false
			}
			if f, ok := w.(http.Flusher); ok {
				f.Flush()
			} else {
				c.AbortWithStatusJSON(http.StatusOK, result.Fail.WithMessage("Unable to flush response"))
				return false
			}
			responseCount++
		}
	})
}
