import os

def verify_titlekeyword_doesnotexist() -> bool:
    files = os.listdir('./web/dist')

    for f in files:
        if ".html" in f:
            with open("./web/dist/" + f, 'r+') as file:
                text = file.read()
                
                if "$title" in text:
                    return False

    return True

if __name__ == '__main__':
    assert verify_titlekeyword_doesnotexist(), "title keyword was found, please run the 'finalize.py' tool found in the website directory"