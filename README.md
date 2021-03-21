# Luban（鲁班）

An project generator aimed to D.R.Y while creating new project based on DDD. visit https://github.com/leopku/luban for more details.

## Get source code

- [Github](https://github.com/leopku/luban)
- [Gitee](https://gitee.com/leopku/luban)

## Features

- [x] Generate model layer(go struct files) from database
  - [x] MySQL
  - [x] PostgreSQL
  - [x] Sqlite3
- [ ] Generate repository layer
  - [x] Interface
  - [x] Implement CRUD
    - [x] squirrel with raw sql
    - [-] queryset with gorm(W.I.P)
  - [x] DAO / ORM
    - [x] sqlx
    - [x] gorm
    - [ ] ~~xorm~~
- [ ] Generate service layer
- [ ] Generate deliver layer
  - [ ] RESTful deliver
    - [ ] gin
    - [ ] gofiber
    - [ ] echo
    - [ ] iris
  - [ ] GraphQL deliver
- [ ] Add go-swagger support
- [ ] Admin dashboard

# Special thanks

This project was heavyly inspired by [gen](https://github.com/smallnest/gen). And `templates/mapping.json` was directly copied from [gen](https://github.com/smallnest/gen) project. Thanks you, [@smallnest](https://github.com/smallnest).

# Credit

Thanks for this projects to build 

- [Cobra](https://github.com/spf13/cobra)
- [jennifer](https://github.com/dave/jennifer)
- [schema](https://github.com/jimsmart/schema)
- [zerolog](https://github.com/rs/zerolog)
- [wire](https://github.com/google/wire)
- [strcase](https://github.com/iancoleman/strcase)
- [xstrings](https://github.com/huandu/xstrings)
- [go-pluralize](https://github.com/gertd/go-pluralize)
- [go-commons-lang](https://github.com/agrison/go-commons-lang)
- [testify](https://github.com/stretchr/testify)
- [go-mysql-server](https://github.com/src-d/go-mysql-server)
- [jsoniter](https://github.com/json-iterator/go)
- [squirrel](github.com/Masterminds/squirrel)
- [sqlx](github.com/jmoiron/sqlx)
- [queryset](github.com/jirfag/go-queryset)

# Feedback&Issue

https://github.com/leopku/luban/issues
