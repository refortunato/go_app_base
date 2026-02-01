-- Database schema for go_app_base

-- Create database if not exists
CREATE DATABASE IF NOT EXISTS go_app_base;
USE go_app_base;

-- Examples table
CREATE TABLE IF NOT EXISTS examples (
    id VARCHAR(36) PRIMARY KEY,
    description TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Insert sample data
INSERT INTO examples (id, description, created_at, updated_at) 
VALUES 
    ('550e8400-e29b-41d4-a716-446655440000', 'First example', NOW(), NOW()),
    ('650e8400-e29b-41d4-a716-446655440001', 'Second example', NOW(), NOW());

-- Products table
CREATE TABLE IF NOT EXISTS products (
    id VARCHAR(40) PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT NOT NULL,
    price DECIMAL(10,2),
    stock INT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;