package xorm

import (
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"
)

/*
CREATE TABLE `userinfo` (
	`id` INT(10) NULL AUTO_INCREMENT,
	`username` VARCHAR(64) NULL,
	`departname` VARCHAR(64) NULL,
	`created` DATE NULL,
	PRIMARY KEY (`uid`)
);
CREATE TABLE `userdeatail` (
	`id` INT(10) NULL,
	`intro` TEXT NULL,
	`profile` TEXT NULL,
	PRIMARY KEY (`uid`)
);
*/

type Userinfo struct {
	Uid        int64  `xorm:"id pk not null autoincr"`
	Username   string `xorm:"unique"`
	Departname string
	Alias      string `xorm:"-"`
	Created    time.Time
	Detail     Userdetail `xorm:"detail_id int(11)"`
	Height     float64
	Avatar     []byte
	IsMan      bool
}

type Userdetail struct {
	Id      int64
	Intro   string `xorm:"text"`
	Profile string `xorm:"varchar(2000)"`
}

func directCreateTable(engine *Engine, t *testing.T) {
	err := engine.DropTables(&Userinfo{}, &Userdetail{})
	if err != nil {
		t.Error(err)
		panic(err)
	}

	err = engine.Sync(&Userinfo{}, &Userdetail{})
	if err != nil {
		t.Error(err)
		panic(err)
	}

	err = engine.DropTables(&Userinfo{}, &Userdetail{})
	if err != nil {
		t.Error(err)
		panic(err)
	}

	err = engine.CreateTables(&Userinfo{}, &Userdetail{})
	if err != nil {
		t.Error(err)
		panic(err)
	}

	err = engine.CreateIndexes(&Userinfo{})
	if err != nil {
		t.Error(err)
		panic(err)
	}

	err = engine.CreateIndexes(&Userdetail{})
	if err != nil {
		t.Error(err)
		panic(err)
	}

	err = engine.CreateUniques(&Userinfo{})
	if err != nil {
		t.Error(err)
		panic(err)
	}

	err = engine.CreateUniques(&Userdetail{})
	if err != nil {
		t.Error(err)
		panic(err)
	}
}

func insert(engine *Engine, t *testing.T) {
	user := Userinfo{0, "xiaolunwen", "dev", "lunny", time.Now(),
		Userdetail{Id: 1}, 1.78, []byte{1, 2, 3}, true}
	cnt, err := engine.Insert(&user)
	fmt.Println(user.Uid)
	if err != nil {
		t.Error(err)
		panic(err)
	}
	if cnt != 1 {
		err = errors.New("insert not returned 1")
		t.Error(err)
		panic(err)
		return
	}
	if user.Uid <= 0 {
		err = errors.New("not return id error")
		t.Error(err)
		panic(err)
	}

	user.Uid = 0
	cnt, err = engine.Insert(&user)
	if err == nil {
		err = errors.New("insert failed but no return error")
		t.Error(err)
		panic(err)
	}
	if cnt != 0 {
		err = errors.New("insert not returned 1")
		t.Error(err)
		panic(err)
		return
	}
}

func testQuery(engine *Engine, t *testing.T) {
	sql := "select * from userinfo"
	results, err := engine.Query(sql)
	if err != nil {
		t.Error(err)
		panic(err)
	}
	fmt.Println(results)
}

func exec(engine *Engine, t *testing.T) {
	sql := "update userinfo set username=? where id=?"
	res, err := engine.Exec(sql, "xiaolun", 1)
	if err != nil {
		t.Error(err)
		panic(err)
	}
	fmt.Println(res)
}

func querySameMapper(engine *Engine, t *testing.T) {
	sql := "select * from `Userinfo`"
	results, err := engine.Query(sql)
	if err != nil {
		t.Error(err)
		panic(err)
	}
	fmt.Println(results)
}

func execSameMapper(engine *Engine, t *testing.T) {
	sql := "update `Userinfo` set `Username`=? where (id)=?"
	res, err := engine.Exec(sql, "xiaolun", 1)
	if err != nil {
		t.Error(err)
		panic(err)
	}
	fmt.Println(res)
}

func insertAutoIncr(engine *Engine, t *testing.T) {
	// auto increment insert
	user := Userinfo{Username: "xiaolunwen2", Departname: "dev", Alias: "lunny", Created: time.Now(),
		Detail: Userdetail{Id: 1}, Height: 1.78, Avatar: []byte{1, 2, 3}, IsMan: true}
	cnt, err := engine.Insert(&user)
	fmt.Println(user.Uid)
	if err != nil {
		t.Error(err)
		panic(err)
	}
	if cnt != 1 {
		err = errors.New("insert not returned 1")
		t.Error(err)
		panic(err)
		return
	}
	if user.Uid <= 0 {
		t.Error(errors.New("not return id error"))
	}
}

