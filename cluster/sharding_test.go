package cluster

import (
	"testing"
	"time"

	"github.com/azhai/gozzo-db/fixture"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
)

type Person struct {
	ID        int
	Name      string
	Height    float64
	BirthDate time.Time `gorm:"column:birth"`
}

func (Person) TableName() string {
	return fixture.TestTableName
}

func (p Person) BaseTableName() string {
	return p.TableName() + "_"
}

func CreateRecords(db *gorm.DB) *gorm.DB {
	query := db.Table(fixture.TestTableName + "_males")
	query.Create(&Person{Name: "Bob", Height: 178, BirthDate: fixture.GetDate("1982-02-28")})
	query.Create(&Person{Name: "David", Height: 181, BirthDate: fixture.GetDate("1984-04-01")})
	query.Create(&Person{Name: "Frank", Height: 169, BirthDate: fixture.GetDate("1986-06-09")})
	query = db.Table(fixture.TestTableName + "_females")
	query.Create(&Person{Name: "Alice", Height: 168, BirthDate: fixture.GetDate("1981-01-08")})
	query.Create(&Person{Name: "Candy", Height: 165, BirthDate: fixture.GetDate("1983-03-15")})
	query.Create(&Person{Name: "Emily", Height: 171, BirthDate: fixture.GetDate("1985-05-01")})
	query.Create(&Person{Name: "Grace", Height: 175, BirthDate: fixture.GetDate("1987-07-05")})
	return db
}

func Test_CountSharding(t *testing.T) {
	db := fixture.InitDB()
	CreateRecords(fixture.TruncateRecords(db))
	query := NewSharding(db, true)
	count := query.CountSharding(Person{})
	assert.Equal(t, int64(7), count)
}

func Test_PaginateSharding(t *testing.T) {
	db := fixture.InitDB()
	CreateRecords(fixture.TruncateRecords(db))
	shr := NewSharding(db.Order("Name"), true)
	var result []Person
	fetch := func(query *gorm.DB) *gorm.DB {
		var rows []Person
		query.Find(&rows)
		for _, row := range rows {
			result = append(result, row)
		}
		return query
	}
	err := shr.PaginateSharding(Person{}, 1, 5, fetch)
	assert.NoError(t, err, "出错了")
	assert.Equal(t, 5, len(result))
	assert.Equal(t, "Bob", result[0].Name)
	assert.Equal(t, float64(181), result[1].Height)
	assert.Equal(t, "1986-06-09", result[2].BirthDate.Format("2006-01-02"))
	assert.Equal(t, "Alice", result[3].Name)
	assert.Equal(t, 2, result[4].ID)
}
