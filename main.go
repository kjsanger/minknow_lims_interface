package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"minknow_lims_interface/minknow/rpc/device"
	"minknow_lims_interface/minknow/rpc/instance"
	"minknow_lims_interface/minknow/rpc/manager"
	"minknow_lims_interface/minknow/rpc/protocol"
)

func main() {
	fqdn := "xxxxxx.sanger.ac.uk"
	port := 9501

	addr := fmt.Sprintf("%s:%d", fqdn, port)

	func() {
		conn, err := grpc.Dial(addr, grpc.WithInsecure())
		if err != nil {
			log.Fatalf("failed to connect: %s", err)
		}

		defer conn.Close()

		client := manager.NewManagerServiceClient(conn)
		response1, err := client.ListDevices(context.Background(), &manager.ListDevicesRequest{})
		if err != nil {
			log.Printf("failed to list devices: %s", err)
		} else {
			for _, d := range response1.GetActive() {
				fmt.Println(fmt.Sprintf("ACTIVE: name: %s", d.GetName()))

				dport := d.GetPorts().InsecureGrpc
				daddr := fmt.Sprintf("%s:%d", fqdn, dport)
				dconn, derr := grpc.Dial(daddr, grpc.WithInsecure())
				if derr != nil {
					log.Printf("failed to connect to device: %s on port %d", err, dport)
				}
				defer dconn.Close()
				log.Printf("connected to device: %s on port %d", d, dport)

				dclient := device.NewDeviceServiceClient(dconn)
				response, err := dclient.GetDeviceInfo(context.Background(), &device.GetDeviceInfoRequest{})
				if err != nil {
					log.Printf("failed to list device info: %s", err)
				} else {
					fmt.Println(fmt.Sprintf("device ID: %s", response.GetDeviceId()))
				}

				response2, err := dclient.GetFlowCellInfo(context.Background(), &device.GetFlowCellInfoRequest{})
				if err != nil {
					log.Printf("failed to list flowcell info: %s", err)
				} else {
					fmt.Println(fmt.Sprintf("flowcell ID: %s", response2.GetFlowCellId()))
				}

				pclient := protocol.NewProtocolServiceClient(dconn)
				response3, err := pclient.GetRunInfo(context.Background(), &protocol.GetRunInfoRequest{})
				if err != nil {
					log.Printf("failed to get run info: %s", err)
				} else {
					fmt.Println(fmt.Sprintf("run ID: %s", response3.GetRunId()))
					fmt.Println(fmt.Sprintf("output path: %s", response3.GetOutputPath()))
					fmt.Println(fmt.Sprintf("flowcell: %s", response3.GetFlowCell()))
					fmt.Println(fmt.Sprintf("user info: %v", response3.GetUserInfo()))
				}





			}
		}

		response2, err := client.GetVersionInfo(context.Background(), &manager.GetVersionInfoRequest{})
		if err != nil {
			log.Printf("failed to get version info: %s", err)
		} else {
			fmt.Println(fmt.Sprintf("MinKNOW version: %s", response2.GetMinknow().GetFull()))
			fmt.Println(fmt.Sprintf("Protocols version: %s", response2.GetProtocols()))
			fmt.Println(fmt.Sprintf("Distribution version %s", response2.GetDistributionVersion()))
		}
	}()

	func() {
		conn, err := grpc.Dial(addr, grpc.WithInsecure())
		if err != nil {
			log.Fatalf("failed to connect: %s", err)
		}

		defer conn.Close()

		client := instance.NewInstanceServiceClient(conn)

		response1, err := client.GetDiskSpaceInfo(context.Background(), &instance.GetDiskSpaceInfoRequest{})
		if err != nil {
			log.Printf("failed to get disk space info: %s", err)
		} else {
			for _, di := range response1.GetFilesystemDiskSpaceInfo() {
				fmt.Println(fmt.Sprintf("disk space: %d", di.GetBytesAvailable()))
			}
		}

		response2, err := client.GetOutputDirectories(context.Background(), &instance.GetOutputDirectoriesRequest{})
		if err != nil {
			log.Printf("failed to list output directories: %s", err)
		} else {
			for _, d := range response2.GetOutput() {
				fmt.Println(fmt.Sprintf("output directory: %s", d))
			}
		}

		response3, err := client.GetDefaultOutputDirectories(context.Background(), &instance.GetDefaultOutputDirectoriesRequest{})
		if err != nil {
			log.Printf("failed to get default output directory: %s", err)
		} else {
			for _, d := range response3.GetOutput() {
				fmt.Println(fmt.Sprintf("default output directory: %s", d))
			}
		}

	}()

	func() {
		conn, err := grpc.Dial(addr, grpc.WithInsecure())
		if err != nil {
			log.Fatalf("failed to connect: %s", err)
		}

		defer conn.Close()

		client := device.NewDeviceServiceClient(conn)
		response1, err := client.GetDeviceInfo(context.Background(), &device.GetDeviceInfoRequest{})
		if err != nil {
			log.Println(fmt.Sprintf("failed to list device info: %s", err))
		} else {
			fmt.Println(fmt.Sprintf("device ID: %s", response1.GetDeviceId()))
		}

		response2, err := client.GetFlowCellInfo(context.Background(), &device.GetFlowCellInfoRequest{})
		if err != nil {
			log.Printf("failed to list flowcell info: %s", err)
		} else {
			fmt.Println(fmt.Sprintf("flowcell ID: %s, channel count: %d",
				response2.GetFlowCellId(),
				response2.GetChannelCount()))
		}
	}()

	func() {
		conn, err := grpc.Dial(addr, grpc.WithInsecure())
		if err != nil {
			log.Fatalf("failed to connect: %s", err)
		}

		defer conn.Close()

		client := protocol.NewProtocolServiceClient(conn)
		response1, err := client.GetRunInfo(context.Background(), &protocol.GetRunInfoRequest{})
		if err != nil {
			log.Printf("failed to list run info: %s", err)
		} else {
			fmt.Println(fmt.Sprintf("run ID: %s", response1.GetRunId()))
		}

	}()



}
