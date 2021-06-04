package dbconn

import (
    "database/sql"
    "log"
    "githubmcom/mebusy/goweb/db"
    "fmt"
)

func getMysqlDB() *sql.DB {
    return db.GetMysqlDB()
}


var sql_create_tbl = []string{
    `
    create table if not exists __pmadmin_users {
        id int NOT NULL AUTO_INCREMENT ,
        user  varchar(16) NOT NULL ,
        password  varchar(16) NOT NULL ,
        PRIMARY KEY (id),
        UNIQUE key ( user )
    }
    `,
}

func PrepareLoginTable() {
    db := getMysqlDB()
    for _, v := range sql_create_tbl {
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

