CREATE USER IF NOT EXISTS 'shuashuati'@'%' IDENTIFIED BY 'shuashuati123';
GRANT ALL PRIVILEGES ON shuashuati.* TO 'shuashuati'@'%';
FLUSH PRIVILEGES;