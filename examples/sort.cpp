// 示例4: 数组排序
// 输入格式（输入区）:
// 第一行: 数组长度 n
// 第二行: n 个整数
// 例如:
// 5
// 64 34 25 12 22
#include <iostream>
#include <vector>
#include <algorithm>
using namespace std;

int main() {
    int n;
    cout << "请输入数组长度: ";
    cin >> n;

    vector<int> arr(n);
    cout << "请输入 " << n << " 个整数: " << endl;
    for (int i = 0; i < n; i++) {
        cin >> arr[i];
    }

    cout << "排序前: ";
    for (int x : arr) cout << x << " ";
    cout << endl;

    sort(arr.begin(), arr.end());

    cout << "排序后: ";
    for (int x : arr) cout << x << " ";
    cout << endl;

    return 0;
}
