/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"sort"

	"github.com/fatih/color"
	"github.com/mattermost/mattermost-perf-stats-cli/app"
	"github.com/rodaine/table"
	"github.com/spf13/cobra"
)

// topApiCmd represents the topApi command
var topApiCmd = &cobra.Command{
	Use:   "top-api",
	Short: "List the top API endpoints on calls and time consumed",
	Long:  `Generate a list of the top API endpoints based on the criteria provided.`,
	Run:   topApiHandler,
}

func init() {
	rootCmd.AddCommand(topApiCmd)
	topApiCmd.Flags().StringP("criteria", "c", "total-time", "Set the criteria to get the top API endpoints (total-time, average-time, count)")
	topApiCmd.Flags().IntP("limit", "l", 20, "Set the top limit of the results")
}

func topApiHandler(cmd *cobra.Command, args []string) {
	a := app.New("http://localhost:9090")
	data, err := a.GetAPIMetrics()
	if err != nil {
		panic(err)
	}
	criteria, err := cmd.Flags().GetString("criteria")
	if err != nil {
		panic(err)
	}
	limit, err := cmd.Flags().GetInt("limit")
	if err != nil {
		panic(err)
	}

	if criteria != "total-time" && criteria != "average-time" && criteria != "count" {
		fmt.Println("Invalid criteria")
		return
	}

	headerFmt := color.New(color.FgGreen, color.Underline).SprintfFunc()
	columnFmt := color.New(color.FgYellow).SprintfFunc()

	tbl := table.New("Handler", "Total time", "Count", "Average time")
	tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)

	dataList := make([]*app.APIEntry, 0, len(data))
	for _, d := range data {
		dataList = append(dataList, d)
	}

	sort.Slice(dataList, func(i, j int) bool {
		if criteria == "total-time" {
			return dataList[i].TotalTime > dataList[j].TotalTime
		} else if criteria == "average-time" {
			return dataList[i].Average > dataList[j].Average
		} else if criteria == "count" {
			return dataList[i].Count > dataList[j].Count
		}
		panic("unreachable code")
	})

	if limit > len(dataList) {
		limit = len(dataList)
	}
	for _, d := range dataList[0:limit] {
		tbl.AddRow(d.Handler, d.TotalTime, d.Count, d.Average)
	}
	tbl.Print()
}
