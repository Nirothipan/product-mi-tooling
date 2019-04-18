/*
* Copyright (c) 2019, WSO2 Inc. (http://www.wso2.org) All Rights Reserved.
*
* WSO2 Inc. licenses this file to you under the Apache License,
* Version 2.0 (the "License"); you may not use this file except
* in compliance with the License.
* You may obtain a copy of the License at
*
*    http://www.apache.org/licenses/LICENSE-2.0
*
* Unless required by applicable law or agreed to in writing,
* software distributed under the License is distributed on an
* "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
* KIND, either express or implied. See the License for the
* specific language governing permissions and limitations
* under the License.
 */

package cmd

import (
	"github.com/lithammer/dedent"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"github.com/wso2/micro-integrator/cmd/utils"
	"os"
)

var inboundEndpointName string

// Show InboundEndpoint command related usage info
const showInboundEndpointCmdLiteral = "inboundEndpoint"
const showInboundEndpointCmdShortDesc = "Get information about the specified Inbound Endpoint"

var showInboundEndpointCmdLongDesc = "Get information about the InboundEndpoint specified by the flag --name, -n\n"

var showInboundEndpointCmdExamples = dedent.Dedent(`
Example:
  ` + utils.ProjectName + ` ` + showCmdLiteral + ` ` + showInboundEndpointCmdLiteral + ` -n TestInboundEndpoint
`)

// InboundEndpointShowCmd represents the Show inboundEndpoint command
var inboundEndpointShowCmd = &cobra.Command{
	Use:   showInboundEndpointCmdLiteral,
	Short: showInboundEndpointCmdShortDesc,
	Long:  showInboundEndpointCmdLongDesc + showInboundEndpointCmdExamples,
	Run: func(cmd *cobra.Command, args []string) {
		utils.Logln(utils.LogPrefixInfo + "Show InboundEndpoint called")
		executeGetInboundEndpointCmd(inboundEndpointName)
	},
}

func init() {
	showCmd.AddCommand(inboundEndpointShowCmd)

	inboundEndpointShowCmd.Flags().StringVarP(&inboundEndpointName, "name", "n", "", "Name of the Inbound Endpoint")
	inboundEndpointShowCmd.MarkFlagRequired("name")
}

func executeGetInboundEndpointCmd(inboundEndpointname string) {

	finalUrl := utils.RESTAPIBase + utils.PrefixInboundEndpoints + "?inboundEndpointName=" + inboundEndpointname

	resp, err := utils.UnmarshalData(finalUrl, &utils.InboundEndpoint{})

	if err == nil {
		// Printing the details of the InboundEndpoint
		inboundEndpoint := resp.(*utils.InboundEndpoint)
		printInboundEndpoint(*inboundEndpoint)
	} else {
		utils.Logln(utils.LogPrefixError+"Getting Information of InboundEndpoint", err)
	}
}

// printInboundEndpointInfo
// @param InboundEndpoint : InboundEndpoint object
func printInboundEndpoint(endpoint utils.InboundEndpoint) {
	table := tablewriter.NewWriter(os.Stdout)

	row := []string{"NAME", "", endpoint.Name}
	table.Append(row)

	row = []string{"PROTOCOL", "", endpoint.Protocol}
	table.Append(row)

	row = []string{"CLASS", "", endpoint.Class}
	table.Append(row)

	row = []string{"SEQUENCE", "", endpoint.Sequence}
	table.Append(row)

	row = []string{"ERROR SEQUENCE", "", endpoint.ErrorSequence}
	table.Append(row)

	for _, param := range endpoint.Parameters {
		row = []string{"PARAMETERS", param.Name, param.Value}
		table.Append(row)
	}

	table.SetBorders(tablewriter.Border{Left: true, Top: true, Right: true, Bottom: false})
	table.SetRowLine(true)
	table.SetAutoMergeCells(true)
	table.Render() // Send output
}