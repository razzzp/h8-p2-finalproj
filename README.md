# h8_p2_finalproj

Car Rental App

- Planning, concepts & ERD can be found in `docs` folder

## Decsription

### User

- Client can register users by providing:
  - By providing:
    - Email
    - Password
    - Desposit
  - Check if email already registered, return error if true
- Client can login a user by providing:
  - By Providing:
    - Email
    - Password
  - Check email and password
    - If ok return JWT token
    - Else return 401
- Client can search for available cars for rent
  - provide start date and time
  - provide end date and time
  - can filter by vehicle type/model
  - returns all available cars
- Client can book a rental for a given car
  - provide start date and time
  - provide end date and time
  - provide vehicle id to rent
  - If car not available return error
  - Returns the payment link to the client
- Callback for payment gateway to update status of rental

### Admin

- Can manage cars avaialble
  - Create, Update and Delete

## Technologies

- Echo
- Postgres
- Xandit