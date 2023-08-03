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

// topDbCmd represents the topDb command
var topDbCmd = &cobra.Command{
	Use:   "top-db",
	Short: "List the top DB methods on calls and time consumed",
	Long:  `Generate a list of the top DB methods based on the criteria provided.`,
	Run:   topDbHandler,
}

func init() {
	rootCmd.AddCommand(topDbCmd)
	topDbCmd.Flags().StringP("criteria", "c", "total-time", "Set the criteria to get the top DB methods (total-time, average-time, count)")
	topDbCmd.Flags().IntP("limit", "l", 20, "Set the top limit of the results")
}

func topDbHandler(cmd *cobra.Command, args []string) {
	a := app.New("http://localhost:9090")
	data, err := a.GetDBMetrics()
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

	tbl := table.New("Method", "Total time", "Count", "Average time")
	tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)

	dataList := make([]*app.DBEntry, 0, len(data))
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
		tbl.AddRow(d.Method, d.TotalTime, d.Count, d.Average)
	}
	tbl.Print()
}
