# Trove Engineering Challenge

## Context
The goal of this challenge is to test your proficiency as a software engineer or developer in the
stack of your choosing as well as familiarity with the financial market or ability to quickly pick
up skills on the fly. These are important technologies which are an integral part of the Trove
Tech stack. Furthermore, this challenge aims to test your coding ability, UI/UX implementation
skills, algorithmic prowess and your ability to dive into unknown territory.
## Situation
Meet Ope. Ope is an investor and stock trader. Ope holds a couple of stocks covering
Tesla(TSLA), Apple(AAPL) and Amazon (AMZN) with a portfolio value worth $10,000.00. Ope
would like to take a loan against his portfolio up to the tune of $3,000.00 for a period of either
6 to 12 months and withdraw the cash into his back account. Ope should be able to pay back
his loan prorated on a monthly basis.

## Challenge
You’re to design an {API, Web or Mobile App} to allow Ope to:
1. Create an account on your service - 10 points
2. Perform basic operations on his account such as updating his personal information and
changing his password - 10 points
3. Pull his portfolio positions - 5 points
4. Get his portfolio value - 5 points
5. Take a loan of up to 60% loan against his total portfolio value over period of 6 - 12
months, period should be decided by Ope - 10 points
6. View his active loan and balance - 5 points
7. View or get a prorated payment schedule over the loan period - 10 points
8. Payback his loan through any payment providers, please note using payment provider
test bed is perfectly acceptable provided it works end to end - 20 points

## Technical requirements
- Go 1.17
- PostgreSQL 13

## API Documentation and Testing
The endpoints exposed by the API are:
- POST /api/v1/users - Register a new user
- PUT /api/v1/users/{id} - update user info
- GET /api/v1/portfolio - get portfolio positions
- GET /api/v1/portfolio/value - get portfolio's total value
- POST /api/v1/loans - apply for a loan
- GET /api/v1/loans  - view loans
- POST /api/v1/auth/token - api authorization

There is a thunder-collection that can be imported with the thunder client in vscode to test the api.

## Running locally
- Clone the repo
- get all dependencies
- run `make run` to start the server