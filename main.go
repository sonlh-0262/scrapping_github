package main

import (
  "github.com/sonlh-0262/scrapping_github/database"
  "github.com/sonlh-0262/scrapping_github/entity"
  "log"
  "fmt"
  "strconv"
  "github.com/go-rod/rod"
)

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

  github := entity.Github{
    Owner: getOwner(page),
    Name: getName(page),
    Star: getStar(page),
    Fork: getFork(page),
    BranchCount: getBranchCount(page),
    TagCount: getTagCount(page),
  }

  channel <- github
}

func getStar(page *rod.Page) string {
  return getText(page, "#js-repo-pjax-container > div.hx_page-header-bg > div.d-flex > ul > li:nth-child(2) > a.social-count")
}

func getFork(page *rod.Page) string {
  return getText(page, "#js-repo-pjax-container > div.hx_page-header-bg > div.d-flex > ul > li:nth-child(3) > a.social-count")
}

func getOwner(page *rod.Page) string {
  return getText(page, "#js-repo-pjax-container > div.hx_page-header-bg > div.d-flex > div.flex-auto > h1 > span > a")
}

func getName(page *rod.Page) string {
  return getText(page, "#js-repo-pjax-container > div.hx_page-header-bg > div.d-flex > div.flex-auto > h1 > strong > a")
}

func getBranchCount(page *rod.Page) int {
  branchCountText := getText(page, "#repo-content-pjax-container > div > div.gutter-condensed > div > div.file-navigation > div.flex-self-center > a > strong")
  branchCount, _ := strconv.Atoi(branchCountText)

  return branchCount
}

func getTagCount(page *rod.Page) int {
  tagCountText := getText(page, "#repo-content-pjax-container > div > div.gutter-condensed > div > div.file-navigation > div.flex-self-center > a:nth-child(2) > strong")
  tagCount, _ := strconv.Atoi(tagCountText)

  return tagCount
}

func getText(page *rod.Page, cssPath string) string {
  return page.MustElement(cssPath).MustText()
}
