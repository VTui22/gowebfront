package dbconn

import (
    "database/sql"
    "fmt"
    "log"
    // "reflect"
    "github.com/mebusy/goweb/db"
    // "github.com/mebusy/goweb/tools"
)


func getMysqlDB() *sql.DB {
    return db.GetMysqlDB()
}


func clearTransaction(tx *sql.Tx) {
    err := tx.Rollback()
    if err != sql.ErrTxDone && err != nil {
        log.Println(err)
    }
}



// =================================================================

var SQL_CREATE_TBL = []string{
`
CREATE TABLE if not exists wechat (
  game varchar(64) NOT NULL,
  access_token varchar(512) NOT NULL,
  expires_by int NOT NULL,
  PRIMARY KEY (game)
) ;
` ,

}

func PrepareMysqlTable() {
    db := getMysqlDB()
    for _, v := range SQL_CREATE_TBL {
        _, err := db.Exec(v)
        if err != nil {
            var code int
            _, err2 := fmt.Sscanf(err.Error(), "Error %d", &code)
            if err2 == nil && code == 1060 {
                log.Println("dup column , ignore this error:", err.Error())
                continue
            }
            log.Fatalln(err.Error())
        }
    }
}

