-- Membuat tabel User
CREATE TABLE "User" (
    ID bigserial PRIMARY KEY,
    Username VARCHAR(255) NOT NULL,
    Email VARCHAR(255) UNIQUE NOT NULL,
    Password VARCHAR(255) NOT NULL,
    Created_At TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    Updated_At TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Membuat tabel Photo
CREATE TABLE "Photo" (
    ID bigserial PRIMARY KEY,
    Title VARCHAR(255),
    Caption TEXT,
    PhotoUrl VARCHAR(255),
    UserID bigint,
    FOREIGN KEY (UserID) REFERENCES "User"(ID) ON DELETE CASCADE
);

-- Membuat trigger untuk mengatur Updated_At pada tabel User
CREATE OR REPLACE FUNCTION update_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.Updated_At = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_user_updated_at
BEFORE UPDATE ON "User"
FOR EACH ROW
EXECUTE FUNCTION update_updated_at();
