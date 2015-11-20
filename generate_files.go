package gasegment

//go:generate curl https://www.googleapis.com/analytics/v3/metadata/ga/columns -o files/columns.json
//go:generate go-bindata -ignore="\.DS_Store" -o asset/asset.go -pkg asset -prefix files files/...
