package database

import (
  "database/sql"
  "log"
  "fmt"
  "github.com/go-sql-driver/mysql"
  "github.com/sonlh-0262/scrapping_github/entity"
)

var db *sql.DB

func init() {
  connectDB()
}

func FetchAllScrappingParameters() ([]entity.ScrappingParameter, error) {
  var scrappingParameters []entity.ScrappingParameter

  rows, err := db.Query("SELECT * from scrapping_parameters")

  if err != nil {
    return nil, fmt.Errorf("Query all with err %v", err)
  }

  defer rows.Close()

  for rows.Next() {
    var sp entity.ScrappingParameter
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

func AddGithubDB(github entity.Github) (int64, error) {
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

func FetchAllGithubs() ([]entity.Github, error) {
  var githubs []entity.Github

  rows, err := db.Query("SELECT * FROM github")

  if err != nil {
    return nil, fmt.Errorf("Query all with err %v", err)
  }

  defer rows.Close()

  for rows.Next() {
    var git entity.Github
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
