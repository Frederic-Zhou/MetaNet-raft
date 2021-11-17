package main

import (
	"fmt"
	"metanet/lanscan"
	"testing"
)

func TestLanscan(t *testing.T) {
	// Scan for hosts listening on tcp port 80.
	// Use 20 threads and timeout after 5 seconds.

	for _, current := range lanscan.LinkLocalAddresses("tcp4") {
		allIPs := lanscan.CalculateSubnetIPs(current)
		fmt.Println(allIPs)
	}
	// ctx, cancel := context.WithTimeout(context.Background(), 1000*time.Millisecond)
	// defer cancel()

	// conn, err := grpc.DialContext(ctx, fmt.Sprintf("%s:%s", "192.168.3.82", "8800"), grpc.WithBlock(), grpc.WithInsecure())
	// if err != nil {
	// 	fmt.Println(err.Error())
	// 	return
	// }
	// defer conn.Close()
}
