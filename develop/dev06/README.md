```bash
echo -e "apple\tbanana\tcherry" | go run main.go -f 1,3
# Результат: apple    cherry

echo -e "a;b;c\nd;e;f" | go run main.go -f 2 -d ";"
# Результат: b
#            e

echo -e "no-delimiter\napple\tbanana\tcherry" | go run main.go -f 1 -s
# Результат: apple

