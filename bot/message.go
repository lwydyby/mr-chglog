package bot

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

var (
	createMessageURL     = "https://open.feishu.cn/open-apis/im/v1/messages"
	getMessageHistoryURL = "https://open.feishu.cn/open-apis/im/v1/messages"
	pinURL               = "https://open.feishu.cn/open-apis/im/v1/pins"
)

func SendAlertMessage(ctx context.Context, token, chatID string, title, text string) error {
	var err error

	var createResp *MessageItem
	var createReq *CreateMessageRequest

	cardContent := constructAlterCard(title, text)
	createReq = genCreateMessageRequest(ctx, chatID, cardContent, "interactive")

	createResp, err = sendMessage(ctx, token, createReq)
	if err != nil {
		panic(err)
	}

	msgID := createResp.MessageID
	pinMessage(token, msgID)
	return nil
}

func sendMessage(ctx context.Context, token string, createReq *CreateMessageRequest) (*MessageItem, error) {
	var err error

	cli := &http.Client{}

	reqBytes, err := json.Marshal(createReq)
	if err != nil {
		panic(err)
	}
	req, err := http.NewRequest("POST", createMessageURL, strings.NewReader(string(reqBytes)))
	if err != nil {
		panic(err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	q := req.URL.Query()
	q.Add("receive_id_type", "chat_id")
	req.URL.RawQuery = q.Encode()
	resp, err := cli.Do(req)
	if err != nil {
		return nil, fmt.Errorf("create message failed, err=%v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	createMessageResp := &CreateMessageResponse{}
	err = json.Unmarshal(body, createMessageResp)
	if err != nil {
		panic(err)
	}
	if createMessageResp.Code != 0 {
		fmt.Println(string(body))
		panic(err)
	}
	fmt.Printf("succeed create message, msg_id: %v", createMessageResp.Data.MessageID)
	return createMessageResp.Data, nil
}

func pinMessage(token string, messageID string) {
	cli := &http.Client{}

	reqBytes, err := json.Marshal(&PinMessageRequest{
		MessageID: messageID,
	})
	if err != nil {
		panic(err)
	}
	req, err := http.NewRequest("POST", pinURL, strings.NewReader(string(reqBytes)))
	if err != nil {
		panic(err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	resp, err := cli.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	pinMessageResp := &PinMessageResponse{}
	err = json.Unmarshal(body, pinMessageResp)
	if err != nil {
		panic(err)
	}
	if pinMessageResp.Code != 0 {
		fmt.Println(string(body))
		panic(err)
	}
}

func genCreateMessageRequest(ctx context.Context, chatID, content, msgType string) *CreateMessageRequest {
	return &CreateMessageRequest{
		ReceiveID: chatID,
		Content:   content,
		MsgType:   msgType,
	}
}

func constructAlterCard(title, desc string) (card string) {
	desc = strings.ReplaceAll(desc, "`", "&#96;")
	cardContent := &CardContent{
		Config: &CardConfig{
			WideScreenMode: true,
		},
		Header: &CardHeader{
			Template: "blue",
			Title: &CardText{
				Tag:     "lark_md",
				Content: title,
			},
		},
	}
	var elements []interface{}
	elements = []interface{}{
		&CardElement{
			Tag: "div",
			Fields: []*CardField{
				{
					IsShort: false,
					Text: &CardText{
						Content: desc,
						Tag:     "lark_md",
					},
				},
			},
		},
	}
	cardContent.Elements = elements

	cardBytes, err := json.Marshal(cardContent)
	if err != nil {
		panic(err)
	}
	return string(cardBytes)
}
