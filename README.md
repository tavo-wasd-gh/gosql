# Sample usage

```go
type Pet struct {
    ID   int    `db:"pet_id"`
    Name string `db:"pet_name"`
}

var pets []Pet

rows := db.Query(`SELECT * FROM pets`)
if err := gosql.Scan(rows, pets) ; err != nil {
    log.Fatalf("Error")
}
```
