# Ticket Allocation Coding Test

## Notes
I've implemented the 3 endpoints using Gin as the http server and postgres as the database engine.

I have included tests that most cover the "service" layer, im not happy with my implementation of the gin request response logic because the code is hard to test (controller package). I would refactor this in reality. An integration test suite would be useful because some logic is in the database engine (checking ticket allocation quantity) which isn't unit testable.

There is a bash script called `runme.sh` that will run the basic commands to setup and run the environment in docker compose

I used a transaction in postgres in an attempt to ensure that tickets and allocations are tied together so if something fails during the process it is all rolled back. 

## Problems

 - I'm unsure how concurrently safe this is, will postgres lock the ticket_option row whilst the rest of the SQL statements are executed and commit is ran on the transaction? Or should I have updated the allocation before creating the purchase and the tickets
 - Not enough code in the controller package to validate inputs to the API
 - Mostly just happy path validation (did test for allocation of ticket errors as this was a key requirement)
