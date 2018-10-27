package main

import (
	"fmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"log"
)

// 创建通道客户端请求参数
type CreateChannelClient struct {
	ConfigFile string
	ChannelID  string
	UserName   string
	Org        string
}

// 内部请求周转
type InterRequest struct {
	ChaincodeID  string            // 链码ID
	Fcn          string            // 函数名称
	Args         [][]byte          // 参数
	TransientMap map[string][]byte // 隐私数据
	Peers        []string          // 请求peers

	InvocationChain []*fab.ChaincodeCall // 暂时用不到
}

type Client struct {
	configFile string
	channelID  string
	userName   string
	org        string

	client *channel.Client
	sdk    *fabsdk.FabricSDK
}

func NewClient(ccc CreateChannelClient) (Client, error) {
	client := Client{}
	/// 赋值
	client.configFile = ccc.ConfigFile
	client.channelID = ccc.ChannelID
	client.userName = ccc.UserName
	client.org = ccc.Org
	err := client.initialize()
	if err != nil {
		return client, err
	}
	return client, nil
}

func (setup *Client) initialize() error {
	// 通过配置文件初始化sdk
	sdk, err := fabsdk.New(config.FromFile(setup.configFile))
	if err != nil {
		return fmt.Errorf("failed to create SDK,Error=%s", err.Error())
	}

	setup.sdk = sdk
	log.Println("SDK created")

	// 创建通道客户端，用于query和invoke
	clientContext := setup.sdk.ChannelContext(setup.channelID, fabsdk.WithOrg(setup.org), fabsdk.WithUser(setup.userName))
	setup.client, err = channel.New(clientContext)
	if err != nil {
		return fmt.Errorf("Failed to create new channel client,Error=%s", err.Error())
	}
	log.Println("Channel client created")
	return nil
}

func (setup *Client) Query(req InterRequest) (channel.Response, error) {
	r := channel.Request{}
	r.Fcn = req.Fcn
	r.Args = req.Args
	r.ChaincodeID = req.ChaincodeID
	r.TransientMap = req.TransientMap

	// Create a request (proposal) and send it
	return setup.client.Query(r, channel.WithTargetEndpoints(req.Peers...))
}

func (setup *Client) Invoke(req InterRequest) (channel.Response, error) {
	r := channel.Request{}
	r.Fcn = req.Fcn
	r.Args = req.Args
	r.ChaincodeID = req.ChaincodeID

	return setup.client.Execute(r, channel.WithTargetEndpoints(req.Peers...))
}
