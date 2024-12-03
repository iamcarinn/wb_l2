echo -e "apple\tbanana\tcherry" | go run task.go -f 1,3
Output:
apple    cherry

echo -e "a;b;c\nd;e;f" | go run task.go -f 2 -d ";"
Output:
b
e

echo -e "no-delimiter\napple\tbanana\tcherry" | go run task.go -f 1 -s
Output:
apple
