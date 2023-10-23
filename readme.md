# Go-lang DPL Implementation
This project is an implementation of DPL (Data Persistence Layer) in Go-lang. It consists of two major components:

## Gateway
The gateway is responsible for routing requests to the appropriate service and performing other tasks, such as tournament simulation.

## Services
The services are responsible for handling specific operations related to players and teams. Services are divided into two main categories:

### Player Service
This service manages all player-related operations, including creating, reading, updating, and deleting players.

#### Sub-Services
##### Cru Service: Handles create, read, and update operations for players.
##### Delete Service: Handles delete operations for players.

### Team Service
This service manages all team-related operations, including creating, reading, updating, and deleting teams.

#### Sub-Services
##### Cru Service: Handles create, read, and update operations for teams.
##### Delete Service: Handles delete operations for teams.


### Example Usage
Here's how you can interact with the services:

##### Creating a New Player
To create a new player, send a POST request to the following endpoint:
Endpoint: http://localhost:8080/createPlayer
Request Body: JSON object containing player information, such as name, email, and phone number.

##### Simulating a Tournament
Simulate a tournament based on the teams present in the database:
Endpoint: http://localhost:8080/simulateTournament

##### Getting All Players
To retrieve a list of all players, send a GET request to the following endpoint:
Endpoint: http://localhost:8080/getAllPlayers
Response: JSON array containing all players.

##### Getting a Specific Player
To retrieve a specific player, send a GET request with the player's ID to the following endpoint:
Endpoint: http://localhost:8080/readPlayer/{playerId}
Replace {playerId} with the ID of the player you want to retrieve.

##### Updating a Player
To update a player's information, send a PUT request with the player's ID and the updated information to the following endpoint:
Endpoint: http://localhost:8080/editPlayer/{playerId}
Replace {playerId} with the ID of the player you want to update.
Request Body: JSON object containing the player's updated information.

##### Deleting a Player
To delete a player, send a DELETE request with the player's ID to the following endpoint:
Endpoint: http://localhost:8080/deletePlayer/{playerId}
Replace {playerId} with the ID of the player you want to delete.
You can use the same set of endpoints to manage teams in a similar fashion.

Note: Replace localhost in the endpoints with your host's IP address.
