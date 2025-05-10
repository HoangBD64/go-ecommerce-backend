-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS `pre_go_acc_user_two_factor_9999` (
    `two_factor_id` INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,                    -- Khóa chính tự động tăng
    `user_id` INT UNSIGNED NOT NULL,                                 -- Khóa ngoại liên kết tới bảng người dùng
    `two_factor_auth_type` ENUM('SMS', 'EMAIL', 'APP') NOT NULL,                -- Loại phương thức 2FA (SMS, Email, Ứng dụng như Google Authenticator)
    `two_factor_auth_secret` VARCHAR(255) NOT NULL,                             -- Thông tin bí mật cho 2FA (ví dụ: mã bí mật TOTP cho ứng dụng 2FA)
    `two_factor_phone` VARCHAR(20) NULL,                                 -- Số điện thoại cho 2FA qua SMS (nếu áp dụng)
    `two_factor_email` VARCHAR(255) NULL,                                   -- Địa chỉ email cho 2FA qua Email (nếu áp dụng)
    `two_factor_is_active` BOOLEAN NOT NULL DEFAULT TRUE,                       -- Trạng thái kích hoạt của phương thức 2FA
    `two_factor_created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,                -- Thời điểm tạo phương thức 2FA
    `two_factor_updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP, -- Thời điểm cập nhật phương thức 2FA

    -- Ràng buộc khóa ngoại
    -- FOREIGN KEY (`user_id`) REFERENCES `pre_go_acc_user_base_9999`(`user_id`) ON DELETE CASCADE,

    -- Chỉ mục để tối ưu hóa truy vấn theo `user_id` và `auth_type`
    INDEX `idx_user_id` (`user_id`),
    INDEX `idx_auth_type` (`two_factor_auth_type`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='pre_go_acc_user_two_factor_9999';

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS `pre_go_acc_user_two_factor_9999`;
-- +goose StatementEnd
