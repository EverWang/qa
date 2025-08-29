UPDATE users SET password = '$2a$10$3c7UK5Wo71yB/sjCxuMzBOB7SB4fqAyDaQ0zeeffsE1wPUPS5aYf2' WHERE username = 'admin';
SELECT id, username, password, role FROM users WHERE username = 'admin';