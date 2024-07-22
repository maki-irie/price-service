# price-service

## Requirements
Develop a REST service with one HTTP GET resource API as:
HTTP GET <base-IP>/price?jwt="7394ytuh.473tirg.489th"
Response {"quality":23,"price":123.34}

Where JWT will be signed JWT and the pub-key will be made available via rest service, body of JWT will contain{"item":"item-name","vat-incl":true,"quantity":124}

pub-key-service via a GET HTTP REST client call to <security-service-IP>/pkey?client="clientname"
Response{"pkey":"dsvfehbouenfdnvl"}

Service will need to call DB to retrieve details for the item

DB name testdb, one "item_table" with following columns:
{
    name varchar(64) primary key,
    bar_code varchar(32),
    country_of_origin varchar(64),
    tot_quantity integer,
    available boolean,
    expiry_date date(dd-mm-yyyy),
    price integer
}

discount-service <disc-service-IP>/discount?quantity=435
Response (overloading the content just to execute more JSON unmarshaling)
{
    "item":"name","quantity":134,"applicable_to_EU":false, "shipping_cost":24, "discount":0, "shipping_time_days":5, "related_items":["item1","item2"]
}

Draft sequence:
Price API -> http get pkey -> chk signature and decode -> query DB for the item price -> http get discount -> conditionally add VAT -> apply discount -> marshal JSON -> send response

Thoughts:
1. configure DB thread pool and http client pool so bounded and tune via config
2. GORM might help for small data set
3. Intenally cache?
4. If DB and discount API are independent we could make use of go routines and channels, parallel run them
5. Just use go http handler for request and call other services via go routine, as seems to be synchronous steps

# Setting up postgres
docker pull postgres
docker run --name my-postgres -e POSTGRES_PASSWORD=mysecretpassword -d postgres
docker exec -it my-postgres psql -U postgres
