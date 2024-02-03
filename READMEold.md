# Multi-Functional Golang Webserver

Golang Webserver 

## Features

- In-depth documentation an description of function an controllers
- Validation middleware for individual assessments implemented for each route 
- Validation handler for chained input validation individually customizable according to your own needs
- Flexible error handler covering 32 HTTP status codes
- Flexible response handler for positive 2** responses covering 10 HTTP status codes
- Central collection, overview and management of configuration values in config folder
- Central Regex Expression library with several predefined and tested expressions
- Setup of MongoDB connection including prepared templates to query, save and delete
- Logging with logrus to daily log files in log folder: Logging error, time, log level.
- Deletion of log files in log folder older than 5 days (default). Can be changed in base_config file with 'DeleteLogsAfterDays'
- AWS and Azure functions to upload and delete images to S3 or Blob and change their names
- Fully-featured cookie handler
- Strict file upload validations: number, sizes, type
- Context extensions enabling to reuse queried data amongst more than one controller within one route an thus reducing number of database requests
- JWT setup needed for authentication and authorization
- Practical examples for testing with Go Testify (Used package: https://github.com/stretchr/testify)
- Cron jobs implementation. Cron scheduler starts in a separate goroutine. (Used package: https://github.com/robfig/cron)
- Transactions implemented. Transactions in MongoDB ensure the atomicity, consistency, isolation, and durability (ACID) properties for multiple database operations, allowing developers to group multiple statements into a single unit of work that either succeeds entirely or fails without leaving the database in an inconsistent state.
- Encryption, Decryption, token and hashing functions 
- Email setup for sending emails (mailtrap as placeholder). Switch via config between testing or production email github.com/paulmuenzner/powerplantmanager
- AWS KMS management (must create credentials file -> see ToDos below)
- Sophisticated csv handling: Reading and writing maps with any nested level to csv files. Non existing csv files will be created automatically by using the provided keys in map for header. If new data is added to an existing csv the keys in a map of this new data are not in line with the header of a csv, an error will be returned.

## Routes

- Authentication routes: Registration, registration verification, login
- File routes: Upload

## Security

- Defends against Slowloris attacks to some extent

## ToDo
- You must create credentials file for aws access to use aws --> '~/.aws/credentials'.  The tilde (~) symbol is typically expanded to the home directory by the operating system. The credentials file typically looks like this: 
``` [default]
aws_access_key_id = YOUR_ACCESS_KEY_ID
aws_secret_access_key = YOUR_SECRET_ACCESS_KEY
region = YOUR_REGION ```