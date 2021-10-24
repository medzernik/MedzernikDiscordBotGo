package covid_slovakia

import (
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"os"
)

func COVIDOutputVaccinatedGraph(response VaccinatedSlovakiaResponse) {
	// create a new bar instance
	bar := charts.NewBar()

	bar.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{Title: "toolbox options"}),
		charts.WithToolboxOpts(opts.Toolbox{
			Show:  true,
			Right: "20%",
			Feature: &opts.ToolBoxFeature{
				SaveAsImage: &opts.ToolBoxFeatureSaveAsImage{
					Show:  true,
					Type:  "png",
					Title: "Anything you want",
				},
				DataView: &opts.ToolBoxFeatureDataView{
					Show:  true,
					Title: "DataView",
					// set the language
					// Chinese version: ["数据视图", "关闭", "刷新"]
					Lang: []string{"data view", "turn off", "refresh"},
				},
			}},
		),
	)

	// set some global options like Title/Legend/ToolTip or anything else
	bar.SetGlobalOptions(charts.WithTitleOpts(opts.Title{
		Title:    "Number of vaccinated people in slovakia for the past 7 days",
		Subtitle: "Blue - first dosage; Green - second dosage",
	}))

	// Put data into instance
	bar.SetXAxis([]string{response.Page[0].PublishedOn, response.Page[1].PublishedOn, response.Page[2].PublishedOn, response.Page[3].PublishedOn, response.Page[4].PublishedOn, response.Page[5].PublishedOn, response.Page[6].PublishedOn}).
		AddSeries("Prva davka", generateBarItems1(response)).
		AddSeries("Druha davka", generateBarItems2(response))
	// Where the magic happens
	f, _ := os.Create("bar.html")
	bar.Render(f)

	return
}

// generate random data for bar chart
func generateBarItems1(response VaccinatedSlovakiaResponse) []opts.BarData {
	items := make([]opts.BarData, 0)
	for i := 0; i < 7; i++ {
		items = append(items, opts.BarData{Name: "Prva davka", Value: response.Page[i].Dose1Count})
	}
	return items
}

// generate random data for bar chart2
func generateBarItems2(response VaccinatedSlovakiaResponse) []opts.BarData {
	items := make([]opts.BarData, 0)
	for i := 0; i < 7; i++ {
		items = append(items, opts.BarData{Name: "Druha davka", Value: response.Page[i].Dose2Count})
	}
	return items
}
