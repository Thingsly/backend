# Data Access Layer (DAL)

## Aggregate Function

Các hàm tổng hợp (aggregate function) giúp phân tích dữ liệu cảm biến theo từng khoảng thời gian (window), rất hữu ích cho dashboard, báo cáo, và cảnh báo trong hệ thống IoT.

### Ý nghĩa các hàm:

- **avg**: Trung bình các giá trị trong mỗi khoảng thời gian.
- **max**: Giá trị lớn nhất trong mỗi khoảng thời gian.
- **min**: Giá trị nhỏ nhất trong mỗi khoảng thời gian.
- **sum**: Tổng các giá trị trong mỗi khoảng thời gian.
- **diff**: Hiệu số giữa giá trị lớn nhất và nhỏ nhất trong mỗi khoảng thời gian.

### Ví dụ minh họa với dữ liệu nhiệt độ

Giả sử bạn có dữ liệu nhiệt độ của một thiết bị trong 1 giờ, lấy mẫu mỗi 10 phút:

| Thời gian (phút) | Nhiệt độ (°C) |
| ---------------- | ------------- |
| 0                | 25            |
| 10               | 27            |
| 20               | 26            |
| 30               | 28            |
| 40               | 29            |
| 50               | 27            |
| 60               | 26            |

Giả sử **window = 30 phút**:

- **avg**:
  - 0-30 phút: (25+27+26+28)/4 = 26.5
  - 30-60 phút: (28+29+27+26)/4 = 27.5
- **max**:
  - 0-30 phút: max(25,27,26,28) = 28
  - 30-60 phút: max(28,29,27,26) = 29
- **min**:
  - 0-30 phút: min(25,27,26,28) = 25
  - 30-60 phút: min(28,29,27,26) = 26
- **sum**:
  - 0-30 phút: 25+27+26+28 = 106
  - 30-60 phút: 28+29+27+26 = 110
- **diff**:
  - 0-30 phút: 28-25 = 3
  - 30-60 phút: 29-26 = 3
