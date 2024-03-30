package main

import (
	"database/sql"
	_ "database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"os"
	"strings"
)

type base struct {
	rack_name       string
	prev_rack       string
	order_number    int
	prev_order      int
	product_name    string
	prev_product    int
	product_id      int
	additional_rack string
	quantity        int
}

var DB *sql.DB

func ScanNPrintdb(base *base, rows *sql.Rows) {
	if !rows.Next() {
		return
	} else {
		base.prev_rack = base.rack_name
		err := rows.Scan(&base.order_number, &base.quantity, &base.product_name, &base.product_id, &base.rack_name, &base.additional_rack)
		if err != nil {
			panic("")
		}
		basePrint(base, rows)
	}
}
func closeDb() error {
	return DB.Close()
}
func openDb() error {
	var err error

	DB, err = sql.Open("postgres", "user=postgres dbname=shop sslmode=disable")
	if err != nil {
		return err
	}
	return nil
}
func PrintAddRack(base *base, rows *sql.Rows) {

	fmt.Printf(base.additional_rack)
	base.additional_rack = ""
	base.prev_product = base.product_id
	base.prev_order = base.order_number
	if rows.Next() {
		err := rows.Scan(&base.order_number, &base.quantity, &base.product_name, &base.product_id, &base.rack_name, &base.additional_rack)
		if err != nil {
			panic("")
		}
		if base.prev_product == base.product_id && base.prev_order == base.order_number {
			fmt.Print(",")
			PrintAddRack(base, rows)
		} else {
			fmt.Printf("\n\n")
			basePrint(base, rows)
		}
	}

}
func basePrint(base *base, rows *sql.Rows) {
	if base.prev_rack != base.rack_name {
		fmt.Printf("===Стеллаж %s\n", base.rack_name)
	}
	fmt.Printf("%s (id=%d)\n", base.product_name, base.product_id)
	fmt.Printf("заказ %d, %d шт\n", base.order_number, base.quantity)
	base.prev_rack = base.rack_name
	if base.additional_rack != "" {
		fmt.Printf("доп стеллаж: ")
		PrintAddRack(base, rows)
	}
	fmt.Printf("\n")
	ScanNPrintdb(base, rows)
}
func orders() string {
	var orders string
	orders = os.Args[1]
	fmt.Printf("=+=+=+=\nСтраница сборки заказов " + orders + "\n")
	str := strings.Replace(orders, ",", " OR order_number = ", -1)
	return str
}
func main() {
	err := openDb()
	defer closeDb()
	if err != nil {
		log.Printf("error connecting to postgres database %v", err)
	}
	rows, err := DB.Query("CREATE TEMPORARY TABLE temp_table AS SELECT order_number, quantity, product_name, product.product_id, prime_rack.rack_name AS prime_rack, additional_rack.rack_name from product\n\t" +
		"JOIN  order_details ON product.product_id = order_details.product_id\n    " +
		"JOIN  orders ON order_details.order_id = orders.order_id  \n    " +
		"JOIN  rack prime_rack ON product.rack_prime_id = prime_rack.rack_id \n\t" +
		"LEFT JOIN  product_racks ON product.product_id = product_racks.product_id\n \t" +
		"LEFT JOIN rack additional_rack ON product_racks.rack_id = additional_rack.rack_id\n    " +
		"WHERE order_number =" + orders() +
		";\nUPDATE temp_table\n\t\tSET rack_name = COALESCE(rack_name,'')\n\t\tWHERE rack_name IS NULL;\n" +
		"SELECT * FROM temp_table " +
		"ORDER BY prime_rack;")
	if err != nil {
		panic("cant select orders form database")
	}
	var base base
	ScanNPrintdb(&base, rows)
}
