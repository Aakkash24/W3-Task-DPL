## DPL Project
This project is a Go-lang implementation of DPL. It consists of two major components:

Gateway: The gateway is responsible for routing requests to the appropriate service and performing other tasks such as tournament simulation.
Services: The services are responsible for handling specific operations, such as creating, reading, updating, and deleting players and teams.
The services are further divided into two sub-folders:

Player service: This service handles all operations related to players, such as creating, reading, updating, and deleting players.
Team service: This service handles all operations related to teams, such as creating, reading, updating, and deleting teams.
Each service is further divided into two sub-folders:

Cru service: This service handles create, read, and update operations.
Delete service: This service handles delete operations.

Example Usage
To create a new player, send a POST request to the following endpoint:

http://localhost:8080/simulateTournament
Simulates the tournament based on the teams present in the DB.

http://localhost:8080/createPlayer
The request body should be a JSON object containing the player's information, such as name, email, and phone number.

To get all players, send a GET request to the following endpoint:

http://localhost:8080/getAllTeams
This will return a JSON array of all players.

To get a specific player, send a GET request to the following endpoint:

http://localhost:8080/readPlayer/{playerId}
Replace {playerId} with the ID of the player you want to retrieve.

To update a player, send a PUT request to the following endpoint:

http://localhost:8080/editPlayer/{playerId}
Replace {playerId} with the ID of the player you want to update. The request body should be a JSON object containing the player's updated information.

To delete a player, send a DELETE request to the following endpoint:

http://localhost:8080/deletePlayer/{playerId}
Replace {playerId} with the ID of the player you want to delete.

You can use the same endpoints to manage teams.


REPLACE localhost with your Host IP
