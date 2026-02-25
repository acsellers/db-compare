# Database Comparison

The goal of this project is to compare the features and code for the different major database mapping libraries
for the Go programming language. 

## Contributing

Not yet, but when I finish setting up the first couple of libraries, I'll start accepting pull requests.

The plan is that new libraries will be added by doing the following:

1. Create new directory in the docs/libraries folder with the files from the docs/libraries/template folder.
2. Fill out the files in that directory with the info about the library.
3. Implement the library actions for the store example for Mysql/Postgres/Sqlite using the template project.
4. Load the verification data into the database for the library.
5. Run the verification code to verify that the library works as expected.
6. Setup the sample list in the docs folder so samples will be displayed.
7. Make a pull request to pull in the new library.
8. I pull in the new library and run the benchmarks on the benchmark machine.
9. New library becomes available on the website.

## TODO

### Docs (initial)

- [x] Examples format
- [x] Samples format
- [x] Report Card format
- [x] Examples Draft
- [x] Rework the documentation info into folders and files.
- [ ] With Example
- [ ] Examples Final
- [x] Bob Description Draft
- [ ] Bob Description Final
- [x] Bob Report Card
- [ ] Bob Sample Mapping
- [x] sqlc Description Draft
- [ ] sqlc Description Final
- [x] sqlc Report Card
- [ ] sqlc Sample Mapping
- [x] Other Libraries Basic Info

### Bob (mysql)

- [x] Get Sale
- [x] Create Sale
- [ ] Customer Update (individual statements)
- [ ] Customer Update (temporary table)
- [x] Daily Revenue
- [x] Customer Sales
- [x] Daily Sold Items
- [x] General Sales Report
- [x] Weekly Sales Report
- [ ] Sale Search
- [x] Bulk Customers
- [ ] JSON Marshal/Unmarshal
- [ ] JSON Query
- [ ] CTE Example

### sqlc (mysql)

- [x] Get Sale
- [x] Create Sale
- [ ] Customer Update (individual statements)
- [ ] Customer Update (temporary table)
- [x] Daily Revenue
- [x] Customer Sales
- [x] Daily Sold Items
- [x] General Sales Report
- [x] Weekly Sales Report
- [ ] Sale Search
- [x] Bulk Customers
- [ ] JSON Marshal/Unmarshal
- [ ] JSON Query
- [ ] CTE Example

### Tests (mysql)

- [x] Get Sale
- [ ] Create Sale
- [ ] Customer Update (individual statements)
- [ ] Customer Update (temporary table)
- [x] Daily Revenue
- [x] Customer Sales
- [x] Daily Sold Items
- [x] General Sales Report
- [x] Weekly Sales Report
- [ ] Sale Search
- [ ] Bulk Customers
- [ ] JSON Marshal/Unmarshal
- [ ] JSON Query
- [ ] CTE Example

### Postgres 

- [ ] Schema.sql
- [ ] Queries.sql
- [ ] Initial Data

### Bob (postgres)

- [ ] Everything

### sqlc (postgres)

- [ ] Everything

### Tests (postgres)

- [ ] Everything

### Website

- [x] Setup everything
- [x] Theme
- [x] Markdown Rendering
- [x] Source Code Highlighting
- [x] Landing Page
- [x] Libraries Page
- [ ] Better Library Infoboxes
- [ ] Library Search and Sorting
- [x] Library Report Cards
- [x] Criteria Page
- [ ] Criteria Page Mobile View (source code wide issue)
- [ ] Samples Page
- [ ] Benchmark Page
- [x] Feature Matrix Page
- [ ] Feature Matrix Improvements
- [x] Coming Soon Page
- [ ] Github Pages Setup

### Bun (mysql/postgres)

- [ ] Everything

### Ent (mysql/postgres)

- [ ] Everything

### GORM (mysql/postgres)

- [ ] Everything

### Jet (mysql/postgres)

- [ ] Everything

### SQLBoiler (mysql/postgres)

- [ ] Everything

### sqlx (mysql/postgres)

- [ ] Everything

### Upper DB (mysql/postgres)

- [ ] Everything

### XORM (mysql/postgres)

- [ ] Everything