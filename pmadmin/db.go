package pmadmin

// PM admin db ops

import (
    "database/sql"
    "log"
    "fmt"
)

var db *sql.DB;

var sql_create_tbl = []string{

    `CREATE TABLE IF NOT EXISTS _pmadmin_users (
      id int(11) NOT NULL AUTO_INCREMENT,
      user varchar(16) NOT NULL,
      password varchar(64) NOT NULL,
      PRIMARY KEY (id),
      UNIQUE KEY (user)
    )
    `,
}

func isValidUser( user,  password string ) bool {
    var id int
    err := db.QueryRow( `select id from _pmadmin_users where user = ? AND password = ? UNION select 0 `, user, password ).Scan( &id )
    if err != nil {
        log.Println( err )
    }
    return id > 0
}

func prepareLoginTable( _db *sql.DB) {
    db = _db

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

