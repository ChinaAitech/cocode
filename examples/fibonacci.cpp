// 示例3: 斐波那契数列
// 输入区填写一个数字，例如: 10
#include <iostream>
using namespace std;

int main() {
    int n;
    cout << "请输入要生成的斐波那契数列长度: ";
    cin >> n;

    if (n <= 0) {
        cout << "请输入正整数" << endl;
        return 1;
    }

    cout << "斐波那契数列: ";

    long long a = 0, b = 1;
    for (int i = 0; i < n; i++) {
        cout << a << " ";
        long long next = a + b;
        a = b;
        b = next;
    }
    cout << endl;

    return 0;
}
