# Coin-Price-Tele 

## Giới thiệu


## Mục tiêu của dự án


## Tính năng chính


## Thông tin về nhóm

| STT | Tên thành viên               | Vai trò      | Mã số sinh viên |     GitHub      |
|-----|------------------------------|--------------|-----------------|-----------------|
| 1   | Trần Nguyễn Thanh Lâm        | Product Owner| 2211822         | [Github](https://github.com/clgslsm) |
| 2   | Nguyễn Trung Tín             | Developer    | 2212548         | [Github](https://github.com/thanhlam2000) |
| 3   | Thái Thành Duy               | Developer    | 2212213         | [Github](https://github.com/thanhlam2000) |
| 4   | Nguyễn Đăng Khoa             | Developer    | 2212823         | [Github](https://github.com/thanhlam2000) |
| 5   | Lê Thành Đạt                 | Developer    | 2210683         | [Github](https://github.com/Skehjhj) |
| 6   | Nguyễn Hữu Đăng Khoa         | Developer    | 2212571         | [Github](https://github.com/thanhlam2000) |
| 7   | Dương Hoàng Long             | Developer    | 2212571         | [Github](https://github.com/thanhlam2000) |

## Hướng dẫn cài đặt
Sử dụng .env.template và thay thế bằng token tele của bạn

Tải xuống các package: ```go mod tidy```

Để khởi động bot, cd vào thư mục src, chạy ```go run main.go```

Để khởi động mock be server, cd vào thư mục mock-be
```
npm install
npm run start:dev
```

## Hướng dẫn local test webhook

### Cài đặt ngrok bằng choco
Tham khảo: https://ngrok.com/docs/getting-started/?os=windows
Chú ý tải choco, chạy Powershell/Cmp bằng quyền admin


### Local test
Tạo URL ngrok bằng lệnh: ```ngrok http 8443```

Copy đường dẫn tại Forwarding URL có dạng: https://ngrokurl.ngrok-free.app, thêm /webhook và copy vào WEBHOOK_URL của file .env

Chạy bot bằng lệnh ```go run main.go```

## Liên hệ


---

> *README này sẽ được cập nhật thường xuyên để cung cấp thêm thông tin chi tiết và hướng dẫn sử dụng.*
