package main

import (
	"context"
	"github.com/mhchlib/mconfig-api/api/v1/sdk"
	grpc "google.golang.org/grpc"
	"log"
	"strconv"
)

func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial("client.u.hcyang.top:31790", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := sdk.NewMConfigClient(conn)
	appId := 1003
	configIds := []string{"1003-100", "1003-103"}
	extreData := map[string]string{
		"ip": "192.168.1.1",
	}
	log.Println("client listen app ", appId, " config ", configIds, " with data ", extreData)
	resp, err := client.GetVStream(context.Background())
	resp.SendMsg(&sdk.GetVRequest{
		AppId: strconv.Itoa(appId),
		Filters: &sdk.ConfigFilters{
			ConfigIds: configIds,
			ExtraData: extreData,
		},
	})
	if err != nil {
		log.Fatal(err)
		return
	}
	defer func() {
		log.Println("close stream")
	}()
	for {
		config := sdk.GetVResponse{}
		//config, err := resp.Recv()
		err := resp.RecvMsg(&config)
		//_, err := resp.Recv()
		if err != nil {
			log.Fatal("err: ", err)
			return
		}
		log.Println(appId, " get msg")
		log.Println(" ------------------- ")
		log.Println(config.Configs)
	}
}
