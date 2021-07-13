package main

import (
  "database/sql"
  "log"
  "fmt"
  "strconv"
  "github.com/go-rod/rod"
  "github.com/go-sql-driver/mysql"
)

type Github struct {
  ID int64
  Owner string
  Name string
  Star string
  Fork string
  BranchCount int
  TagCount int
}

type ScrappingParameter struct {
  ID int64
  Url string
  Parameter string
}

var db *sql.DB

func main() {
  connectDB()

  scrappingParameters, err := fetchAllScrappingParameters()
  if err != nil {
    log.Fatal(err)
  }

  for i := 0; i < len(scrappingParameters); i++ {
    github := scrappingPage(scrappingParameters[i].Url)
    fmt.Println(github)
    addGithubDB(github)
  }

  githubs, err := fetchAllGithubs()
  if err != nil {
    log.Fatal(err)
  }
  fmt.Println("All Githubs: %v", githubs)
}

func fetchAllScrappingParameters() ([]ScrappingParameter, error) {
  var scrappingParameters []ScrappingParameter

  rows, err := db.Query("SELECT * from scrapping_parameters")

  if err != nil {
    return nil, fmt.Errorf("Query all with err %v", err)
  }

  defer rows.Close()

  for rows.Next() {
    var sp ScrappingParameter
    if err := rows.Scan(&sp.ID, &sp.Url, &sp.Parameter); err != nil {
      return nil, fmt.Errorf("Error %v", err)
    }

    scrappingParameters = append(scrappingParameters, sp)
  }

  if err := rows.Err(); err != nil {
    return nil, fmt.Errorf("Error %v", err)
  }
  return scrappingParameters, nil
}

func addGithubDB(github Github) (int64, error) {
  result, err := db.Exec(
    "INSERT INTO github (owner, name, star, fork, branch_count, tag_count) VALUES (?, ?, ?, ?, ?, ?)",
    github.Owner, github.Name, github.Star, github.Fork, github.BranchCount, github.TagCount,
  )

  if err != nil {
    return 0, fmt.Errorf("addGithub: %v", err)
  }

  id, err := result.LastInsertId()
  if err != nil {
    return 0, fmt.Errorf("addGithub: %v", err)
  }

  return id, nil
}

func fetchAllGithubs() ([]Github, error) {
  var githubs []Github

  rows, err := db.Query("SELECT * FROM github")

  if err != nil {
    return nil, fmt.Errorf("Query all with err %v", err)
  }

  defer rows.Close()

  for rows.Next() {
    var git Github
    if err := rows.Scan(&git.ID, &git.Owner, &git.Name, &git.Star, &git.Fork, &git.BranchCount, &git.TagCount); err != nil {
      return nil, fmt.Errorf("Error %v", err)
    }

    githubs = append(githubs, git)
  }

  if err := rows.Err(); err != nil {
    return nil, fmt.Errorf("Error %v", err)
  }
  return githubs, nil
}

func connectDB() {
  config := mysql.Config{
    User: "root",
    Passwd: "",
    Net: "tcp",
    Addr: "localhost:3306",
    DBName: "scrapping_dev",
    AllowNativePasswords: true,
  }

  var err error
  db, err = sql.Open("mysql", config.FormatDSN())
  if err != nil {
    log.Fatal(err)
  }

  pingErr := db.Ping()
  if pingErr != nil {
    log.Fatal(pingErr)
  }

  fmt.Println("Connected!")
}

func scrappingPage(url string) Github {
  page := rod.New().MustConnect().MustPage(url).MustWindowFullscreen()

  return Github{
    Owner: getOwner(page),
    Name: getName(page),
    Star: getStar(page),
    Fork: getFork(page),
    BranchCount: getBranchCount(page),
    TagCount: getTagCount(page),
  }
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