func insertMulti(engine *Engine, t *testing.T) {
	//engine.InsertMany = true
	users := []Userinfo{
		{Username: "xlw", Departname: "dev", Alias: "lunny2", Created: time.Now()},
		{Username: "xlw2", Departname: "dev", Alias: "lunny3", Created: time.Now()},
		{Username: "xlw11", Departname: "dev", Alias: "lunny2", Created: time.Now()},
		{Username: "xlw22", Departname: "dev", Alias: "lunny3", Created: time.Now()},
	}
	cnt, err := engine.Insert(&users)
	if err != nil {
		t.Error(err)
		panic(err)
	}
	if cnt != int64(len(users)) {
		err = errors.New("insert not returned 1")
		t.Error(err)
		panic(err)
		return
	}

	users2 := []*Userinfo{
		&Userinfo{Username: "1xlw", Departname: "dev", Alias: "lunny2", Created: time.Now()},
		&Userinfo{Username: "1xlw2", Departname: "dev", Alias: "lunny3", Created: time.Now()},
		&Userinfo{Username: "1xlw11", Departname: "dev", Alias: "lunny2", Created: time.Now()},
		&Userinfo{Username: "1xlw22", Departname: "dev", Alias: "lunny3", Created: time.Now()},
	}

	cnt, err = engine.Insert(&users2)
	if err != nil {
		t.Error(err)
		panic(err)
	}

	if cnt != int64(len(users2)) {
		err = errors.New(fmt.Sprintf("insert not returned %v", len(users2)))
		t.Error(err)
		panic(err)
		return
	}
}

func insertTwoTable(engine *Engine, t *testing.T) {
	userdetail := Userdetail{Id: 1, Intro: "I'm a very beautiful women.", Profile: "sfsaf"}
	userinfo := Userinfo{Username: "xlw3", Departname: "dev", Alias: "lunny4", Created: time.Now(), Detail: userdetail}

	cnt, err := engine.Insert(&userinfo, &userdetail)
	if err != nil {
		t.Error(err)
		panic(err)
	}

	if userinfo.Uid <= 0 {
		err = errors.New("not return id error")
		t.Error(err)
		panic(err)
	}

	if userdetail.Id <= 0 {
		err = errors.New("not return id error")
		t.Error(err)
		panic(err)
	}

	if cnt != 2 {
		err = errors.New("insert not returned 2")
		t.Error(err)
		panic(err)
		return
	}
}

type Condi map[string]interface{}

func update(engine *Engine, t *testing.T) {
	// update by id
	user := Userinfo{Username: "xxx", Height: 1.2}
	cnt, err := engine.Id(1).Update(&user)
	if err != nil {
		t.Error(err)
		panic(err)
	}
	if cnt != 1 {
		err = errors.New("update not returned 1")
		t.Error(err)
		panic(err)
		return
	}

	condi := Condi{"username": "zzz", "height": 0.0, "departname": ""}
	cnt, err = engine.Table(&user).Id(1).Update(&condi)
	if err != nil {
		t.Error(err)
		panic(err)
	}
	if cnt != 1 {
		err = errors.New("update not returned 1")
		t.Error(err)
		panic(err)
		return
	}

	cnt, err = engine.Update(&Userinfo{Username: "yyy"}, &user)
	if err != nil {
		t.Error(err)
		panic(err)
	}
	total, err := engine.Count(&user)
	if err != nil {
		t.Error(err)
		panic(err)
	}

	if cnt != total {
		err = errors.New("insert not returned 1")
		t.Error(err)
		panic(err)
		return
	}
}

func updateSameMapper(engine *Engine, t *testing.T) {
	// update by id
	user := Userinfo{Username: "xxx", Height: 1.2}
	cnt, err := engine.Id(1).Update(&user)
	if err != nil {
		t.Error(err)
		panic(err)
	}
	if cnt != 1 {
		err = errors.New("update not returned 1")
		t.Error(err)
		panic(err)
		return
	}

	condi := Condi{"Username": "zzz", "Height": 0.0, "Departname": ""}
	cnt, err = engine.Table(&user).Id(1).Update(&condi)
	if err != nil {
		t.Error(err)
		panic(err)
	}

	if cnt != 1 {
		err = errors.New("insert not returned 1")
		t.Error(err)
		panic(err)
		return
	}

	cnt, err = engine.Update(&Userinfo{Username: "yyy"}, &user)
	if err != nil {
		t.Error(err)
		panic(err)
	}

	if cnt != 1 {
		err = errors.New("insert not returned 1")
		t.Error(err)
		panic(err)
		return
	}
}

func testdelete(engine *Engine, t *testing.T) {
	user := Userinfo{Uid: 1}
	cnt, err := engine.Delete(&user)
	if err != nil {
		t.Error(err)
		panic(err)
	}
	if cnt != 1 {
		err = errors.New("delete not returned 1")
		t.Error(err)
		panic(err)
		return
	}

	user.Uid = 0
	has, err := engine.Id(3).Get(&user)
	if err != nil {
		t.Error(err)
		panic(err)
	}

	if has {
		//var tt time.Time
		//user.Created = tt
		cnt, err := engine.UseBool().Delete(&user)
		if err != nil {
			t.Error(err)
			panic(err)
		}
		if cnt != 1 {
			t.Error(errors.New("delete failed"))
			panic(err)
		}
	}
}

func get(engine *Engine, t *testing.T) {
	user := Userinfo{Uid: 2}

	has, err := engine.Get(&user)
	if err != nil {
		t.Error(err)
		panic(err)
	}
	if has {
		fmt.Println(user)
	} else {
		fmt.Println("no record id is 2")
	}
}

func cascadeGet(engine *Engine, t *testing.T) {
	user := Userinfo{Uid: 11}

	has, err := engine.Get(&user)
	if err != nil {
		t.Error(err)
		panic(err)
	}
	if has {
		fmt.Println(user)
	} else {
		fmt.Println("no record id is 2")
	}
}

