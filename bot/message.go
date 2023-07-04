package bot

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"
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
	if resp.Body != nil {
		defer resp.Body.Close()
	}

	body, err := io.ReadAll(resp.Body)
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
	if resp.Body != nil {
		defer resp.Body.Close()
	}

	body, err := io.ReadAll(resp.Body)
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
	// desc = strings.ReplaceAll(desc, "`", "&#96;")
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
	// var elements []interface{}
	// elements = []interface{}{
	// 	&CardElement{
	// 		Tag: "div",
	// 		Fields: []*CardField{
	// 			{
	// 				IsShort: false,
	// 				Text: &CardText{
	// 					Content: desc,
	// 					Tag:     "lark_md",
	// 				},
	// 			},
	// 		},
	// 	},
	// }

	node, err := parseMarkdown([]byte(desc))
	if err != nil {
		panic(err)
	}
	elements := convertNodesToJSON(node, []byte(desc))
	cardContent.Elements = elements

	cardBytes, err := json.Marshal(cardContent)
	if err != nil {
		panic(err)
	}
	return string(cardBytes)
}

func parseMarkdown(markdown []byte) (ast.Node, error) {
	md := goldmark.New(goldmark.WithExtensions())
	reader := text.NewReader(markdown)
	return md.Parser().Parse(reader), nil
}

func convertNodesToJSON(astRoot ast.Node, source []byte) []*CardElement {
	elements := []*CardElement{}

	ast.Walk(astRoot, func(node ast.Node, entering bool) (ast.WalkStatus, error) {
		if entering {
			switch n := node.(type) {
			case *ast.Heading:
				if n.Level == 2 {
					elements = append(elements, &CardElement{
						Tag: "div",
						Text: &CardText{
							Tag:     "lark_md",
							Content: "<font color='green'>" + extractTextAndLinks(n, source) + "</font><at id=all></at>",
						},
					})
				} else {
					elements = append(elements, &CardElement{
						Tag: "div",
						Text: &CardText{
							Tag:     "lark_md",
							Content: "**" + extractTextAndLinks(n, source) + "**",
						},
					})
				}

			case *ast.List:
				listContent := ""
				i := 1
				ast.Walk(n, func(child ast.Node, entering bool) (ast.WalkStatus, error) {
					if entering {
						switch listItem := child.(type) {
						case *ast.ListItem:
							listContent += strconv.Itoa(i) + ". " + strings.TrimSpace(extractTextAndLinks(listItem, source)) + "\n"
							i++
						}
					}
					return ast.WalkContinue, nil
				})
				if listContent != "" {
					elements = append(elements, &CardElement{
						Tag: "div",
						Fields: []*CardField{
							{
								Text: &CardText{
									Tag:     "lark_md",
									Content: listContent,
								},
							},
						},
					})
				}
			case *ast.FencedCodeBlock:
				code := ""
				lines := n.Lines()
				for i := 0; i < lines.Len(); i++ {
					line := lines.At(i)
					code += string(line.Value(source))
				}
				if code != "" {
					elements = append(elements, &CardElement{Tag: "hr"}, &CardElement{
						Tag: "div",
						Fields: []*CardField{
							{
								Text: &CardText{
									Tag:     "lark_md",
									Content: code,
								},
							},
						},
					}, &CardElement{Tag: "hr"})
				}
			}
		}
		return ast.WalkContinue, nil
	})

	return elements
}

func extractTextAndLinks(node ast.Node, source []byte) string {
	var buffer bytes.Buffer
	ast.Walk(node, func(child ast.Node, entering bool) (ast.WalkStatus, error) {
		if entering {
			switch childNode := child.(type) {
			case *ast.Text:
				buffer.WriteString(string(childNode.Segment.Value(source)))
			case *ast.Link:
				destination := string(childNode.Destination)
				linkText := string(childNode.Text(source))
				buffer.WriteString(fmt.Sprintf("[%s](%s)", linkText, destination))
				return ast.WalkSkipChildren, nil
			}
		}
		return ast.WalkContinue, nil
	})
	return buffer.String()
}
