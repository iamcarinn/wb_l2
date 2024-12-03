echo -e "apple\tbanana\tcherry" | go run task.go -f 1,3 \n
Output:\n
apple    cherry

echo -e "a;b;c\nd;e;f" | go run task.go -f 2 -d ";" \n
Output: \n
b \n
e \n

echo -e "no-delimiter\napple\tbanana\tcherry" | go run task.go -f 1 -s \n
Output: \n 
apple \n
