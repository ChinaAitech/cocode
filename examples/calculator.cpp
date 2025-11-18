// 示例2: 简单计算器
// 输入区填写两个数字，例如: 10 20
#include <iostream>
using namespace std;

int main() {
    int a, b;
    cout << "请输入两个整数: " << endl;
    cin >> a >> b;

    cout << "a + b = " << (a + b) << endl;
    cout << "a - b = " << (a - b) << endl;
    cout << "a * b = " << (a * b) << endl;

    if (b != 0) {
        cout << "a / b = " << (a / b) << endl;
    } else {
        cout << "除数不能为0" << endl;
    }

    return 0;
}
