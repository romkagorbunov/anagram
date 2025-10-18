#include <bits/stdc++.h>
#include <locale>


using namespace std;

const int cmax = 256;

struct Node {
    Node* go[cmax];
    bool isTerminal;
};

Node* root;

Node* NewNode() {
    return new Node{};
}

Node* go(Node* current, char c) {
    if (current->go[c] == nullptr) {
        current->go[c] = NewNode();
    }
    return current->go[c];
}

bool tryGet(string s) {
    Node* current = root;
    for (auto c : s) {
        if (c == ' ') {
            if (!current->isTerminal) {
                return false;
            }
            current = root;
        } else {
            if (current->go[c] == nullptr) {
                return false;
            }
            current = current->go[c];
        }
    }
    return true;
}

int main() {
    string filename = "russian.txt";
    std::ifstream file(filename); // путь к файлу
    cout << "Ready" << endl;

    if (!file.is_open()) {
        std::cerr << "Не удалось открыть файл!" << std::endl;
        return 1;
    }
    cout << "Ready" << endl;

    vector<string> words;
    std::string line;
    while (std::getline(file, line)) {
        words.push_back(line);
    }
    cout << "Ready" << endl;

    root = NewNode();
    for (auto w : words) {
        Node* current = root;
        cout << "trying " << w << endl;
        for (auto c : w) {
            current = go(current, c);
        }
        current->isTerminal = true;
    }
    cout << "Ready" << endl;
    string name;
    getline(cin, name);
    string initial = name;
    name = next_permutation(name.begin(), name.end());
    while (name != initial) {
        if (tryGet(name + " ")) {
            cout << name << endl;
        }
    }
    file.close(); 
}
