# minesweeper
Minesweeper game in go

To start the API locally : docker-compose up --build
To view the Swagger API Docs: http://localhost:8080/api/v1/swagger/swagger.html

#Deployed in Now (https://zeit.co/) 
url: https://minesweeper.javieriseruk.now.sh/api/v1

The version deployed in Zeit does not support the swagger doc yet.

The persistence supported are:
 
sqlite (for Zeit Now deployment)
in-memory (for testing purposes)

The Rest API is developed with the Gin Gonic Framework


