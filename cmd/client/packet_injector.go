/*
 * Copyright (C) 2016 Red Hat, Inc.
 *
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *  http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 *
 */

package client

import (
	"os"

	"github.com/skydive-project/skydive/api/client"
	api "github.com/skydive-project/skydive/api/types"
	"github.com/skydive-project/skydive/logging"
	"github.com/skydive-project/skydive/validator"

	"github.com/spf13/cobra"
)

var (
	srcNode    string
	dstNode    string
	srcIP      string
	srcMAC     string
	srcPort    int64
	dstPort    int64
	dstIP      string
	dstMAC     string
	packetType string
	payload    string
	id         int64
	count      int64
	interval   int64
	uuid       string
)

// PacketInjectorCmd skydive inject-packet root command
var PacketInjectorCmd = &cobra.Command{
	Use:          "inject-packet",
	Short:        "Inject packets",
	Long:         "Inject packets",
	SilenceUsage: false,
}

var PacketInjectionCreate = &cobra.Command{
	Use:          "create",
	Short:        "create packet injection",
	Long:         "create packet injection",
	SilenceUsage: false,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := client.NewCrudClientFromConfig(&AuthenticationOpts)
		if err != nil {
			logging.GetLogger().Critical(err.Error())
			os.Exit(1)
		}

		packet := &api.PacketParamsReq{
			Src:      srcNode,
			Dst:      dstNode,
			SrcIP:    srcIP,
			SrcMAC:   srcMAC,
			SrcPort:  srcPort,
			DstIP:    dstIP,
			DstMAC:   dstMAC,
			DstPort:  dstPort,
			Type:     packetType,
			Payload:  payload,
			ICMPID:   id,
			Count:    count,
			Interval: interval,
		}

		if err = validator.Validate(packet); err != nil {
			logging.GetLogger().Error(err.Error())
			os.Exit(1)
		}

		if err := client.Create("injectpacket", &packet); err != nil {
			logging.GetLogger().Error(err.Error())
			os.Exit(1)
		}

		printJSON(packet)
	},
}

var PacketInjectionGet = &cobra.Command{
	Use:   "get",
	Short: "get packet injection",
	Long:  "get packet injection",
	PreRun: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Usage()
			os.Exit(1)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		var injection api.PacketParamsReq
		client, err := client.NewCrudClientFromConfig(&AuthenticationOpts)
		if err != nil {
			logging.GetLogger().Critical(err.Error())
			os.Exit(1)
		}

		if err := client.Get("injectpacket", args[0], &injection); err != nil {
			logging.GetLogger().Error(err.Error())
			os.Exit(1)
		}
		printJSON(&injection)
	},
}

var PacketInjectionList = &cobra.Command{
	Use:   "list",
	Short: "list packet injections",
	Long:  "list packet injections",
	Run: func(cmd *cobra.Command, args []string) {
		var injections map[string]api.PacketParamsReq
		client, err := client.NewCrudClientFromConfig(&AuthenticationOpts)
		if err != nil {
			logging.GetLogger().Critical(err.Error())
			os.Exit(1)
		}

		if err := client.List("injectpacket", &injections); err != nil {
			logging.GetLogger().Error(err.Error())
			os.Exit(1)
		}
		printJSON(injections)
	},
}

var PacketInjectionDelete = &cobra.Command{
	Use:   "delete [injection]",
	Short: "Delete injection",
	Long:  "Delete packet injection",
	PreRun: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Usage()
			os.Exit(1)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		client, err := client.NewCrudClientFromConfig(&AuthenticationOpts)
		if err != nil {
			logging.GetLogger().Critical(err.Error())
			os.Exit(1)
		}

		if err := client.Delete("injectpacket", args[0]); err != nil {
			logging.GetLogger().Error(err.Error())
			os.Exit(1)
		}
	},
}

func addInjectPacketFlags(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&srcNode, "src", "", "", "source node gremlin expression (mandatory)")
	cmd.Flags().StringVarP(&dstNode, "dst", "", "", "destination node gremlin expression")
	cmd.Flags().StringVarP(&srcIP, "srcIP", "", "", "source node IP")
	cmd.Flags().StringVarP(&dstIP, "dstIP", "", "", "destination node IP")
	cmd.Flags().StringVarP(&srcMAC, "srcMAC", "", "", "source node MAC")
	cmd.Flags().StringVarP(&dstMAC, "dstMAC", "", "", "destination node MAC")
	cmd.Flags().Int64VarP(&srcPort, "srcPort", "", 0, "source port for TCP packet")
	cmd.Flags().Int64VarP(&dstPort, "dstPort", "", 0, "destination port for TCP packet")
	cmd.Flags().StringVarP(&packetType, "type", "", "icmp4", "packet type: icmp4, icmp6, tcp4, tcp6, udp4 and udp6")
	cmd.Flags().StringVarP(&payload, "payload", "", "", "payload")
	cmd.Flags().Int64VarP(&id, "id", "", 0, "ICMP identification")
	cmd.Flags().Int64VarP(&count, "count", "", 1, "number of packets to be generated")
	cmd.Flags().Int64VarP(&interval, "interval", "", 1000, "wait interval milliseconds between sending each packet")
}

func init() {
	PacketInjectorCmd.AddCommand(PacketInjectionList)
	PacketInjectorCmd.AddCommand(PacketInjectionGet)
	PacketInjectorCmd.AddCommand(PacketInjectionDelete)
	PacketInjectorCmd.AddCommand(PacketInjectionCreate)

	addInjectPacketFlags(PacketInjectionCreate)
}
