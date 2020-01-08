# pkg-viewer

# How to run

Run ```go run .``` from root of the repository.

# Endpoints

## /

Reads the status file and intializes the app. Other endpoints don't work before this endpoint has been visited once.

## /index

Lists all the software packages.

## /packages/{package}

Provides information, such as description, packages it depends on and packages that are dependant of it, for a single package.
