package flamedb

import (
	"encoding/json"
	"math"
	"math/rand"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/jinzhu/gorm/dialects/postgres"

	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	testDB     *gorm.DB
	testResult int
	letters    = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
)

func init() {
	rand.Seed(time.Now().UnixNano())

	var err error
	testDB, err = Connect(Settings{
		Host:     "127.0.0.1",
		Port:     5432,
		User:     "test",
		Password: "test",
		Name:     "test",
		SSLMode:  SSLModeDisabled,
	})
	if err != nil {
		panic(err)
	}

	if err := testDB.DropTableIfExists(&Record{}).Error; err != nil {
		panic(err)
	}
	if err := testDB.AutoMigrate(&Record{}).Error; err != nil {
		panic(err)
	}
}

func jsonb(obj interface{}) postgres.Jsonb {
	js, err := json.Marshal(obj)
	if err != nil {
		panic(err)
	}
	return postgres.Jsonb{RawMessage: js}
}

func timeDelta(t1, t2 time.Time) float64 {
	dt := t1.Sub(t2).Seconds()
	return math.Abs(dt)
}

func testRecordExists(id string) bool {
	_, err := NewFlame(testDB).Get(Key{ID: id})
	return err == nil
}

func testClearRecords() {
	db := NewFlame(testDB)
	records, err := db.List(Query{Limit: 10000, Prefix: "/", OrderBy: "path"})
	if err != nil {
		panic(err)
	}
	for _, r := range records {
		if err := db.Delete(r); err != nil {
			panic(err)
		}
	}
}

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func testRandPath(len int) string {
	parts := make([]string, len)
	for i := range parts {
		parts[i] = randSeq(len)
	}
	return "/" + strings.Join(parts, "/")
}

func getRecordPaths(records []*Record) []string {
	res := make([]string, len(records))
	for i, blob := range records {
		res[i] = blob.Path
	}
	return res
}

func TestEmptyDatabase(t *testing.T) {
	db := NewFlame(testDB)
	defer testClearRecords()
	records, err := db.List(Query{OrderBy: "path", Parent: "/"})
	require.Nil(t, err)
	require.Equal(t, 0, len(records))
}

func TestCreateRecord(t *testing.T) {
	db := NewFlame(testDB)
	defer testClearRecords()
	blob := &Record{
		ID:   NewID(),
		Path: "/boom.wav",
	}
	require.Nil(t, db.Save(blob))
	assert.Equal(t, "/", blob.Parent)

	records, err := db.List(Query{OrderBy: "path", Parent: "/"})
	require.Nil(t, err)
	require.Equal(t, 1, len(records))

	r1 := records[0]
	assert.Equal(t, blob.ID, r1.ID)
	assert.Equal(t, blob.Path, r1.Path)
	assert.Equal(t, blob.Properties, r1.Properties)
	assert.Equal(t, blob.CreatedBy, r1.CreatedBy)
	assert.Equal(t, blob.UpdatedBy, r1.UpdatedBy)
	assert.True(t, timeDelta(blob.CreatedAt, r1.CreatedAt) < 0.01)
	assert.True(t, timeDelta(blob.UpdatedAt, r1.UpdatedAt) < 0.01)
}

func TestCreateNoPath(t *testing.T) {
	db := NewFlame(testDB)
	defer testClearRecords()

	err := db.Save(&Record{ID: "123"})
	require.NotNil(t, err)
	require.Equal(t, "Invalid record path: empty", err.Error())
}

func TestGetRecord(t *testing.T) {
	db := NewFlame(testDB)
	defer testClearRecords()

	r1 := &Record{ID: NewID(), Path: "/a/b"}
	require.Nil(t, db.Save(r1))

	b2 := &Record{ID: NewID(), Path: "/c/d"}
	require.Nil(t, db.Save(b2))

	got, err := db.Get(Key{ID: r1.ID})
	require.Nil(t, err)
	assert.Equal(t, r1.ID, got.ID)
	assert.Equal(t, r1.Path, got.Path)

	got, err = db.Get(Key{Path: b2.Path})
	require.Nil(t, err)
	assert.Equal(t, b2.ID, got.ID)
	assert.Equal(t, b2.Path, got.Path)

	_, err = db.Get(Key{})
	require.NotNil(t, err)
	assert.Equal(t, "Invalid key", err.Error())

	_, err = db.Get(Key{ID: "123"})
	require.NotNil(t, err)
	assert.Equal(t, "record not found", err.Error())

	_, err = db.Get(Key{Path: "/1/2/3"})
	require.NotNil(t, err)
	assert.Equal(t, "record not found", err.Error())
}

