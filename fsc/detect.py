import sys
if __name__ == "__main__":
    file = open(sys.path[0] + "/result.txt", "r")
    print(''.join(file.readlines()))
    