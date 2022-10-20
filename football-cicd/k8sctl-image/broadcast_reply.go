package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"joynova.com/greenly/greenly-master/src/services/reply/api_reply"
	"joynova.com/joynova/joymicro/registry"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func newK8sClient(configPath string) (*kubernetes.Clientset, error) {
	RestConfig, err := clientcmd.BuildConfigFromFlags("", configPath)
	if err != nil {
		errStr := fmt.Sprintf("Build Config FromFlags is err:%v", err)
		fmt.Fprintln(os.Stderr, errStr)
		os.Exit(1)
		return nil, err
	}

	ip := strings.Split(RestConfig.Host[8:], ":")[0]
	fmt.Printf("host ip:%v\n", ip)

	// clientSet := kube.ClientSet(KubernetesConfigFlags)
	clientSet, err := kubernetes.NewForConfig(RestConfig)
	if err != nil {
		errStr := fmt.Sprintf("new restclient from rest config is err:%v", err)
		fmt.Fprintln(os.Stderr, errStr)
		return nil, err
	}

	return clientSet, nil
}

func main() {
	namespace := flag.String("ns", "", "k8s命名空间，trunk|release|...")
	msg := flag.String("msg", "no msg", "发给reply的消息")
	imageVersion := flag.String("version", "", "镜像版本，要和被通知的战斗服匹配")

	etcd := flag.String("etcd", "192.168.1.22:2379", "etcd地址")
	flag.Parse()

	fmt.Printf("ns:%v, etcd:%v, msg:%v, version:%v\n", *namespace, *etcd, *msg, *imageVersion)

	registry.SetNameSpace("/football/football-" + *namespace)
	replyClient := api_reply.NewReplyServiceInstance([]string{*etcd}, 20*time.Second, true, false)
	// selector := new(ReplyAddrsCollector)
	// selector.replys = new(sync.Map)
	// api_reply.SetReplyServiceSelector(replyClient, selector)
	ctx := context.WithValue(context.Background(), "to_all_nodes", true)
	_, err := replyClient.BroadCast(ctx, &api_reply.BroadCastReq{Msg: *msg, ImageVersion: *imageVersion})
	if err != nil {
		fmt.Printf("广播reply返回错误：%v\n", err)
		os.Exit(1)
	}
	fmt.Printf("通知命名空间%v下所有版本为%v的战斗服消息%v成功！\n", *namespace, *imageVersion, *msg)
}

//
// // PeerSelector 点对点选择器
// type ReplyAddrsCollector struct {
// 	servers []string
// 	// replys  map[string][2]string
// 	replys *sync.Map // 所有战斗服的公网ip，用来客户端登录时候获取到地址做延迟计算
// }
//
// // Select 根据context里的select_key选择匹配的服务器进行调用
// func (ms *ReplyAddrsCollector) Select(ctx context.Context, servicePath, serviceMethod string, args interface{}) string {
// 	if len(ms.servers) <= 0 {
// 		return ""
// 	}
//
// 	return ms.servers[rand.Intn(len(ms.servers))]
// }
//
// // UpdateServer 更新服务器
// func (ms *ReplyAddrsCollector) UpdateServer(servers map[string]string) {
// 	ss := make([]string, 0, len(servers))
//
// 	// 清空所有key
// 	ms.replys.Range(func(k, v interface{}) bool {
// 		ms.replys.Delete(k)
// 		return true
// 	})
//
// 	for k, v := range servers {
// 		values, err := url.ParseQuery(v)
// 		if err != nil {
// 			panic(err)
// 		}
//
// 		publicTcpAddr := values.Get("public_tcp_ip")
// 		publicUdpAddr := values.Get("public_udp_ip")
// 		state := values.Get("state")
//
// 		if state == "stop" {
// 			fmt.Printf("server:%v, state stop:%v\n", k, values.Encode())
// 			continue
// 		}
//
// 		addrs := [2]string{publicTcpAddr, publicUdpAddr}
//
// 		ms.replys.Store(k, addrs)
//
// 		ss = append(ss, k)
// 	}
//
// 	ms.servers = ss
// }