func find(engine *Engine, t *testing.T) {
	users := make([]Userinfo, 0)

	err := engine.Find(&users)
	if err != nil {
		t.Error(err)
		panic(err)
	}
	for _, user := range users {
		fmt.Println(user)
	}
}

func find2(engine *Engine, t *testing.T) {
	users := make([]*Userinfo, 0)

	err := engine.Find(&users)
	if err != nil {
		t.Error(err)
		panic(err)
	}
	for _, user := range users {
		fmt.Println(user)
	}
}

func findMap(engine *Engine, t *testing.T) {
	users := make(map[int64]Userinfo)

	err := engine.Find(&users)
	if err != nil {
		t.Error(err)
		panic(err)
	}
	for _, user := range users {
		fmt.Println(user)
	}
}

func findMap2(engine *Engine, t *testing.T) {
	users := make(map[int64]*Userinfo)

	err := engine.Find(&users)
	if err != nil {
		t.Error(err)
		panic(err)
	}
	for id, user := range users {
		fmt.Println(id, user)
	}
}

func count(engine *Engine, t *testing.T) {
	user := Userinfo{Departname: "dev"}
	total, err := engine.Count(&user)
	if err != nil {
		t.Error(err)
		panic(err)
	}
	fmt.Printf("Total %d records!!!\n", total)
}

func where(engine *Engine, t *testing.T) {
	users := make([]Userinfo, 0)
	err := engine.Where("(id) > ?", 2).Find(&users)
	if err != nil {
		t.Error(err)
		panic(err)
	}
	fmt.Println(users)

	err = engine.Where("(id) > ?", 2).And("(id) < ?", 10).Find(&users)
	if err != nil {
		t.Error(err)
		panic(err)
	}
	fmt.Println(users)
}

func in(engine *Engine, t *testing.T) {
	users := make([]Userinfo, 0)
	err := engine.In("(id)", 1, 2, 3).Find(&users)
	if err != nil {
		t.Error(err)
		panic(err)
	}
	fmt.Println(users)

	ids := []interface{}{1, 2, 3}
	err = engine.Where("(id) > ?", 2).In("(id)", ids...).Find(&users)
	if err != nil {
		t.Error(err)
		panic(err)
	}
	fmt.Println(users)
}

func limit(engine *Engine, t *testing.T) {
	users := make([]Userinfo, 0)
	err := engine.Limit(2, 1).Find(&users)
	if err != nil {
		t.Error(err)
		panic(err)
	}
	fmt.Println(users)
}

func order(engine *Engine, t *testing.T) {
	users := make([]Userinfo, 0)
	err := engine.OrderBy("id desc").Find(&users)
	if err != nil {
		t.Error(err)
		panic(err)
	}
	fmt.Println(users)

	users2 := make([]Userinfo, 0)
	err = engine.Asc("id", "username").Desc("height").Find(&users2)
	if err != nil {
		t.Error(err)
		panic(err)
	}
	fmt.Println(users2)
}

func join(engine *Engine, t *testing.T) {
	users := make([]Userinfo, 0)
	err := engine.Join("LEFT", "userdetail", "userinfo.id=userdetail.id").Find(&users)
	if err != nil {
		t.Error(err)
		panic(err)
	}
}

func having(engine *Engine, t *testing.T) {
	users := make([]Userinfo, 0)
	err := engine.GroupBy("username").Having("username='xlw'").Find(&users)
	if err != nil {
		t.Error(err)
		panic(err)
	}
	fmt.Println(users)
}

func orderSameMapper(engine *Engine, t *testing.T) {
	users := make([]Userinfo, 0)
	err := engine.OrderBy("(id) desc").Find(&users)
	if err != nil {
		t.Error(err)
		panic(err)
	}
	fmt.Println(users)

	users2 := make([]Userinfo, 0)
	err = engine.Asc("(id)", "Username").Desc("Height").Find(&users2)
	if err != nil {
		t.Error(err)
		panic(err)
	}
	fmt.Println(users2)
}

func joinSameMapper(engine *Engine, t *testing.T) {
	users := make([]Userinfo, 0)
	err := engine.Join("LEFT", `"Userdetail"`, `"Userinfo"."id"="Userdetail"."Id"`).Find(&users)
	if err != nil {
		t.Error(err)
		panic(err)
	}
}

func havingSameMapper(engine *Engine, t *testing.T) {
	users := make([]Userinfo, 0)
	err := engine.GroupBy("Username").Having(`"Username"='xlw'`).Find(&users)
	if err != nil {
		t.Error(err)
		panic(err)
	}
	fmt.Println(users)
}

func transaction(engine *Engine, t *testing.T) {
	counter := func() {
		total, err := engine.Count(&Userinfo{})
		if err != nil {
			t.Error(err)
		}
		fmt.Printf("----now total %v records\n", total)
	}

	counter()
	defer counter()
	session := engine.NewSession()
	defer session.Close()

	err := session.Begin()
	if err != nil {
		t.Error(err)
		panic(err)
	}
	//session.IsAutoRollback = false
	user1 := Userinfo{Username: "xiaoxiao", Departname: "dev", Alias: "lunny", Created: time.Now()}
	_, err = session.Insert(&user1)
	if err != nil {
		session.Rollback()
		t.Error(err)
		panic(err)
	}

	user2 := Userinfo{Username: "yyy"}
	_, err = session.Where("(id) = ?", 0).Update(&user2)
	if err != nil {
		session.Rollback()
		fmt.Println(err)
		//t.Error(err)
		return
	}

	_, err = session.Delete(&user2)
	if err != nil {
		session.Rollback()
		t.Error(err)
		panic(err)
	}

	err = session.Commit()
	if err != nil {
		t.Error(err)
		panic(err)
	}
	panic(err)
}

