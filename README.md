# allytrak

# Vision

Our vision is to create a cutting-edge partner management platform that serves as the ultimate solution for businesses to seamlessly connect, collaborate, and cultivate strong partnerships. We envision a platform that empowers organizations to effortlessly onboard and enable partners, fostering transparency, trust, and alignment. Through advanced features, intuitive user experience, and data-driven insights, our platform will revolutionize partner management, driving accelerated growth, expanding market reach, and unlocking unparalleled opportunities for success in the dynamic business landscape.

# Docker & Postgres

$docker compose up

# Env:

Server: db POSTGRES_USER: rick POSTGRES_PASSWORD: picklerick POSTGRES_DB: GalacticFederation Adminer: localhost:3333 PSQL: docker

compose exec -it db psql -U rick -d GalacticFederation

# connection string for goose:

host=localhost port=5432 user=rick password=picklerick dbname=GalacticFederation sslmode=disable

# Goose install:

go install github.com/pressly/goose/v3/cmd/goose@latest export GOPATH=$HOME/alex/git export PATH=$PATH:/home/alex/go:$GOPATH/bi
