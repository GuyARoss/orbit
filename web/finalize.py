import os

files = os.listdir('./dist')

for f in files:
    if ".html" in f:
        with open("./dist/" + f, 'r+') as file:
            text = file.read()
            
            text = text.replace("$title", "Orbit - " + f.replace(".html", ""))

            file.seek(0)
            file.write(text)
            file.truncate()
        