func combineTransaction(engine *Engine, t *testing.T) {
	counter := func() {
		total, err := engine.Count(&Userinfo{})
		if err != nil {
			t.Error(err)
		}
		fmt.Printf("----now total %v records\n", total)
	}

	counter()
	defer counter()
	session := engine.NewSession()
	defer session.Close()

	err := session.Begin()
	if err != nil {
		t.Error(err)
		panic(err)
	}
	//session.IsAutoRollback = false
	user1 := Userinfo{Username: "xiaoxiao2", Departname: "dev", Alias: "lunny", Created: time.Now()}
	_, err = session.Insert(&user1)
	if err != nil {
		session.Rollback()
		t.Error(err)
		panic(err)
	}
	user2 := Userinfo{Username: "zzz"}
	_, err = session.Where("id = ?", 0).Update(&user2)
	if err != nil {
		session.Rollback()
		t.Error(err)
		panic(err)
	}

	_, err = session.Exec("delete from userinfo where username = ?", user2.Username)
	if err != nil {
		session.Rollback()
		t.Error(err)
		panic(err)
	}

	err = session.Commit()
	if err != nil {
		t.Error(err)
		panic(err)
	}
}

func combineTransactionSameMapper(engine *Engine, t *testing.T) {
	counter := func() {
		total, err := engine.Count(&Userinfo{})
		if err != nil {
			t.Error(err)
		}
		fmt.Printf("----now total %v records\n", total)
	}

	counter()
	defer counter()
	session := engine.NewSession()
	defer session.Close()

	err := session.Begin()
	if err != nil {
		t.Error(err)
		panic(err)
	}
	//session.IsAutoRollback = false
	user1 := Userinfo{Username: "xiaoxiao2", Departname: "dev", Alias: "lunny", Created: time.Now()}
	_, err = session.Insert(&user1)
	if err != nil {
		session.Rollback()
		t.Error(err)
		panic(err)
	}
	user2 := Userinfo{Username: "zzz"}
	_, err = session.Where("(id) = ?", 0).Update(&user2)
	if err != nil {
		session.Rollback()
		t.Error(err)
		panic(err)
	}

	_, err = session.Exec("delete from `Userinfo` where `Username` = ?", user2.Username)
	if err != nil {
		session.Rollback()
		t.Error(err)
		panic(err)
	}

	err = session.Commit()
	if err != nil {
		t.Error(err)
		panic(err)
	}
}

func table(engine *Engine, t *testing.T) {
	err := engine.DropTables("user_user")
	if err != nil {
		t.Error(err)
		panic(err)
	}

	err = engine.Table("user_user").CreateTable(&Userinfo{})
	if err != nil {
		t.Error(err)
		panic(err)
	}
}

func createMultiTables(engine *Engine, t *testing.T) {
	session := engine.NewSession()
	defer session.Close()

	user := &Userinfo{}
	err := session.Begin()
	if err != nil {
		t.Error(err)
		panic(err)
	}
	for i := 0; i < 10; i++ {
		tableName := fmt.Sprintf("user_%v", i)

		err = session.DropTable(tableName)
		if err != nil {
			session.Rollback()
			t.Error(err)
			panic(err)
		}

		err = session.Table(tableName).CreateTable(user)
		if err != nil {
			session.Rollback()
			t.Error(err)
			panic(err)
		}
	}
	err = session.Commit()
	if err != nil {
		t.Error(err)
		panic(err)
	}
}

func tableOp(engine *Engine, t *testing.T) {
	user := Userinfo{Username: "tablexiao", Departname: "dev", Alias: "lunny", Created: time.Now()}
	tableName := fmt.Sprintf("user_%v", len(user.Username))
	cnt, err := engine.Table(tableName).Insert(&user)
	if err != nil {
		t.Error(err)
		panic(err)
	}
	if cnt != 1 {
		err = errors.New("insert not returned 1")
		t.Error(err)
		panic(err)
		return
	}

	has, err := engine.Table(tableName).Get(&Userinfo{Username: "tablexiao"})
	if err != nil {
		t.Error(err)
		panic(err)
	}
	if !has {
		err = errors.New("Get has return false")
		t.Error(err)
		panic(err)
		return
	}

	users := make([]Userinfo, 0)
	err = engine.Table(tableName).Find(&users)
	if err != nil {
		t.Error(err)
		panic(err)
	}

	id := user.Uid
	cnt, err = engine.Table(tableName).Id(id).Update(&Userinfo{Username: "tableda"})
	if err != nil {
		t.Error(err)
		panic(err)
	}

	_, err = engine.Table(tableName).Id(id).Delete(&Userinfo{})
	if err != nil {
		t.Error(err)
		panic(err)
	}
}

func testCharst(engine *Engine, t *testing.T) {
	err := engine.DropTables("user_charset")
	if err != nil {
		t.Error(err)
		panic(err)
	}

	err = engine.Charset("utf8").Table("user_charset").CreateTable(&Userinfo{})
	if err != nil {
		t.Error(err)
		panic(err)
	}
}

func testStoreEngine(engine *Engine, t *testing.T) {
	err := engine.DropTables("user_store_engine")
	if err != nil {
		t.Error(err)
		panic(err)
	}

	err = engine.StoreEngine("InnoDB").Table("user_store_engine").CreateTable(&Userinfo{})
	if err != nil {
		t.Error(err)
		panic(err)
	}
}

