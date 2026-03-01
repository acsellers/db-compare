package main

var samples = []string{
	"bob",
	"bob2",
	"bun",
	"ent",
	"gorm",
	"jet",
	"sqlboiler",
	"sqlc",
	"sqlx",
	"upper db",
	"xorm",
}

type Samples struct {
	Store StoreSamples `json:"store"`
}
type StoreSamples struct {
	MySQL    StoreDBSamples `json:"mysql"`
	Postgres StoreDBSamples `json:"postgres"`
	SQLite   StoreDBSamples `json:"sqlite"`
}
type StoreDBSamples struct {
	GetSale          SampleRefs `json:"get_sale"`
	CreateSale       SampleRefs `json:"create_sale"`
	CustomerUpdate   SubSamples `json:"customer_update"`
	BasicGrouping    SubSamples `json:"basic_grouping"`
	AdvancedGrouping SubSamples `json:"advanced_grouping"`
	WithQueries      SubSamples `json:"with_queries"`
	SaleSearch       SampleRefs `json:"sale_search"`
	BulkCustomers    SampleRefs `json:"bulk_customers"`
	JSON             SampleRefs `json:"json"`
}
type SampleRefs struct {
	Go  SampleFile `json:"go"`
	Sql SampleFile `json:"sql"`
}
type SampleFile struct {
	File       string   `json:"file"`
	Lines      [][2]int `json:"lines"`
	Highlights [][2]int `json:"highlights"`
}
type SubSamples map[string]SampleFile

func main() {

}
