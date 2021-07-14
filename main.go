package main

import (
  "github.com/sonlh-0262/scrapping_github/database"
  "github.com/sonlh-0262/scrapping_github/entity"
  "log"
  "fmt"
  "strconv"
  "github.com/go-rod/rod"
)

type PageStruct struct {
  Page *rod.Page
}

func main() {
  // Fetch Scrapping parameter: it includes URl and parameter (if required)
  scrappingParameters, err := database.FetchAllScrappingParameters()
  if err != nil {
    log.Fatal(err)
  }
  fmt.Println("parameters: %v", scrappingParameters)

  // Fetch old Github data
  githubs, err := database.FetchAllGithubs()
  if err != nil {
    log.Fatal(err)
  }
  fmt.Println("All Githubs: %v", githubs)

  // Make concurrency scrapping
  channel := make(chan entity.Github)

  for i := 0; i < len(scrappingParameters); i++ {
    go scrappingPage(scrappingParameters[i].Url, channel)
  }

  // Add data to DB
  for github := range channel {
    fmt.Println(github)
    database.AddGithubDB(github)
  }
}

func scrappingPage(url string, channel chan entity.Github) {
  page := rod.New().MustConnect().MustPage(url).MustWindowFullscreen()
  site := PageStruct{
    Page: page,
  }

  github := entity.Github{
    Owner: site.getOwner(),
    Name: site.getName(),
    Star: site.getStar(),
    Fork: site.getFork(),
    BranchCount: site.getBranchCount(),
    TagCount: site.getTagCount(),
  }

  channel <- github
}

func (site PageStruct) getStar() string {
  return site.getText("#js-repo-pjax-container > div.hx_page-header-bg > div.d-flex > ul > li:nth-child(2) > a.social-count")
}

func (site PageStruct) getFork() string {
  return site.getText("#js-repo-pjax-container > div.hx_page-header-bg > div.d-flex > ul > li:nth-child(3) > a.social-count")
}

func (site PageStruct) getOwner() string {
  return site.getText("#js-repo-pjax-container > div.hx_page-header-bg > div.d-flex > div.flex-auto > h1 > span > a")
}

func (site PageStruct) getName() string {
  return site.getText("#js-repo-pjax-container > div.hx_page-header-bg > div.d-flex > div.flex-auto > h1 > strong > a")
}

func (site PageStruct) getBranchCount() int {
  branchCountText := site.getText("#repo-content-pjax-container > div > div.gutter-condensed > div > div.file-navigation > div.flex-self-center > a > strong")
  branchCount, _ := strconv.Atoi(branchCountText)

  return branchCount
}

func (site PageStruct) getTagCount() int {
  tagCountText := site.getText("#repo-content-pjax-container > div > div.gutter-condensed > div > div.file-navigation > div.flex-self-center > a:nth-child(2) > strong")
  tagCount, _ := strconv.Atoi(tagCountText)

  return tagCount
}

func (site PageStruct) getText(cssPath string) string {
  return site.Page.MustElement(cssPath).MustText()
}
