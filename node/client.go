package node

import (
	context "context"
	"encoding/json"
	"fmt"
	"log"
	"metanet/rpc"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type Client = Node

func (c *Client) ClientRequestCall(cmd []byte, to string, md map[string]string) (result *rpc.ClientResults, err error) {

	result = &rpc.ClientResults{}

	//链接到节点
	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", to, PORT), grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
		return
	}

	defer conn.Close()

	nodeclient := rpc.NewNodeClient(conn)

	// ctx, cancel := context.WithTimeout(context.Background(), MinTimeout*time.Millisecond)
	// defer cancel()

	ctx := metadata.NewOutgoingContext(context.Background(), metadata.New(md))

	result, err = nodeclient.ClientRequest(ctx, &rpc.ClientArguments{Data: cmd})

	return
}

func (c *Client) Command(cmd string, args ...string) (r string, err error) {

	switch cmd {
	case "ll":
		d, err := json.Marshal(c.Log)
		r = string(d)
		return fmt.Sprintf("%d:%s", len(c.Log), r), err
	case "lc":
		d, err := json.Marshal(c.NodesConfig)
		r = string(d)
		return fmt.Sprintf("%d:%s", len(c.NodesConfig), r), err
	case "if":
		return fmt.Sprintf("Term %v ,ROLE: %v, ID: %v, LeaderID: %v, commit %v", c.CurrentTerm, c.CurrentRole, c.ID, c.LeaderID, c.CommitIndex), err
	case CMD_JOIN:

		if len(args) > 1 {
			c.CurrentTerm = 0
			r, err := c.ClientRequestCall([]byte(CMD_JOIN), args[1], nil)
			if err != nil {
				return err.Error(), err
			}
			return []string{"失败", "OK"}[r.State], err
		} else {
			err = fmt.Errorf("提供一个地址参数")
			return
		}

	default:
		logrus.Info("default:", cmd)

		r, err := c.ClientRequestCall([]byte(cmd), c.LeaderID, nil)
		if err != nil {
			return err.Error(), err
		}
		return []string{"失败", "OK"}[r.State], err

	}

}
