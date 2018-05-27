cd ../../migrations
goose postgres "host=gis.cesfozorknmw.us-west-2.rds.amazonaws.com user=hesidoryn password=hesidoryn dbname=gis sslmode=disable" up
cd ../deploy/yalcha-server
up