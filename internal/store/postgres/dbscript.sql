create table books (
	id SERIAL PRIMARY KEY,
	title VARCHAR(255) NOT NULL UNIQUE,
	available_copies INT NOT NULL CHECK (available_copies >= 0)
)

INSERT INTO books (title, available_copies) VALUES ('Alchemist', 3);
INSERT INTO books (title, available_copies) VALUES ('Atomic Habbits', 4);
INSERT INTO books (title, available_copies) VALUES ('Sapiens', 7);
INSERT INTO books (title, available_copies) VALUES ('Mocking Bird', 5);
INSERT INTO books (title, available_copies) VALUES ('Animal Farm', 10);

select * from  books;

create table loans (
	id SERIAL PRIMARY KEY,
	title VARCHAR(256) NOT NULL,
	name_of_borrower VARCHAR(256) NOT NULL,
	loan_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	return_date TIMESTAMP NOT NULL,
	status VARCHAR(100) NOT NULL
)

select * from loans;