type tempUser struct {
	Id       int64
	Username string
}

func testCols(engine *Engine, t *testing.T) {
	users := []Userinfo{}
	err := engine.Cols("id, username").Find(&users)
	if err != nil {
		t.Error(err)
		panic(err)
	}

	fmt.Println(users)

	tmpUsers := []tempUser{}
	err = engine.NoCache().Table("userinfo").Cols("id, username").Find(&tmpUsers)
	if err != nil {
		t.Error(err)
		panic(err)
	}
	fmt.Println(tmpUsers)

	user := &Userinfo{Uid: 1, Alias: "", Height: 0}
	affected, err := engine.Cols("departname, height").Id(1).Update(user)
	if err != nil {
		t.Error(err)
		panic(err)
	}
	fmt.Println("===================", user, affected)
}

func testColsSameMapper(engine *Engine, t *testing.T) {
	users := []Userinfo{}
	err := engine.Cols("(id), Username").Find(&users)
	if err != nil {
		t.Error(err)
		panic(err)
	}

	fmt.Println(users)

	tmpUsers := []tempUser{}
	err = engine.Table("Userinfo").Cols("(id), Username").Find(&tmpUsers)
	if err != nil {
		t.Error(err)
		panic(err)
	}
	fmt.Println(tmpUsers)

	user := &Userinfo{Uid: 1, Alias: "", Height: 0}
	affected, err := engine.Cols("Departname, Height").Update(user)
	if err != nil {
		t.Error(err)
		panic(err)
	}
	fmt.Println("===================", user, affected)
}

type tempUser2 struct {
	tempUser   `xorm:"extends"`
	Departname string
}

func testExtends(engine *Engine, t *testing.T) {
	err := engine.DropTables(&tempUser2{})
	if err != nil {
		t.Error(err)
		panic(err)
	}

	err = engine.CreateTables(&tempUser2{})
	if err != nil {
		t.Error(err)
		panic(err)
	}

	tu := &tempUser2{tempUser{0, "extends"}, "dev depart"}
	_, err = engine.Insert(tu)
	if err != nil {
		t.Error(err)
		panic(err)
	}

	tu2 := &tempUser2{}
	_, err = engine.Get(tu2)
	if err != nil {
		t.Error(err)
		panic(err)
	}

	tu3 := &tempUser2{tempUser{0, "extends update"}, ""}
	_, err = engine.Id(tu2.Id).Update(tu3)
	if err != nil {
		t.Error(err)
		panic(err)
	}
}

type allCols struct {
	Bit       int   `xorm:"BIT"`
	TinyInt   int8  `xorm:"TINYINT"`
	SmallInt  int16 `xorm:"SMALLINT"`
	MediumInt int32 `xorm:"MEDIUMINT"`
	Int       int   `xorm:"INT"`
	Integer   int   `xorm:"INTEGER"`
	BigInt    int64 `xorm:"BIGINT"`

	Char       string `xorm:"CHAR(12)"`
	Varchar    string `xorm:"VARCHAR(54)"`
	TinyText   string `xorm:"TINYTEXT"`
	Text       string `xorm:"TEXT"`
	MediumText string `xorm:"MEDIUMTEXT"`
	LongText   string `xorm:"LONGTEXT"`
	Binary     []byte `xorm:"BINARY(23)"`
	VarBinary  []byte `xorm:"VARBINARY(12)"`

	Date      time.Time `xorm:"DATE"`
	DateTime  time.Time `xorm:"DATETIME"`
	Time      time.Time `xorm:"TIME"`
	TimeStamp time.Time `xorm:"TIMESTAMP"`

	Decimal float64 `xorm:"DECIMAL"`
	Numeric float64 `xorm:"NUMERIC"`

	Real   float32 `xorm:"REAL"`
	Float  float32 `xorm:"FLOAT"`
	Double float64 `xorm:"DOUBLE"`

	TinyBlob   []byte `xorm:"TINYBLOB"`
	Blob       []byte `xorm:"BLOB"`
	MediumBlob []byte `xorm:"MEDIUMBLOB"`
	LongBlob   []byte `xorm:"LONGBLOB"`
	Bytea      []byte `xorm:"BYTEA"`

	Bool bool `xorm:"BOOL"`

	Serial int `xorm:"SERIAL"`
	//BigSerial int64 `xorm:"BIGSERIAL"`
}

