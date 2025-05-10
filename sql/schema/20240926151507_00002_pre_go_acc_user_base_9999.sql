-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS pre_go_acc_user_base_9999 (
    user_id INT AUTO_INCREMENT PRIMARY KEY,             -- User ID
    user_account VARCHAR(255) NOT NULL,                 -- User account (used to verify identity)
    user_password VARCHAR(255) NOT NULL,                -- User password
    user_salt VARCHAR(255) NOT NULL,                    -- Salt used for password encryption
    -- isTwoFactorEnabled
    user_login_time TIMESTAMP NULL DEFAULT NULL,        -- Last login time
    user_logout_time TIMESTAMP NULL DEFAULT NULL,       -- Last logout time
    user_login_ip VARCHAR(45) NULL,                     -- Login IP address (45 characters to support IPv6)

    user_created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, -- Record creation time
    user_updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP, -- Record update time

    -- Ensure user_account is unique
    UNIQUE KEY unique_user_account (user_account)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='pre_go_acc_user_base_9999';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS `pre_go_acc_user_base_9999`;
-- +goose StatementEnd
