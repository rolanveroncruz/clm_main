# clm_main
An Automated Certificate Lifecycle Managment system

This is the main server for the Certificate Lifecycle Management System. 
It will support the following:
## Agentless Discovery of External Certificates

## Agent Discovery, Renewal, and Application of Internal Certificates


---
# Journal
### May 28, 2025 11:34 am
I can already pull certificates from a public web server.
Need to store them in the database so the front end can display them. Need to define db schema.

### Jun 6, 2025 10:05am
I can already perform logging in, and have solved the CORS problem using a preflight middleware.
Now need to work on querying websites for their certs, storing these in the database, and allowing querying the db 
via http.

### Jun 10, 2025 8:51pm
Have already implemented:
-checking if server for discovery was already previously requested by user.
-saving discovered certs into DB.
TODO: handler for querying certs discovered.
