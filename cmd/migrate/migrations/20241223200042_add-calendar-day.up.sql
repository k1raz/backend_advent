CREATE TABLE IF NOT EXISTS calendar_days (
    id INT AUTO_INCREMENT PRIMARY KEY,
    user_id INT,
    day_date DATE,
    title VARCHAR(255),
    image_url VARCHAR(255),
    content TEXT,
    FOREIGN KEY (user_id) REFERENCES users(id)
);