func testColTypes(engine *Engine, t *testing.T) {
	err := engine.DropTables(&allCols{})
	if err != nil {
		t.Error(err)
		panic(err)
	}

	err = engine.CreateTables(&allCols{})
	if err != nil {
		t.Error(err)
		panic(err)
	}

	ac := &allCols{
		1,
		4,
		8,
		16,
		32,
		64,
		128,

		"123",
		"fafdafa",
		"fafafafdsafdsafdaf",
		"fdsafafdsafdsaf",
		"fafdsafdsafdsfadasfsfafd",
		"fadfdsafdsafasfdasfds",
		[]byte("fdafsafdasfdsafsa"),
		[]byte("fdsafsdafs"),

		time.Now(),
		time.Now(),
		time.Now(),
		time.Now(),

		1.34,
		2.44302346,

		1.3344,
		2.59693523,
		3.2342523543,

		[]byte("fafdasf"),
		[]byte("fafdfdsafdsafasf"),
		[]byte("faffadsfdsdasf"),
		[]byte("faffdasfdsadasf"),
		[]byte("fafasdfsadffdasf"),

		true,

		21,
	}

	cnt, err := engine.Insert(ac)
	if err != nil {
		t.Error(err)
		panic(err)
	}
	if cnt != 1 {
		err = errors.New("insert return not 1")
		t.Error(err)
		panic(err)
	}
	newAc := &allCols{}
	has, err := engine.Get(newAc)
	if err != nil {
		t.Error(err)
		panic(err)
	}
	if !has {
		err = errors.New("error no ideas")
		t.Error(err)
		panic(err)
	}

	// don't use this type as query condition
	newAc.Real = 0
	newAc.Float = 0
	newAc.Double = 0
	cnt, err = engine.Delete(newAc)
	if err != nil {
		t.Error(err)
		panic(err)
	}
	if cnt != 1 {
		err = errors.New(fmt.Sprintf("delete error, deleted counts is %v", cnt))
		t.Error(err)
		panic(err)
	}
}

type MyInt int
type MyUInt uint
type MyFloat float64
type MyString string

/*func (s *MyString) FromDB(data []byte) error {
	reflect.
	s MyString(data)
	return nil
}

func (s *MyString) ToDB() ([]byte, error) {
	return []byte(string(*s)), nil
}*/

type MyStruct struct {
	Type      MyInt
	U         MyUInt
	F         MyFloat
	S         MyString
	IA        []MyInt
	UA        []MyUInt
	FA        []MyFloat
	SA        []MyString
	NameArray []string
	Name      string
	UIA       []uint
	UIA8      []uint8
	UIA16     []uint16
	UIA32     []uint32
	UIA64     []uint64
	UI        uint
	//C64       complex64
	MSS map[string]string
}

func testCustomType(engine *Engine, t *testing.T) {
	err := engine.DropTables(&MyStruct{})
	if err != nil {
		t.Error(err)
		panic(err)
		return
	}

	err = engine.CreateTables(&MyStruct{})
	i := MyStruct{Name: "Test", Type: MyInt(1)}
	i.U = 23
	i.F = 1.34
	i.S = "fafdsafdsaf"
	i.UI = 2
	i.IA = []MyInt{1, 3, 5}
	i.UIA = []uint{1, 3}
	i.UIA16 = []uint16{2}
	i.UIA32 = []uint32{4, 5}
	i.UIA64 = []uint64{6, 7, 9}
	i.UIA8 = []uint8{1, 2, 3, 4}
	i.NameArray = []string{"ssss", "fsdf", "lllll, ss"}
	i.MSS = map[string]string{"s": "sfds,ss", "x": "lfjljsl"}
	cnt, err := engine.Insert(&i)
	if err != nil {
		t.Error(err)
		panic(err)
		return
	}
	if cnt != 1 {
		err = errors.New("insert not returned 1")
		t.Error(err)
		panic(err)
		return
	}

	fmt.Println(i)
	has, err := engine.Get(&i)
	if err != nil {
		t.Error(err)
		panic(err)
	} else if !has {
		t.Error(errors.New("should get one record"))
		panic(err)
	}

	ss := []MyStruct{}
	err = engine.Find(&ss)
	if err != nil {
		t.Error(err)
		panic(err)
	}
	fmt.Println(ss)

	sss := MyStruct{}
	has, err = engine.Get(&sss)
	if err != nil {
		t.Error(err)
		panic(err)
	}
	fmt.Println(sss)

	if has {
		cnt, err := engine.Delete(&sss)
		if err != nil {
			t.Error(err)
			panic(err)
		}
		if cnt != 1 {
			t.Error(errors.New("delete error"))
			panic(err)
		}
	}
}

type UserCU struct {
	Id      int64
	Name    string
	Created time.Time `xorm:"created"`
	Updated time.Time `xorm:"updated"`
}

func testCreatedAndUpdated(engine *Engine, t *testing.T) {
	u := new(UserCU)
	err := engine.DropTables(u)
	if err != nil {
		t.Error(err)
		panic(err)
	}

	err = engine.CreateTables(u)
	if err != nil {
		t.Error(err)
		panic(err)
	}

	u.Name = "sss"
	cnt, err := engine.Insert(u)
	if err != nil {
		t.Error(err)
		panic(err)
	}
	if cnt != 1 {
		err = errors.New("insert not returned 1")
		t.Error(err)
		panic(err)
		return
	}

	u.Name = "xxx"
	cnt, err = engine.Id(u.Id).Update(u)
	if err != nil {
		t.Error(err)
		panic(err)
	}
	if cnt != 1 {
		err = errors.New("update not returned 1")
		t.Error(err)
		panic(err)
		return
	}

	u.Id = 0
	u.Created = time.Now().Add(-time.Hour * 24 * 365)
	u.Updated = u.Created
	fmt.Println(u)
	cnt, err = engine.NoAutoTime().Insert(u)
	if err != nil {
		t.Error(err)
		panic(err)
	}
	if cnt != 1 {
		err = errors.New("insert not returned 1")
		t.Error(err)
		panic(err)
		return
	}
}

type IndexOrUnique struct {
	Id        int64
	Index     int `xorm:"index"`
	Unique    int `xorm:"unique"`
	Group1    int `xorm:"index(ttt)"`
	Group2    int `xorm:"index(ttt)"`
	UniGroup1 int `xorm:"unique(lll)"`
	UniGroup2 int `xorm:"unique(lll)"`
}

