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
  - provide start date
  - provide end date
  - provide vehicle id to rent
  - If car not available return error
  - Returns the payment link to the client
- Client can make payment at the payment gateway
  - Callback for payment gateway to update status of rental
- Client can top up his/her deposit
- Will be notifed by email on registration, booking & payment

### Admin

- Can manage cars avaialble
  - Create, Update and Delete

## Technologies

- Echo
- Gorm
- Postgres
- Xendit
- Gmail SMTP
- Testify

## Environment Variables

```
DB_HOST=
DB_PORT=
DB_NAME=
DB_USER=
DB_PASS=
PORT=
JWT_KEY=
XENDIT_API_KEY=
XENDIT_WEBHOOK_TOKEN=
XENDIT_INVOICE_CALLBACK=
SMTP_HOST=
SMTP_PORT=
SMTP_USER=
SMTP_PASS=
```