func TestDeleteRecord(t *testing.T) {
	db := NewFlame(testDB)
	defer testClearRecords()
	blob := &Record{ID: NewID(), Path: "/foo/baz"}
	require.Nil(t, db.Save(blob))
	require.True(t, testRecordExists(blob.ID))
	require.Nil(t, db.Delete(blob))
	require.False(t, testRecordExists(blob.ID))
}

func TestDeleteErrors(t *testing.T) {
	db := NewFlame(testDB)
	defer testClearRecords()
	blob := &Record{Path: "/foo/baz"}
	err := db.Delete(blob)
	require.NotNil(t, err)
	require.Equal(t, "Invalid id: empty", err.Error())
}

func TestListRecords(t *testing.T) {
	db := NewFlame(testDB)
	defer testClearRecords()

	require.Nil(t, db.Save(&Record{ID: NewID(), Path: "/b/c/d"}))
	require.Nil(t, db.Save(&Record{ID: NewID(), Path: "/b/b/c/d"}))
	require.Nil(t, db.Save(&Record{ID: NewID(), Path: "/x/y/z"}))
	require.Nil(t, db.Save(&Record{ID: NewID(), Path: "/aaa/aa/a"}))

	// OrderBy: path
	records, err := db.List(Query{OrderBy: "path", Prefix: "/"})
	require.Nil(t, err)
	require.Len(t, records, 4)
	assert.Equal(t, "/aaa/aa/a", records[0].Path)
	assert.Equal(t, "/b/b/c/d", records[1].Path)
	assert.Equal(t, "/b/c/d", records[2].Path)
	assert.Equal(t, "/x/y/z", records[3].Path)

	records, err = db.List(Query{OrderBy: "id", Prefix: "/"})
	require.Nil(t, err)
	require.Len(t, records, 4)
	assert.Equal(t, "/b/c/d", records[0].Path)
	assert.Equal(t, "/b/b/c/d", records[1].Path)
	assert.Equal(t, "/x/y/z", records[2].Path)
	assert.Equal(t, "/aaa/aa/a", records[3].Path)

	records, err = db.List(Query{
		Limit:   2,
		OrderBy: "path",
		Prefix:  "/",
	})
	require.Nil(t, err)
	require.Len(t, records, 2)
	assert.Equal(t, "/aaa/aa/a", records[0].Path)
	assert.Equal(t, "/b/b/c/d", records[1].Path)

	records, err = db.List(Query{
		Offset:  3,
		Limit:   2,
		OrderBy: "path",
		Prefix:  "/",
	})
	require.Nil(t, err)
	require.Len(t, records, 1)
	assert.Equal(t, "/x/y/z", records[0].Path)

	records, err = db.List(Query{OrderBy: "path", Prefix: "/b"})
	require.Nil(t, err)
	require.Len(t, records, 2)
	assert.Equal(t, "/b/b/c/d", records[0].Path)
	assert.Equal(t, "/b/c/d", records[1].Path)

	records, err = db.List(Query{OrderBy: "path", Prefix: "/b/c/"})
	require.Nil(t, err)
	require.Len(t, records, 1)
	assert.Equal(t, "/b/c/d", records[0].Path)

	records, err = db.List(Query{OrderBy: "path", Prefix: "/b/c/d/e"})
	require.Nil(t, err)
	require.Len(t, records, 0)
}

func TestListErrors(t *testing.T) {
	db := NewFlame(testDB)
	defer testClearRecords()

	_, err := db.List(Query{OrderBy: "path", Offset: -1})
	require.NotNil(t, err)
	require.Equal(t, "Invalid negative offset", err.Error())

	_, err = db.List(Query{OrderBy: "path", Limit: -1})
	require.NotNil(t, err)
	require.Equal(t, "Invalid negative limit", err.Error())
}

