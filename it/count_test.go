package it

import (
	"context"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestCount(t *testing.T) {

	Convey("统计测试", t, func() {
		Convey("正常", func() {
			g0 := newEnvelopeGoods(1)
			rs, err := db.InsertContext(context.Background(), &g0)
			So(err, ShouldBeNil)
			So(rs, ShouldNotBeNil)
			id, err := rs.LastInsertId()
			So(err, ShouldBeNil)
			So(id > 0, ShouldBeTrue)
			rows, err := rs.RowsAffected()
			So(err, ShouldBeNil)
			So(rows, ShouldEqual, 1)

			query := "select count(1) from red_envelope_goods where id>=?"
			var x int64
			row := db.QueryRow(query, id)
			err = row.Scan(&x)

			So(err, ShouldBeNil)
			So(x >= 1, ShouldBeTrue)
			fmt.Println(x)
		})
	})
}

func TestCountInt32(t *testing.T) {

	Convey("统计测试", t, func() {
		Convey("正常", func() {
			g0 := newEnvelopeGoods(1)
			rs, err := db.InsertContext(context.Background(), &g0)
			So(err, ShouldBeNil)
			So(rs, ShouldNotBeNil)
			id, err := rs.LastInsertId()
			So(err, ShouldBeNil)
			So(id > 0, ShouldBeTrue)
			rows, err := rs.RowsAffected()
			So(err, ShouldBeNil)
			So(rows, ShouldEqual, 1)

			query := "select count(1) from red_envelope_goods where id>=?"
			x, err := db.GetInt32(query, id)
			So(err, ShouldBeNil)
			So(x >= 1, ShouldBeTrue)
			fmt.Println(x)
		})
	})
}

func TestCountInt64(t *testing.T) {

	Convey("统计测试", t, func() {
		Convey("正常", func() {
			g0 := newEnvelopeGoods(1)
			rs, err := db.InsertContext(context.Background(), &g0)
			So(err, ShouldBeNil)
			So(rs, ShouldNotBeNil)
			id, err := rs.LastInsertId()
			So(err, ShouldBeNil)
			So(id > 0, ShouldBeTrue)
			rows, err := rs.RowsAffected()
			So(err, ShouldBeNil)
			So(rows, ShouldEqual, 1)

			query := "select count(1) from red_envelope_goods where id>= ?"
			x, err := db.GetInt64(query, id)
			So(err, ShouldBeNil)
			So(x >= 1, ShouldBeTrue)
			fmt.Println(x)
		})
	})
}