func testIndexAndUnique(engine *Engine, t *testing.T) {
	err := engine.CreateTables(&IndexOrUnique{})
	if err != nil {
		t.Error(err)
		//panic(err)
	}

	err = engine.DropTables(&IndexOrUnique{})
	if err != nil {
		t.Error(err)
		//panic(err)
	}

	err = engine.CreateTables(&IndexOrUnique{})
	if err != nil {
		t.Error(err)
		//panic(err)
	}

	err = engine.CreateIndexes(&IndexOrUnique{})
	if err != nil {
		t.Error(err)
		//panic(err)
	}

	err = engine.CreateUniques(&IndexOrUnique{})
	if err != nil {
		t.Error(err)
		//panic(err)
	}
}

type IntId struct {
	Id   int
	Name string
}

type Int32Id struct {
	Id   int32
	Name string
}

func testIntId(engine *Engine, t *testing.T) {
	err := engine.DropTables(&IntId{})
	if err != nil {
		t.Error(err)
		panic(err)
	}

	err = engine.CreateTables(&IntId{})
	if err != nil {
		t.Error(err)
		panic(err)
	}

	_, err = engine.Insert(&IntId{Name: "test"})
	if err != nil {
		t.Error(err)
		panic(err)
	}
}

func testInt32Id(engine *Engine, t *testing.T) {
	err := engine.DropTables(&Int32Id{})
	if err != nil {
		t.Error(err)
		panic(err)
	}

	err = engine.CreateTables(&Int32Id{})
	if err != nil {
		t.Error(err)
		panic(err)
	}

	_, err = engine.Insert(&Int32Id{Name: "test"})
	if err != nil {
		t.Error(err)
		panic(err)
	}
}

func testMetaInfo(engine *Engine, t *testing.T) {
	tables, err := engine.DBMetas()
	if err != nil {
		t.Error(err)
		panic(err)
	}

	for _, table := range tables {
		fmt.Println(table.Name)
		for _, col := range table.Columns {
			fmt.Println(col.String(engine.dialect))
		}

		for _, index := range table.Indexes {
			fmt.Println(index.Name, index.Type, strings.Join(index.Cols, ","))
		}
	}
}

func testIterate(engine *Engine, t *testing.T) {
	err := engine.Omit("is_man").Iterate(new(Userinfo), func(idx int, bean interface{}) error {
		user := bean.(*Userinfo)
		fmt.Println(idx, "--", user)
		return nil
	})

	if err != nil {
		t.Error(err)
		panic(err)
	}
}

type StrangeName struct {
	Id_t int64 `xorm:"pk autoincr"`
	Name string
}

func testStrangeName(engine *Engine, t *testing.T) {
	err := engine.DropTables(new(StrangeName))
	if err != nil {
		t.Error(err)
	}

	err = engine.CreateTables(new(StrangeName))
	if err != nil {
		t.Error(err)
	}

	_, err = engine.Insert(&StrangeName{Name: "sfsfdsfds"})
	if err != nil {
		t.Error(err)
	}

	beans := make([]StrangeName, 0)
	err = engine.Find(&beans)
	if err != nil {
		t.Error(err)
	}
}

type Version struct {
	Id   int64
	Name string
	Ver  int `xorm:"version"`
}

func testVersion(engine *Engine, t *testing.T) {
	err := engine.DropTables(new(Version))
	if err != nil {
		t.Error(err)
		return
	}

	err = engine.CreateTables(new(Version))
	if err != nil {
		t.Error(err)
		return
	}

	ver := &Version{Name: "sfsfdsfds"}
	_, err = engine.Cols("name").Insert(ver)
	if err != nil {
		t.Error(err)
		return
	}

	newVer := new(Version)
	has, err := engine.Id(ver.Id).Get(newVer)
	if err != nil {
		t.Error(err)
		return
	}

	if !has {
		t.Error(errors.New(fmt.Sprintf("no version id is %v", ver.Id)))
		return
	}

	newVer.Name = "-------"
	_, err = engine.Id(ver.Id).Update(newVer, &Version{Ver: newVer.Ver})
	if err != nil {
		t.Error(err)
		return
	}

	has, err = engine.Id(ver.Id).Get(newVer)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(ver)

	newVer.Name = "-------"
	_, err = engine.Id(ver.Id).Update(newVer, &Version{Ver: newVer.Ver})
	if err != nil {
		t.Error(err)
		return
	}
}

func testDistinct(engine *Engine, t *testing.T) {
	users := make([]Userinfo, 0)
	err := engine.Distinct("departname").Find(&users)
	if err != nil {
		t.Error(err)
		panic(err)
	}
	if len(users) != 1 {
		t.Error(err)
		panic(errors.New("should be one record"))
	}

	fmt.Println(users)

	type Depart struct {
		Departname string
	}

	users2 := make([]Depart, 0)
	err = engine.Distinct("departname").Table(new(Userinfo)).Find(&users2)
	if err != nil {
		t.Error(err)
		panic(err)
	}
	if len(users2) != 1 {
		t.Error(err)
		panic(errors.New("should be one record"))
	}
	fmt.Println(users2)
}