func TestListRecordsOrdering(t *testing.T) {
	db := NewFlame(testDB)
	defer testClearRecords()

	r1 := &Record{
		ID:   NewID(),
		Path: "/1/d.wav",
		Properties: jsonb(map[string]interface{}{
			"name":  "d",
			"count": 1,
		}),
	}
	r2 := &Record{
		ID:   NewID(),
		Path: "/9/c.wav",
		Properties: jsonb(map[string]interface{}{
			"name":  "c",
			"count": -1,
		}),
	}
	r3 := &Record{
		ID:   NewID(),
		Path: "/7/a.wav",
		Properties: jsonb(map[string]interface{}{
			"name":  "a",
			"count": 3,
		}),
	}
	r4 := &Record{
		ID:   NewID(),
		Path: "/8/b.wav",
		Properties: jsonb(map[string]interface{}{
			"name":  "b",
			"count": 11,
		}),
	}

	for _, r := range []*Record{r1, r2, r3, r4} {
		require.Nil(t, db.Save(r))
	}

	records, _ := db.List(Query{
		Limit:               100,
		OrderByProperty:     "name",
		OrderByPropertyDesc: false,
	})
	require.Len(t, records, 4)
	expPaths := []string{"/7/a.wav", "/8/b.wav", "/9/c.wav", "/1/d.wav"}
	assert.Equal(t, expPaths, getRecordPaths(records))

	records, _ = db.List(Query{
		Limit:               100,
		OrderByProperty:     "name",
		OrderByPropertyDesc: true,
	})
	require.Len(t, records, 4)
	expPaths = []string{"/1/d.wav", "/9/c.wav", "/8/b.wav", "/7/a.wav"}
	assert.Equal(t, expPaths, getRecordPaths(records))

	records, _ = db.List(Query{
		Limit:               100,
		OrderByProperty:     "count",
		OrderByPropertyDesc: false,
	})
	require.Len(t, records, 4)
	expPaths = []string{"/9/c.wav", "/1/d.wav", "/7/a.wav", "/8/b.wav"}
	assert.Equal(t, expPaths, getRecordPaths(records))

	records, _ = db.List(Query{
		Limit:               100,
		OrderByProperty:     "count",
		OrderByPropertyDesc: true,
	})
	require.Len(t, records, 4)
	expPaths = []string{"/8/b.wav", "/7/a.wav", "/1/d.wav", "/9/c.wav"}
	assert.Equal(t, expPaths, getRecordPaths(records))
}

func TestListRecordsPropertyFilters(t *testing.T) {
	db := NewFlame(testDB)
	defer testClearRecords()

	r1 := &Record{
		ID:   NewID(),
		Path: "/1/d.wav",
		Properties: jsonb(map[string]interface{}{
			"name":   "d",
			"age":    2,
			"folder": "folderA",
		}),
	}
	r2 := &Record{
		ID:   NewID(),
		Path: "/9/c.wav",
		Properties: jsonb(map[string]interface{}{
			"name":   "c",
			"age":    1,
			"folder": "folderA",
		}),
	}
	r3 := &Record{
		ID:   NewID(),
		Path: "/7/a.wav",
		Properties: jsonb(map[string]interface{}{
			"name":   "a",
			"age":    2,
			"folder": "folderB",
		}),
	}

	for _, r := range []*Record{r1, r2, r3} {
		require.Nil(t, db.Save(r))
	}

	records, _ := db.List(Query{
		PropertyFilter:  map[string]string{"folder": "folderA"},
		OrderByProperty: "name",
	})
	require.Len(t, records, 2)
	expPaths := []string{"/9/c.wav", "/1/d.wav"}
	assert.Equal(t, expPaths, getRecordPaths(records))

	records, _ = db.List(Query{
		PropertyFilter: map[string]string{"name": "a"},
	})
	require.Len(t, records, 1)
	expPaths = []string{"/7/a.wav"}
	assert.Equal(t, expPaths, getRecordPaths(records))

	records, _ = db.List(Query{
		PropertyFilter:      map[string]string{"age": "2"},
		OrderByProperty:     "name",
		OrderByPropertyDesc: true,
	})
	require.Len(t, records, 2)
	expPaths = []string{"/1/d.wav", "/7/a.wav"}
	assert.Equal(t, expPaths, getRecordPaths(records))
}

