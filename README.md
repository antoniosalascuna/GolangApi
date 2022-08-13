# Instructions

Install Docker https://www.docker.com/products/docker-desktop

Install Docker compose

Run docker-compose build

Run docker-compose up

In visual studio code attach to running container.

# Recommend Tree

Reference https://github.com/golang-standards/project-layout

exitus/
├── cmd
│   ├── authtest
│   │   └── main.go
│   ├── backend
│   │   └── main.go
│   └── client
│       └── main.go
├── dev
│   ├── add_migration.sh
│   └── docker-compose.yml
├── Dockerfile
├── go.mod
├── go.sum
│   ├── 20190721131113_extensions.down.sql
│   ├── 20190721131113_extensions.up.sql
│   ├── 20190723044115_customer_projects.down.sql
│   ├── 20190723044115_customer_projects.up.sql
│   ├── 20190726175158_issues.down.sql
│   ├── 20190726175158_issues.up.sql
│   ├── 20190726201649_comments.down.sql
│   ├── 20190726201649_comments.up.sql
│   ├── bindata.go
│   ├── gen.go
│   ├── migrations_test.go
│   └── README.md
├── pkg
│   ├── api
│   │   ├── exitus.gen.go
│   │   ├── exitus.yml
│   │   └── gen.go
│   ├── auth
│   │   ├── scopes.go
│   │   └── user.go
│   ├── conf
│   │   ├── conf.go
│   │   └── conf_test.go
│   ├── db
│   │   ├── db.go
│   │   ├── dbtesting.go
│   │   ├── migrate.go
│   │   ├── sqlhooks.go
│   │   └── transactions.go
│   ├── env
│   │   └── env.go
│   ├── healthz
│   │   ├── healthz.go
│   │   └── healthz_test.go
│   ├── jwt
│   │   └── jwt.go
│   ├── metrics
│   │   └── metrics.go
│   ├── middleware
│   │   ├── jwt.go
│   │   └── middleware.go
│   ├── oidc
│   │   └── client.go
│   ├── server
│   │   ├── reflect.go
│   │   └── server.go
│   └── store
│       ├── comments.go
│       ├── comments_test.go
│       ├── customers.go
│       ├── customers_test.go
│       ├── issues.go
│       ├── issues_test.go
│       ├── migrate_test.go
│       ├── projects.go
│       ├── projects_test.go
│       └── store.go
└── README.md