func testUseBool(engine *Engine, t *testing.T) {
	cnt1, err := engine.Count(&Userinfo{})
	if err != nil {
		t.Error(err)
		panic(err)
	}

	users := make([]Userinfo, 0)
	err = engine.Find(&users)
	if err != nil {
		t.Error(err)
		panic(err)
	}
	var fNumber int64
	for _, u := range users {
		if u.IsMan == false {
			fNumber += 1
		}
	}

	cnt2, err := engine.UseBool().Update(&Userinfo{IsMan: true})
	if err != nil {
		t.Error(err)
		panic(err)
	}
	if fNumber != cnt2 {
		fmt.Println("cnt1", cnt1, "fNumber", fNumber, "cnt2", cnt2)
		/*err = errors.New("Updated number is not corrected.")
		t.Error(err)
		panic(err)*/
	}

	_, err = engine.Update(&Userinfo{IsMan: true})
	if err == nil {
		err = errors.New("error condition")
		t.Error(err)
		panic(err)
	}
}

func testBool(engine *Engine, t *testing.T) {
	_, err := engine.UseBool().Update(&Userinfo{IsMan: true})
	if err != nil {
		t.Error(err)
		panic(err)
	}
	users := make([]Userinfo, 0)
	err = engine.Find(&users)
	if err != nil {
		t.Error(err)
		panic(err)
	}
	for _, user := range users {
		if !user.IsMan {
			err = errors.New("update bool or find bool error")
			t.Error(err)
			panic(err)
		}
	}

	_, err = engine.UseBool().Update(&Userinfo{IsMan: false})
	if err != nil {
		t.Error(err)
		panic(err)
	}
	users = make([]Userinfo, 0)
	err = engine.Find(&users)
	if err != nil {
		t.Error(err)
		panic(err)
	}
	for _, user := range users {
		if user.IsMan {
			err = errors.New("update bool or find bool error")
			t.Error(err)
			panic(err)
		}
	}
}

func testAll(engine *Engine, t *testing.T) {
	fmt.Println("-------------- directCreateTable --------------")
	directCreateTable(engine, t)
	fmt.Println("-------------- insert --------------")
	insert(engine, t)
	fmt.Println("-------------- query --------------")
	testQuery(engine, t)
	fmt.Println("-------------- exec --------------")
	exec(engine, t)
	fmt.Println("-------------- insertAutoIncr --------------")
	insertAutoIncr(engine, t)
	fmt.Println("-------------- insertMulti --------------")
	insertMulti(engine, t)
	fmt.Println("-------------- insertTwoTable --------------")
	insertTwoTable(engine, t)
	fmt.Println("-------------- update --------------")
	update(engine, t)
	fmt.Println("-------------- testdelete --------------")
	testdelete(engine, t)
	fmt.Println("-------------- get --------------")
	get(engine, t)
	fmt.Println("-------------- cascadeGet --------------")
	cascadeGet(engine, t)
	fmt.Println("-------------- find --------------")
	find(engine, t)
	fmt.Println("-------------- find2 --------------")
	find2(engine, t)
	fmt.Println("-------------- findMap --------------")
	findMap(engine, t)
	fmt.Println("-------------- findMap2 --------------")
	findMap2(engine, t)
	fmt.Println("-------------- count --------------")
	count(engine, t)
	fmt.Println("-------------- where --------------")
	where(engine, t)
	fmt.Println("-------------- in --------------")
	in(engine, t)
	fmt.Println("-------------- limit --------------")
	limit(engine, t)
	fmt.Println("-------------- order --------------")
	order(engine, t)
	fmt.Println("-------------- join --------------")
	join(engine, t)
	fmt.Println("-------------- having --------------")
	having(engine, t)
}

func testAll2(engine *Engine, t *testing.T) {
	fmt.Println("-------------- combineTransaction --------------")
	combineTransaction(engine, t)
	fmt.Println("-------------- table --------------")
	table(engine, t)
	fmt.Println("-------------- createMultiTables --------------")
	createMultiTables(engine, t)
	fmt.Println("-------------- tableOp --------------")
	tableOp(engine, t)
	fmt.Println("-------------- testCols --------------")
	testCols(engine, t)
	fmt.Println("-------------- testCharst --------------")
	testCharst(engine, t)
	fmt.Println("-------------- testStoreEngine --------------")
	testStoreEngine(engine, t)
	fmt.Println("-------------- testExtends --------------")
	testExtends(engine, t)
	fmt.Println("-------------- testColTypes --------------")
	testColTypes(engine, t)
	fmt.Println("-------------- testCustomType --------------")
	testCustomType(engine, t)
	fmt.Println("-------------- testCreatedAndUpdated --------------")
	testCreatedAndUpdated(engine, t)
	fmt.Println("-------------- testIndexAndUnique --------------")
	testIndexAndUnique(engine, t)
	fmt.Println("-------------- testIntId --------------")
	//testIntId(engine, t)
	fmt.Println("-------------- testInt32Id --------------")
	//testInt32Id(engine, t)
	fmt.Println("-------------- testMetaInfo --------------")
	testMetaInfo(engine, t)
	fmt.Println("-------------- testIterate --------------")
	testIterate(engine, t)
	fmt.Println("-------------- testStrangeName --------------")
	testStrangeName(engine, t)
	fmt.Println("-------------- testVersion --------------")
	testVersion(engine, t)
	fmt.Println("-------------- testDistinct --------------")
	testDistinct(engine, t)
	fmt.Println("-------------- testUseBool --------------")
	testUseBool(engine, t)
	fmt.Println("-------------- testBool --------------")
	testBool(engine, t)
	fmt.Println("-------------- transaction --------------")
	transaction(engine, t)
}