func TestCreateDeleteIndex(t *testing.T) {
	db := NewFlame(testDB)
	defer testClearRecords()

	_, err := db.HasIndex(Index{})
	require.NotNil(t, err)
	require.Equal(t, "Must specify an index name", err.Error())

	exists, err := db.HasIndex(Index{Name: "name_idx"})
	require.Nil(t, err)
	require.False(t, exists)

	require.Nil(t, db.CreateIndex(Index{Name: "name_idx", Property: "name"}))

	exists, err = db.HasIndex(Index{Name: "name_idx"})
	require.Nil(t, err)
	require.True(t, exists)

	require.Nil(t, db.DeleteIndex(Index{Name: "name_idx"}))

	exists, err = db.HasIndex(Index{Name: "name_idx"})
	require.Nil(t, err)
	require.False(t, exists)
}

func TestPartialIndex(t *testing.T) {
	db := NewFlame(testDB)
	defer testClearRecords()

	index := Index{
		Name:     "sound_names",
		Property: "name",
		Parent:   "/sounds",
	}

	defer db.DeleteIndex(index)

	exists, err := db.HasIndex(index)
	require.Nil(t, err)
	require.False(t, exists)

	require.Nil(t, db.CreateIndex(index))

	dbStruct := db.(*flame)
	indexes, _ := dbStruct.getPostgresIndexes()
	foundIndex := indexes.Get("sound_names")
	require.NotNil(t, foundIndex)
	indexDef := foundIndex.IndexDef

	expectedDef := "CREATE INDEX sound_names ON public.records "
	expectedDef += "USING btree (((properties ->> 'name'::text))) "
	expectedDef += "WHERE ((parent)::text = '/sounds'::text)"
	assert.Equal(t, expectedDef, indexDef)
}

func TestIndexErrors(t *testing.T) {
	db := NewFlame(testDB)
	defer testClearRecords()

	err := db.CreateIndex(Index{Property: "name"})
	require.NotNil(t, err)
	require.Equal(t, "Must specify an index name and property", err.Error())

	err = db.DeleteIndex(Index{Property: "name"})
	require.NotNil(t, err)
	require.Equal(t, "Must specify an index name", err.Error())
}

func TestGetIndexes(t *testing.T) {
	db := NewFlame(testDB)
	flame := db.(*flame)
	defer testClearRecords()
	indexes, err := flame.getPostgresIndexes()
	require.Nil(t, err)
	var names []string
	for _, index := range indexes {
		require.Equal(t, "records", index.TableName)
		names = append(names, index.IndexName)
	}
	sort.Strings(names)
	expected := []string{
		"idx_records_parent",
		"records_pkey",
		"uix_records_id",
		"uix_records_path",
	}
	assert.Equal(t, expected, names)
}

func BenchmarkCreateRecord(b *testing.B) {
	db := NewFlame(testDB)
	defer testClearRecords()
	for n := 0; n < b.N; n++ {
		p := testRandPath(4)
		db.Save(&Record{ID: NewID(), Path: p})
	}
}

func BenchmarkListRecord(b *testing.B) {
	db := NewFlame(testDB)
	defer testClearRecords()
	for i := 0; i < 1000; i++ {
		p := testRandPath(4)
		db.Save(&Record{ID: NewID(), Path: p})
	}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_, err := db.List(Query{Limit: 100, OrderBy: "path", Parent: "/"})
		if err != nil {
			panic(err)
		}
	}
	b.StopTimer()
}

func BenchmarkQueryRecordName(b *testing.B) {
	db := NewFlame(testDB)
	defer testClearRecords()
	for i := 0; i < 1000; i++ {
		p := testRandPath(4)
		props := map[string]interface{}{
			"name": randSeq(8),
		}
		propJSON, err := json.Marshal(props)
		if err != nil {
			panic(err)
		}
		db.Save(&Record{
			ID:         NewID(),
			Path:       p,
			Properties: postgres.Jsonb{RawMessage: propJSON},
		})
	}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_, err := db.List(Query{Limit: 100, OrderBy: "path", Parent: "/"})
		if err != nil {
			panic(err)
		}
	}
	b.StopTimer()
}

func TestBuildWhere(t *testing.T) {
	filter := buildWhere(map[string]string{"name": "frank"})
	assert.Equal(t, "properties->>'name' = 'frank'", filter)
}
