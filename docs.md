# Documentation

The project was implemented using Golang and a minimal frontend was done with html, css and javascript which is accessible at the index page of this app. This API tries to follow REST patterns and Golang standards.

Users are able to sign up, update their details, view their portfolio positions,request loans, view their loan balance.
The project is structured in the following formats:
- The root package contains domain/application data structures (models) and how to use them.
- The postgres package implements the interfaces defined in the domain types for working with a postgresql database. This enables the decoupling of our data models from the database and we could easily swap databases with little modification and eases testing
- The mock package contains stubs for testing, this helps with mocking out implementation specific details, e.g mocking a database call.
- The app package contains the core handlers and business logic for the application.
- The ui package contains the templates for the user interface
- The cmd holds the main package that ties the different parts together.

This structure helps with layering the codebase, so that the code remains readable, maintainable and avoids cyclic dependencies. 

## API Endpoints

### Authentication
Authentication is done using stateful tokens. The user must first register, then navigate to the api endpoint `/api/v1/auth/token` with their email and password to generate a new token, which should be passed in the `Authorization` header of subsequent api requests that require authentication.

### User Registration
Users can register using the following endpoint `/api/v1/users` providing their basic information including first_name, last_name, email and password.

### User profile update
Users can update their profile by making a PUT request to `/api/v1/users/:id` with the data to be updated.

### Portfolio
The `/api/v1/portfolio` endpoint returns the users portfolio.

### Portfolio Value
The `api/v1/portfolio/value` returns the total value of all positions held by the user

### Loan
The `/api/v1/loans` is used to create or fetch all loans based on the http method used (POST or GET). A user can request multiple loans, given that the total amount do not exceed 60% of the portfoilo value. An interest rate of 15% is applied to each loan request.
When retireving all loans, the payment due date and a prorated payment amount is provided.

## Testing
Although the project does not include extensive unit tests, a simple test was written for testing the user registration and update handlers in the /app/user_test.go file.

There is a [thunder client](/thunder-collection_Trove.json) (for vscode) collection to test available endpoints in the github repo.

## Running Locally
- Install the Go compiler
- Make sure to have postgresql >=13 installed
- clone the repo
- set the env vars as described in the .env.example file
- in the root directory, run `make run`

Although, not all features in the specifications are implemented, these are some areas that have not been implemented or could be improved upon: 
- Repayment of loans
- Integrating a 3rd party payment processor
- Allowing users to add and update their stocks
- Email notification