package entity

// Parameter is used when want to collect params from other table when scrapping
type ScrappingParameter struct {
  ID int64
  Url string
  Parameter string
}
