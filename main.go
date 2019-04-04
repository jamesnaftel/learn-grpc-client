package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	pb "github.com/jamesnaftel/learn-grpc/api"
	"google.golang.org/grpc"
)

func main() {
	host := flag.String("host", "localhost", "Server host")
	port := flag.String("port", "3001", "Server port")
	flag.Parse()

	//todo: change this code to be CLI driven once done tinkering
	conn, err := grpc.Dial(fmt.Sprintf("%s:%s", *host, *port), grpc.WithInsecure())
	if err != nil {
		fmt.Printf("Error dialing %s: %v\n", fmt.Sprintf("%s:%s", *host, *port), err)
		os.Exit(1)
	}
	defer conn.Close()

	c := pb.NewPodcastsClient(conn)
	ctx := context.Background()

	//Request a specific podcast - return empty
	req := pb.PodcastRequest{Name: "SE Daily"}
	podcast, err := c.GetPodcast(ctx, &req)
	if err != nil {
		fmt.Printf("%v\n", err)
	}

	fmt.Printf("GetPodcast: %s\n", podcast.Podcast.GetName())

	//Request a specific podcast
	req = pb.PodcastRequest{Name: "SE Daily: GRPC"}
	podcast, err = c.GetPodcast(ctx, &req)
	if err != nil {
		fmt.Printf("%v\n", err)
	}

	fmt.Printf("GetPodcast: %s\n", podcast.Podcast.String())

	empty := pb.Empty{}
	podcasts, err := c.GetPodcasts(ctx, &empty)
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}

	fmt.Printf("GetPodcasts: %v\n", podcasts)

}
