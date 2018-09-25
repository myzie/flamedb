package database

import (
	"errors"
	"fmt"
	"path"
	"strings"

	gorm "github.com/jinzhu/gorm"
)

type flame struct {
	gormDB *gorm.DB
}

// NewFlame returns an interface to a FlameDB
func NewFlame(gormDB *gorm.DB) Flame {
	return &flame{gormDB: gormDB}
}

func (f *flame) DropTable() error {
	if err := f.gormDB.DropTableIfExists(&Record{}).Error; err != nil {
		return fmt.Errorf("Failed to drop table: %s", err.Error())
	}
	return nil
}

func (f *flame) Migrate() error {
	if err := f.gormDB.AutoMigrate(&Record{}).Error; err != nil {
		return fmt.Errorf("Failed to migrate database: %s", err.Error())
	}
	return nil
}

func (f *flame) getPostgresIndexes() (PostgresIndexes, error) {

	rows, err := f.gormDB.
		Table("pg_indexes").
		Where("tablename = ?", "records").
		Select("tablename, indexname, indexdef").
		Rows()

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []*PostgresIndex

	for rows.Next() {
		var tableName, indexName, indexDef string
		if err := rows.Scan(&tableName, &indexName, &indexDef); err != nil {
			return nil, err
		}
		result = append(result, &PostgresIndex{
			TableName: tableName,
			IndexName: indexName,
			IndexDef:  indexDef,
		})
	}
	return result, nil
}

func (f *flame) GetIndexes() ([]Index, error) {

	pIndexes, err := f.getPostgresIndexes()
	if err != nil {
		return nil, fmt.Errorf("Failed to get pgsql indexes: %s", err.Error())
	}

	results := make([]Index, len(pIndexes))
	for i := range pIndexes {
		results[i] = Index{
			Name:     pIndexes[i].IndexName,
			Parent:   "", // todo
			Property: "", // todo
		}
	}
	return results, nil
}

func (f *flame) CreateIndex(index Index) error {

	if index.Name == "" || index.Property == "" {
		return errors.New("Must specify an index name and property")
	}

	cmd := fmt.Sprintf("CREATE INDEX %s ON records ((properties->>'%s'))",
		index.Name, index.Property)
	if index.Parent != "" {
		cmd = fmt.Sprintf("%s WHERE parent = '%s';", cmd, index.Parent)
	} else {
		cmd += ";"
	}

	return f.gormDB.Exec(cmd).Error
}

func (f *flame) DeleteIndex(index Index) error {

	if index.Name == "" {
		return errors.New("Must specify an index name")
	}

	cmd := fmt.Sprintf("DROP INDEX %s;", index.Name)
	return f.gormDB.Exec(cmd).Error
}

func (f *flame) HasIndex(index Index) (bool, error) {

	if index.Name == "" {
		return false, errors.New("Must specify an index name")
	}

	indexes, err := f.getPostgresIndexes()
	if err != nil {
		return false, err
	}

	for _, curIndex := range indexes {
		if curIndex.IndexName == index.Name {
			return true, nil
		}
	}
	return false, nil
}

func (f *flame) Get(key Key) (*Record, error) {

	if key.ID == "" && key.Path == "" {
		return nil, errors.New("Invalid key")
	}
	where := Record{ID: key.ID, Path: key.Path}

	r := Record{}
	if err := f.gormDB.Where(where).First(&r).Error; err != nil {
		return nil, err
	}
	return &r, nil
}

func (f *flame) Save(r *Record) error {
	if r.ID == "" {
		r.ID = NewID()
	}
	return f.gormDB.Save(r).Error
}

func (f *flame) Delete(r *Record) error {
	if r.ID == "" {
		return errors.New("Invalid id: empty")
	}
	return f.gormDB.Delete(r).Error
}

func (f *flame) List(q Query) ([]*Record, error) {

	if q.Offset < 0 {
		return nil, errors.New("Invalid negative offset")
	}
	if q.Limit < 0 {
		return nil, errors.New("Invalid negative limit")
	}
	if q.Limit == 0 {
		q.Limit = 100
	}

	query := f.gormDB

	if q.OrderBy != "" {
		if q.OrderByDesc {
			query = query.Order(fmt.Sprintf("%s desc", q.OrderBy))
		} else {
			query = query.Order(q.OrderBy)
		}
	}

	if q.OrderByProperty != "" {
		if q.OrderByPropertyDesc {
			query = query.Order(fmt.Sprintf("(properties->'%s') desc", q.OrderByProperty))
		} else {
			query = query.Order(fmt.Sprintf("(properties->'%s')", q.OrderByProperty))
		}
	}

	if q.Offset > 0 {
		query = query.Offset(q.Offset)
	}

	query = query.Limit(q.Limit)

	if q.Parent != "" {
		query = query.Where("parent = ?", q.Parent)
	} else if q.Prefix != "" {
		prefix := path.Join(q.Prefix, "%")
		query = query.Where("path LIKE ?", prefix)
	}

	if len(q.PropertyFilter) > 0 {
		query = query.Where(buildWhere(q.PropertyFilter))
	}

	var records []*Record
	if err := query.Find(&records).Error; err != nil {
		return nil, err
	}
	return records, nil
}

func buildWhere(filter map[string]string) string {
	var props []string
	for k, v := range filter {
		props = append(props, fmt.Sprintf("properties->>'%s' = '%s'", k, v))
	}
	return strings.Join(props, " AND ")
}
